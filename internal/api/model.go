package api

type Response struct {
	Status    string `json:"status"`
	ReportURL string `json:"reportUrl"`
	ReportID  int    `json:"reportId,omitempty"`
	Error     string `json:"error,omitempty"`
}
