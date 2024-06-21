# POC for generating signed JWTs with a key secured by aws kms

Simple proof-of-concept that securely creates signed json web tokens using a
crypto signer backed by AWS KMS.

The JWT key ID is generated by hashing the public key, and is therefor not random.
This allows multiple signers to use the same KMS key without producing different
values for this header.

As a bonus, I also produce a JWKS that can be used by verifiers to validate the
signatures.

**See [code](main.go) for detailed comments**

## Example Output
```
# ensure that you change the KMS key ARN before executing

$ go run main.go
kms keyId -> arn:aws:kms:ca-central-1:392058377485:key/b50a55b8-14f2-4f1f-97cc-2766b787a499

Public JWKS:
{
  "keys": [
    {
      "alg": "RS256",
      "e": "AQAB",
      "kid": "28a2ef0153883879ebd343d997fedf2d7ea56f3282eff4bbaa4a83f5942a2a77",
      "kty": "RSA",
      "n": "xGCzQXx5Q1Q44VxgmO0vAo1C2iMSMCtC3kpdf9bd2mHtDk4-wjXC5say_P_ajWX8M8s_vJS1me2X8A1acVI4IneiEY0Tjl7zH_AlG8AzTBJT-v0STHPPtfDtJ-FlSaWsblNorgV8TaknOR27sTnXQ4rUCdM28W6_dcbKlNEPwlYtxKyxl2moug3fxQiy9kun7aIrez3C0F4F6kFEjPqcPJVohpMqDIxVFk9B98MJLpUwdFMsqUD2TH44PZBMZCGEHroBV8gnVmFY09KMPmbXTjDSYRTGwZJwmq4LsUpywJX1KHJwCXljdScUCq0yyX10t_-r6VIWkx4quF-SIQPtgQ",
      "use": "sig"
    }
  ]
}

Signed JWT:
eyJhbGciOiJSUzI1NiIsImtpZCI6IjI4YTJlZjAxNTM4ODM4NzllYmQzNDNkOTk3ZmVkZjJkN2VhNTZmMzI4MmVmZjRiYmFhNGE4M2Y1OTQyYTJhNzciLCJ0eXAiOiJKV1QifQ.eyJhdWQiOlsibGludXhjdGwuY29tIl0sImlhdCI6MTY0NDQ2MjQxNywiaXNzIjoiZ2JvbG8ifQ.fttKtDev-kFzD1cp0XvFyr8u5nZUfhPaUqKZ_hnh6fawnR2uJLEwCFaJIVc0uO6CNa3SijkyD1hUXVoZQFzbepyG9wsqHet_repHRHWsiBcvAkBxp1SonGkTJ7l2LfpwxYF1JK9b22Xxy8p7YH4O5YjKkyBZIFHcqyLcORliMk8fNUIJwDS5gaV9PONbrm3pWqiWenejz2Iw0wpPQ-ent9jv2ftxU1-Ny7Mt9oDzS5e0NllbdONvXYK-q1pRj8OVaRn_hkAWrAGjXh7fFpqEGtsAHziPM16vZR6Kn97MT7KMgmR5LN-1ZWPD6x_-CITm1N4ZgrT-Muz51GFA2Fze_w
```