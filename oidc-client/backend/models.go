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
