package storage

// current interface
type Store interface {
	AddService(description string) *Service
	AddServiceWithID(id, description string) *Service
	RemoveService(id string)
	UpdateService(service *Service)
	ListServices() map[string]*Service
	GetServiceByID(id string) *Service
	ValidateApiKeyForServiceID(id, rawKey string) bool
}

// WIP interface
type tokenProvider interface {
	// init provider
	Init(options interface{}) (err error)

	// service management
	AddService(description string) (id string, err error)
	AddServiceWithID(id, description string) (err error)
	GetService(id string) (service *Service, err error)
	RemoveService(id string) (err error)
	UpdateService(id, description string) (err error)
	ListServices() (services []*Service, err error)

	// token management
	GenerateTokenForService(serviceID, description string, ttlMins int) (rawToken string, err error)
	ValidateTokenForService(serviceID, rawToken string) (valid bool, err error)
	RevokeTokenForService(serviceID, tokenID string) (err error)
	UpdateTokenForService(serviceID, tokenID, description string) (err error)
}
