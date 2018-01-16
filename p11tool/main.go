package main

import (
	pw "github.com/gbolo/go-util/p11tool/pkcs11wrapper"
	//de "github.com/gbolo/go-util/lib/debugging"
	"flag"
	"fmt"
	"github.com/miekg/pkcs11"
	"os"
	"encoding/hex"
)

func exitWhenError(err error) {
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func main() {

	// get flags
	pkcs11Library := flag.String("lib", "/usr/lib/softhsm/libsofthsm2.so", "Location of pkcs11 library")
	slotLabel := flag.String("slot", "ForFabric", "Slot Label")
	slotPin := flag.String("pin", "98765432", "Slot PIN")
	action := flag.String("action", "list", "list,import,generateAndImport")
	keyFile := flag.String("keyFile", "/some/dir/key.pem", "path to key you want to import")

	flag.Parse()

	// initialize pkcs11
	p11w := pw.Pkcs11Wrapper{
		Library: pw.Pkcs11Library{
			Path: *pkcs11Library,
		},
		SlotLabel: *slotLabel,
		SlotPin:   *slotPin,
	}

	err := p11w.InitContext()
	exitWhenError(err)

	err = p11w.InitSession()
	exitWhenError(err)

	err = p11w.Login()
	exitWhenError(err)

	// defer cleanup
	defer p11w.Context.Destroy()
	defer p11w.Context.Finalize()
	defer p11w.Context.CloseSession(p11w.Session)
	defer p11w.Context.Logout(p11w.Session)

	// complete actions
	switch *action {

	case "importEC":
		err = p11w.ImportECKeyFromFile(*keyFile)
		exitWhenError(err)

	case "importRSA":
		err = p11w.ImportRSAKeyFromFile(*keyFile)
		exitWhenError(err)

	case "generateAndImport":
		ec := pw.EcdsaKey{}
		// TODO: fix non working curves (P-521)
		ec.Generate("P-256")
		p11w.ImportECKey(ec)

	case "test":

		message := "Some Test Message"

		// test SW ecdsa sign and verify
		ec := pw.EcdsaKey{}
		ec.ImportPrivKeyFromFile("contrib/testfiles/key.pem")
		sig, err := ec.SignMessage(message)
		exitWhenError(err)
		fmt.Println("Signature:", sig)
		verified := ec.VerifySignature(message, sig)
		fmt.Println("Verified:", verified)


		// test PKCS11 ecdsa sign and verify
		// Find object
		id, err := hex.DecodeString("018f389d200e48536367f05b99122f355ba33572009bd2b8b521cdbbb717a5b5")
		exitWhenError(err)

		o, _, err := p11w.FindObjects([]*pkcs11.Attribute{
			pkcs11.NewAttribute(pkcs11.CKA_KEY_TYPE, pkcs11.CKK_EC),
			pkcs11.NewAttribute(pkcs11.CKA_LABEL, "BCPRV1"),
			pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_PRIVATE_KEY),
			pkcs11.NewAttribute(pkcs11.CKA_ID, id),
			},
			2,
		)

		exitWhenError(err)

		sig, err = p11w.SignMessage(message, o[0])
		exitWhenError(err)
		fmt.Println("pkcs11 Signature:", sig)
		verified = ec.VerifySignature(message, sig)
		fmt.Println("Verified:", verified)

		// test pkcs11 verify
		o, _, err = p11w.FindObjects([]*pkcs11.Attribute{
			pkcs11.NewAttribute(pkcs11.CKA_KEY_TYPE, pkcs11.CKK_EC),
			pkcs11.NewAttribute(pkcs11.CKA_LABEL, "BCPUB1"),
			pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_PUBLIC_KEY),
			pkcs11.NewAttribute(pkcs11.CKA_ID, id),
		},
			2,
		)

		verified, err = p11w.VerifySignature(message, sig, o[0])
		exitWhenError(err)
		fmt.Println("pkcs11 Verified:", verified)

		// derive test
		ec2 := pw.EcdsaKey{}
		ec2.Generate("P-256")

		secret, err := ec.DeriveSharedSecret(ec2.PubKey)
		exitWhenError(err)
		fmt.Printf("shared secret: %x\n", secret)

		secret, err = ec2.DeriveSharedSecret(ec.PubKey)
		exitWhenError(err)
		fmt.Printf("shared secret: %x\n", secret)

	case "testRsa":
		message := "Some Test Message"

		rsa := pw.RsaKey{}
		//rsa.Generate(2048)
		err = rsa.ImportPrivKeyFromFile("contrib/testfiles/key.rsa.pem")
		exitWhenError(err)
		rsa.GenSKI()

		err = p11w.ImportRSAKey(rsa)
		exitWhenError(err)

		sig, err := rsa.SignMessage(message, 256)
		exitWhenError(err)

		fmt.Println("Signature:", sig)

		// test PKCS11 ecdsa sign and verify
		// Find object
		id, err := hex.DecodeString("0344ae0121e025d998f5923174e9e4d69b899144ac79bfdf01c065bd4d99d6cb")
		exitWhenError(err)

		o, _, err := p11w.FindObjects([]*pkcs11.Attribute{
			pkcs11.NewAttribute(pkcs11.CKA_KEY_TYPE, pkcs11.CKK_RSA),
			pkcs11.NewAttribute(pkcs11.CKA_LABEL, "TLSPRVKEY"),
			pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_PRIVATE_KEY),
			pkcs11.NewAttribute(pkcs11.CKA_ID, id),
		},
			2,
		)
		exitWhenError(err)

		sig, err = p11w.SignMessageAdvanced([]byte(message), o[0], pkcs11.NewMechanism(pkcs11.CKM_SHA256_RSA_PKCS, nil))
		exitWhenError(err)

		fmt.Println("pkcs11 Signature:", sig)



	default:
		p11w.ListObjects(
			[]*pkcs11.Attribute{},
			50,
		)

	}

}
