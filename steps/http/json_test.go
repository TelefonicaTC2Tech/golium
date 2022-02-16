package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/TelefonicaTC2Tech/golium"
	"github.com/stretchr/testify/assert"
)

func TestGetParamFromJSON(t *testing.T) {
	t.Run("Should return selected value from JSON file", func(t *testing.T) {

		var JSONhttpFileValues = `
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
		var ctx = context.Background()
		golium.GetConfig().Dir.Schemas = "./schemas"

		os.MkdirAll("./schemas", os.ModePerm)
		ioutil.WriteFile("./schemas/http.json", []byte(JSONhttpFileValues), os.ModePerm)
		defer os.RemoveAll("./schemas/")

		var fileName = "http"
		var code = "example1"
		var param = "response"
		var JSON = `{
            "boolean": false, 
            "empty": "", 
            "list": [
                { "attribute": "attribute0", "value": "value0"},
                { "attribute": "attribute1", "value": "value1"},
                { "attribute": "attribute2", "value": "value2"}
            ]
        }`
		var expectedParam interface{}
		if err := json.Unmarshal([]byte(fmt.Sprint(JSON)), &expectedParam); err != nil {
			t.Error("error Unmarshaling expected response body: %w", err)
		}

		// Call function to test
		resultParam, err := GetParamFromJSON(ctx, fileName, code, param)
		if err != nil {
			t.Errorf("error loading parameter from file %s due to error: %v", fileName, err)
		}

		assert.True(t,
			reflect.DeepEqual(resultParam, expectedParam),
			fmt.Sprintf("expected JSON parameter does not match response JSON parametr, \n%v\n vs \n%s", resultParam, JSON))
	})
}

func TestFindValueByCode(t *testing.T) {
	var JSONhttpFileValues = `
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

	var JSON = `{
		"boolean": false, 
		"empty": "", 
		"list": [
			{ "attribute": "attribute0", "value": "value0"},
			{ "attribute": "attribute1", "value": "value1"},
			{ "attribute": "attribute2", "value": "value2"}
		]
	}`
	var expectedValue interface{}
	if err := json.Unmarshal([]byte(fmt.Sprint(JSON)), &expectedValue); err != nil {
		t.Error("error Unmarshaling expected response body: %w", err)
	}

	dataStruct := []map[string]interface{}{}
	if err := json.Unmarshal([]byte(fmt.Sprint(JSONhttpFileValues)), &dataStruct); err != nil {
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
				if err.Error() != fmt.Sprintf("value for param: '%s' with code: '%s' not found", tc.param, tc.code) {
					t.Errorf("error not expected with param '%s' and code '%s':\n%v", tc.param, tc.code, err)
				}
			}

			if !JSONEquals(value, tc.expectedValue) {
				t.Errorf("value %v for param %s and code %s is not expected: %v", value, tc.param, tc.code, tc.expectedValue)
			}
		})
	}

}

func TestLoadJSONData(t *testing.T) {

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

	var message = "error unmarshalling JSON data due to error: json: cannot unmarshal object into Go value of type []map[string]interface {}"
	formatError := fmt.Errorf(message)
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
	if err := json.Unmarshal([]byte(fmt.Sprint(expectedString)), &expected); err != nil {
		t.Error("error Unmarshaling expected response body: %w", err)
	}

	var different interface{}
	if err := json.Unmarshal([]byte(fmt.Sprint(differentString)), &different); err != nil {
		t.Error("error Unmarshaling expected response body: %w", err)
	}

	var current interface{}
	if err := json.Unmarshal([]byte(fmt.Sprint(currentString)), &current); err != nil {
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
