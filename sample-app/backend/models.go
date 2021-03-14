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

type Client struct {
	ID   string `json:"id" valid:"length(1|64)" gorm:"primaryKey"`
	Name string `json:"name" valid:"length(1|64)"`
	URL  string `json:"url" valid:"url"`
}
