package main

import (
	"path/filepath"

	"github.com/openrelayxyz/xplugeth"
	"github.com/openrelayxyz/xplugeth/types"
	"github.com/openrelayxyz/xplugeth/hooks/apis"
	"github.com/openrelayxyz/xplugeth/hooks/initialize"
	
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/log"

	"github.com/urfave/cli/v2"
)

func pluginRegisterBackend(backend types.Backend, ethObj *eth.Ethereum) {
	if err := xplugeth.StoreSingleton[types.Backend](backend); err != nil {
		log.Error("Error storing backend singleton", "err", err.Error())
	}
	xplugeth.StoreSingleton[*eth.Ethereum](ethObj)
}

func pluginRegisterStack(stack *node.Node) {
	xplugeth.StoreSingleton[*node.Node](stack)
}

func pluginInitialize(ctx *cli.Context) {
	if !ctx.IsSet(utils.DisablePluginsFlag.Name) {
		var pluginsConfigDir string
		if ctx.IsSet(utils.PluginsConfigDirFlag.Name) {
			pluginsConfigDir = ctx.String(utils.PluginsConfigDirFlag.Name)
		} else {
			pluginsConfigDir = filepath.Join(ctx.String(utils.DataDirFlag.Name), "pluginsconfig")
		}
		xplugeth.Initialize(pluginsConfigDir)
		log.Info("xplugeth initialized")
	}
}

func pluginInitializeNode() {
	stack, ok := xplugeth.GetSingleton[*node.Node]()
	if !ok {
		panic("*node.Node singleton not set, xplugeth InitializeNode")
	}
	backend, ok := xplugeth.GetSingleton[types.Backend]()
	if !ok {
		panic("types.Backend singleton not set, xplugeth Initializenode")
	}
	for _, init := range xplugeth.GetModules[initialize.Initializer]() {
		init.InitializeNode(stack, backend)
	}
}

func pluginGetAPIs() []rpc.API {
	result := []rpc.API{}

	stack, ok := xplugeth.GetSingleton[*node.Node]()
	if !ok {
		panic("*node.Node singleton not set, xplugeth GetAPIs")
	}
	backend, ok := xplugeth.GetSingleton[types.Backend]()
	if !ok {
		panic("types.Backend singleton not set xplugeth GetAPIs")
	}

	for _, a := range xplugeth.GetModules[apis.GetAPIs]() {
		result = append(result, a.GetAPIs(stack, backend)...)
	}

	return result
}

func pluginOnShutdown() {
	for _, shutdown := range xplugeth.GetModules[initialize.Shutdown]() {
		shutdown.Shutdown()
	}
}

func pluginBlockchain() {
	for _, b := range xplugeth.GetModules[initialize.Blockchain]() {
		b.Blockchain()
	}
}
