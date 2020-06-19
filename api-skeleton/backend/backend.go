package backend

import "github.com/asaskevich/govalidator"

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
	//govalidator.SetNilPtrAllowedByRequired(false)
}

// StartBackendDeamon a blocking function that starts the backend processes
func StartBackendDeamon(cfgFile string) {

	// init the config
	ConfigInit(cfgFile, true)

	// TODO init the backend here for now...
	//if err != nil {
	//	log.Fatalf("could not init provider: %v", err)
	//}

	// start the server. block for now...
	startHTTPServer()
}
