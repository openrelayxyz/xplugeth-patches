package core

import(
	"testing"
	
	"github.com/ethereum/go-ethereum/common"

	"github.com/openrelayxyz/xplugeth"
	"github.com/openrelayxyz/xplugeth/hooks/triecommit"
)

func init() {
	xplugeth.RegisterModule[hashdbTestModule]("hashdbTestModule")
}

type hashdbTestModule struct {}

func TestHashDBInjections(t *testing.T) {
	xplugeth.Initialize("")

	TestHeaderVerification(t)

	if len(injections) > 0 {
		var uncalledInjections []string
		for k, _ := range injections {
			uncalledInjections = append(uncalledInjections, k)
		}
		t.Fatalf("injections not called %v", uncalledInjections)
	} 
}

func (*hashdbTestModule) PreTrieCommit(node common.Hash) {
	delete(injections, "PreTrieCommit")
}

func (*hashdbTestModule) PostTrieCommit(node common.Hash) {
	delete(injections, "PostTrieCommit")
}

var injections map[string]struct{} = map[string]struct{}{
	"PreTrieCommit":struct{}{},
	"PostTrieCommit":struct{}{},
}

var (
	_ triecommit.PreTrieCommit = (*hashdbTestModule)(nil)
	_ triecommit.PostTrieCommit = (*hashdbTestModule)(nil)
)