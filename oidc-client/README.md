# OIDC Test Client for use with OIDC DAC Flow

## Requirements
In order for the flow to work you need the following:
- you need a reachable URL so that the oidc server can reach your jwks URL. Set this in config file: `external_self_baseurl`
- you need to generate an RSA key (although a test one is already provided)
- you need to set the `oidc.discovery_url` to the well-known URL of the OIDC provider
- you will need to be onboarded with the following configuration:
  - `oidc.client_id` must match configuration
  - jwks URL must be onboarded as `<external_self_baseurl>/api/v1/jwks`
  - redirect URL must be onboarded as `<external_self_baseurl>/api/v1/callback`
    
## Begin flow
start your OIDC client server like: `go run cmd/oidc-client/main.go`.
Then navigate to: `<external_self_baseurl>/api/v1/auth`