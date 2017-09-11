package main

import (
	"fmt"
	"os"

	"github.com/gbolo/go-util/pkcs11-test/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}