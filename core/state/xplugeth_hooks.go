package state

import (
	"bytes"
	"time"

	"github.com/openrelayxyz/xplugeth"
	"github.com/openrelayxyz/xplugeth/hooks/stateupdates"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state/snapshot"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/metrics"
)

var (
	acctCheckTimer = metrics.NewRegisteredTimer("plugeth/statedb/accounts/checks", nil)
)

type acctChecker struct {
	snap snapshot.Snapshot
	trie Trie
	hasSnap bool
}

type acctByHasher interface {
	GetAccountByHash(common.Hash) (*types.StateAccount, error)
}

func (ac *acctChecker) exists(k common.Hash) bool {
	if hadStorage, ok := ac.snapExists(k); ok {
		return hadStorage
	}
	return ac.trieExists(k)
}

func (ac *acctChecker) snapExists(k common.Hash) (bool, bool) {
	if ac.hasSnap {
		acct, err := ac.snap.AccountRLP(k)
		if err != nil {
			return false, false
		}
		return len(acct) != 0, true
	}
	return false, false
}

func (ac *acctChecker) trieExists(k common.Hash) bool {
	trie, ok := ac.trie.(acctByHasher)
	if !ok {
		log.Warn("Couldn't check trie existence, wrong trie type")
		return true
	}
	acct, _ := trie.GetAccountByHash(k)
	return acct != nil
}

func (ac *acctChecker) updated(k common.Hash, v []byte) bool {
	if updated, ok := ac.snapUpdated(k, v); ok {
		return updated
	}
	return ac.trieUpdated(k, v)
}

func pluginStateUpdate(blockRoot, parentRoot common.Hash, tree *snapshot.Tree, trie Trie, destructs map[common.Hash]struct{}, accounts map[common.Hash][]byte, storage map[common.Hash]map[common.Hash][]byte, codeUpdates map[common.Address]contractCode) {
	var snap snapshot.Snapshot
	if tree != nil {
		snap = tree.Snapshot(parentRoot)
	}
	checker := &acctChecker{snap:snap, trie:trie, hasSnap: snap != nil}

	filteredDestructs := make(map[common.Hash]struct{})
	for k, v := range destructs {
		if _, ok := accounts[k]; ok {
			if !checker.hadStorage(k) {
				// If an account is in both destructs and accounts, that means it was
				// "destroyed" and recreated in the same block. Especially post-cancun,
				// that generally means that an account that had ETH but no code got
				// replaced by an account that had code.
				//
				// If there's data in the accounts map, we only need to process this
				// account if there are storage slots we need to clear out, so we
				// check the account storage for the empty root. If it's empty, we can
				// skip this destruct. We need this check to normalize parallel blocks
				// with serial blocks, because they report destructs differently.
				continue
			}
		} else {
			if !checker.exists(k) {
				// If we have a destruct for an account that didn't exist at the previous
				// block, we don't need to destruct it. This normalizes the use case of
				// CREATE and CREATE2 deployments to empty addresses that revert.
				continue
			}
		}
		filteredDestructs[k] = v
	}

	filteredAccounts := make(map[common.Hash][]byte)
	start := time.Now()
	for k, v := range accounts {
		if checker.updated(k, v) {
			filteredAccounts[k] = v
		}
	}
	acctCheckTimer.UpdateSince(start)

	cleanCodeUpdates := make(map[common.Hash][]byte)
	for _, code := range codeUpdates {
		if len(code.blob) > 0 {
			// Empty code may or may not be included in this list by Geth based on
			// factors I'm not certain of, so to normalize our pending batches,
			// exclude empty blobs.
			cleanCodeUpdates[code.hash] = code.blob
		}
	}

	for _, m := range xplugeth.GetModules[stateupdates.StateUpdatePlugin]() {
		m.StateUpdate(blockRoot, parentRoot, filteredDestructs, filteredAccounts, storage, cleanCodeUpdates)
	}
}

func (ac *acctChecker) hadStorage(k common.Hash) bool {
	if hadStorage, ok := ac.snapHadStorage(k); ok {
		return hadStorage
	}
	return ac.trieHadStorage(k)
}

func (ac *acctChecker) trieHadStorage(k common.Hash) bool {
	trie, ok := ac.trie.(acctByHasher)
	if !ok {
		log.Warn("Couldn't check trie updates, wrong trie type")
		return true
	}
	acct, err := trie.GetAccountByHash(k)
	if err != nil {
		return true
	}
	if acct == nil {
		// No account existed at this hash, so it wouldn't have had storage
		return false
	}
	return acct.Root != types.EmptyRootHash
}

func (ac *acctChecker) snapHadStorage(k common.Hash) (bool, bool) {
	if ac.hasSnap {
		acct, err := ac.snap.Account(k)
		if err != nil {
			return false, false
		}
		if acct == nil {
			// No account existed at this hash, so it wouldn't have had storage
			return false, true
		}
		if len(acct.Root) > 0 && !bytes.Equal(acct.Root, types.EmptyRootHash.Bytes()) {
			return true, true
		}
		return false, true
	}
	return false, false
}

func (ac *acctChecker) snapUpdated(k common.Hash, v []byte) (bool, bool) {
	if ac.hasSnap {
		acct, err := ac.snap.AccountRLP(k)
		if err != nil {
			return false, false
		}
		if len(acct) == 0 {
			return false, false
		}
		return !bytes.Equal(acct, v), true
	} 
	return false, false
}

func (ac *acctChecker) trieUpdated(k common.Hash, v []byte) bool {
	trie, ok := ac.trie.(acctByHasher)
	if !ok {
		log.Warn("Couldn't check trie updates, wrong trie type")
		return true
	}
	oAcct, err := trie.GetAccountByHash(k)
	if err != nil {
		return true
	}
	return !bytes.Equal(v, types.SlimAccountRLP(*oAcct))
}
