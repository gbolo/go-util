package backend

import (
	"strings"

	"github.com/spf13/viper"
)

const EnvConfigPrefix = "APPNAME"

// ConfigInit instantiates and validates the configuration options
// optionally it can print out a configuration summary
func ConfigInit(cfgFile string, printConfig bool) {

	// init viper
	initViper(cfgFile)

	// Print config if required
	if printConfig {
		printConfigSummary()
	}

	// Sanity checks
	sanityChecks()

	return
}

// setup viper
func initViper(cfgFile string) {

	// Set some defaults
	viper.SetDefault("log_level", "DEBUG")
	viper.SetDefault("server.bind_address", "127.0.0.1")
	viper.SetDefault("server.bind_port", "8080")
	viper.SetDefault("server.access_log", true)

	// Configuring and pulling overrides from environmental variables
	viper.SetEnvPrefix(EnvConfigPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// set default config name and paths to look for it
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath("./")
	viper.AddConfigPath("./testdata/sampleconfig")

	// if the user provides a config file in a flag, lets use it
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}

	// If a config file is found, read it in.
	err := viper.ReadInConfig()

	// Kick-off the logging module
	loggingInit(viper.GetString("log_level"))

	// output version
	log.Noticef("initialized: %s", getAppInfo().summary())

	// configuration file is mandatory for now...
	if err != nil {
		log.Fatalf("unable to load configuration file: %v", err)
	}
	log.Infof("using config file: %s", viper.ConfigFileUsed())
}

// prints the config options
func printConfigSummary() {

	log.Debugf("Configuration:\n")
	for _, c := range []string{
		"log_level",
		"server.bind_address",
		"server.bind_port",
		"server.tls.enabled",
		"server.access_log",
		"server.compression",
	} {
		log.Debugf("%s: %s\n", c, viper.GetString(c))
	}
}

// checks that the config is correctly defined
func sanityChecks() {
	// these values cannot be empty
	keysThatCannotBeEmpty := []string{}
	for _, key := range keysThatCannotBeEmpty {
		if viper.GetString(key) == "" {
			log.Fatalf("%s cannot be empty", key)
		}
	}
}
