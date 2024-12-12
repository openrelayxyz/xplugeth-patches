package main

import (
	"os"

	"github.com/openrelayxyz/xplugeth"

	"github.com/ethereum/go-ethereum/internal/cli"
	"github.com/ethereum/go-ethereum/log"
)

func main() {
	//begin xplugeth injection
	if !disablePlugins() {
		xplugeth.Initialize(pluginsConfig())
		log.Info("xplugeth initialized")
	}

	if i, ok := xplugeth.ParseCommands(os.Args); ok {
		os.Exit(cli.Run(append([]string{"server"}, os.Args[1:i]...)))
	} else {
		os.Exit(cli.Run(os.Args[1:]))
	}
	//end xplugeth injection
}
