# pkcs11 helper tool

Meant to help import keys for use in fabric.

```
# prepare slot
softhsm2-util --init-token --label someLabel --pin somePin --free --so-pin 1234

# build tool

./p11tool --help
Usage of ./p11tool:
  -action string
    	list,import,generateAndImport (default "list")
  -keyFile string
    	path to key you want to import (default "/some/dir/key.pem")
  -keyType string
    	Type of key (EC,RSA) (default "EC")
  -lib string
    	Location of pkcs11 library
  -pin string
    	Slot PIN (default "98765432")
  -slot string
    	Slot Label (default "ForFabric")

# import ec key
./p11tool -slot someLabel -pin somePin -action import -keyFile contrib/testfiles/key.pem -keyType EC
PKCS11 provider found specified slot label: someLabel (slot: 0, index: 0)
Object was imported with CKA_LABEL:BCPUB1 CKA_ID:018f389d200e48536367f05b99122f355ba33572009bd2b8b521cdbbb717a5b5
Object was imported with CKA_LABEL:BCPRV1 CKA_ID:018f389d200e48536367f05b99122f355ba33572009bd2b8b521cdbbb717a5b5

# import rsa key (not supported by fabric BCCSP)
./p11tool -slot someLabel -pin somePin -action import -keyFile contrib/testfiles/key.rsa.pem -keyType RSA
PKCS11 provider found specified slot label: someLabel (slot: 0, index: 0)
Object was imported with CKA_LABEL:TLSPUBKEY CKA_ID:0344ae0121e025d998f5923174e9e4d69b899144ac79bfdf01c065bd4d99d6cb
Object was imported with CKA_LABEL:TLSPRVKEY CKA_ID:0344ae0121e025d998f5923174e9e4d69b899144ac79bfdf01c065bd4d99d6cb

# list objects
./p11tool -slot someLabel -pin somePin -action list
PKCS11 provider found specified slot label: someLabel (slot: 0, index: 0)
+-------+-----------------+-----------+------------------------------------------------------------------+
| COUNT |    CKA CLASS    | CKA LABEL |                              CKA ID                              |
+-------+-----------------+-----------+------------------------------------------------------------------+
|   001 | CKO_PRIVATE_KEY | BCPRV1    | 018f389d200e48536367f05b99122f355ba33572009bd2b8b521cdbbb717a5b5 |
|   002 | CKO_PUBLIC_KEY  | TLSPUBKEY | 0344ae0121e025d998f5923174e9e4d69b899144ac79bfdf01c065bd4d99d6cb |
|   003 | CKO_PRIVATE_KEY | TLSPRVKEY | 0344ae0121e025d998f5923174e9e4d69b899144ac79bfdf01c065bd4d99d6cb |
|   004 | CKO_PUBLIC_KEY  | BCPUB1    | 018f389d200e48536367f05b99122f355ba33572009bd2b8b521cdbbb717a5b5 |
+-------+-----------------+-----------+------------------------------------------------------------------+
Total objects found (max 50): 4
```