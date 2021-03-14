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

	// connect to database
	if err := openDatabase(); err != nil {
		log.Fatalf("unable to connect to database: %v", err)
	}
	// ensure we have the proper schema
	if err := migrateDatabaseSchema(); err != nil {
		log.Fatalf("unable to migrate database schema: %v", err)
	}

	// start the server. block for now...
	startHTTPServer()
}
