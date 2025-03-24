package state

import (
	"sync"

	"github.com/openrelayxyz/xplugeth"
	"github.com/openrelayxyz/xplugeth/hooks/stateupdates"

	"github.com/ethereum/go-ethereum/common"
)

var rootMap sync.Map

func writeToRootMap(root common.Hash, deletes map[common.Hash]*accountDelete) {
	destructs := make(map[common.Hash]struct{})
	for addrHash, _ := range deletes {
		destructs[addrHash] = struct{}{}
	}
	rootMap.Store(root, destructs)
}

func pluginStateUpdate(blockRoot, parentRoot common.Hash, accounts map[common.Hash][]byte, storage map[common.Hash]map[common.Hash][]byte, codeUpdates map[common.Address]contractCode) {
	cleanCodeUpdates := make(map[common.Hash][]byte)
	for _, code := range codeUpdates {
		if len(code.blob) > 0 {
			// Empty code may or may not be included in this list by Geth based on
			// factors I'm not certain of, so to normalize our pending batches,
			// exclude empty blobs.
			cleanCodeUpdates[code.hash] = code.blob
		}
	}

	destructs := make(map[common.Hash]struct{})
	if d, ok := rootMap.Load(blockRoot); ok {
		destructs = d.(map[common.Hash]struct{})
		rootMap.Delete(blockRoot)
	}

	accountsCopy := make(map[common.Hash][]byte)
	for k, v := range accounts {
		accountsCopy[k] = v
	}

	for acct := range destructs {
		if _, present := accountsCopy[acct]; present {
			delete(accountsCopy, acct)
		}
	}

	for _, m := range xplugeth.GetModules[stateupdates.StateUpdatePlugin]() {
		m.StateUpdate(blockRoot, parentRoot, destructs, accountsCopy, storage, cleanCodeUpdates)
	}
}
