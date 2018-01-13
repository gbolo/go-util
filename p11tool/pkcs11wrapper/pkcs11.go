package pkcs11wrapper

import (
	"errors"
	"fmt"
	"os"

	"github.com/miekg/pkcs11"
	"github.com/olekukonko/tablewriter"
	//"encoding/asn1"
	"crypto/elliptic"
	//"encoding/hex"
	//"encoding/asn1"

)

type Pkcs11Library struct {

	Path string
	Info pkcs11.Info

}

type Pkcs11Wrapper struct {

	// Context
	Library Pkcs11Library
	Context *pkcs11.Ctx

	// Session Handler
	SlotLabel string
	Session pkcs11.SessionHandle

	// Optional Slot Login
	SlotPin string

}

type Pkcs11Object struct {

	ObjectHandle pkcs11.ObjectHandle

	// Some human readable attributes
	Count string
	CKA_CLASS string
	CKA_LABEL string
	CKA_ID string

}

// Initialize pkcs11 context
func (p11w *Pkcs11Wrapper) InitContext() (err error) {

	// check if lib file exists
	if _, err = os.Stat(p11w.Library.Path); os.IsNotExist(err) {
		return
	}

	// try to initialize
	p11w.Context = pkcs11.New(p11w.Library.Path)
	err = p11w.Context.Initialize()
	if err != nil {
		return
	}

	// get library info
	p11w.Library.Info, err = p11w.Context.GetInfo()

	return

}

// Initialize Session to slot
func (p11w *Pkcs11Wrapper) InitSession() (err error) {

	// Look for provided slot
	slot, _, err := FindSlotByLabel(p11w.Context, p11w.SlotLabel)
	if err != nil {
		return
	}

	// Create session for matching slot
	p11w.Session, err = p11w.Context.OpenSession(slot, pkcs11.CKF_SERIAL_SESSION|pkcs11.CKF_RW_SESSION)

	return

}

// This will perform a login
func (p11w *Pkcs11Wrapper) Login() (err error) {

	err = p11w.Context.Login(p11w.Session, pkcs11.CKU_USER, p11w.SlotPin)

	return
}

// This should return a list of object handlers and true if more than max
func (p11w *Pkcs11Wrapper) FindObjects(template []*pkcs11.Attribute, max int) (p11ObjHandlers []pkcs11.ObjectHandle, moreThanMax bool, err error) {

	// start the search for object
	err = p11w.Context.FindObjectsInit(
		p11w.Session,
		template,
	)
	if err != nil {
		return
	}

	// continue the search, get object handlers
	p11ObjHandlers, moreThanMax, err = p11w.Context.FindObjects(p11w.Session, max)
	if err != nil {
		return
	}

	// finishes the search
	err = p11w.Context.FindObjectsFinal(p11w.Session)
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

		//fmt.Printf("PKCS11 provider found %d slots\n", len(slots))

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

// List content of slot
func (p11w *Pkcs11Wrapper) ListObjects(template []*pkcs11.Attribute, max int) {

	// do an object search
	objects, _, err := p11w.FindObjects(template, max)

	if err != nil {
		fmt.Println("Could not find any objects:", err)
	} else {

		// prepare table headers
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"COUNT", "CKA_CLASS", "CKA_LABEL", "CKA_ID"})
		table.SetCaption(true, fmt.Sprintf("Total objects found (max %d): %d",max, len(objects)))

		// populate table data
		for i, k := range objects {
			al, err := p11w.Context.GetAttributeValue(
				p11w.Session,
				k,
				[]*pkcs11.Attribute{
					pkcs11.NewAttribute(pkcs11.CKA_LABEL, nil),
					pkcs11.NewAttribute(pkcs11.CKA_ID, nil),
					pkcs11.NewAttribute(pkcs11.CKA_CLASS, nil),
				},
			)

			if err != nil {
				panic(err)
			}
			table.Append(
				[]string{
					fmt.Sprintf("%03d", i+1),
					DecodeCKACLASS(al[2].Value[0]),
					fmt.Sprintf("%s", al[0].Value),
					fmt.Sprintf("%x", al[1].Value),
				},
			)
		}

		// render table
		table.Render()

	}
}

func DecodeCKACLASS(b byte) string {

	switch b {
		case 0:
			return "CKO_DATA"
		case 1:
			return "CKO_CERTIFICATE"
		case 2:
			return "CKO_PUBLIC_KEY"
		case 3:
			return "CKO_PRIVATE_KEY"
		case 4:
			return "CKO_SECRET_KEY"
		default:
			return "UNKNOWN"
	}

}

func (p11w *Pkcs11Wrapper) ImportECKey(file string) (err error) {

	ec := EcdsaKey{}
	err = ec.ImportPrivKeyFromFile(file)
	if err != nil {
		return
	}

	ec.GenSKI()

	marshaledOID, err := GetECParamMarshaled(ec.PrivKey.Params().Name)
	if err != nil {
		return
	}

	// pubkey import
	ecPt := elliptic.Marshal(ec.PubKey.Curve, ec.PubKey.X, ec.PubKey.Y)
	// Add DER encoding for the CKA_EC_POINT
	ecPt = append([]byte{0x04, byte(len(ecPt))}, ecPt...)

	keyTemplate := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_KEY_TYPE, pkcs11.CKK_EC),
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_PUBLIC_KEY),
		pkcs11.NewAttribute(pkcs11.CKA_TOKEN, true),
		pkcs11.NewAttribute(pkcs11.CKA_VERIFY, true),
		pkcs11.NewAttribute(pkcs11.CKA_EC_PARAMS, marshaledOID),

		pkcs11.NewAttribute(pkcs11.CKA_ID, ec.SKI.Sha256Bytes),
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, "BCPUB1"),
		pkcs11.NewAttribute(pkcs11.CKA_EC_POINT, ecPt),
		pkcs11.NewAttribute(pkcs11.CKA_PRIVATE, false),
	}

	_, err = p11w.Context.CreateObject(p11w.Session, keyTemplate)
	if err != nil {
		return
	} else {
		fmt.Printf("Object was imported with CKA_LABEL:%s CKA_ID:%x\n", "BCPUB1", ec.SKI.Sha256Bytes)
	}

	keyTemplate = []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_KEY_TYPE, pkcs11.CKK_EC),
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_PRIVATE_KEY),
		pkcs11.NewAttribute(pkcs11.CKA_TOKEN, true),
		pkcs11.NewAttribute(pkcs11.CKA_PRIVATE, false),
		pkcs11.NewAttribute(pkcs11.CKA_SIGN, true),
		pkcs11.NewAttribute(pkcs11.CKA_EC_PARAMS, marshaledOID),

		pkcs11.NewAttribute(pkcs11.CKA_ID, ec.SKI.Sha256Bytes),
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, "BCPRV1"),
		pkcs11.NewAttribute(pkcs11.CKR_ATTRIBUTE_SENSITIVE, false),
		pkcs11.NewAttribute(pkcs11.CKA_EXTRACTABLE, true),
		pkcs11.NewAttribute(pkcs11.CKA_PRIVATE, false),
		pkcs11.NewAttribute(pkcs11.CKA_VALUE, ec.PrivKey.D.Bytes()),
	}

	_, err = p11w.Context.CreateObject(p11w.Session, keyTemplate)
	if err == nil {
		fmt.Printf("Object was imported with CKA_LABEL:%s CKA_ID:%x\n", "BCPRV1", ec.SKI.Sha256Bytes)
	}
	return

}