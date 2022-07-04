package validator

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/TelefonicaTC2Tech/golium/steps/http"
	"github.com/stretchr/testify/require"
)

const (
	fieldName           = "field_name"
	fieldType           = "field_type"
	mapComparissonError = "\nmap[details:map[field_name:map[code:incorrect_type message:" +
		"Wrong message Field 'field_name' has invalid type, expected field_type]] message:" +
		"Validation error. status_code:400]\n vs \nmap[details:" +
		"map[field_name:map[code:incorrect_type message:" +
		"Field 'field_name' has invalid type, expected field_type]] message:" +
		"Validation error. status_code:%!s(float64=400)]"
	replaceMapStringFile = `
		[
			{
				"code": "example1",
				"body": {
					"title": "title_to_replace",
					"body": "bar1",
					"userId": 1
				},
				"response": {
					"message": "Validation error.",
					"details": {
					  "field_to_replace": {
						"message": "Field 'field_to_replace' has invalid type, expected type_to_replace",
						"code": "incorrect_type"
					  }
					},
					"status_code": 400
				}
			},
			{
				"code": "example2",
				"body": {
					"title": "title_to_replace",
					"body": "bar1",
					"userId": 1
				},
				"response": {
					"message": "Validation error.",
					"details": {
					  "field_to_replace": {
						"message": "Wrong message Field 'field_to_replace' has invalid type, expected` +
		` type_to_replace",
						"code": "incorrect_type"
					  }
					},
					"status_code": 400
				}
			}
		]
		`
	mapResponseBody = `{
			"message": "Validation error.",
			"details": {
				"field_name": {
				"message": "Field 'field_name' has invalid type, expected field_type",
				"code": "incorrect_type"
				}
			},
			"status_code": 400
			}
		`
	mapResponseBodyError = `{
			"message": "Validation error.",
			"details": {
				"field_name": {
				"message": "Field 'field_name' has invalid type, expected field_type",
				"code": "incorrect_type"
				}
			},
			"status_code": 400
			},
		`
	replaceStringFile = `
		[
			{
				"code": "example1",
				"body": {
					"title": "title_to_replace",
					"body": "bar1",
					"userId": 1
				},
				"response": "field_to_replace and type_to_replace has been replaced"
			},
			{
				"code": "example2",
				"body": {
					"title": "title_to_replace",
					"body": "bar1",
					"userId": 1
				},
				"response": "field_to_replace and type_to_replace has been replaced with match error"
			}
		]
		`
	stringError = "\nfield_name and field_type has been replaced with match error\n " +
		"vs \nfield_name and field_type has been replaced"
	schemasPath = "./schemas"
	logsPath    = "./logs"
)

func TestReplaceMapStringResponse(t *testing.T) {
	os.MkdirAll(logsPath, os.ModePerm)
	defer os.RemoveAll(logsPath)

	os.MkdirAll(schemasPath, os.ModePerm)
	os.WriteFile("./schemas/test.json", []byte(replaceMapStringFile), os.ModePerm)
	defer os.RemoveAll(schemasPath)

	paramsInput := make(map[string]interface{})
	paramsInput["field"] = fieldName
	paramsInput["type"] = fieldType

	tcs := []struct {
		name         string
		code         string
		file         string
		responseBody string
		marshalMock  bool
		expectedErr  error
		params       map[string]interface{}
	}{
		{
			name:         "not_error",
			code:         "example1",
			file:         "test",
			params:       paramsInput,
			marshalMock:  false,
			responseBody: mapResponseBody,
			expectedErr:  nil,
		},
		{
			name:         "not_equal_error",
			code:         "example2",
			file:         "test",
			params:       paramsInput,
			marshalMock:  false,
			responseBody: mapResponseBody,
			expectedErr:  fmt.Errorf("expected JSON does not match actual, %v", mapComparissonError),
		},
		{
			name:         "unmarshal_error",
			code:         "example2",
			file:         "test",
			params:       paramsInput,
			marshalMock:  true,
			responseBody: mapResponseBodyError,
			expectedErr:  fmt.Errorf("error unmarshalling response body: %w", fmt.Errorf("unmarshal error")),
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ctx := http.InitializeContext(context.Background())
			j := JSONService{}

			if tc.marshalMock {
				originalUnmarshal := unmarshal
				defer func() { unmarshal = originalUnmarshal }()
				unmarshal = func(data []byte, v interface{}) error {
					return fmt.Errorf("unmarshal error")
				}
			}
			respBody, _ := io.ReadAll(bytes.NewBufferString(tc.responseBody))
			bodyContent, _ := http.GetParamFromJSON(ctx, tc.file, tc.code, "response")

			err := j.ReplaceMapStringResponse(respBody, bodyContent, tc.params)
			require.Equal(t, tc.expectedErr, err)
		})
	}
}
func TestReplaceStringResponse(t *testing.T) {
	os.MkdirAll(logsPath, os.ModePerm)
	defer os.RemoveAll(logsPath)

	var stringResponseBody = `field_name and field_type has been replaced`

	os.MkdirAll(schemasPath, os.ModePerm)
	os.WriteFile("./schemas/test.json", []byte(replaceStringFile), os.ModePerm)
	defer os.RemoveAll(schemasPath)

	paramsInput := make(map[string]interface{})
	paramsInput["field"] = fieldName
	paramsInput["type"] = fieldType

	tcs := []struct {
		name         string
		code         string
		request      string
		file         string
		responseBody string
		expectedErr  error
		params       map[string]interface{}
	}{
		{
			name:         "not_error",
			code:         "example1",
			request:      "string_test",
			file:         "test",
			responseBody: stringResponseBody,
			params:       paramsInput,
			expectedErr:  nil,
		},
		{
			name:         "not_equal_err",
			code:         "example2",
			request:      "string_test",
			file:         "test",
			responseBody: stringResponseBody,
			params:       paramsInput,
			expectedErr:  fmt.Errorf("received body does not match expected, %v", stringError),
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ctx := http.InitializeContext(context.Background())
			j := JSONService{}

			respBody, _ := io.ReadAll(bytes.NewBufferString(tc.responseBody))
			bodyContent, _ := http.GetParamFromJSON(ctx, tc.file, tc.code, "response")

			err := j.ReplaceStringResponse(respBody, fmt.Sprint(bodyContent), tc.params)
			require.Equal(t, tc.expectedErr, err)
		})
	}
}
