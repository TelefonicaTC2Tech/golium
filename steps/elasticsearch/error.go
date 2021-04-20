package elasticsearch

// ResponseError represents the error information in elasticsearch response
type ResponseError struct {
	Type   string `json:"type"`
	Reason string `json:"reason"`
}

// Response represents the body of an elasticsearch response
type Response struct {
	Error ResponseError `json:"error"`
}
