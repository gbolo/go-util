package backend

type versionInfo struct {
	Version   string `json:"version"`
	CommitSHA string `json:"build_ref"`
	BuildDate string `json:"build_date"`
}

type errorResponse struct {
	Error string `json:"error"`
}

type successResponse struct {
	Message string `json:"message"`
}

type addService struct {
	ID          string `json:"id,omitempty"`
	Description string `json:"description"`
}

type generateAPIKey struct {
	Name      string `json:"name"`
	ServiceID string `json:"service_id"`
}

type revokeAPIKey struct {
	Prefix    string `json:"prefix"`
	ServiceID string `json:"service_id"`
}
