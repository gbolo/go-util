package storage

import (
	"crypto/sha256"
	"fmt"
	"regexp"
)

const (
	apiKeyLength                = 42
	apiKeyLengthWithDeliminator = apiKeyLength + 1
	apiKeyPrefixLength          = 8
	zbase32CharSet              = "ybndrfg8ejkmcpqxot1uwisza345h769"
)

var (
	// regex to match the api key format
	// example api key: 4bqgeysp.n7swohxky7imgq5jcuy8iue8kf3csmua65
	validKeyFormat = regexp.MustCompile(
		fmt.Sprintf(
			"^[%s]{%d}\\.[%s]{%d}$",
			zbase32CharSet,
			apiKeyPrefixLength,
			zbase32CharSet,
			apiKeyLength-apiKeyPrefixLength,
		),
	)
)

func ValidateKeyFormat(apiKey string) bool {
	return validKeyFormat.MatchString(apiKey)
}

type ApiKeyRaw string

func GenerateApiKey() (ApiKeyRaw, error) {
	key, err := generateRandomStringForKey(apiKeyLength)
	return ApiKeyRaw(key[:apiKeyPrefixLength] + "." + key[apiKeyPrefixLength:]), err
}

func (a ApiKeyRaw) HasValidFormat() bool {
	return ValidateKeyFormat(string(a))
}

func (a ApiKeyRaw) GetPrefix() (prefix string) {
	return string(a)[:apiKeyPrefixLength]
}

func (a ApiKeyRaw) GetHash() (hash string) {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(a)))
}

type ApiKeyStored struct {
	Prefix    string `json:"prefix"`
	Name      string `json:"name"`
	Hash      string `json:"hash"`
	CreatedAt int64  `json:"created_at"`
	LastUsed  int64  `json:"last_used,omitempty"`
	Revoked   bool   `json:"revoked"`
	RevokedAt int64  `json:"revoked_at,omitempty"`
}

func (a *ApiKeyStored) Validate(raw string) (valid bool) {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(raw))) == a.Hash
}
