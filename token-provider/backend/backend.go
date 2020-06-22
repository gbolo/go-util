package backend

import (
	"token-provider/storage"

	"github.com/asaskevich/govalidator"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
	//govalidator.SetNilPtrAllowedByRequired(false)
}

var store storage.Store

// StartBackendDeamon a blocking function that starts the backend processes
func StartBackendDeamon(cfgFile string) {

	// init the config
	ConfigInit(cfgFile, true)

	// init storage provider
	store = storage.NewMemoryStore()

	// start the server. block for now...
	startHTTPServer()
}
