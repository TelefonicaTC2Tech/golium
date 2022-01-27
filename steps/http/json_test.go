package http

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/Telefonica/golium"
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
		os.WriteFile("./schemas/http.json", []byte(JSONhttpFileValues), os.ModePerm)
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
		resultParam, paramFromJSONErr := GetParamFromJSON(ctx, fileName, code, param)
		if paramFromJSONErr != nil {
			t.Errorf("error loading parameter from file %s due to error: %v", fileName, paramFromJSONErr)
		}

		assert.True(t, 
			reflect.DeepEqual(resultParam, expectedParam), 
			fmt.Sprintf("expected JSON parameter does not match response JSON parametr, \n%v\n vs \n%s", resultParam, JSON))
	})
}

func TestGetParamFromJSONFile(t *testing.T) {
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

		os.MkdirAll("./schemas", os.ModePerm)
		os.WriteFile("./schemas/http.json", []byte(JSONhttpFileValues), os.ModePerm)
		defer os.RemoveAll("./schemas/")

		var fileName = "http"
		var code = "example1"
		var param = "body"
		var JSON = `{
            "boolean": false, 
            "empty": "", 
            "list": [
                { "attribute": "attribute0", "value": "value0"},
                { "attribute": "attribute1", "value": "value1"},
                { "attribute": "attribute2", "value": "value2"}
            ]
        }`
		var ctx = context.Background()
		var expectedParam interface{}
		if err := json.Unmarshal([]byte(fmt.Sprint(JSON)), &expectedParam); err != nil {
			t.Error("error Unmarshaling expected response body: %w", err)
		}

		// Call function to test
		resultParam, paramFromJSONErr := GetParamFromJSONFile(ctx, fileName, code, param)
		if paramFromJSONErr != nil {
			t.Errorf("error loading parameter from file %s due to error: %v", fileName, paramFromJSONErr)
		}

		assert.True(t, 
			reflect.DeepEqual(resultParam, expectedParam), 
			fmt.Sprintf("expected JSON parameter does not match response JSON parametr, \n%v\n vs \n%s", resultParam, JSON))
	})
}
