package issue

type IssueRequest struct {
	Issue string `json:"issue"`
}

type ClassificationResponse struct {
	Severity string `json:"severity"`
}
