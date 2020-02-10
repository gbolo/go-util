package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"

	"github.com/spf13/viper"
)

var (
	// PCI compliance as of Jun 30, 2018: anything under TLS 1.1 must be disabled
	// we bump this up to TLS 1.2 so we can support best possible ciphers
	tlsMinVersion = uint16(tls.VersionTLS12)
	// allowed ciphers when in hardened mode
	// disable CBC suites (Lucky13 attack) this means TLS 1.1 can't work (no GCM)
	// only use perfect forward secrecy ciphers
	tlsCiphers = []uint16{
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		// these ciphers require go 1.8+
		tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
		tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
	}
	// EC curve preference when in hardened mode
	// curve reference: http://safecurves.cr.yp.to/
	tlsCurvePreferences = []tls.CurveID{
		// this curve is a non-NIST curve with no NSA influence. Prefer this over all others!
		// this curve required go 1.8+
		tls.X25519,
		// These curves are provided by NIST; prefer in descending order
		tls.CurveP521,
		tls.CurveP384,
		tls.CurveP256,
	}
)

// createTLSConfig creates the required TLS configuration that we need to establish a TLS connection
func createTLSConfig() (tlsConfig *tls.Config) {
	var err error
	tlsConfig = &tls.Config{}

	// use stronger TLS settings
	tlsConfig.InsecureSkipVerify = false
	tlsConfig.MinVersion = tlsMinVersion
	tlsConfig.CipherSuites = tlsCiphers
	tlsConfig.CurvePreferences = tlsCurvePreferences

	// always load system trust store
	tlsConfig.RootCAs, err = x509.SystemCertPool()
	if err != nil {
		log.Errorf("unable to load system trust store: %s", err)
		return
	}

	// attempt to load additional ca certs if specified
	// DO NOT FAIL if unable to. Just log an error
	if viper.GetString("tls.ca_cert") != "" {
		log.Infof("attempting to load additional CA pem file: %s", viper.GetString("tls.ca_cert"))
		var certBytes []byte
		certBytes, err = ioutil.ReadFile(viper.GetString("tls.ca_cert"))
		if err != nil {
			log.Errorf("unable to read CA pem file from disk: %s", err)
			return
		}
		if !tlsConfig.RootCAs.AppendCertsFromPEM(certBytes) {
			err = fmt.Errorf("failed to load CA certificate(s): %s", viper.GetString("tls.ca_cert"))
		}
	}

	// load a client certificate and key if mutual TLS is enabled
	// DO NOT FAIL if unable to. Just log an error
	if viper.GetBool("tls.client_auth_enabled") {
		log.Infof("TLS client authentication is enabled. Loading client cert and key")
		var clientCert tls.Certificate
		clientCert, err = tls.LoadX509KeyPair(viper.GetString("tls.client_cert"), viper.GetString("tls.client_key"))
		// we should fail if unable to load the keypair since the user intended mutual authentication
		if err != nil {
			log.Errorf("unable to load client cert and/or key: %s", err)
			return
		}
		// according to TLS spec (RFC 5246 appendix F.1.1) the certificate message
		// must provide a valid certificate chain leading to an acceptable certificate authority.
		// We will make this optional; the client cert pem file can contain more than one certificate
		tlsConfig.Certificates = []tls.Certificate{clientCert}
	}
	return
}
