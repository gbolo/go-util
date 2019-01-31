// simple tool to get Subject Key Identifier in various formats
package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

var (
	PubKey         *ecdsa.PublicKey
	PrivKey        *ecdsa.PrivateKey
	SKISha1        string
	SKISha1Bytes   []byte
	SKISha256      string
	SKISha256Bytes []byte
)

func main() {

	file := flag.String("file", "/some/dir/cert.pem", "path to pem file")
	flag.Parse()

	err := importPubKeyFromCertFile(*file)
	if err != nil {
		fmt.Println("Error importPubKeyFromCertFile:", err)
		os.Exit(1)
	}

	genSKI()
	fmt.Printf(" file: %s\n SKI-sha1: %s\n SKI-sha256: %s\n", *file, SKISha1, SKISha256)

}

func importPubKeyFromCertFile(file string) (err error) {

	certFile, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}

	certBlock, _ := pem.Decode(certFile)
	x509Cert, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		return
	}

	PubKey = x509Cert.PublicKey.(*ecdsa.PublicKey)

	return
}

func genSKI() {
	if PubKey == nil {
		return
	}

	// Marshall the public key
	raw := elliptic.Marshal(PubKey.Curve, PubKey.X, PubKey.Y)

	// Hash it
	hash := sha256.New()
	hash.Write(raw)
	SKISha256Bytes = hash.Sum(nil)
	SKISha256 = hex.EncodeToString(SKISha256Bytes)

	hash = sha1.New()
	hash.Write(raw)
	SKISha1Bytes = hash.Sum(nil)
	SKISha1 = hex.EncodeToString(SKISha1Bytes)

	return
}
