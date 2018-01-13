package pkcs11wrapper

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"io/ioutil"
)

type EcdsaKey struct {
	PubKey  *ecdsa.PublicKey
	PrivKey *ecdsa.PrivateKey
	SKI     SubjectKeyIdentifier
}

type SubjectKeyIdentifier struct {
	Sha1        string
	Sha1Bytes   []byte
	Sha256      string
	Sha256Bytes []byte
}

// SKI returns the subject key identifier of this key.
func (k *EcdsaKey) GenSKI() (ski []byte) {
	if k.PubKey == nil {
		return nil
	}

	// Marshall the public key
	raw := elliptic.Marshal(k.PubKey.Curve, k.PubKey.X, k.PubKey.Y)

	// Hash it
	hash := sha256.New()
	hash.Write(raw)
	k.SKI.Sha256Bytes = hash.Sum(nil)
	k.SKI.Sha256 = hex.EncodeToString(k.SKI.Sha256Bytes)

	hash = sha1.New()
	hash.Write(raw)
	k.SKI.Sha1Bytes = hash.Sum(nil)
	k.SKI.Sha1 = hex.EncodeToString(k.SKI.Sha1Bytes)

	return
}

func (k *EcdsaKey) Generate(namedCurve string) (err error) {

	// generate private key
	switch namedCurve {
	case "P-224":
		k.PrivKey, err = ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
	case "P-256":
		k.PrivKey, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	case "P-384":
		k.PrivKey, err = ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	case "P-521":
		k.PrivKey, err = ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	default:
		k.PrivKey, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	}

	// store public key
	k.PubKey = &k.PrivKey.PublicKey

	return
}

func (k *EcdsaKey) ImportPubKeyFromPubKeyFile(file string) (err error) {
	return
}

func (k *EcdsaKey) ImportPubKeyFromCertFile(file string) (err error) {

	certFile, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}

	certBlock, _ := pem.Decode(certFile)
	x509Cert, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		return
	}

	k.PubKey = x509Cert.PublicKey.(*ecdsa.PublicKey)

	return
}

func (k *EcdsaKey) ImportPrivKeyFromFile(file string) (err error) {

	keyFile, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}

	keyBlock, _ := pem.Decode(keyFile)
	key, err := x509.ParsePKCS8PrivateKey(keyBlock.Bytes)
	if err != nil {
		return
	}

	k.PrivKey = key.(*ecdsa.PrivateKey)
	k.PubKey = &k.PrivKey.PublicKey

	return
}

/* returns value for CKA_EC_PARAMS */
func GetECParamMarshaled(namedCurve string) (ecParamMarshaled []byte, err error) {

	ecParamOID := asn1.ObjectIdentifier{}

	switch namedCurve {
	case "P-224":
		ecParamOID = asn1.ObjectIdentifier{1, 3, 132, 0, 33}
	case "P-256":
		ecParamOID = asn1.ObjectIdentifier{1, 2, 840, 10045, 3, 1, 7}
	case "P-384":
		ecParamOID = asn1.ObjectIdentifier{1, 3, 132, 0, 34}
	case "P-521":
		ecParamOID = asn1.ObjectIdentifier{1, 3, 132, 0, 35}
	}

	if len(ecParamOID) == 0 {
		err = fmt.Errorf("Error with curve name: %s", namedCurve)
		return
	}

	ecParamMarshaled, err = asn1.Marshal(ecParamOID)
	return
}
