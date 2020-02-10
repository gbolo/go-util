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

// configure TLS as defined in configuration
func configureTLS() (tlsConfig tls.Config, err error) {

	if !viper.GetBool("server.tls.enabled") {
		log.Debug("TLS not enabled, skipping TLS config")
		return
	}

	// attempt to load configured cert/key
	log.Info("TLS enabled, loading cert and key")
	log.Debugf("loading TLS cert and key: %s %s", viper.GetString("server.tls.cert_chain"), viper.GetString("server.tls.private_key"))
	cert, err := tls.LoadX509KeyPair(viper.GetString("server.tls.cert_chain"), viper.GetString("server.tls.private_key"))
	if err != nil {
		return
	}

	// configure hardened TLS settings
	tlsConfig.Certificates = []tls.Certificate{cert}
	tlsConfig.MinVersion = tlsMinVersion
	tlsConfig.InsecureSkipVerify = false
	tlsConfig.PreferServerCipherSuites = true
	tlsConfig.CurvePreferences = tlsCurvePreferences
	tlsConfig.CipherSuites = tlsCiphers

	// configure client authentication if enabled
	if viper.GetBool("server.tls.client_auth_enabled") {
		log.Infof("client authentication enabled. Loading CA pem file: %s", viper.GetString("server.tls.client_auth_ca"))
		tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert

		caCertPool := x509.NewCertPool()
		var certBytes []byte
		certBytes, err = ioutil.ReadFile(viper.GetString("server.tls.client_auth_ca"))
		if err != nil {
			log.Errorf("unable to read CA pem file from disk: %s", err)
			return
		}
		if !caCertPool.AppendCertsFromPEM(certBytes) {
			err = fmt.Errorf("failed to load CA certificate(s): %s", viper.GetString("server.tls.client_auth_ca"))
			return
		}
		tlsConfig.ClientCAs = caCertPool
	}

	return
}
