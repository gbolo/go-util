package oidc

import (
	"crypto/rsa"
	"io/ioutil"
	"oidc-client/util"
	"time"

	"github.com/dgrijalva/jwt-go"
	cv "github.com/nirasan/go-oauth-pkce-code-verifier"
	"github.com/spf13/viper"
	jose "gopkg.in/square/go-jose.v2"
)

var (
	jwtSigningKey *rsa.PrivateKey
	jwks          jose.JSONWebKeySet
	codeVerifier  *cv.CodeVerifier
)

func initJwtSigningKey() {
	// read in file
	keyFile := viper.GetString("jwt.signing_key")
	log.Infof("reading in JWT signing key: %v", keyFile)
	jwtSigningKeyBytes, err := ioutil.ReadFile(keyFile)
	if err != nil {
		log.Fatalf("could not load jwt signing key: %v", err)
	}

	// load it as an RSA key
	jwtSigningKey, err = jwt.ParseRSAPrivateKeyFromPEM(jwtSigningKeyBytes)
	if err != nil {
		log.Fatalf("could not parse jwt signing key: %v", err)
	}

	// init the jwks
	jwks.Keys = []jose.JSONWebKey{
		{
			Key:   jwtSigningKey.Public(),
			KeyID: viper.GetString("jwt.kid"),
			Use:   "sig",
		},
	}

	// init pkce
	codeVerifier, err = cv.CreateCodeVerifier()
	if err != nil {
		log.Fatalf("could not init PKCE code verifier: %v", err)
	}
}

func GetJwks() jose.JSONWebKeySet {
	return jwks
}

// jwt used during initial auth request
func createRequestJWT(state, redirectURL string) (jwtString string) {

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"aud":                   discoveryCache.Issuer,
		"iss":                   viper.GetString("oidc.client_id"),
		"code_challenge_method": "S256",
		"code_challenge":        codeVerifier.CodeChallengeS256(),
		"scope":                 "openid foundation_profile",
		"response_type":         "code",
		"redirect_uri":          redirectURL,
		"state":                 state,
		"iat":                   time.Now().Unix(),
		"ui_locales":            "en-CA",
	})

	// has to match our advertised jwks
	token.Header["kid"] = viper.GetString("jwt.kid")

	// Sign and get the complete encoded token as a string using the secret
	jwtString, _ = token.SignedString(jwtSigningKey)
	return
}

// jwt used during access token request
func createAssertionJWT() (jwtString string) {

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"aud": discoveryCache.TokenEndpoint,
		"sub": viper.GetString("oidc.client_id"),
		"iss": viper.GetString("oidc.client_id"),
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Minute * 10).Unix(),
		// random string
		"jti": util.GenerateRandomString(16),
	})

	// has to match our advertised jwks
	token.Header["kid"] = viper.GetString("jwt.kid")

	// Sign and get the complete encoded token as a string using the secret
	jwtString, _ = token.SignedString(jwtSigningKey)
	return
}
