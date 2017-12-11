package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/miekg/pkcs11"
	"github.com/gbolo/go-util/pkcs11-test/p11"
	"github.com/spf13/viper"
)

// aesImport represents the aes command
var aesImport = &cobra.Command{
	Use:   "aes-import",
	Short: "Import an AES key in HexString format",
	Long: `Import an AES key in HexString format`,
	Run: func(cmd *cobra.Command, args []string) {

		setGlobalFlagValues()
		PrintPkcs11Settings()
		p, session, sindex := LoginPkcs11()
		defer p.Destroy()
		defer p.Finalize()
		defer p.CloseSession(session)
		defer p.Logout(session)

		ImportAESKey(p, session, sindex)
	},
}


func init() {
	RootCmd.AddCommand(aesImport)

	aesImport.PersistentFlags().StringP( "object-label", "o", "testkeyobject", "Label of Object to import")
	aesImport.PersistentFlags().Bool( "non-ephemeral",  false, "Sets CKA_TOKEN to true")
	aesImport.PersistentFlags().StringP( "hex-value", "x", "", "Hex value of key to import")
	viper.BindPFlag("aes.label", aesImport.PersistentFlags().Lookup("object-label"))
	viper.BindPFlag("aes.non-ephemeral", aesImport.PersistentFlags().Lookup("non-ephemeral"))
	viper.BindPFlag("aes.hexvalue", aesImport.PersistentFlags().Lookup("hex-value"))

}

// Creates an AES key object then tests mechanism CKM_SHA256_HMAC with it
func ImportAESKey(p *pkcs11.Ctx, session pkcs11.SessionHandle, sindex int) {

	// Set aes variables
	keyLabel := viper.GetString("aes.label")
	nonEphemeral := viper.GetBool("aes.non-ephemeral")
	keyHexString := viper.GetString("aes.hexvalue")

	// Get library info
	//pkcs11LibInfo, _ := p.GetInfo()

	_, err := p11.ImportAesKey(p, session, keyLabel, !nonEphemeral, keyHexString)
	if err != nil {
		ExitWithMessage(fmt.Sprintf("importing aes key with label: %s", keyLabel), err)
	}

	fmt.Printf("\nSuccessfully imported AES key with label: %s and value %s \n",
		keyLabel,
		keyHexString,
	)

	// Exit nicely if we reached this point
	os.Exit(0)

}


