package p11

import (
	"encoding/asn1"
	"fmt"

	"github.com/miekg/pkcs11"
)

/* return a set of attributes that we require for our ecdsa keypair */
func GetECDSAPkcs11Template(objectLabel string, namedCurve string, ephemeral bool) (pubKeyTemplate []*pkcs11.Attribute, privKeyTemplate []*pkcs11.Attribute, err error) {

	// get ec params
	ecParam, err := GetECParamMarshaled(namedCurve)
	if err != nil {
		return
	}

	// spec taken from fabric
	pubKeyTemplate = []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_KEY_TYPE, pkcs11.CKK_EC),
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_PUBLIC_KEY),
		pkcs11.NewAttribute(pkcs11.CKA_TOKEN, !ephemeral), /* session only. destroy later */
		pkcs11.NewAttribute(pkcs11.CKA_VERIFY, true),
		pkcs11.NewAttribute(pkcs11.CKA_EC_PARAMS, ecParam),
		pkcs11.NewAttribute(pkcs11.CKA_PRIVATE, false),
		pkcs11.NewAttribute(pkcs11.CKA_ID, []byte(objectLabel)),
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, objectLabel),
	}

	// spec taken from fabric
	privKeyTemplate = []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_KEY_TYPE, pkcs11.CKK_EC),
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_PRIVATE_KEY),
		pkcs11.NewAttribute(pkcs11.CKA_TOKEN, !ephemeral), /* session only. destroy later */
		pkcs11.NewAttribute(pkcs11.CKA_PRIVATE, true),
		pkcs11.NewAttribute(pkcs11.CKA_SIGN, true),
		pkcs11.NewAttribute(pkcs11.CKA_ID, []byte(objectLabel)),
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, objectLabel),
		pkcs11.NewAttribute(pkcs11.CKA_EXTRACTABLE, false),
	}

	fmt.Println("PKCS11 Attributes Required:")
	fmt.Println(" - CKA_KEY_TYPE:", "CKK_EC")
	fmt.Println(" - CKA_LABEL:", objectLabel)
	fmt.Println(" - CKA_EC_PARAMS:", namedCurve)
	fmt.Println(" - CKA_TOKEN:", !ephemeral)
	fmt.Println(" - CKA_SIGN:", true)

	return
}

/* returns value for CKA_EC_PARAMS */
func GetECParamMarshaled(namedCurve string) (ecParamMarshaled []byte, err error) {

	ecParamOID := asn1.ObjectIdentifier{}

	switch namedCurve {
	case "P224":
		ecParamOID = asn1.ObjectIdentifier{1, 3, 132, 0, 33}
	case "P256":
		ecParamOID = asn1.ObjectIdentifier{1, 2, 840, 10045, 3, 1, 7}
	case "P384":
		ecParamOID = asn1.ObjectIdentifier{1, 3, 132, 0, 34}
	case "P521":
		ecParamOID = asn1.ObjectIdentifier{1, 3, 132, 0, 35}
	}

	if len(ecParamOID) == 0 {
		err = fmt.Errorf("Error with curve name: %s", namedCurve)
		return
	}

	ecParamMarshaled, err = asn1.Marshal(ecParamOID)
	return
}

/* Create an ECDSA keypair with required template */
func CreateECDSAKeyPair(p *pkcs11.Ctx, session pkcs11.SessionHandle, objectLabel string, namedCurve string, ephemeral bool) (ecdsaPrivKey pkcs11.ObjectHandle, ecdsaPubKey pkcs11.ObjectHandle, err error) {

	// get the required attributes
	pubKeyTemplate, privKeyTemplate, err := GetECDSAPkcs11Template(objectLabel, namedCurve, ephemeral)

	if err != nil {
		return
	}

	// generate the ecdsa key
	ecdsaPubKey, ecdsaPrivKey, err = p.GenerateKeyPair(session,
		[]*pkcs11.Mechanism{
			pkcs11.NewMechanism(pkcs11.CKM_EC_KEY_PAIR_GEN, nil),
		},
		pubKeyTemplate,
		privKeyTemplate,
	)

	return
}


/* This should verify that our key has the correct attributes */
func VerifyECDSAKey(p *pkcs11.Ctx, session pkcs11.SessionHandle, oLabel string, namedCurve string, ephemeral bool) (verified bool, err error) {

	// get the required attributes for priv key
	_, privKey_requiredAttributes, err := GetECDSAPkcs11Template(oLabel, namedCurve, ephemeral)

	if err != nil {
		return
	}

	// Search for objects which have ALL these attributes
	oHs, moreThanOne, err := FindObjects(p, session, privKey_requiredAttributes, 1)

	// object is verified if there is exactly 1 match and no errors
	if len(oHs) == 1 && !moreThanOne && err == nil {
		verified = true
	}

	return
}