package model

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddBody(t *testing.T) {
	var messageJSON = map[string]interface{}{}
	tests := []struct {
		name    string
		message interface{}
	}{
		{
			name:    "Message is a string",
			message: "string message",
		},
		{
			name:    "Message is a json",
			message: messageJSON,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Request{}
			r.AddBody(tt.message)
		})
	}
}

func TestAddAuthorization(t *testing.T) {
	tests := []struct {
		name     string
		apiKey   string
		jwtValue string
	}{
		{
			name:     "Empty values",
			apiKey:   "",
			jwtValue: "",
		},
		{
			name:     "With values",
			apiKey:   "apikey",
			jwtValue: "jwt",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Request{
				Headers: make(map[string][]string),
			}
			r.AddAuthorization(tt.apiKey, tt.jwtValue)
			if tt.apiKey != "" {
				require.Equal(t, r.Headers["X-API-KEY"], []string{tt.apiKey})
			} else {
				_, ok := r.Headers["X-API-KEY"]
				require.Equal(t, ok, false)
			}
			if tt.jwtValue != "" {
				require.Equal(t, r.Headers["Authorization"], []string{fmt.Sprintf("Bearer %s", tt.jwtValue)})
			} else {
				_, ok := r.Headers["Authorization"]
				require.Equal(t, ok, false)
			}
		})
	}
}

func TestAddJSONHeaders(t *testing.T) {
	tests := []struct {
		name    string
		headers map[string][]string
	}{
		{
			name: "With empty headers",
		},
		{
			name: "With non empty headers",
			headers: map[string][]string{
				"header": {"value"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Request{
				Headers: tt.headers,
			}
			r.AddJSONHeaders()
			require.Equal(t, r.Headers["Content-Type"], []string{"application/json"})
		})
	}
}

func TestNormalizeEndpoint(t *testing.T) {
	tests := []struct {
		name      string
		endpoint  string
		backslash bool
		want      string
	}{
		{
			name:      "Without backslash and endpoint with",
			endpoint:  "test/",
			backslash: false,
			want:      "test",
		},
		{
			name:      "Without backslash and endpoint without",
			endpoint:  "test",
			backslash: false,
			want:      "test",
		},
		{
			name:      "With backslash and endpoint with",
			endpoint:  "test/",
			backslash: true,
			want:      "test/",
		},
		{
			name:      "With backslash and endpoint without",
			endpoint:  "test",
			backslash: true,
			want:      "test/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NormalizeEndpoint(tt.endpoint, tt.backslash); got != tt.want {
				t.Errorf("normalizeEndpoint() = %v, want %v", got, tt.want)
			}
		})
	}
}
