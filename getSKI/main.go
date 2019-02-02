// simple tool to get Subject Key Identifier (SKI) in various formats
// supports pem encoded certificates and pk1/pk8 keys (unencrypted)
package main

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

var (
	publicKey crypto.PublicKey
	// SKISha1 is the raw public key hashed with SHA1 (hex encoded)
	SKISha1 string
	// SKISha256 is the raw public key hashed with SHA256 (hex encoded)
	SKISha256 string
)

func exitWithErr(i ...interface{}) {
	fmt.Printf("Error: %v\n", i)
	os.Exit(1)
}

func main() {
	filePath := flag.String("pem", "", "path to pem file (supports x509 certs and unencrypted pkcs8/pkcs1 private keys)")
	flag.Parse()

	if *filePath == "" {
		exitWithErr("pem file path was not specified (-pem). See --help")
	}

	fileBytes, err := ioutil.ReadFile(*filePath)
	if err != nil {
		exitWithErr(err)
	}

	pemBlock, _ := pem.Decode(fileBytes)
	if pemBlock == nil {
		exitWithErr("invalid pem file")
	}

	switch pemBlock.Type {
	case "CERTIFICATE":
		fmt.Println("pem encoded CERTIFICATE")
		x509Cert, err := x509.ParseCertificate(pemBlock.Bytes)
		if err != nil {
			exitWithErr(err)
		}
		switch pkey := x509Cert.PublicKey.(type) {
		case *ecdsa.PublicKey:
			publicKey = *pkey
		case *rsa.PublicKey:
			publicKey = *pkey
		default:
			exitWithErr(fmt.Sprintf("unsupported public key type: %T", pkey))
		}
	// PKCS8 key
	case "PRIVATE KEY":
		fmt.Println("pem encoded PKCS8 PRIVATE KEY")
		parsePKCS8PrivateKey(pemBlock.Bytes)
	// PKCS1 key
	case "EC PRIVATE KEY", "RSA PRIVATE KEY":
		fmt.Println("pem encoded PKCS1 PRIVATE KEY")
		parsePKCS1PrivateKey(pemBlock.Bytes)
	default:
		exitWithErr(fmt.Sprintf("unsupported pem file type: %v", pemBlock.Type))
	}

	calculateSKI()
	fmt.Printf(" File Path:  %s\n SKI-sha1:   %s\n SKI-sha256: %s\n", *filePath, SKISha1, SKISha256)
}

func calculateSKI() {
	switch pubKey := publicKey.(type) {
	case ecdsa.PublicKey:
		fmt.Printf("EC Key (%v)\n", pubKey.Curve.Params().Name)
		raw := elliptic.Marshal(pubKey.Curve, pubKey.X, pubKey.Y)
		SKISha1 = fmt.Sprintf("%x", sha1.Sum(raw))
		SKISha256 = fmt.Sprintf("%x", sha256.Sum256(raw))
	case rsa.PublicKey:
		fmt.Printf("RSA Key (%v bits)\n", pubKey.N.BitLen())
		raw := pubKey.N.Bytes()
		SKISha1 = fmt.Sprintf("%x", sha1.Sum(raw))
		SKISha256 = fmt.Sprintf("%x", sha256.Sum256(raw))
	}
}

func parsePKCS1PrivateKey(der []byte) {
	// try RSA
	if key, err := x509.ParsePKCS1PrivateKey(der); err == nil {
		publicKey = key.PublicKey
		return
	}
	// try EC
	if key, err := x509.ParseECPrivateKey(der); err == nil {
		publicKey = key.PublicKey
		return
	}
	// if neither RSA or EC, we should fail
	exitWithErr("PKCS1 key is neither RSA or EC")
}

func parsePKCS8PrivateKey(der []byte) {
	key, err := x509.ParsePKCS8PrivateKey(der)
	if err == nil {
		switch key := key.(type) {
		case *rsa.PrivateKey:
			publicKey = key.PublicKey
			return
		case *ecdsa.PrivateKey:
			publicKey = key.PublicKey
			return
		default:
			exitWithErr("PKCS8 key is neither RSA or EC")
		}
	}

	// err must have a value
	exitWithErr(err)
}

// another way of getting SKI. requires pointers (ie: *ecdsa.publicKey)
// not used for now...
func calculateSKIAlternative() {
	fmt.Printf("pk: %T %v\n", publicKey, publicKey)

	encodedPub, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		exitWithErr("MarshalPKIXPublicKey", err)
	}

	var subPKI subjectPublicKeyInfo
	_, err = asn1.Unmarshal(encodedPub, &subPKI)
	if err != nil {
		exitWithErr(err)
	}

	SKISha1 = fmt.Sprintf("%x", sha1.Sum(subPKI.SubjectPublicKey.Bytes))
	SKISha256 = fmt.Sprintf("%x", sha256.Sum256(subPKI.SubjectPublicKey.Bytes))
}

type subjectPublicKeyInfo struct {
	Algorithm        pkix.AlgorithmIdentifier
	SubjectPublicKey asn1.BitString
}
