package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/ethereum/go-ethereum/log"
)

func isValidPath(path string) bool {
    cleanPath := filepath.Clean(path)

    fileInfo, err := os.Stat(cleanPath)
    if err != nil {
        return false
    }
    return fileInfo.IsDir()
}

func disablePlugins() bool {
    val := os.Getenv("DISABLE_PLUGINS")
    if len(val) > 0 {
    	truthyValues := []string{"true", "1", "yes", "y"}
    	for _, truthy := range truthyValues {
        	if strings.EqualFold(val, truthy) {
            	return true
        	}
    	}
	}
    return false
}


func pluginsConfig() string {
	pluginsConfigEnv := os.Getenv("PLUGINS_CONFIG")

	if pluginsConfigEnv != "" && isValidPath(pluginsConfigEnv) {
		log.Info("plugins config path provided", "path", pluginsConfigEnv)
		return pluginsConfigEnv
	} 

	log.Info("plugins config path not provided or invalid")
	return ""
}