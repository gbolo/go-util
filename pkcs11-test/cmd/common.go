package cmd

import (
	"fmt"
	"os"

	"github.com/miekg/pkcs11"
	"github.com/gbolo/go-util/pkcs11-test/p11"
	"github.com/spf13/viper"
)

// Sets the global pkcs11 variables
func setGlobalFlagValues() {
	pkcs11Lib = viper.GetString("pkcs11.library")
	pkcs11SlotLabel = viper.GetString("pkcs11.label")
	pkcs11SlotPin = viper.GetString("pkcs11.pin")
}

// Loads the PKCS11 library and performs a login
func LoginPkcs11() (p *pkcs11.Ctx, session pkcs11.SessionHandle, sindex int) {

	// Initialize Library
	p, err := p11.InitPkcs11(pkcs11Lib)
	if err != nil {
		ExitWithMessage(fmt.Sprintf("Could not load pkcs11 library: %s", pkcs11Lib), err)
	}
	//defer p.Destroy()
	//defer p.Finalize()

	// line break for readability
	fmt.Printf("\n")

	// Look for provided slot
	slot, sindex, err := p11.FindSlotByLabel(p, pkcs11SlotLabel)
	if err != nil {
		ExitWithMessage(fmt.Sprintf("PKCS11 provider slot label not found: %s", pkcs11SlotLabel), err)
	}

	// Create session for matching slot
	session, err = p.OpenSession(slot, pkcs11.CKF_SERIAL_SESSION|pkcs11.CKF_RW_SESSION)
	if err != nil {
		ExitWithMessage("Creating session", err)
	}
	//defer p.CloseSession(session)

	// Login to access private objects
	fmt.Printf("PKCS11 provider attempting login to slot labeled: %s\n", pkcs11SlotLabel)
	err = p.Login(session, pkcs11.CKU_USER, pkcs11SlotPin)
	if err != nil {
		ExitWithMessage("Login", err)
	}
	//defer p.Logout(session)

	return

}

/* Exit with message and code 1 */
func ExitWithMessage(message string, err error) {

	if err == nil {
		fmt.Printf("\nFatal Error: %s\n", message)
	} else {
		fmt.Printf("\nFatal Error: %s\n%s\n", message, err)
	}
	os.Exit(1)
}