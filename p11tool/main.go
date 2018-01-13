package main

import (
	pw "github.com/gbolo/go-util/p11tool/pkcs11wrapper"
	//de "github.com/gbolo/go-util/lib/debugging"
	"flag"
	"fmt"
	"github.com/miekg/pkcs11"
	"os"
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
	action := flag.String("action", "list", "list,import")
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

	case "import":
		err = p11w.ImportECKeyFromFile(*keyFile)
		exitWhenError(err)

	case "generate":
		ec := pw.EcdsaKey{}
		ec.Generate("P-256")
		p11w.ImportECKey(ec)

	default:
		p11w.ListObjects(
			[]*pkcs11.Attribute{},
			50,
		)

	}

}
