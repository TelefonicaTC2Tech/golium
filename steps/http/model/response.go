package model

import "net/http"

// Response information of the session.
type Response struct {
	// HTTP response
	HTTPResponse *http.Response
	// Response body as slice of bytes
	ResponseBody []byte
}
