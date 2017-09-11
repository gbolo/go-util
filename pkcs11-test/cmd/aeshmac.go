// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

var keyType string
var keyLabel string
var nonEphemeral bool

// aesHmac represents the aes command
var aesHmac = &cobra.Command{
	Use:   "aes-hmac",
	Short: "Creates an AES key then tests HMAC signing with it",
	Long: `Creates an AES key then tests HMAC signing with it`,
	Run: func(cmd *cobra.Command, args []string) {

		setCommonFlagValues()
		PrintPkcs11Settings()
		CreateKey()
	},
}


func init() {
	RootCmd.AddCommand(aesHmac)

	aesHmac.PersistentFlags().StringP( "key-type", "t", "aes", "Type of Key. Supported values: aes, ecdsa")
	aesHmac.PersistentFlags().IntP( "aes-keylength", "k", 32, "Length of AES Key")
	aesHmac.PersistentFlags().StringP( "object-label", "o", "testkeyobject", "Label of Object to use")
	aesHmac.PersistentFlags().Bool( "non-ephemeral",  false, "Sets CKA_TOKEN to true")
	viper.BindPFlag("aes.keylength", aesHmac.PersistentFlags().Lookup("aes-keylength"))
	viper.BindPFlag("create.label", aesHmac.PersistentFlags().Lookup("object-label"))
	viper.BindPFlag("create.type", aesHmac.PersistentFlags().Lookup("key-type"))
	viper.BindPFlag("create.non-ephemeral", aesHmac.PersistentFlags().Lookup("non-ephemeral"))

}


func setCommonFlagValues() {
	pkcs11Lib = viper.GetString("pkcs11.library")
	pkcs11SlotLabel = viper.GetString("pkcs11.label")
	pkcs11SlotPin = viper.GetString("pkcs11.pin")
	keyType = viper.GetString("create.type")
	keyLabel = viper.GetString("create.label")
	nonEphemeral = viper.GetBool("create.non-ephemeral")
}

func displaySettings() {
	fmt.Printf(
		"\nObject Settings:\n - type: %s\n - label: %s\n - nonEphemeral: %t\n",
		keyType,
		keyLabel,
		nonEphemeral,
	)
}

func CreateKey() {

	// Initialize Library
	p, err := p11.InitPkcs11(pkcs11Lib)
	if err != nil {
		ExitWithMessage(fmt.Sprintf("Could not load pkcs11 library: %s", pkcs11Lib), err)
	}
	defer p.Destroy()
	defer p.Finalize()

	// line break for readability
	fmt.Printf("\n")

	// Look for provided slot
	slot, sindex, err := p11.FindSlotByLabel(p, pkcs11SlotLabel)
	if err != nil {
		ExitWithMessage(fmt.Sprintf("PKCS11 provider slot label not found: %s", pkcs11SlotLabel), err)
	}

	// Create session for matching slot
	session, err := p.OpenSession(slot, pkcs11.CKF_SERIAL_SESSION|pkcs11.CKF_RW_SESSION)
	if err != nil {
		ExitWithMessage("Creating session", err)
	}
	defer p.CloseSession(session)

	// Login to access private objects
	fmt.Printf("PKCS11 provider attempting login to slot labeled: %s\n", pkcs11SlotLabel)
	err = p.Login(session, pkcs11.CKU_USER, pkcs11SlotPin)
	if err != nil {
		ExitWithMessage("Login", err)
	}
	defer p.Logout(session)

	// output the settings
	displaySettings()

	if c.CaseInsensitiveEquals(keyType, "aes") {
		CreateAESKey(p, session, sindex)
	} else {
		fmt.Println("non aes keytype not implemented yet.")
	}

}

func CreateAESKey(p *pkcs11.Ctx, session pkcs11.SessionHandle, sindex int) {

	// Get library info
	pkcs11LibInfo, _ := p.GetInfo()

	AesKeyLength := viper.GetInt("aes.keylength")

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

			if keyVerified {
				fmt.Printf("Successfully verified key attributes for key labeled: %s\n", ObjLabel)
			} else {
				ExitWithMessage(fmt.Sprintf("existing key with label: %s has incorrect attribute(s) set", ObjLabel), nil)
			}

		}

		// Test signing with mechanism CKM_SHA256_HMAC
		testMsg := []byte("someRandomString")
		hmac, err := p11.SignHmacSha256(p, session, aesKey, testMsg)
		if err != nil {
			ExitWithMessage("Error signing with CKM_SHA256_HMAC", err)
		}
		fmt.Printf("Successfully tested CKM_SHA256_HMAC on key with label: %s \n HMAC %x\n", ObjLabel, hmac)
	}

	// Exit nicely if we reached this point
	os.Exit(0)

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
