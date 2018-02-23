package pkcs11wrapper

import (
	"crypto"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"

	"crypto/rsa"
	"crypto/sha512"
	"encoding/pem"
	"io/ioutil"
)

type RsaKey struct {
	PubKey  *rsa.PublicKey
	PrivKey *rsa.PrivateKey
	SKI     SubjectKeyIdentifier
}

// SKI returns the subject key identifier of this key.
func (k *RsaKey) GenSKI() {
	if k.PubKey == nil {
		return
	}

	// Marshall the public key
	raw, err := x509.MarshalPKIXPublicKey(k.PubKey)
	if err != nil {
		return
	}

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

func (k *RsaKey) Generate(bits int) (err error) {

	// generate private key
	k.PrivKey, err = rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return
	}

	// store public key
	k.PubKey = &k.PrivKey.PublicKey

	return
}

func (k *RsaKey) ImportPrivKeyFromFile(file string) (err error) {

	keyFile, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}

	keyBlock, _ := pem.Decode(keyFile)
	k.PrivKey, err = x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
	if err != nil {
		return
	}

	// store public key
	k.PubKey = &k.PrivKey.PublicKey

	return
}

func (k *RsaKey) SignMessage(message string, shaSize int) (signature string, err error) {

	var digest []byte
	var hash crypto.Hash
	switch shaSize {

	case 256:
		d := sha256.Sum256([]byte(message))
		digest = d[:]
		hash = crypto.SHA256

	case 384:
		d := sha512.Sum384([]byte(message))
		digest = d[:]
		hash = crypto.SHA384

	case 512:
		d := sha512.Sum512([]byte(message))
		digest = d[:]
		hash = crypto.SHA512

	default:
		d := sha256.Sum256([]byte(message))
		digest = d[:]
		hash = crypto.SHA256

	}

	// sign the hash
	signatureBytes, err := rsa.SignPKCS1v15(rand.Reader, k.PrivKey, hash, digest)
	if err != nil {
		return
	}

	signature = hex.EncodeToString(signatureBytes)

	return
}
