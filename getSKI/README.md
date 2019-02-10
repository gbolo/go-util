# getSKI

Output the Subject Key Identifier (SKI) for a specified pem encoded x509 certficate or unencrypted pkcs #1/8 pem encoded private key (EC or RSA).
Outputs both the SHA1 and SHA256 hex-encoded SKI.

```
# Install
go get github.com/gbolo/go-util/getSKI

# Usage Help
getSKI --help
Usage of ./getSKI:
  -pem string
    	path to pem file (supports x509 certs and unencrypted pkcs8/pkcs1 private keys)

# Output Example
getSKI -pem /tmp/server_testServerCert.pem
pem encoded CERTIFICATE
key type: EC (P-384)
 File Path:  /tmp/server_testServerCert.pem
 SKI-sha1:   596b81cdf0fd8c37908e40650f6f94f2518ffa15
 SKI-sha256: 8bdf32b4e12ac909128aff390c12041d2e76bbc4e1b11544bcbf2f66fa23592a

```
