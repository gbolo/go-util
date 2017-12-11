package p11

import (
	"errors"
	"fmt"
	"os"

	"github.com/miekg/pkcs11"
)

/* This should return a context and print out lib information */
func InitPkcs11(pkcs11LibraryPath string) (p *pkcs11.Ctx, err error) {

	// check if lib file exists
	if _, err = os.Stat(pkcs11LibraryPath); os.IsNotExist(err) {
		return
	}

	// Initialize Library
	fmt.Println("Using PKCS11 provider:", pkcs11LibraryPath)
	p = pkcs11.New(pkcs11LibraryPath)
	err = p.Initialize()
	if err != nil {
		//ExitWithMessage(fmt.Sprintf("Could not load pkcs11 lib:", pkcs11LibraryPath), err)
		return
	}

	// print some info about the library
	pkcs11LibInfo, err := p.GetInfo()
	if err == nil {
		fmt.Printf(
			" - Manufacturer: %s\n - Description: %s\n - Lib Version: %s\n - Cryptoki version: %s\n",
			pkcs11LibInfo.ManufacturerID,
			pkcs11LibInfo.LibraryDescription,
			fmt.Sprintf("%d.%d", pkcs11LibInfo.LibraryVersion.Major, pkcs11LibInfo.LibraryVersion.Minor),
			fmt.Sprintf("%d.%d", pkcs11LibInfo.CryptokiVersion.Major, pkcs11LibInfo.CryptokiVersion.Minor),
		)
	} else {
		fmt.Println("Could not retrieve additional information about this pkcs11 lib")
	}

	return
}

/* This should return a list of object handlers and true if more than max */
func FindObjects(p *pkcs11.Ctx, session pkcs11.SessionHandle, template []*pkcs11.Attribute, max int) (oHs []pkcs11.ObjectHandle, moreThanMax bool, err error) {

	// start the search for object
	err = p.FindObjectsInit(
		session,
		template,
	)
	if err != nil {
		return
	}

	// continue the search, get object handlers
	oHs, moreThanMax, err = p.FindObjects(session, max)
	if err != nil {
		return
	}

	// finishes the search
	err = p.FindObjectsFinal(session)
	if err != nil {
		return
	}

	return
}

/* Return the slotID of token label */
func FindSlotByLabel(p *pkcs11.Ctx, slotLabel string) (slot uint, index int, err error) {

	var slotFound bool

	// Get list of slots
	slots, err := p.GetSlotList(true)
	if err == nil {

		fmt.Printf("PKCS11 provider found %d slots\n", len(slots))

		// Look for matching slot label
		for i, s := range slots {
			tInfo, errGt := p.GetTokenInfo(s)
			if errGt != nil {
				//ExitWithMessage(fmt.Sprintf("getting TokenInfo slot: %d", s), err)
				err = errGt
				return
			}
			if slotLabel == tInfo.Label {
				slotFound = true
				slot = s
				index = i
				fmt.Printf("PKCS11 provider found specified slot label: %s (slot: %d, index: %d)\n", slotLabel, slot, i)
				break
			}
		}
	}

	// set error if slot not found
	if !slotFound {
		err = errors.New(fmt.Sprintf("Could not find slot with label: %s", slotLabel))
	}

	return
}
