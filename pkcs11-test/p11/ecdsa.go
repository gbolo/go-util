package p11

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/asn1"
	"fmt"
	"io/ioutil"

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
		pkcs11.NewAttribute(pkcs11.CKA_ID, []byte(objectLabel)),
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, objectLabel),
		// public key should be easily accessed
		pkcs11.NewAttribute(pkcs11.CKA_PRIVATE, false),
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
		// TODO: make these options configurable...
		pkcs11.NewAttribute(pkcs11.CKA_EXTRACTABLE, false),
		// support key derivation by default for now...
		pkcs11.NewAttribute(pkcs11.CKA_DERIVE, true),
		// pkcs11.NewAttribute(pkcs11.CKR_ATTRIBUTE_SENSITIVE, false),
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

/* This should return the public key in PEM format */
func GetPublicKey(p *pkcs11.Ctx, session pkcs11.SessionHandle, objectLabel string) (pubKeyPem string, err error) {

	// search for the public key
	pubKeyTemplate := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_KEY_TYPE, pkcs11.CKK_EC),
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_PUBLIC_KEY),
		pkcs11.NewAttribute(pkcs11.CKA_PRIVATE, false),
		pkcs11.NewAttribute(pkcs11.CKA_ID, []byte(objectLabel)),
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, objectLabel),
	}

	// Search for objects with this template
	oHs, moreThanOne, err := FindObjects(p, session,
		pubKeyTemplate,
		1,
	)
	if err != nil {
		return
	}

	// If we got more than 1, we should exit with this information!
	if moreThanOne {
		err = fmt.Errorf("more than 1 key found")
		return
	}

	// get CKA_VALUE
	fmt.Println("Found a pub key:", oHs[0])
	pubKeyHandle := oHs[0]
	template := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_EC_PARAMS, nil),
		pkcs11.NewAttribute(pkcs11.CKA_EC_POINT, nil),
	}
	pubKeyAttrValues, err := p.GetAttributeValue(session, pubKeyHandle, template)

	if err != nil {
		return
	}

	// according to: http://docs.oasis-open.org/pkcs11/pkcs11-curr/v2.40/csprd02/pkcs11-curr-v2.40-csprd02.html#_Toc387327769
	// DER-encoding of an ANSI X9.62 Parameters value
	fmt.Println("CKA_EC_PARAMS:", pubKeyAttrValues[0].Value)
	// DER-encoding of ANSI X9.62 ECPoint value Q
	fmt.Println("CKA_EC_POINT:", pubKeyAttrValues[1].Value)

	err = ioutil.WriteFile("ecpoint.der", pubKeyAttrValues[1].Value, 0644)
	if err != nil {
		panic(err)
	}

	var ecp []byte
	_, err1 := asn1.Unmarshal(pubKeyAttrValues[1].Value, &ecp)
	if err1 != nil {
		fmt.Printf("Failed to decode ASN.1 encoded CKA_EC_POINT (%s)", err1.Error())
	}
	fmt.Println("ecp:", ecp)

	pubKey, err := getPublic(ecp)
	if err != nil {
		fmt.Printf("Failed to decode public key (%s)", err.Error())
		return
	}

	fmt.Printf("Public key: %#v", pubKey)

	pubKeyPem = ""
	return
}

func getPublic(point []byte) (pub crypto.PublicKey, err error) {
	var ecdsaPub ecdsa.PublicKey

	ecdsaPub.Curve = elliptic.P256()
	pointLenght := ecdsaPub.Curve.Params().BitSize/8*2 + 1
	if len(point) != pointLenght {
		err = fmt.Errorf("CKA_EC_POINT (%d) does not fit used curve (%d)", len(point), pointLenght)
		return
	}
	ecdsaPub.X, ecdsaPub.Y = elliptic.Unmarshal(ecdsaPub.Curve, point[:pointLenght])
	if ecdsaPub.X == nil {
		err = fmt.Errorf("Failed to decode CKA_EC_POINT")
		return
	}
	if !ecdsaPub.Curve.IsOnCurve(ecdsaPub.X, ecdsaPub.Y) {
		err = fmt.Errorf("Public key is not on Curve")
		return
	}

	pub = &ecdsaPub
	return
}
