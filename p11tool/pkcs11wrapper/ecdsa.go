package pkcs11wrapper

import (
	"crypto/elliptic"
	"crypto/sha256"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/x509"
	"io/ioutil"
	"encoding/pem"
	"encoding/hex"
	"crypto/sha1"
	"encoding/asn1"
	"fmt"
)

type EcdsaKey struct {
	PubKey *ecdsa.PublicKey
	PrivKey *ecdsa.PrivateKey
	SKIsha256 string
	SKIsha1 string
	SKIsha256Bytes []byte
	SKIsha1Bytes []byte
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
	k.SKIsha256Bytes = hash.Sum(nil)
	k.SKIsha256 = hex.EncodeToString(hash.Sum(nil))

	hash = sha1.New()
	hash.Write(raw)
	k.SKIsha1Bytes = hash.Sum(nil)
	k.SKIsha1 = hex.EncodeToString(hash.Sum(nil))

	return
}

func (k *EcdsaKey) Generate() (err error) {

	k.PrivKey, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
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

