package state

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/openrelayxyz/xplugeth"
)

// This function is being brought over from foundation to enable our producer plugin to work agnostically across networks.
func (s *StateDB) GetTrie() Trie {
	return s.trie
}

type stateUpdatePlugin interface {
	StateUpdate(blockRoot, parentRoot common.Hash, destructs map[common.Hash]struct{}, accounts map[common.Hash][]byte, storage map[common.Hash]map[common.Hash][]byte, codeUpdates map[common.Hash][]byte)
}

func init() {
	xplugeth.RegisterHook[stateUpdatePlugin]()

}

func pluginStateUpdate(blockRoot, parentRoot common.Hash, destructs map[common.Hash]struct{}, accounts map[common.Hash][]byte, storage map[common.Hash]map[common.Hash][]byte, codeUpdates map[common.Hash][]byte) {
	for _, m := range xplugeth.GetModules[stateUpdatePlugin]() {
		m.StateUpdate(blockRoot, parentRoot, destructs, accounts, storage, codeUpdates)
	}
}
