# OIDC Test Client for use with OIDC DAC Flow

# MOVED LOCATION
This tool has been moved to: https://github.com/gbolo/muggle-oidc

## DISCLAIMER
This application was written to help me understand the details of a specific OIDC flow.
This is NOT a general purpose OIDC client, and such does not support many standard options.
This application skips many steps and does not use any oauth2/oidc client libraries.
PLEASE USE A REAL LIBRARY WHEN DEVELOPING YOUR OIDC APPLICATION.

**!! THIS APPLICATION SHOULD NOT BE USED FOR ANY OTHER PURPOSE OTHER THAN TESTING AND EDUCATION !!**

## Requirements
In order for the flow to work you need the following:
- you need a reachable URL so that the oidc server can reach your jwks URL. Set this in config file: `external_self_baseurl`
- you need to generate an RSA key (although a test one is already provided)
- you need to set the `oidc.discovery_url` to the well-known URL of the OIDC provider
- you will need to be onboarded with the following configuration:
  - `oidc.client_id` must match configuration
  - jwks URL must be onboarded as `<external_self_baseurl>/api/v1/jwks`
  - redirect URL must be onboarded as `<external_self_baseurl>/api/v1/callback`

**NOTE** the default configuration file is located in: [testdata/sampleconfig/config.yaml](testdata/sampleconfig/config.yaml)

## Begin flow
start your OIDC client server like: `go run cmd/oidc-client/main.go`.
Then navigate to: `<external_self_baseurl>/api/v1/auth`
