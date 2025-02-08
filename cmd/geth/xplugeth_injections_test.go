package main

import(
	"testing"
	"flag"
	"fmt"
	"os"
	"strings"
	"strconv"
	"syscall"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/openrelayxyz/xplugeth"
	"github.com/openrelayxyz/xplugeth/hooks/initialize"
	"github.com/openrelayxyz/xplugeth/hooks/apis"
	"github.com/openrelayxyz/xplugeth/types"
)

func init() {
	xplugeth.RegisterModule[gethTestModule]("gethTestModule")
}

func TestGethPkgInjections(t *testing.T) {
	pseudoProtocols := strings.Join(func() (strings []string) {
		for _, s := range ethconfig.Defaults.ProtocolVersions {
			strings = append(strings, strconv.Itoa(int(s)))
		}
		return
	}(), ",")

	testDataDir := "./xp-test-data"
	defer os.RemoveAll(testDataDir)

	flagSet := flag.NewFlagSet("test", flag.ContinueOnError)
	
	flagSet.String("gcmode", "full", "Blockchain garbage collection mode")
	flagSet.String("crypto.kzg", "gokzg", "KZG library implementation to use; gokzg (recommended) or ckzg")
	flagSet.String("eth.protocols",pseudoProtocols, "Sets the Ethereum Protocol versions (first is primary)")
	flagSet.String("datadir", "", "Data directory for the databases and keystore")
	
	if err := flagSet.Parse([]string{"--datadir", testDataDir}); err != nil {
		t.Fatalf("Failed to parse flag set: %v", err)
	}

	app := &cli.App{
		Flags: nodeFlags,
	}
	
	ctx := cli.NewContext(app, flagSet, nil)

	go func() {
		if err := geth(ctx); err != nil {
			t.Fatalf("err calling geth, err %v", err)
		}
	}()

	time.Sleep(2 * time.Second)

	if err := syscall.Kill(syscall.Getpid(), syscall.SIGINT); err != nil {
		t.Fatalf("failed to send SIGINT: %v", err)
	}

	time.Sleep(2 * time.Second)
	
	if len(injections) > 0 {
		var uncalledInjections []string
		for k, _ := range injections {
			uncalledInjections = append(uncalledInjections, k)
		}
		t.Fatalf("test failed, injections not called %v", uncalledInjections)
	} 
}

type gethTestModule struct {}

func (*gethTestModule) InitializeNode(s *node.Node, b types.Backend) {
	if s != nil && b != nil {
		delete(injections, "InitializeNode")
	} else {
		fmt.Printf("get singleton error stack %v, backend %v", s, b)
	}
}

func (*gethTestModule) GetAPIs(*node.Node, types.Backend) []rpc.API {
	delete(injections, "GetAPIs")
	return nil
}

func (*gethTestModule) Blockchain() {
	delete(injections, "BlockChain")
}

func (*gethTestModule) Shutdown() {
	delete(injections, "Shutdown")
}

var injections map[string]struct{} = map[string]struct{}{
	"InitializeNode":struct{}{},
	"GetAPIs":struct{}{},
	"BlockChain":struct{}{},
	"Shutdown":struct{}{},
}


var (
	_ initialize.Initializer = (*gethTestModule)(nil)
	_ initialize.Blockchain = (*gethTestModule)(nil)
	_ initialize.Shutdown = (*gethTestModule)(nil)
	_ apis.GetAPIs = (*gethTestModule)(nil)
)