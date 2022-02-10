package main

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	awssigner "github.com/jwx-go/crypto-signer/aws"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jws"
	"github.com/lestrrat-go/jwx/jwt"
	"strings"
	"time"
)

var keyId = "arn:aws:kms:ca-central-1:392058377485:key/b50a55b8-14f2-4f1f-97cc-2766b787a499"

func handleError(err error, hint string) {
	if err != nil {
		panic(fmt.Sprintf("Error Hint: %s\n Error Message: %v", hint, err))
	}
}

func main() {
	fmt.Printf("kms keyId -> %s\n\n", keyId)

	// create our aws-kms client
	awscfg, err := config.LoadDefaultConfig(
		context.Background(),
	)
	if err != nil {
		handleError(err, "cannot connect to aws")
	}
	kmsClient := kms.NewFromConfig(awscfg)

	// validate the key
	err = validateKmsKey(kmsClient, keyId)
	handleError(err, "key failed validation")

	// generate a JWKS for the public key soi others can validate our signatures
	jwks, kid, err := generatePublicJwks(kmsClient, keyId)
	handleError(err, "failed to generate jwks")
	fmt.Println("Public JWKS:")
	prettyPrint(jwks)

	// create a signer for this key
	signer := awssigner.NewRSA(kms.NewFromConfig(awscfg)).
		WithAlgorithm(types.SigningAlgorithmSpecRsassaPkcs1V15Sha256).
		WithKeyID(keyId).
		WithContext(context.Background())

	// create a jwt
	token := jwt.New()
	token.Set(jwt.AudienceKey, "linuxctl.com")
	token.Set(jwt.IssuerKey, "gbolo")
	token.Set(jwt.IssuedAtKey, time.Now().Unix())
	// we need to include the kid in the header so others can validate our signatures from our public jwks
	headers := jws.NewHeaders()
	headers.Set(jws.KeyIDKey, kid)
	signedToken, err := jwt.Sign(token, jwa.RS256, signer, jwt.WithHeaders(headers))
	handleError(err, "failed to sign token")
	fmt.Printf("\nSigned JWT:\n%s\n", signedToken)
}

// validateKmsKey ensures that this kms key meets our expectations
func validateKmsKey(c *kms.Client, keyArn string) (err error) {
	// ensure that we can access the key
	output, err := c.DescribeKey(context.Background(), &kms.DescribeKeyInput{
		KeyId: &keyArn,
	})
	if err != nil {
		return
	}

	// ensure that the key usage allows signing
	if output.KeyMetadata.KeyUsage != "SIGN_VERIFY" {
		return fmt.Errorf("key usage is missing: SIGN_VERIFY")
	}

	// ensure that the key is enabled
	if output.KeyMetadata.KeyState != "Enabled" {
		return fmt.Errorf("key is not in an ENABLED state")
	}

	// ensure we are using an RSA key (we can add support for EC later)
	if !strings.HasPrefix(string(output.KeyMetadata.KeySpec), "RSA_") {
		return fmt.Errorf("key is not of type RSA")
	}

	return
}

// generate a JWKS for our public key so that external parties can validate our signatures
func generatePublicJwks(c *kms.Client, keyArn string) (pubJwks jwk.Set, kid string, err error) {
	pubOut, err := c.GetPublicKey(context.Background(), &kms.GetPublicKeyInput{
		KeyId: &keyArn,
	})
	if err != nil {
		return
	}

	pubKey, err := x509.ParsePKIXPublicKey(pubOut.PublicKey)
	if err != nil {
		return
	}

	// construct the jwks
	jwKey, err := jwk.New(pubKey)
	if err != nil {
		return
	}
	// the key ID should be computed from this public key (not be random)
	kid = calculateKeyID(pubKey)
	jwKey.Set(jwk.KeyIDKey, kid)
	jwKey.Set(jwk.KeyUsageKey, jwk.ForSignature)
	jwKey.Set(jwk.AlgorithmKey, jwa.RS256)
	pubJwks = jwk.NewSet()
	pubJwks.Add(jwKey)
	return
}

// calculateKeyID will return the sha256 hash sum of the public key
func calculateKeyID(pub interface{}) (keyId string) {
	switch pubKey := pub.(type) {
	case *ecdsa.PublicKey:
		//fmt.Printf("key type: EC (%v)\n", pubKey.Curve.Params().Name)
		raw := elliptic.Marshal(pubKey.Curve, pubKey.X, pubKey.Y)
		keyId = fmt.Sprintf("%x", sha256.Sum256(raw))
	case *rsa.PublicKey:
		//fmt.Printf("key type: RSA Key (%v bits)\n", pubKey.N.BitLen())
		raw := pubKey.N.Bytes()
		keyId = fmt.Sprintf("%x", sha256.Sum256(raw))
	default:
		panic("unsupported key type")
	}
	return
}

// print a struct in human readable form
func prettyPrint(v interface{}) {
	b, _ := json.MarshalIndent(v, "", "  ")
	println(string(b))
}
