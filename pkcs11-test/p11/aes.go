package p11

import (
	"os"

	"github.com/miekg/pkcs11"
	c "github.com/gbolo/go-util/lib/common"
	"fmt"
	"strings"
	"encoding/hex"
)

/* return a set of attributes that we require for our aes key */
func GetAesPkcs11Template(objectLabel string, AesKeyLength int, pkcs11LibInfo pkcs11.Info, ephemeral bool) (AesPkcs11Template []*pkcs11.Attribute) {

	// default CKA_KEY_TYPE
	pkcs11_keytype := pkcs11.NewAttribute(pkcs11.CKA_KEY_TYPE, pkcs11.CKK_AES)
	HR_keytype := "CKK_AES"

	// Overrides env first, then autodetect from vendor
	switch {
	case c.CaseInsensitiveContains(os.Getenv("SECURITY_PROVIDER_CONFIG_KEYTYPE"), "CKK_GENERIC_SECRET"):
		pkcs11_keytype = pkcs11.NewAttribute(pkcs11.CKA_KEY_TYPE, pkcs11.CKK_GENERIC_SECRET)
		HR_keytype = "CKK_GENERIC_SECRET"
	case c.CaseInsensitiveContains(pkcs11LibInfo.ManufacturerID, "softhsm") &&
		pkcs11LibInfo.LibraryVersion.Major > 1 &&
		pkcs11LibInfo.LibraryVersion.Minor > 2:
		// matches softhsm versions greater than 2.2 (scott patched 2.3)
		pkcs11_keytype = pkcs11.NewAttribute(pkcs11.CKA_KEY_TYPE, pkcs11.CKK_GENERIC_SECRET)
		HR_keytype = "CKK_GENERIC_SECRET"
	case c.CaseInsensitiveContains(pkcs11LibInfo.ManufacturerID, "ncipher"):
		pkcs11_keytype = pkcs11.NewAttribute(pkcs11.CKA_KEY_TYPE, pkcs11.CKK_SHA256_HMAC)
		HR_keytype = "CKK_SHA256_HMAC"
	case c.CaseInsensitiveContains(pkcs11LibInfo.ManufacturerID, "SafeNet"):
		pkcs11_keytype = pkcs11.NewAttribute(pkcs11.CKA_KEY_TYPE, pkcs11.CKK_GENERIC_SECRET)
		HR_keytype = "CKK_GENERIC_SECRET"
	}

	// Scott's Reference
	// default template common to all manufactures
	AesPkcs11Template = []*pkcs11.Attribute{
		// common to all
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, objectLabel),      /* Name of Key */
		pkcs11.NewAttribute(pkcs11.CKA_TOKEN, !ephemeral),             /* This key should persist */
		pkcs11.NewAttribute(pkcs11.CKA_VALUE_LEN, AesKeyLength), /* KeyLength */
		pkcs11.NewAttribute(pkcs11.CKA_SIGN, true),
		// vendor specific override
		pkcs11_keytype,
	}


	fmt.Println("PKCS11 Attributes Required:")
	fmt.Println(" - CKA_KEY_TYPE:", HR_keytype)
	fmt.Println(" - CKA_LABEL:", objectLabel)
	fmt.Println(" - CKA_VALUE_LEN:", AesKeyLength)
	fmt.Println(" - CKA_TOKEN:", !ephemeral)
	fmt.Println(" - CKA_SIGN:", true)

	return
}

/* This should verify that our key has the correct attributes */
func VerifyAesKey(p *pkcs11.Ctx, session pkcs11.SessionHandle, oLabel string, AesKeyLength int, ephemeral bool) (verified bool, err error) {

	// get lib info
	pkcs11LibInfo, err := p.GetInfo()
	if err != nil {
		verified = false
		return
	}

	// get the required attributes
	requiredAttributes := GetAesPkcs11Template(oLabel, AesKeyLength, pkcs11LibInfo, ephemeral)

	// Search for objects which have ALL these attributes
	oHs, moreThanOne, err := FindObjects(p, session, requiredAttributes, 1)

	// object is verified if there is exactly 1 match and no errors
	if len(oHs) == 1 && !moreThanOne && err == nil {
		verified = true
	}

	return
}

/* Create an AES key with required template */
func CreateAesKey(p *pkcs11.Ctx, session pkcs11.SessionHandle, objectLabel string, AesKeyLength int, ephemeral bool) (aesKey pkcs11.ObjectHandle, err error) {

	// get lib info
	pkcs11LibInfo, err := p.GetInfo()
	if err != nil {
		return
	}

	// default mech CKM_AES_KEY_GEN
	pkcs11_mech := pkcs11.NewMechanism(pkcs11.CKM_AES_KEY_GEN, nil)

	// Overrides env first, then autodetect from vendor
	switch {
	case c.CaseInsensitiveContains(os.Getenv("SECURITY_PROVIDER_CONFIG_MECH"), "CKM_GENERIC_SECRET_KEY_GEN"):
		pkcs11_mech = pkcs11.NewMechanism(pkcs11.CKM_GENERIC_SECRET_KEY_GEN, nil)
	case c.CaseInsensitiveContains(pkcs11LibInfo.ManufacturerID, "SafeNet"):
		pkcs11_mech = pkcs11.NewMechanism(pkcs11.CKM_GENERIC_SECRET_KEY_GEN, nil)
	}

	// get the required attributes
	requiredAttributes := GetAesPkcs11Template(objectLabel, AesKeyLength, pkcs11LibInfo, ephemeral)

	// generate the aes key
	aesKey, err = p.GenerateKey(
		session,
		[]*pkcs11.Mechanism{
			// vendor specific
			pkcs11_mech,
		},
		requiredAttributes,
	)

	return
}

/* test CKM_SHA256_HMAC signing */
func SignHmacSha256(p *pkcs11.Ctx, session pkcs11.SessionHandle, o pkcs11.ObjectHandle, message []byte) (hmac []byte, err error) {

	// start the signing
	err = p.SignInit(
		session,
		[]*pkcs11.Mechanism{
			pkcs11.NewMechanism(pkcs11.CKM_SHA256_HMAC, nil),
		},
		o,
	)
	if err != nil {
		return
	}

	// do the signing
	hmac, err = p.Sign(session, message)
	if err != nil {
		return
	}

	return
}

func getBytes(input string) ([]byte, bool) {
	var result []byte
	isHex := false
	if strings.HasPrefix(input, "0x") {
		d, err := hex.DecodeString(input[2:])
		if err != nil {
			panic(err)
		}
		result = d
		isHex = true
	} else {
		result = []byte(input)
	}
	return result, isHex
}


func ImportAesKey(p *pkcs11.Ctx, session pkcs11.SessionHandle, objectLabel string, ephemeral bool, hexKey string) (aesKey pkcs11.ObjectHandle, err error) {

	// for now, lets only support softhsm for this operation
	pkcs11LibInfo, err := p.GetInfo()
	if err != nil {
		return
	}

	if !c.CaseInsensitiveContains(pkcs11LibInfo.ManufacturerID, "SoftHSM") {
		err = fmt.Errorf("only SoftHSM is supported for key import")
		return
	}

	// now we can try to import the key
	b, isHex := getBytes(hexKey)
	if !isHex {
		err = fmt.Errorf("value provided for key is not hex: %s", hexKey)
		return
	}
	aesKey, err = p.CreateObject(session,
		[]*pkcs11.Attribute{
			pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_SECRET_KEY),
			pkcs11.NewAttribute(pkcs11.CKA_KEY_TYPE, pkcs11.CKK_AES),
			pkcs11.NewAttribute(pkcs11.CKA_TOKEN, !ephemeral),
			pkcs11.NewAttribute(pkcs11.CKA_LABEL, objectLabel),
			pkcs11.NewAttribute(pkcs11.CKA_SIGN, true),
			//pkcs11.NewAttribute(pkcs11.CKA_VERIFY, true),
			pkcs11.NewAttribute(pkcs11.CKA_VALUE, b),
		})

	return
}