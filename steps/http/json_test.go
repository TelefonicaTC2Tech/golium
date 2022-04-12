package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/TelefonicaTC2Tech/golium"
	"github.com/stretchr/testify/assert"
)

const (
	schemasDir         = "./schemas"
	JSONhttpFileValues = `
	[
		{
			"code": "example1",
			"body": {
				"empty": "",
				"boolean": false,
				"list": [
				{ "attribute": "attribute0", "value": "value0"},
				{ "attribute": "attribute1", "value": "value1"},
				{ "attribute": "attribute2", "value": "value2"}
				]
			},
			"response": {
				"boolean": false, 
				"empty": "", 
				"list": [
					{ "attribute": "attribute0", "value": "value0"},
					{ "attribute": "attribute1", "value": "value1"},
					{ "attribute": "attribute2", "value": "value2"}
				]
			}
		}
	]
	`

	JSONhttpResponse = `{
	"code": "example1",
	"body": {
		"empty": "",
		"boolean": false,
		"list": [
		{ "attribute": "attribute0", "value": "value0"},
		{ "attribute": "attribute1", "value": "value1"},
		{ "attribute": "attribute2", "value": "value2"}
		]
	},
	"response": {
		"boolean": false, 
		"empty": "", 
		"list": [
			{ "attribute": "attribute0", "value": "value0"},
			{ "attribute": "attribute1", "value": "value1"},
			{ "attribute": "attribute2", "value": "value2"}
		]
	}
	}`

	JSON = `{
		"boolean": false, 
		"empty": "", 
		"list": [
			{ "attribute": "attribute0", "value": "value0"},
			{ "attribute": "attribute1", "value": "value1"},
			{ "attribute": "attribute2", "value": "value2"}
		]
	}`

	JSONhttpFileBadFormat = `
	[
		{
			"code": "example1",
			"body": {
	`
)

func TestGetParamFromJSON(t *testing.T) {
	var expectedParam interface{}
	if err := json.Unmarshal([]byte(JSON), &expectedParam); err != nil {
		t.Error("error Unmarshaling expected response body: %w", err)
	}

	var dataStruct interface{}
	if err := json.Unmarshal([]byte(JSONhttpFileValues), &dataStruct); err != nil {
		t.Error("error Unmarshaling expected response body: %w", err)
	}

	golium.GetConfig().Dir.Schemas = schemasDir

	os.MkdirAll("./schemas", os.ModePerm)
	os.WriteFile("./schemas/httpBadFormat.json", []byte(JSONhttpFileBadFormat), os.ModePerm)
	os.WriteFile("./schemas/http.json", []byte(JSONhttpFileValues), os.ModePerm)
	defer os.RemoveAll("./schemas/")

	tcs := []struct {
		name          string
		fileName      string
		code          string
		param         string
		expectedErr   string
		expectedValue interface{}
	}{
		{
			name:          "Should return selected value from JSON file",
			fileName:      "http",
			code:          "example1",
			param:         "response",
			expectedErr:   "",
			expectedValue: expectedParam,
		},
		{
			name:          "Should return a error loading file",
			fileName:      "httpNotExist",
			code:          "example1",
			param:         "response",
			expectedErr:   "error loading file at httpNotExist due to error:",
			expectedValue: nil,
		},
		{
			name:          "Should return a error unmarsharlling JSON file",
			fileName:      "httpBadFormat",
			code:          "example1",
			param:         "response",
			expectedErr:   "error unmarsharlling JSON file at httpBadFormat due to error:",
			expectedValue: nil,
		},
		{
			name:     "Should return a error param value not found",
			fileName: "http",
			code:     "non-existing-code",
			param:    "response",
			expectedErr: fmt.Sprintf("param value: 'response' not found in '%v' due to error:",
				dataStruct),
			expectedValue: nil,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			var ctx = context.Background()
			resultParam, err := GetParamFromJSON(ctx, tc.fileName, tc.code, tc.param)
			if err != nil {
				assert.Containsf(t, err.Error(), tc.expectedErr, "error message %s", "formatted")
			}
			if !JSONEquals(resultParam, tc.expectedValue) {
				t.Errorf("value %v for param %s and code %s is not expected: %v",
					resultParam, tc.param, tc.code, tc.expectedValue)
			}
		})
	}
}

func TestFindValueByCode(t *testing.T) {
	var expectedValue interface{}
	if err := json.Unmarshal([]byte(JSON), &expectedValue); err != nil {
		t.Error("error Unmarshaling expected response body: %w", err)
	}

	dataStruct := []map[string]interface{}{}
	if err := json.Unmarshal([]byte(JSONhttpFileValues), &dataStruct); err != nil {
		t.Error("error Unmarshaling expected response body: %w", err)
	}

	tcs := []struct {
		name          string
		code          string
		param         string
		expectedValue interface{}
	}{
		{
			name:          "value found with code and param",
			code:          "example1",
			param:         "response",
			expectedValue: expectedValue,
		},
		{
			name:          "value not found due non existing param",
			code:          "example1",
			param:         "non-existing-param",
			expectedValue: nil,
		},
		{
			name:          "value not found due non existing code",
			code:          "non-existing-code",
			param:         "response",
			expectedValue: nil,
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			value, err := FindValueByCode(dataStruct, tc.code, tc.param)
			if err != nil {
				if err.Error() != fmt.Sprintf("value for param: '%s' with code: '%s' not found",
					tc.param, tc.code) {
					t.Errorf("error not expected with param '%s' and code '%s':\n%v",
						tc.param, tc.code, err)
				}
			}

			if !JSONEquals(value, tc.expectedValue) {
				t.Errorf("value %v for param %s and code %s is not expected: %v",
					value, tc.param, tc.code, tc.expectedValue)
			}
		})
	}
}

func TestLoadJSONData(t *testing.T) {
	tcs := []struct {
		name        string
		fileName    string
		expectedErr string
	}{
		{
			name:        "Should return data json file",
			fileName:    "http",
			expectedErr: "",
		},
		{
			name:        "Should return error reading file",
			fileName:    "httpNotExistsFile",
			expectedErr: "error reading file",
		},
	}

	golium.GetConfig().Dir.Schemas = schemasDir
	os.MkdirAll("./schemas", os.ModePerm)
	os.WriteFile("./schemas/http.json", []byte(JSONhttpFileValues), os.ModePerm)
	defer os.RemoveAll("./schemas/")

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			_, err := LoadJSONData(tc.fileName)
			if err != nil {
				assert.Containsf(t, err.Error(), tc.expectedErr, "error message %s", "formatted")
				fmt.Printf(err.Error(), tc.expectedErr)
			}
		})
	}
}

func TestUnmarshalJSONData(t *testing.T) {
	var expectedString = `[
		{
			"boolean": false, 
			"empty": "", 
			"list": [
				{ "attribute": "attribute0", "value": "value0"},
				{ "attribute": "attribute1", "value": "value1"},
				{ "attribute": "attribute2", "value": "value2"}
			]
		}
	]`

	var current = `[
		{
			"boolean": false, 
			"empty": "", 
			"list": [
				{ "attribute": "attribute0", "value": "value0"},
				{ "attribute": "attribute1", "value": "value1"},
				{ "attribute": "attribute2", "value": "value2"}
			]
		}
	]`

	var incorrect = `
		{
			"boolean": false, 
			"empty": ""
		}`

	var message = "error unmarshalling JSON data due to error: json: cannot unmarshal " +
		"object into Go value of type []map[string]interface {}"
	formatError := errors.New(message)
	var expected interface{}
	if err := json.Unmarshal([]byte(expectedString), &expected); err != nil {
		t.Error("error Unmarshaling expected response body: %w", err)
	}

	tcs := []struct {
		name     string
		expected interface{}
		current  string
		err      error
	}{
		{
			name:     "equals JSON values from structure",
			expected: expected,
			current:  current,
			err:      nil,
		},
		{
			name:     "equals JSON values from structure",
			expected: expected,
			current:  incorrect,
			err:      formatError,
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			unmarshalled, err := UnmarshalJSONData([]byte(tc.current))
			if err != nil {
				if err.Error() != tc.err.Error() {
					t.Errorf("unexpected error unmarshalling data:\n%v\nexpected:\n%v", err, tc.err)
				}
			}
			if JSONEquals(tc.expected, unmarshalled) {
				t.Errorf("expected unmarshalled data error:\n%v", err)
			}
		})
	}
}

func TestJSONEquals(t *testing.T) {
	var expectedString = `{
		"boolean": false, 
		"empty": "", 
		"list": [
			{ "attribute": "attribute0", "value": "value0"},
			{ "attribute": "attribute1", "value": "value1"},
			{ "attribute": "attribute2", "value": "value2"}
		]
	}`

	var differentString = `{
		"boolean": true, 
		"empty": "", 
		"list": [
			{ "attribute": "attribute0", "value": "value0"},
			{ "attribute": "attribute1", "value": "value1"},
			{ "attribute": "attribute2", "value": "value2"}
		]
	}`

	var currentString = `{
		"boolean": false, 
		"empty": "", 
		"list": [
			{ "attribute": "attribute0", "value": "value0"},
			{ "attribute": "attribute1", "value": "value1"},
			{ "attribute": "attribute2", "value": "value2"}
		]
	}`
	var expected interface{}
	if err := json.Unmarshal([]byte(expectedString), &expected); err != nil {
		t.Error("error Unmarshaling expected response body: %w", err)
	}

	var different interface{}
	if err := json.Unmarshal([]byte(differentString), &different); err != nil {
		t.Error("error Unmarshaling expected response body: %w", err)
	}

	var current interface{}
	if err := json.Unmarshal([]byte(currentString), &current); err != nil {
		t.Error("error Unmarshaling expected response body: %w", err)
	}

	tcs := []struct {
		name     string
		expected interface{}
		current  interface{}
		equals   bool
	}{
		{
			name:     "equals JSON values from structure",
			expected: expected,
			current:  current,
			equals:   true,
		},
		{
			name:     "not equals JSON values from structure",
			expected: different,
			current:  current,
			equals:   false,
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			if tc.equals != JSONEquals(tc.expected, tc.current) {
				t.Errorf("expected JSON comparison should be %t \n%v\n vs \n%v", tc.equals, tc.expected,
					tc.current)
			}
		})
	}
}
