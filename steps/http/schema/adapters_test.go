package schema

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/TelefonicaTC2Tech/golium"
	"github.com/cucumber/godog"
	"github.com/stretchr/testify/assert"
)

const (
	schemasPath                   = "./schemas"
	fileName                      = "health"
	validateModifyingResponseFile = `
		[
			{
				"code": "example1",
				"body": {
					"title": "foo1",
					"body": "bar1",
					"userId": 1
				},
				"response": {
					"id": 101,
					"title": "foo2",
					"body": "bar1",
					"userId": 1
				}
			},
			{
				"code": "example2",
				"body": {
					"title": "foo",
					"body": "bar1",
					"userId": 1
				},
				"response": {
					"id": 1,
					"name": "Leanne Graham",
					"username": "Bret",
					"email": "Sincere@april.biz",
					"address": {
					  "street": "Kulas Light",
					  "suite": "Apt. 556",
					  "city": "Gwenborough",
					  "zipcode": "92998-3874",
					  "geo": {
						"lat": "-37.3159",
						"lng": "81.1496"
					  }
					},
					"phone": "1-770-736-8031 x56442",
					"website": "hildegard.org",
					"company": {
					  "name": "to replace",
					  "catchPhrase": "Multi-layered client-server neural-net",
					  "bs": "harness real-time e-markets"
					}
				}
			}
		]
		`
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

	JSONFile = `{
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

func TestModifyResponse(t *testing.T) {
	os.MkdirAll(schemasPath, os.ModePerm)
	os.WriteFile("./schemas/posts.json", []byte(validateModifyingResponseFile), os.ModePerm)
	os.WriteFile("./schemas/users.json", []byte(validateModifyingResponseFile), os.ModePerm)
	defer os.RemoveAll(schemasPath)
	type args struct {
		code string
		file string
		t    *godog.Table
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Error getting response from file",
			args: args{
				code: "example3",
				file: "posts",
				t: golium.NewTable([][]string{
					{"parameter", "value"},
					{"title", "foo1"},
				}),
			},
			wantErr: true,
		},
		{
			name: "Error converting table",
			args: args{
				code: "example1",
				file: "posts",
				t: golium.NewTable([][]string{
					{"title", "foo1"},
				}),
			},
			wantErr: true,
		},
		{
			name: "Error not present nested key",
			args: args{
				code: "example2",
				file: "posts",
				t: golium.NewTable([][]string{
					{"parameter", "value"},
					{"not-present.name", "Romaguera-Crona"},
				}),
			},
			wantErr: true,
		},
		{
			name: "Error not present simple key",
			args: args{
				code: "example1",
				file: "posts",
				t: golium.NewTable([][]string{
					{"parameter", "value"},
					{"wrong_key", "true"},
				}),
			},
			wantErr: true,
		},
		{
			name: "Valid nested key",
			args: args{
				code: "example2",
				file: "users",
				t: golium.NewTable([][]string{
					{"parameter", "value"},
					{"company.name", "Romaguera-Crona"},
				}),
			},
			wantErr: false,
		},
		{
			name: "Valid simple key",
			args: args{
				code: "example1",
				file: "posts",
				t: golium.NewTable([][]string{
					{"parameter", "value"},
					{"title", "foo1"},
				}),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ModifyResponse(context.Background(), tt.args.code, tt.args.file, tt.args.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("ModifyResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetParam(t *testing.T) {
	var expectedParam interface{}
	if err := json.Unmarshal([]byte(JSONFile), &expectedParam); err != nil {
		t.Error("error Unmarshaling expected response body: %w", err)
	}

	var dataStruct interface{}
	if err := json.Unmarshal([]byte(JSONhttpFileValues), &dataStruct); err != nil {
		t.Error("error Unmarshaling expected response body: %w", err)
	}

	golium.GetConfig().Dir.Schemas = schemasPath

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
			resultParam, err := GetParam(tc.fileName, tc.code, tc.param)
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
	if err := json.Unmarshal([]byte(JSONFile), &expectedValue); err != nil {
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

func TestLoadData(t *testing.T) {
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

	golium.GetConfig().Dir.Schemas = schemasPath
	os.MkdirAll("./schemas", os.ModePerm)
	os.WriteFile("./schemas/http.json", []byte(JSONhttpFileValues), os.ModePerm)
	defer os.RemoveAll("./schemas/")

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			_, err := LoadData(tc.fileName)
			if err != nil {
				assert.Containsf(t, err.Error(), tc.expectedErr, "error message %s", "formatted")
				fmt.Printf(err.Error(), tc.expectedErr)
			}
		})
	}
}

func TestUnmarshalData(t *testing.T) {
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
			unmarshalled, err := UnmarshalData([]byte(tc.current))
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

func TestGetBody(t *testing.T) {
	os.MkdirAll(schemasPath, os.ModePerm)
	os.WriteFile("./schemas/health.json", []byte(JSONhttpFileValues), os.ModePerm)
	defer os.RemoveAll(schemasPath)
	type args struct {
		code string
		file string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Error getting from file",
			args: args{
				code: "not_valid_code",
				file: fileName,
			},
			wantErr: true,
		},
		{
			name: "Happy Path",
			args: args{
				code: "example1",
				file: fileName,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetBody(context.Background(), tt.args.code, tt.args.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBody() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestModifyBody(t *testing.T) {
	os.MkdirAll(schemasPath, os.ModePerm)
	os.WriteFile("./schemas/health.json", []byte(JSONhttpFileValues), os.ModePerm)
	defer os.RemoveAll(schemasPath)
	type args struct {
		code string
		file string
		t    *godog.Table
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Error converting table",
			args: args{
				code: "not_valid_code",
				file: fileName,
				t: golium.NewTable([][]string{
					{"parameter"},
				}),
			},
			wantErr: true,
		},
		{
			name: "Error getting from file",
			args: args{
				code: "not_valid_code",
				file: fileName,
				t: golium.NewTable([][]string{
					{"parameter", "value"},
					{"boolean", "true"},
				}),
			},
			wantErr: true,
		},
		{
			name: "Error with wrong body param",
			args: args{
				code: "example1",
				file: fileName,
				t: golium.NewTable([][]string{
					{"parameter", "value"},
					{"wrong_key", "true"},
				}),
			},
			wantErr: true,
		},
		{
			name: "Happy Path",
			args: args{
				code: "example1",
				file: fileName,
				t: golium.NewTable([][]string{
					{"parameter", "value"},
					{"boolean", "true"},
				}),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ModifyBody(context.Background(), tt.args.code, tt.args.file, tt.args.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("ModifyBody() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestDeleteBodyFields(t *testing.T) {
	os.MkdirAll(schemasPath, os.ModePerm)
	os.WriteFile("./schemas/health.json", []byte(JSONhttpFileValues), os.ModePerm)
	defer os.RemoveAll(schemasPath)
	type args struct {
		code string
		file string
		t    *godog.Table
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Error converting table",
			args: args{
				code: "not_valid_code",
				file: fileName,
				t: golium.NewTable([][]string{
					{"parameter"},
				}),
			},
			wantErr: true,
		},
		{
			name: "Error getting from file",
			args: args{
				code: "not_valid_code",
				file: fileName,
				t: golium.NewTable([][]string{
					{"parameter"},
					{"boolean"},
				}),
			},
			wantErr: true,
		},
		{
			name: "Happy Path",
			args: args{
				code: "example1",
				file: fileName,
				t: golium.NewTable([][]string{
					{"parameter"},
					{"boolean"},
				}),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := DeleteBodyFields(context.Background(), tt.args.code, tt.args.file, tt.args.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteBodyFields() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestDeleteResponseFields(t *testing.T) {
	os.MkdirAll("./schemas", os.ModePerm)
	os.WriteFile("./schemas/health.json", []byte(JSONhttpFileValues), os.ModePerm)
	defer os.RemoveAll("./schemas/")

	type args struct {
		code string
		file string
		t    *godog.Table
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Error getting response",
			args: args{
				code: "wrong_code",
				file: fileName,
				t: golium.NewTable([][]string{
					{"parameter"},
					{"boolean"},
				}),
			},
			wantErr: true,
		},
		{
			name: "Error converting table",
			args: args{
				code: "example1",
				file: fileName,
				t: golium.NewTable([][]string{
					{"boolean"},
				}),
			},
			wantErr: true,
		},
		{
			name: "Happy Path",
			args: args{
				code: "example1",
				file: fileName,
				t: golium.NewTable([][]string{
					{"parameter"},
					{"boolean"},
				}),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := DeleteResponseFields(context.Background(), tt.args.code, tt.args.file, tt.args.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteResponseFields() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
