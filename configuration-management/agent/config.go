package main

import (
	"github.com/spf13/viper"
)

// setup viper
func initViper(cfgFile string) {

	// Set some defaults
	viper.SetDefault("log_level", "DEBUG")
	viper.SetDefault("server.bind_address", "0.0.0.0")
	viper.SetDefault("server.bind_port", "8080")
	viper.SetDefault("server.access_log", false)

	// set default config name and paths to look for it
	viper.SetConfigType("yaml")
	viper.SetConfigName("agent")
	viper.AddConfigPath("./testdata")

	// if the user provides a config file in a flag, lets use it
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}

	// If a config file is found, read it in.
	err := viper.ReadInConfig()

	// Kick-off the logging module
	loggingInit(viper.GetString("log_level"))

	if err == nil {
		log.Infof("using config file: %s", viper.ConfigFileUsed())
	} else {
		log.Warningf("no config file found: using environment variables and hard-coded defaults: %v", err)
	}
}
