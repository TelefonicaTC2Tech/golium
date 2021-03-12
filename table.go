// Copyright 2021 Telefonica Cybersecurity & Cloud Tech SL
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package golium

import (
	"context"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/cucumber/godog"
)

// ConvertTableToMap converts a godog table with 2 columns into a map[string]interface{}.
func ConvertTableToMap(ctx context.Context, t *godog.Table) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	if len(t.Rows) == 0 {
		return m, nil
	}
	for i := 0; i < len(t.Rows); i++ {
		cells := t.Rows[i].Cells
		if len(cells) != 2 {
			return m, fmt.Errorf("Table must have 2 columns")
		}
		propKey := cells[0].Value
		propValue := cells[1].Value
		m[propKey] = Value(ctx, propValue)
	}
	return m, nil
}

// ConvertTableToMultiMap converts a godog table with 2 columns into a map[string][]string.
// The multimap is using url.Values.
// The multimap is useful to support multiple values for the same key (e.g. for query parameters
// or HTTP headers).
func ConvertTableToMultiMap(ctx context.Context, t *godog.Table) (map[string][]string, error) {
	m := url.Values{}
	if len(t.Rows) == 0 {
		return m, nil
	}
	for i := 0; i < len(t.Rows); i++ {
		cells := t.Rows[i].Cells
		if len(cells) != 2 {
			return m, fmt.Errorf("Table must have 2 columns")
		}
		propKey := ValueAsString(ctx, cells[0].Value)
		propValue := ValueAsString(ctx, cells[1].Value)
		m.Add(propKey, propValue)
	}
	return m, nil
}

// ConvertTableWithHeaderToStructSlice converts a godog table, where the first row is a header row,
// and the rest of rows are the values.
// For each column, the header value specifies the property in the struct (case insensitive).
// The first argument is the godog table, and the second one is a pointer to a slice of structs.
// The array is filled via reflection.
//
// With the following example table:
//		Scenario: Table to struct test
//			Given I configure my struct slice
//				| Name      | Value |
//				| example 1 | 1     |
//				| example 2 | 10    |
// And the code:
//		scenCtx.Step(`^I configure my struct slice$`, func(table *godog.Table) error {
//			type TestElement struct {
//				Name  string
//      		Value int
//    		}
//			testSlice := []TestElement{}
//			err := golium.ConvertTableWithHeaderToStructSlice(ctx, table, &testSlice)
// 		})
// It will be equivalent to:
//		testSlice := []TestElement{
//			TestElement{Name: "example 1", Value: 1},
//			TestElement{Name: "example 2", Value: 10},
//		}
func ConvertTableWithHeaderToStructSlice(ctx context.Context, t *godog.Table, slicePtr interface{}) error {
	if len(t.Rows) == 0 {
		return fmt.Errorf("Table requires at least 1 row with the header")
	}
	if len(t.Rows) == 1 {
		// No data
		return nil
	}

	if reflect.TypeOf(slicePtr).Kind() != reflect.Ptr {
		return fmt.Errorf("Expected a pointer to an slice of structs")
	}
	slicePtrValue := reflect.ValueOf(slicePtr)
	sliceValue := slicePtrValue.Elem()
	sliceElemType := sliceValue.Type().Elem()

	header := t.Rows[0].Cells
	for i := 1; i < len(t.Rows); i++ {
		elemValue := reflect.New(sliceElemType).Elem()
		for n, cell := range t.Rows[i].Cells {
			if err := assignFieldInStruct(elemValue, header[n].Value, Value(ctx, cell.Value)); err != nil {
				return fmt.Errorf("Error setting element '%s' in struct of type '%s'. %s",
					header[n].Value, sliceElemType, err)
			}
		}
		sliceValue.Set(reflect.Append(sliceValue, elemValue))
	}

	return nil
}

// ConvertTableWithoutHeaderToStruct converts a godog table with two columns into a struct.
// The first column of the table corresponds to the struct property, and the seconds column
// to the value to be assigned.
//
// With the following example table:
//		Scenario: Table to struct test
//			Given I configure my struct
//				| Name  | 1         |
//				| Value | example 1 |
// And the code:
//		scenCtx.Step(`^I configure my struct$`, func(table *godog.Table) error {
//			type TestElement struct {
//				Name  string
//      		Value int
//    		}
//			testElement := TestElement{}
//			err := golium.ConvertTableWithoutHeaderToStruct(ctx, table, &testElement)
// 		})
// It will be equivalent to:
//		testElement := TestElement{Name: "example 1", Value: 1}
func ConvertTableWithoutHeaderToStruct(ctx context.Context, t *godog.Table, v interface{}) error {
	if len(t.Rows) == 0 {
		return nil
	}
	ptrValue := reflect.ValueOf(v)
	value := ptrValue.Elem()
	for i := 0; i < len(t.Rows); i++ {
		cells := t.Rows[i].Cells
		if len(cells) != 2 {
			return fmt.Errorf("Table must have 2 columns")
		}
		propKey := cells[0].Value
		propValue := cells[1].Value
		if err := assignFieldInStruct(value, propKey, Value(ctx, propValue)); err != nil {
			return fmt.Errorf("Error setting element '%s' in struct of type '%s'. %s",
				propKey, value.Type(), err)
		}
	}
	return nil
}

func assignFieldInStruct(value reflect.Value, fieldName string, fieldValue interface{}) error {
	if value.Kind() == reflect.Ptr {
		value = reflect.Indirect(value)
	}
	if value.Kind() != reflect.Struct {
		return fmt.Errorf("orig must be a struct")
	}
	f := value.FieldByNameFunc(func(n string) bool {
		return strings.ToLower(n) == strings.ToLower(fieldName)
	})
	if !f.IsValid() {
		return fmt.Errorf("Field %s is not valid", fieldName)
	}
	if !f.CanSet() {
		return fmt.Errorf("Field %s cannot be set", fieldName)
	}
	if f.Kind() == reflect.Ptr {
		if fieldValue == nil {
			f.SetPointer(nil)
			return nil
		}
		fv := reflect.New(f.Type().Elem())
		f.Set(fv)
		f = fv.Elem()
	}
	fieldValueStr := fmt.Sprintf("%v", fieldValue)
	if f.Kind() == reflect.String {
		f.SetString(fieldValueStr)
		return nil
	}
	if f.Kind() == reflect.Bool {
		v, err := strconv.ParseBool(fieldValueStr)
		if err != nil {
			return fmt.Errorf("Error parsing to boolean the field '%s' with value '%s'", fieldName, fieldValueStr)
		}
		f.SetBool(v)
		return nil
	}
	if f.Kind() == reflect.Int || f.Kind() == reflect.Int8 || f.Kind() == reflect.Int16 || f.Kind() == reflect.Int32 || f.Kind() == reflect.Int64 {
		v, err := strconv.ParseInt(fieldValueStr, 10, 64)
		if err != nil {
			return fmt.Errorf("Error parsing to integer the field '%s' with value '%s'", fieldName, fieldValueStr)
		}
		f.SetInt(v)
		return nil
	}
	if f.Kind() == reflect.Uint || f.Kind() == reflect.Uint8 || f.Kind() == reflect.Uint16 || f.Kind() == reflect.Uint32 || f.Kind() == reflect.Uint64 {
		v, err := strconv.ParseUint(fieldValueStr, 10, 64)
		if err != nil {
			return fmt.Errorf("Error parsing to unsigned integer the field '%s' with value '%s'", fieldName, fieldValueStr)
		}
		f.SetUint(v)
		return nil
	}
	if f.Kind() == reflect.Float32 || f.Kind() == reflect.Float64 {
		v, err := strconv.ParseFloat(fieldValueStr, 64)
		if err != nil {
			return fmt.Errorf("Error parsing to float the field '%s' with value '%s'", fieldName, fieldValueStr)
		}
		f.SetFloat(v)
		return nil
	}
	if f.Kind() == reflect.Complex64 || f.Kind() == reflect.Complex128 {
		v, err := strconv.ParseComplex(fieldValueStr, 128)
		if err != nil {
			return fmt.Errorf("Error parsing to complex the field '%s' with value '%s'", fieldName, fieldValueStr)
		}
		f.SetComplex(v)
		return nil
	}
	return nil
}
