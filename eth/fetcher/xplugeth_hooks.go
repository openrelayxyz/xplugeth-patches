package fetcher

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/openrelayxyz/xplugeth"
	"github.com/openrelayxyz/xplugeth/hooks/fetcher"
)

func pluginPeerEval(id string, headers []*types.Header) {
	for _, m := range xplugeth.GetModules[fetcher.PeerEvalPlugin]() {
		m.PeerEval(id, headers)
	}
}
