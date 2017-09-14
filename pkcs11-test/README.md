# PKCS11 Test Utility

A simple test utility that can:

  - Create an AES key object then test mechanism CKM_SHA256_HMAC with it
  - Create an EC key object then test mechanism CKM_ECDSA with it

## Installation

  * Install [Go 1.7+](https://golang.org/dl/)
  * **Install libtool (or libltdl7)**
  * run `go get -u github.com/gbolo/go-util/pkcs11-test`
  * run `pkcs11-test --help`

## SoftHSM2 Setup (testing pkcs11)

  * Install softhsm2
  * Initialize a slot: `softhsm2-util --init-token --label someLabel --pin somePin --free --so-pin 1234`

## Configuration
Configuration can be applied via file, ENV vars, and CLI arguments. Override order:
    1. Configuration file (default ./pkcs11-config.yaml)
    2. Environment Variables
    3. CLI Arguments (use --help for list)

### Configuration File
default ./pkcs11-config.yaml

```
pkcs11:
  library: /usr/lib/softhsm/libsofthsm2.so
  label: DlbpAdapter
  pin: securekey

aes:
  keylength: 32
  message: "Some Important Message"
  non-ephemeral: false
  label: aes_testkey_01

ecdsa:
  curve: P256
  message: "Some Important Message"
  non-ephemeral: false
  label: ec_testkey_01
```

### Environment Variables
the above configuration file can have any of it's values overridden by environment variables that are all in CAPS and with underscores (_) between maps.
Example:

```
export PKCS11_LIBRARY=/usr/lib/softhsm/libsofthsm2.so
export PKCS11_LABEL=somelabel
export PKCS11_PIN=somepin
export AES_KEYLENGTH=16
export AES_LABEL=test_aeskey007
...
```

### CLI ARGUMENTS
cli arguments can be seen with `--help` option

## COMMANDS
the following commands are supported:

### AES-HMAC
```
pkcs11-test aes-hmac --help
Creates an AES key object then tests mechanism CKM_SHA256_HMAC with it

Usage:
  pkcs11-test aes-hmac [flags]

Flags:
  -k, --aes-keylength int     Length of AES Key (default 32)
  -h, --help                  help for aes-hmac
      --message string        Raw message to sign (default "FooBar")
      --non-ephemeral         Sets CKA_TOKEN to true
  -o, --object-label string   Label of Object to use (default "testkeyobject")

Global Flags:
  -c, --config string    optional config file (default is ./pkcs11-config.yaml)
  -l, --label string     Label of Slot to Use
  -m, --library string   Location of PKCS11 Library
  -p, --pin string       PIN Required for Login to Slot
```

### ECDSA
```
pkcs11-test ecdsa --help
Creates an EC key object then tests mechanism CKM_ECDSA with it

Usage:
  pkcs11-test ecdsa [flags]

Flags:
  -k, --curve string          Named Curve to Use. (P224, P256, P384, P521) (default "P224")
  -h, --help                  help for ecdsa
      --message string        Raw message to sign (default "FooBar")
      --non-ephemeral         Sets CKA_TOKEN to true
  -o, --object-label string   Label of Object to use (default "testkeyobject")

Global Flags:
  -c, --config string    optional config file (default is ./pkcs11-config.yaml)
  -l, --label string     Label of Slot to Use
  -m, --library string   Location of PKCS11 Library
  -p, --pin string       PIN Required for Login to Slot
```

# Example Usage

**test AES+HMAC signing:**

```
pkcs11-test aes-hmac
Using config file: /opt/gopath/src/github.com/gbolo/go-util/pkcs11-test/pkcs11-config.yaml

PKCS11 Settings:
 - lib: /usr/lib/softhsm/libsofthsm2.so
 - label: someLabel
 - pin: somePin

Using PKCS11 provider: /usr/lib/softhsm/libsofthsm2.so
 - Manufacturer: SoftHSM
 - Description: Implementation of PKCS11
 - Lib Version: 2.2
 - Cryptoki version: 2.30

PKCS11 provider found 4 slots
PKCS11 provider found specified slot label: someLabel (slot: 220785084, index: 0)
PKCS11 provider attempting login to slot labeled: someLabel

Object Settings:
 - type: AES
 - label: aes_testkey_01
 - length: 32
 - nonEphemeral: false

Key not found with the label: aes_testkey_01. Attempting to create it...
PKCS11 Attributes Required:
 - CKA_KEY_TYPE: CKK_AES
 - CKA_LABEL: aes_testkey_01
 - CKA_VALUE_LEN: 32
 - CKA_TOKEN: false
 - CKA_SIGN: true
Successfully created key with label: aes_testkey_01 on slot: someLabel
Successfully tested CKM_SHA256_HMAC on key with label: aes_testkey_01
 MESSAGE: Some Important Message
 HMAC: ef3c36b16b52976cabe347831645a498e6aea3ebab168ecff727e2b0afd95471
```

**test ECDSA signing:**

```
pkcs11-test ecdsa
Using config file: /opt/gopath/src/github.com/gbolo/go-util/pkcs11-test/pkcs11-config.yaml

PKCS11 Settings:
 - lib: /usr/lib/softhsm/libsofthsm2.so
 - label: someLabel
 - pin: somePin

Using PKCS11 provider: /usr/lib/softhsm/libsofthsm2.so
 - Manufacturer: SoftHSM
 - Description: Implementation of PKCS11
 - Lib Version: 2.2
 - Cryptoki version: 2.30

PKCS11 provider found 4 slots
PKCS11 provider found specified slot label: someLabel (slot: 220785084, index: 0)
PKCS11 provider attempting login to slot labeled: someLabel

Object Settings:
 - type: ECDSA
 - label: ec_testkey_01
 - curve: P256
 - nonEphemeral: false

Key not found with the label: ec_testkey_01. Attempting to create it...
PKCS11 Attributes Required:
 - CKA_KEY_TYPE: CKK_EC
 - CKA_LABEL: ec_testkey_01
 - CKA_EC_PARAMS: P256
 - CKA_TOKEN: false
 - CKA_SIGN: true
Successfully created key with label: ec_testkey_01 on slot: someLabel
Successfully tested CKM_ECDSA on key with label: ec_testkey_01
 MESSAGE: Some Important Message
 SIGNATURE: 80781a9620dbf7f3d1cb400dd0a10b8402ebf2a49a3b9ae645e8cab449207c552d687e8c61dc8627eab6603eee56ec3fc316fb3b23b6ae21149e40ddb86c8c0d
```
