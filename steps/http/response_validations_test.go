package http

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"syscall"
	"testing"

	"github.com/TelefonicaTC2Tech/golium"
	"github.com/TelefonicaTC2Tech/golium/steps/http/validator"
	"github.com/cucumber/godog"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

const (
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
	replaceMapResponseBody = `{
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
	replaceStringResponseBody       = `field_name and field_type has been replaced`
	jsonPlaceURL                    = "https://jsonplaceholder.typicode.com/"
	validateModifyingWrongCodeError = "param value: 'response' not found in" +
		" '[map[body:map[body:bar1 " +
		"title:foo1 userId:1] code:example1 response:map[body:bar1 id:101 title" +
		":foo2 userId:1]] map[body:map[body:bar1 title:foo userId:1] code:example2" +
		" response:map[address:map[city:Gwenborough geo:map[lat:-37.3159 lng:81.1496]" +
		" street:Kulas Light suite:Apt. 556 zipcode:92998-3874] company:map[bs:" +
		"harness real-time e-markets catchPhrase:Multi-layered client-server " +
		"neural-net name:to replace] email:Sincere@april.biz id:1 name:Leanne " +
		"Graham phone:1-770-736-8031 x56442 username:Bret website:hildegard.org]]]'"
	validateModifyNestedError = "expected JSON does not match real response, " +
		"\nmap[address:map[city:Gwenborough geo:map[lat:fake lng:81.1496] " +
		"street:Kulas Light suite:Apt. 556 zipcode:92998-3874] company:map[bs:harness real-time" +
		" e-markets catchPhrase:Multi-layered client-server neural-net name:to replace] email:Sin" +
		"cere@april.biz id:1 name:Leanne Graham phone:1-770-736-8031 x56442 username:Bret website:" +
		"hildegard.org]\n vs \nmap[address:map[city:Gwenborough geo:map[lat:-37.3159 lng:81.1496]" +
		" street:Kulas Light suite:Apt. 556 zipcode:92998-3874] company:map[bs:harness real-time " +
		"e-markets catchPhrase:Multi-layered client-server neural-net name:Romaguera-Crona]" +
		" email:Sincere@april.biz id:%!s(float64=1) name:Leanne Graham phone:1-770-736-8031" +
		" x56442 username:Bret website:hildegard.org]"
)

type replaceTestTable struct {
	name               string
	code               string
	request            string
	file               string
	expectedErr        error
	jsonValidator      validator.JSONFunctions
	jsonValidatorError string
	table              *godog.Table
}

type validateModifyingTestTable struct {
	name        string
	code        string
	request     string
	path        string
	method      string
	table       *godog.Table
	expectedErr error
}

func TestValidateResponseBodyJSONFileModifying(t *testing.T) {
	os.MkdirAll(logsPath, os.ModePerm)
	defer os.RemoveAll(logsPath)

	os.MkdirAll(schemasPath, os.ModePerm)
	os.WriteFile("./schemas/posts.json", []byte(validateModifyingResponseFile), os.ModePerm)
	os.WriteFile("./schemas/users.json", []byte(validateModifyingResponseFile), os.ModePerm)
	defer os.RemoveAll(schemasPath)

	paramsInput := make(map[string]interface{})
	paramsInput["title"] = "foo1"

	paramsWrongInput := make(map[string]interface{})
	paramsWrongInput["wrong_key"] = true

	tcs := validateModifyingTestCases()

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ctx, s := setRequestToTestJSONPlaceHolderContext(tc.request)
			request := tc.request
			s.SendRequestWithPathUsingJSON(ctx, tc.method, request, tc.path, tc.code)
			err := s.ValidateResponseBodyJSONFileModifying(ctx, tc.code, request, tc.table)
			require.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestValidateErrorBodyJSONFileReplaceString(t *testing.T) {
	httpmock.Activate()
	loadReplaceHTTPResponseMock()
	defer httpmock.DeactivateAndReset()

	os.MkdirAll(logsPath, os.ModePerm)
	defer os.RemoveAll(logsPath)
	os.MkdirAll(schemasPath, os.ModePerm)
	os.WriteFile("./schemas/test.json", []byte(JSONpostsfile), os.ModePerm)
	defer os.RemoveAll(schemasPath)
	filePathError := loadReplaceArgs()
	tcs := []replaceTestTable{
		{
			name:               "valid map response",
			code:               "example1",
			request:            "map_test",
			file:               "test",
			jsonValidator:      validator.JSONService{},
			jsonValidatorError: "",
			table:              generateReplaceTable(),
			expectedErr:        nil,
		},
		{
			name:               "valid string response",
			code:               "example5",
			request:            "string_test",
			file:               "test",
			jsonValidator:      validator.JSONService{},
			jsonValidatorError: "",
			table:              generateReplaceTable(),
			expectedErr:        nil,
		},
		{
			name:               "file error",
			code:               "example2",
			request:            "map_test",
			file:               "file_error",
			jsonValidator:      validator.JSONService{},
			jsonValidatorError: "",
			table:              generateReplaceTable(),
			expectedErr: fmt.Errorf("error getting parameter from json: %w",
				fmt.Errorf("error loading file at file_error due to error: %w",
					fmt.Errorf("%w", filePathError))),
		},
		{
			name:               "default case error",
			code:               "example4",
			request:            "map_test",
			jsonValidator:      validator.JSONService{},
			jsonValidatorError: "",
			file:               "test",
			table:              generateReplaceTable(),
			expectedErr:        fmt.Errorf("body content should be string or map: %v", "1"),
		},
		{
			name:               "error map response replace",
			code:               "example1",
			request:            "map_test",
			file:               "test",
			jsonValidator:      validator.JSONServiceMock{},
			jsonValidatorError: "error in replace map string function",
			table:              generateReplaceTable(),
			expectedErr:        fmt.Errorf("error in replace map string function"),
		},
		{
			name:               "error string response replace",
			code:               "example6",
			request:            "string_test",
			file:               "test",
			jsonValidator:      validator.JSONServiceMock{},
			jsonValidatorError: "error in replace string function",
			table:              generateReplaceTable(),
			expectedErr:        fmt.Errorf("error in replace string function"),
		},
		{
			name:               "input table error",
			code:               "example1",
			request:            "map_test",
			file:               "test",
			jsonValidator:      validator.JSONService{},
			jsonValidatorError: "",
			table:              generateWrongReplaceTable(),
			expectedErr: fmt.Errorf("cannot remove header: %v",
				fmt.Errorf("table must have at least one header and one useful row")),
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ctx, s := setRequestToTestLocalhostContext(tc.request, tc.jsonValidator, tc.jsonValidatorError)

			s.SendRequestTo(ctx, http.MethodGet, tc.request)

			err := s.ValidateErrorBodyJSONFileReplace(ctx, tc.code, tc.file, tc.table)
			if tc.expectedErr != nil && err == nil {
				t.Errorf("expecting error but not found")
			}
		})
	}
}

func validateModifyingTestCases() []validateModifyingTestTable {
	tcs := []validateModifyingTestTable{
		{
			name:    "valid replace",
			code:    "example1",
			request: "posts",
			method:  http.MethodPost,
			path:    "",
			table: golium.NewTable([][]string{
				{"parameter", "value"},
				{"title", "foo1"},
			}),
			expectedErr: nil,
		},
		{
			name:    "valid replace nested",
			code:    "example2",
			request: "users",
			method:  http.MethodGet,
			path:    "1",
			table: golium.NewTable([][]string{
				{"parameter", "value"},
				{"company.name", "Romaguera-Crona"},
			}),
			expectedErr: nil,
		},
		{
			name:    "not present nested key",
			code:    "example2",
			request: "users",
			method:  http.MethodGet,
			path:    "1",
			table: golium.NewTable([][]string{
				{"parameter", "value"},
				{"not-present.name", "Romaguera-Crona"},
			}),
			expectedErr: fmt.Errorf("error processing params: error modifying" +
				" nested param: param not-present does not exists"),
		},
		{
			name:    "validation nested key",
			code:    "example2",
			request: "users",
			method:  http.MethodGet,
			path:    "1",
			table: golium.NewTable([][]string{
				{"parameter", "value"},
				{"address.geo.lat", "fake"},
			}),
			expectedErr: errors.New(validateModifyNestedError),
		},
		{
			name:    "wrong code",
			code:    "example3",
			request: "posts",
			method:  http.MethodPost,
			path:    "",
			table: golium.NewTable([][]string{
				{"parameter", "value"},
				{"title", "foo1"},
			}),
			expectedErr: fmt.Errorf("error getting parameter from json: %w",
				fmt.Errorf(validateModifyingWrongCodeError+" due to error: %w",
					fmt.Errorf("value for param: 'response' with code: 'example3' not found"))),
		},
		{
			name:    "wrong param",
			code:    "example1",
			request: "posts",
			method:  http.MethodPost,
			path:    "",
			table: golium.NewTable([][]string{
				{"parameter", "value"},
				{"wrong_key", "true"},
			}),
			expectedErr: fmt.Errorf("error modifying param: param wrong_key does not exists"),
		},
		{
			name:    "input table error",
			code:    "example1",
			request: "posts",
			method:  http.MethodPost,
			path:    "",
			table: golium.NewTable([][]string{
				{"title", "foo1"},
			}),
			expectedErr: fmt.Errorf("cannot remove header: %v",
				fmt.Errorf("table must have at least one header and one useful row")),
		},
	}
	return tcs
}

func setRequestToTestLocalhostContext(endpoint string,
	jsonValidator validator.JSONFunctions,
	jsonValidatorError string,
) (context.Context, *Session) {
	ValuesAsString = map[string]string{
		"[CONF:url]":                                "http://localhost/",
		"[CTXT:url]":                                NilString,
		"[CONF:endpoints.map_test.api-endpoint]":    endpoint,
		"[CONF:endpoints.string_test.api-endpoint]": endpoint,
	}
	validator.JSON = jsonValidator
	if jsonValidatorError != "" {
		validator.ErrorResponse = jsonValidatorError
	}
	ctx, s := setGoliumContextAndService()
	return ctx, s
}

func setRequestToTestJSONPlaceHolderContext(request string) (context.Context, *Session) {
	ValuesAsString = map[string]string{
		"[CONF:url]":                          "https://jsonplaceholder.typicode.com/",
		"[CTXT:url]":                          NilString,
		"[CONF:endpoints.posts.api-endpoint]": request,
		"[CONF:endpoints.users.api-endpoint]": request,
	}
	FakeResponse = ""

	ctx, s := setGoliumContextAndService()
	return ctx, s
}

func loadReplaceHTTPResponseMock() {
	httpmock.RegisterResponder(http.MethodGet, "http://localhost/map_test/",
		httpmock.NewStringResponder(400, replaceMapResponseBody))
	httpmock.RegisterResponder(http.MethodGet, "http://localhost/string_test/",
		httpmock.NewStringResponder(400, replaceStringResponseBody))
}

func loadReplaceArgs() error {
	absPath, _ := filepath.Abs("./schemas/file_error.json")

	fsError := &fs.PathError{
		Op:   "open",
		Path: absPath,
		Err:  syscall.ENOENT,
	}
	filePathError := fmt.Errorf(
		"error reading file at %v due to error: %v",
		absPath, fsError.Error())
	return filePathError
}

func generateReplaceTable() *godog.Table {
	return golium.NewTable([][]string{
		{"parameter", "value"},
		{"field", "field_name"},
		{"type", "field_type"},
	})
}

func generateWrongReplaceTable() *godog.Table {
	return golium.NewTable([][]string{
		{"field", "field_name"},
	})
}
