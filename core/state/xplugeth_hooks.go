package state

import (
	"github.com/openrelayxyz/xplugeth"
	"github.com/openrelayxyz/xplugeth/hooks/stateupdates"

	"github.com/ethereum/go-ethereum/common"
)

var xpDestructs = make(map[common.Hash]struct{})

func pluginStateUpdate(blockRoot, parentRoot common.Hash, destructs map[common.Hash]struct{}, accounts map[common.Hash][]byte, storage map[common.Hash]map[common.Hash][]byte, codeUpdates map[common.Address]contractCode) {
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
		m.StateUpdate(blockRoot, parentRoot, destructs, accounts, storage, cleanCodeUpdates)
	}
	clear(xpDestructs)
}
