package elasticsearch

type ResponseError struct {
	Type   string `json:"type"`
	Reason string `json:"reason"`
}

type Response struct {
	Error ResponseError `json:"error"`
}
