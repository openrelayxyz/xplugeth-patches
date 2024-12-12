package fetcher

import (
	"testing"

	gtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/openrelayxyz/xplugeth"
	"github.com/openrelayxyz/xplugeth/hooks/fetcher"
)

type fetcherTestModule struct{}

func init() {
	xplugeth.RegisterModule[fetcherTestModule]("fetcherTestModule")
}

func TestFetcherPkgInjections(t *testing.T) {
	xplugeth.Initialize("")
	testSequentialAnnouncements(t, false)

	if len(injections) > 0 {
		var uncalledInjections []string
		for k, _ := range injections {
			uncalledInjections = append(uncalledInjections, k)
		}
		t.Fatalf("injections not called %v", uncalledInjections)
	}
}

func (p *fetcherTestModule) PeerEval(id string, headers []*gtypes.Header) {
	delete(injections, "PeerEval")
}

var injections map[string]struct{} = map[string]struct{}{
	"PeerEval": struct{}{},
}

var (
	_ fetcher.PeerEvalPlugin = (*fetcherTestModule)(nil)
)
