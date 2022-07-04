package http

import (
	"context"
	"fmt"
	"strings"

	"github.com/TelefonicaTC2Tech/golium"
	"github.com/TelefonicaTC2Tech/golium/steps/http/validator"
	"github.com/cucumber/godog"
)

// ValidateResponseBodyJSONFileModifying
// Validates the response body against the JSON in File modifying params.
func (s *Session) ValidateResponseBodyJSONFileModifying(
	ctx context.Context,
	code, file string,
	t *godog.Table,
) error {
	jsonResponseBody, err := GetParamFromJSON(ctx, file, code, "response")
	if err != nil {
		return fmt.Errorf("error getting parameter from json: %w", err)
	}
	jsonResponseBodyMap, _ := jsonResponseBody.(map[string]interface{})

	params, err := golium.ConvertTableToMap(ctx, t)
	if err != nil {
		return err
	}
	for key, value := range params {
		if strings.Contains(key, ".") {
			err = processNestedParams(jsonResponseBodyMap, key, value)
			if err != nil {
				return fmt.Errorf("error processing params: %v", err)
			}
		} else {
			_, present := jsonResponseBodyMap[key]
			if !present {
				return fmt.Errorf("error modifying param: param %v does not exists", key)
			}
			jsonResponseBodyMap[key] = value
		}
	}
	return s.ValidateResponseFromJSONFile(jsonResponseBody, "")
}

// processNestedParams Replace nested params split by "." for modified validation.
func processNestedParams(
	jsonResponseBodyMap map[string]interface{},
	key string, value interface{},
) error {
	var aux, lastAux map[string]interface{}
	aux = jsonResponseBodyMap
	keys := strings.Split(key, ".")
	var lastKey string
	for _, key := range keys {
		lastAux = aux
		_, present := aux[key]
		if !present {
			return fmt.Errorf("error modifying nested param: param %v does not exists", key)
		}
		aux, _ = aux[key].(map[string]interface{})
		lastKey = key
	}
	lastAux[lastKey] = value
	return nil
}

// ValidateResponseBodyJSONFileReplace
// Validates the response body against the JSON in File replacing values.
func (s *Session) ValidateErrorBodyJSONFileReplace(
	ctx context.Context,
	code, file string,
	t *godog.Table,
) error {
	replaceValues, err := golium.ConvertTableToMap(ctx, t)
	if err != nil {
		return err
	}
	respBody := s.Response.ResponseBody
	bodyContent, err := GetParamFromJSON(ctx, file, code, "response")
	if err != nil {
		return fmt.Errorf("error getting parameter from json: %w", err)
	}
	switch bC := bodyContent.(type) {
	case string:
		if err := validator.JSON.ReplaceStringResponse(
			respBody,
			fmt.Sprint(bodyContent),
			replaceValues,
		); err != nil {
			return err
		}

	case map[string]interface{}:
		if err := validator.JSON.ReplaceMapStringResponse(
			respBody,
			bodyContent,
			replaceValues,
		); err != nil {
			return err
		}
	default:
		return fmt.Errorf("body content should be string or map: %v", bC)
	}
	return nil
}
