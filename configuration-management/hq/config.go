package main

import (
	"github.com/spf13/viper"
)

// setup viper
func initViper(cfgFile string) {

	// Set some defaults
	viper.SetDefault("log_level", "DEBUG")

	// set default config name and paths to look for it
	viper.SetConfigType("yaml")
	viper.SetConfigName("dsl")
	viper.AddConfigPath("./testdata")

	// if the user provides a config file in a flag, lets use it
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}

	// If a config file is found, read it in.
	err := viper.ReadInConfig()

	// Kick-off the logging module
	loggingInit(viper.GetString("log_level"))

	if err != nil {
		log.Fatalf("unable to load required config file: %v", err)
	}
}

func getTargets() []string {
	return viper.GetStringSlice("targets")
}

func getTasks() (tasks []task) {
	err := viper.UnmarshalKey("tasks", &tasks)
	if err != nil {
		log.Fatalf("unable to read tasks: %v", err)
	}
	return
}
