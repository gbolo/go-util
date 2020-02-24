fabric-ca-reenroll-request
==========================
simple wrapper to capture a `REENROLL` request made by the `fabric-ca-client` for use when executer cannot reach the fabric-ca server

## Example
```
Usage of ./fabric-ca-reenroll-request:
-b string
    path to fabric-ca-client binary (default "./fabric-ca-client")
-d	enable debug
-label string
    name of pkcs11 label/slot
-lib string
    path to pkcs11 library
-m string
    base folder of msp directory (default "./testdata")
-n string
    CA instance name
-p string
    CA profile to use
-pin string
    pin for pkcs11 label/slot
-pkcs11
    enable pkcs11
```

Successful Usage
```
./fabric-ca-reenroll-request -b /tmp/fabric-ca-client-v1.4.3
using bin: /tmp/fabric-ca-client-v1.4.3
fabric-ca-client:
 Version: 1.4.3
 Go version: go1.11.5
 OS/Arch: linux/amd64

Starting HTTP server...
Using MSP path: ./testdata/msp
Expecting Request...

POST /reenroll HTTP/1.1
Host: 127.0.0.1:21700
Accept-Encoding: gzip
Authorization: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUNVVENDQWZpZ0F3SUJBZ0lVTmMvMXN4aVlxZm5XY2xZRVZyOGJLZi90QmFBd0NnWUlLb1pJemowRUF3SXcKZnpFTE1Ba0dBMVVFQmhNQ1ZWTXhFekFSQmdOVkJBZ1RDa05oYkdsbWIzSnVhV0V4RmpBVUJnTlZCQWNURFZOaApiaUJHY21GdVkybHpZMjh4SHpBZEJnTlZCQW9URmtsdWRHVnlibVYwSUZkcFpHZGxkSE1zSUVsdVl5NHhEREFLCkJnTlZCQXNUQTFkWFZ6RVVNQklHQTFVRUF4TUxaWGhoYlhCc1pTNWpiMjB3SGhjTk1Ua3dPVEV3TVRrek5qQXcKV2hjTk1qQXdPVEE1TVRrME1UQXdXakJkTVFzd0NRWURWUVFHRXdKVlV6RVhNQlVHQTFVRUNCTU9UbTl5ZEdnZwpRMkZ5YjJ4cGJtRXhGREFTQmdOVkJBb1RDMGg1Y0dWeWJHVmtaMlZ5TVE4d0RRWURWUVFMRXdaamJHbGxiblF4CkRqQU1CZ05WQkFNVEJXRmtiV2x1TUZrd0V3WUhLb1pJemowQ0FRWUlLb1pJemowREFRY0RRZ0FFVTlXYlNXcTAKZ3VDVzNQVk9FeXBOKzZJSk96YVJsS3lOOHNaV3hkK25Ec2ExZG1JbGU0dURHcWFzamkvbXlGQ1dlWnNDNE8wMwpOWDVSV2U5N3VsYWNpS04wTUhJd0RnWURWUjBQQVFIL0JBUURBZ2VBTUF3R0ExVWRFd0VCL3dRQ01BQXdIUVlEClZSME9CQllFRkZhenBxbWx0NG56b3Y4Y1QrQmtZUkN0STJRc01COEdBMVVkSXdRWU1CYUFGQmRuUWoycW5vSS8KeE1VZG4xdkRtZEcxbkVnUU1CSUdBMVVkRVFRTE1BbUNCMkpwWjJKaGJtY3dDZ1lJS29aSXpqMEVBd0lEUndBdwpSQUlnQ1g2b2J4WXU0RGhBTmVpK05xMnZ2K0J4a00xQjhGMktOZUxjcCt1c0lmY0NJR1ZDMnhxcGRrcS9reHFvClQyUlF5bExQWVE4cDBVYlJLR2JsNXh0akh0cVUKLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=.MEUCIQDVzSEQ4z8WuIROQHEHej7MCiM4vKcHt+D+/i1A8CZlYQIgdrBceGGUzx3+P+kMQqxOf8pPodinuPTpTJ2HgJeWPQk=
Content-Length: 683
User-Agent: Go-http-client/1.1

{"hosts":["bigbang"],"certificate_request":"-----BEGIN CERTIFICATE REQUEST-----\nMIIBPDCB5AIBADBdMQswCQYDVQQGEwJVUzEXMBUGA1UECBMOTm9ydGggQ2Fyb2xp\nbmExFDASBgNVBAoTC0h5cGVybGVkZ2VyMQ8wDQYDVQQLEwZGYWJyaWMxDjAMBgNV\nBAMTBWFkbWluMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAERrNWwZMZnskX9K4e\nSk1JjnfGGn/CoFGvTdobbP0wdfg3LgOAljhNEr+BTR5EKIEaw6L/0n/ij8C12jvG\nXhb496AlMCMGCSqGSIb3DQEJDjEWMBQwEgYDVR0RBAswCYIHYmlnYmFuZzAKBggq\nhkjOPQQDAgNHADBEAiBEPbeVv42uB9W/vRSt0iqheBhhrFzG9rymFeySRR70FAIg\nVMgxJI/dxhIUdjNWGL6IV/YEVApEerVragdXpDkMJfQ=\n-----END CERTIFICATE REQUEST-----\n","profile":"","crl_override":"","label":"","NotBefore":"0001-01-01T00:00:00Z","NotAfter":"0001-01-01T00:00:00Z","CAName":""}
```
