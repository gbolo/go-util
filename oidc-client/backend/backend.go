package backend

import (
	"oidc-client/oidc"

	"github.com/asaskevich/govalidator"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
	//govalidator.SetNilPtrAllowedByRequired(false)
}

// StartBackendDeamon a blocking function that starts the backend processes
func StartBackendDeamon(cfgFile string) {

	// init the config
	ConfigInit(cfgFile, true)

	// TODO init the backend here for now...
	err := oidc.InitProvider()
	if err != nil {
		log.Fatalf("could not init oidc provider: %v", err)
	}

	// start the server. block for now...
	startHTTPServer()
}
