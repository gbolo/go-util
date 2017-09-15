package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/miekg/pkcs11"
	"github.com/gbolo/go-util/lib/crypto/p11"
	c "github.com/gbolo/go-util/lib/common"
	"github.com/spf13/viper"
)

// aesHmac represents the aes command
var ecdsaCmd = &cobra.Command{
	Use:   "ecdsa",
	Short: "Creates an EC key object then tests mechanism CKM_ECDSA with it",
	Long: `Creates an EC key object then tests mechanism CKM_ECDSA with it`,
	Run: func(cmd *cobra.Command, args []string) {

		setGlobalFlagValues()
		PrintPkcs11Settings()
		p, session, sindex := LoginPkcs11()
		defer p.Destroy()
		defer p.Finalize()
		defer p.CloseSession(session)
		defer p.Logout(session)

		CreateECDSAKey(p, session, sindex)
	},
}


func init() {
	RootCmd.AddCommand(ecdsaCmd)

	ecdsaCmd.PersistentFlags().StringP( "curve", "k", "P224", "Named Curve to Use. (P224, P256, P384, P521)")
	ecdsaCmd.PersistentFlags().StringP( "object-label", "o", "testkeyobject", "Label of Object to use")
	ecdsaCmd.PersistentFlags().Bool( "non-ephemeral",  false, "Sets CKA_TOKEN to true")
	ecdsaCmd.PersistentFlags().String( "message", "FooBar", "Raw message to sign")
	viper.BindPFlag("ecdsa.curve", ecdsaCmd.PersistentFlags().Lookup("curve"))
	viper.BindPFlag("ecdsa.label", ecdsaCmd.PersistentFlags().Lookup("object-label"))
	viper.BindPFlag("ecdsa.non-ephemeral", ecdsaCmd.PersistentFlags().Lookup("non-ephemeral"))
	viper.BindPFlag("ecdsa.message", ecdsaCmd.PersistentFlags().Lookup("message"))

}

// Prints out the object settings
func displayECDSASettings(keyLabel string, curve string, nonEphemeral bool) {
	fmt.Printf(
		"\nObject Settings:\n - type: %s\n - label: %s\n - curve: %s\n - nonEphemeral: %t\n",
		"ECDSA",
		keyLabel,
		curve,
		nonEphemeral,
	)
}

// Creates an ECDSA key object then tests mechanism CKM_SHA256_HMAC with it
func CreateECDSAKey(p *pkcs11.Ctx, session pkcs11.SessionHandle, sindex int) {

	// Set ecdsa variables
	keyLabel := viper.GetString("ecdsa.label")
	nonEphemeral := viper.GetBool("ecdsa.non-ephemeral")
	curve := viper.GetString("ecdsa.curve")
	messageToSign := viper.GetString("ecdsa.message")

	// output the settings
	displayECDSASettings(keyLabel, curve, nonEphemeral)

	// Get library info
	pkcs11LibInfo, _ := p.GetInfo()

	// pkcs11 object labels to look for
	pkcs11ObjLabels := []string{keyLabel}

	for _, ObjLabel := range pkcs11ObjLabels {

		// line break for readability
		fmt.Printf("\n")

		// Do we have a SINGLE object with this JUST this label?
		oHs, moreThanOne, err := p11.FindObjects(p, session,
			[]*pkcs11.Attribute{pkcs11.NewAttribute(pkcs11.CKA_LABEL, ObjLabel)},
			1,
		)
		if err != nil {
			ExitWithMessage(fmt.Sprintf("finding key with label: %s", ObjLabel), err)
		}

		// If we got more than 1, we should exit with this information!
		if moreThanOne {
			ExitWithMessage(fmt.Sprintf("found more than 1 key matching the label: %s", ObjLabel), nil)
		}

		var ecdsaKey pkcs11.ObjectHandle
		// If there are no keys with this label we should create it...
		if len(oHs) == 0 && !c.CaseInsensitiveContains(pkcs11LibInfo.ManufacturerID, "ncipher") {
			fmt.Printf("Key not found with the label: %s. Attempting to create it...\n", ObjLabel)
			ecdsaKey, _, err = p11.CreateECDSAKeyPair(p, session, ObjLabel, curve, !nonEphemeral)
			if err != nil {
				ExitWithMessage(fmt.Sprintf("Error creating key with label: %s on slot: %s", ObjLabel, pkcs11SlotLabel), err)
			} else {
				fmt.Printf("Successfully created key with label: %s on slot: %s\n", ObjLabel, pkcs11SlotLabel)
			}
		} else if len(oHs) == 0 && c.CaseInsensitiveContains(pkcs11LibInfo.ManufacturerID, "ncipher") {
			ExitWithMessage(
				fmt.Sprintf(
					"Key not found with the label: %s. PKCS11 manufacturer is %s which requires vendor key creation.\nPlease create key with vendors generatekey binary and start again.\nEXAMLE:\n\n /opt/nfast/bin/generatekey -g -s %d pkcs11 protect=softcard recovery=yes type=HMACSHA256 size=256 plainname=%s nvram=no",
					ObjLabel,
					pkcs11LibInfo.ManufacturerID,
					sindex,
					ObjLabel,
				),
				errors.New(fmt.Sprintf("Cannot Create Key for Vendor %s", pkcs11LibInfo.ManufacturerID)),
			)
		} else {
			// If we found a key with this label, lets set it to first (and only) one
			ecdsaKey = oHs[0]

			// We need to verify that our key has the correct pkcs11 attributes
			keyVerified, err := p11.VerifyECDSAKey(p, session, ObjLabel, curve, !nonEphemeral)
			if err != nil {
				ExitWithMessage(fmt.Sprintf("finding key with label: %s", ObjLabel), err)
			}

			if keyVerified {
				fmt.Printf("Successfully verified key attributes for key labeled: %s\n", ObjLabel)
			} else {
				ExitWithMessage(fmt.Sprintf("existing key with label: %s has incorrect attribute(s) set", ObjLabel), nil)
			}

		}

		//TODO: create function for this
		err = p.SignInit(session, []*pkcs11.Mechanism{pkcs11.NewMechanism(pkcs11.CKM_ECDSA, nil)}, ecdsaKey)
		if err != nil {
			ExitWithMessage("Error using CKM_ECDSA with key", err)
		}
		// Test signing with mechanism CKM_ECDSA
		testMsg := []byte(messageToSign)
		sig, err := p.Sign(session, testMsg)
		if err != nil {
			ExitWithMessage("Error signing with key", err)
		}
		fmt.Printf("Successfully tested CKM_ECDSA on key with label: %s \n MESSAGE: %s\n SIGNATURE: %x\n",
			ObjLabel,
			messageToSign,
			sig,
		)
	}

	// Exit nicely if we reached this point
	os.Exit(0)

}


