package oidc

func InitProvider() (err error) {
	initJwtSigningKey()
	err = pollDiscovery()
	return
}
