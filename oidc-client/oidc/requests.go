package oidc

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/spf13/viper"
)

// this is the initial request we direct the user's browser to.
// It is sent to the oidc auth endpoint and starts the whole flow.
func GenerateAuthURL(state, redirectURL string) (authUrl string) {
	base, _ := url.Parse(discoveryCache.AuthEndpoint)
	// construct query params
	params := url.Values{}
	params.Set("response_type", "code")
	params.Set("scope", "openid foundation_profile")
	params.Set("client_id", viper.GetString("oidc.client_id"))
	params.Set("state", state)
	params.Set("request", createRequestJWT(state, redirectURL))
	base.RawQuery = params.Encode()

	log.Debugf("auth URL was constructed with state: %s", state)
	return base.String()
}

type tokenResponse struct {
	AccessToken string `json:"access_token"`
}

func AccessTokenRequest(code, redirectURL string) (accessToken string, err error) {
	base, _ := url.Parse(discoveryCache.TokenEndpoint)
	// construct query params
	params := url.Values{}
	params.Set("grant_type", "authorization_code")
	params.Set("client_assertion_type", "urn:ietf:params:oauth:client-assertion-type:jwt-bearer")
	params.Set("client_assertion", createAssertionJWT())
	params.Set("redirect_uri", redirectURL)
	params.Set("code", code)
	params.Set("client_id", viper.GetString("oidc.client_id"))
	params.Set("code_verifier", codeVerifier.String())

	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}
	resp := doHttpCall("POST", base.String(), headers, []byte(params.Encode()))

	// handle any errors
	if resp.err != nil {
		err = resp.err
		return
	}
	if resp.StatusCode != 200 {
		err = fmt.Errorf("token url returned a non-200 status: %v with body: %s", resp.StatusCode, resp.BodyBytes)
		return
	}

	// parse the access token
	var tr tokenResponse
	err = json.Unmarshal(resp.BodyBytes, &tr)
	if err == nil {
		accessToken = tr.AccessToken
	}
	return
}

func UserInfoRequest(accessToken string) (responseBody []byte, err error) {
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", accessToken),
	}
	resp := doHttpCall("GET", discoveryCache.UserinfoEndpoint, headers, nil)

	// handle any errors
	if resp.err != nil {
		err = resp.err
		return
	}
	if resp.StatusCode != 200 {
		err = fmt.Errorf("userinfo url returned a non-200 status: %v with body: %s", resp.StatusCode, resp.BodyBytes)
		return
	}
	responseBody = resp.BodyBytes
	return
}
