package issue

type IssueRequest struct {
	Issue string `json:"issue"`
}

type ClassificationResponse struct {
	Severity string `json:"severity"`
}

type Instance struct {
	Content string `json:"content"`
}

type Params struct {
	Temperature     float64 `json:"temperature"`
	MaxOutputTokens int     `json:"maxOutputTokens"`
	TopP            float64 `json:"topP"`
	TopK            float64 `json:"topK"`
}
type VertexAIRequest struct {
	Instances  []Instance `json:"instances"`
	Parameters Params     `json:"parameters"`
}

type Prediction struct {
	Content string `json:"content"`
}
type VertexAIResponse struct {
	Predictions []Prediction `json:"predictions"`
}
