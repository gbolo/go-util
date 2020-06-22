package main

import (
	"flag"
	"os"

	"token-provider/backend"
)

func main() {
	// parse flags
	cfgFileFromFlag := flag.String("config", "", "path to config file")
	outputVersion := flag.Bool("version", false, "prints version then exits")
	flag.Parse()

	// allow config file to be specified via environment variable
	cfgFileFromEnv := os.Getenv(backend.EnvConfigPrefix + "_CONFIG")

	// print version and exit if flag is present
	if *outputVersion {
		backend.PrintVersion()
		os.Exit(0)
	}

	// config file is taken from flag, unless flag is empty
	cfgFile := *cfgFileFromFlag
	if cfgFile == "" {
		cfgFile = cfgFileFromEnv
	}

	// start the backend
	backend.StartBackendDeamon(cfgFile)
}
