package api

type Response struct {
	Status    string `json:"status"`
	ReportUrl string `json:"reportUrl"`
	ReportId  int    `json:"reportId,omitempty"`
	Error     string `json:"error,omitempty"`
}
