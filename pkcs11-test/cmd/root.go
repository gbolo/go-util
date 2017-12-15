package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
)

var cfgFile string
var pkcs11Lib string
var pkcs11SlotLabel string
var pkcs11SlotPin string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "pkcs11-test",
	Short: "A Simple PKCS11 Utility/Tool for Testing PKCS11.",
	Long: `This tool is meant to be able to create AES and ECDSA objects for testing purposes.
	`,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global Flags
	RootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "optional config file (default is ./pkcs11-config.yaml)")
	RootCmd.PersistentFlags().StringP("library", "m", "", "Location of PKCS11 Library")
	RootCmd.PersistentFlags().StringP("pin", "p", "", "PIN Required for Login to Slot")
	RootCmd.PersistentFlags().StringP("label", "l", "", "Label of Slot to Use")
	viper.BindPFlag("pkcs11.library", RootCmd.PersistentFlags().Lookup("library"))
	viper.BindPFlag("pkcs11.pin", RootCmd.PersistentFlags().Lookup("pin"))
	viper.BindPFlag("pkcs11.label", RootCmd.PersistentFlags().Lookup("label"))

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	// Configuring and pulling overrides from environmental variables
	//viper.SetEnvPrefix("PKCS11")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// from config file
	viper.SetConfigName("pkcs11-config")
	viper.AddConfigPath("./")
	viper.AddConfigPath("$HOME/")

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	// set defaults
	viper.SetDefault("pkcs11.library", "/usr/lib/softhsm/libsofthsm2.so")
}

func PrintPkcs11Settings() {
	fmt.Printf("\nPKCS11 Settings:\n - lib: %s\n - label: %s\n - pin: %s\n\n",
		pkcs11Lib,
		pkcs11SlotLabel,
		pkcs11SlotPin,
	)
}
