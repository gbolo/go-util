package cmd

import (
	"errors"
	"fmt"
	"os"

	c "github.com/gbolo/go-util/lib/common"
	"github.com/gbolo/go-util/pkcs11-test/p11"
	"github.com/miekg/pkcs11"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// aesHmac represents the aes command
var aesHmac = &cobra.Command{
	Use:   "aes-hmac",
	Short: "Creates an AES key object then tests mechanism CKM_SHA256_HMAC with it",
	Long:  `Creates an AES key object then tests mechanism CKM_SHA256_HMAC with it`,
	Run: func(cmd *cobra.Command, args []string) {

		setGlobalFlagValues()
		PrintPkcs11Settings()
		p, session, sindex := LoginPkcs11()
		defer p.Destroy()
		defer p.Finalize()
		defer p.CloseSession(session)
		defer p.Logout(session)

		CreateAESKey(p, session, sindex)
	},
}

func init() {
	RootCmd.AddCommand(aesHmac)

	aesHmac.PersistentFlags().IntP("aes-keylength", "k", 32, "Length of AES Key")
	aesHmac.PersistentFlags().StringP("object-label", "o", "testkeyobject", "Label of Object to use")
	aesHmac.PersistentFlags().Bool("non-ephemeral", false, "Sets CKA_TOKEN to true")
	aesHmac.PersistentFlags().Bool("skip-verify", false, "Skips verification of pkcs11 object attributes")
	aesHmac.PersistentFlags().String("message", "FooBar", "Raw message to sign")
	viper.BindPFlag("aes.keylength", aesHmac.PersistentFlags().Lookup("aes-keylength"))
	viper.BindPFlag("aes.label", aesHmac.PersistentFlags().Lookup("object-label"))
	viper.BindPFlag("aes.non-ephemeral", aesHmac.PersistentFlags().Lookup("non-ephemeral"))
	viper.BindPFlag("aes.skip-verify", aesHmac.PersistentFlags().Lookup("skip-verify"))
	viper.BindPFlag("aes.message", aesHmac.PersistentFlags().Lookup("message"))

}

// Prints out the object settings
func displayAesSettings(keyLabel string, AesKeyLength int, nonEphemeral bool) {
	fmt.Printf(
		"\nObject Settings:\n - type: %s\n - label: %s\n - length: %d\n - nonEphemeral: %t\n",
		"AES",
		keyLabel,
		AesKeyLength,
		nonEphemeral,
	)
}

// Creates an AES key object then tests mechanism CKM_SHA256_HMAC with it
func CreateAESKey(p *pkcs11.Ctx, session pkcs11.SessionHandle, sindex int) {

	// Set aes variables
	keyLabel := viper.GetString("aes.label")
	nonEphemeral := viper.GetBool("aes.non-ephemeral")
	skipVerify := viper.GetBool("aes.skip-verify")
	AesKeyLength := viper.GetInt("aes.keylength")
	messageToSign := viper.GetString("aes.message")

	// output the settings
	displayAesSettings(keyLabel, AesKeyLength, nonEphemeral)

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

		var aesKey pkcs11.ObjectHandle
		// If there are no keys with this label we should create it...
		if len(oHs) == 0 && !c.CaseInsensitiveContains(pkcs11LibInfo.ManufacturerID, "ncipher") {
			fmt.Printf("Key not found with the label: %s. Attempting to create it...\n", ObjLabel)
			aesKey, err = p11.CreateAesKey(p, session, ObjLabel, AesKeyLength, !nonEphemeral)
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
			aesKey = oHs[0]

			// We need to verify that our key has the correct pkcs11 attributes
			keyVerified, err := p11.VerifyAesKey(p, session, ObjLabel, AesKeyLength, !nonEphemeral)
			if err != nil {
				ExitWithMessage(fmt.Sprintf("finding key with label: %s", ObjLabel), err)
			}

			if skipVerify {
				fmt.Println("!! Skipping verification of pkcs11 attributes !!")
			} else if keyVerified {
				fmt.Printf("Successfully verified key attributes for key labeled: %s\n", ObjLabel)
			} else {
				ExitWithMessage(fmt.Sprintf("existing key with label: %s has incorrect attribute(s) set", ObjLabel), nil)
			}

		}

		// Test signing with mechanism CKM_SHA256_HMAC
		testMsg := []byte(messageToSign)
		hmac, err := p11.SignHmacSha256(p, session, aesKey, testMsg)
		if err != nil {
			ExitWithMessage("Error signing with CKM_SHA256_HMAC", err)
		}
		fmt.Printf("Successfully tested CKM_SHA256_HMAC on key with label: %s \n MESSAGE: %s\n HMAC: %x\n",
			ObjLabel,
			messageToSign,
			hmac,
		)
	}

	// Exit nicely if we reached this point
	os.Exit(0)

}
