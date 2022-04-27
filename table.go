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
	"errors"
	"fmt"
	"net/url"
	"reflect"

	"github.com/cucumber/godog"
	"github.com/cucumber/messages-go/v16"
)

// Remove headers form table
func RemoveHeaders(t *godog.Table) error {
	if len(t.Rows) < 2 {
		return errors.New("cannot remove header: table must have at least one header and one useful row")
	}
	t.Rows = t.Rows[1:]
	return nil
}

func ColumnsChecker(cells []*messages.PickleTableCell, n int) error {
	if len(cells) != n {
		return fmt.Errorf("table must have %d columns", n)
	}
	return nil
}

// NewTable Aux function that creates a new table
// from string matrix for testing purposes.
func NewTable(src [][]string) *godog.Table {
	rows := make([]*messages.PickleTableRow, len(src))

	for i, row := range src {
		cells := make([]*messages.PickleTableCell, len(row))

		for j, value := range row {
			cells[j] = &messages.PickleTableCell{Value: value}
		}

		rows[i] = &messages.PickleTableRow{Cells: cells}
	}
	return &godog.Table{Rows: rows}
}

// ConvertTableToMap converts a godog table with 2 columns into a map[string]interface{}.
func ConvertTableToMap(ctx context.Context, t *godog.Table) (map[string]interface{}, error) {
	err := RemoveHeaders(t)
	if err != nil {
		return nil, err
	}

	m := make(map[string]interface{})
	if len(t.Rows) == 0 {
		return m, nil
	}

	err = ColumnsChecker(t.Rows[0].Cells, 2)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(t.Rows); i++ {
		cells := t.Rows[i].Cells
		propKey := cells[0].Value
		propValue := cells[1].Value
		m[propKey] = Value(ctx, propValue)
	}
	return m, nil
}

// ConvertTableColumnToArray converts a godog table with 1 column into a []string.
func ConvertTableColumnToArray(ctx context.Context, t *godog.Table) ([]string, error) {
	err := RemoveHeaders(t)
	if err != nil {
		return nil, err
	}
	m := []string{}
	if len(t.Rows) == 0 {
		return m, nil
	}

	err = ColumnsChecker(t.Rows[0].Cells, 1)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(t.Rows); i++ {
		cells := t.Rows[i].Cells
		propKey := cells[0].Value
		m = append(m, propKey)
	}
	return m, nil
}

// ConvertTableToMultiMap converts a godog table with 2 columns into a map[string][]string.
// The multimap is using url.Values.
// The multimap is useful to support multiple values for the same key (e.g. for query parameters
// or HTTP headers).
func ConvertTableToMultiMap(ctx context.Context, t *godog.Table) (map[string][]string, error) {
	err := RemoveHeaders(t)
	if err != nil {
		return nil, err
	}
	m := url.Values{}
	if len(t.Rows) == 0 {
		return m, nil
	}

	err = ColumnsChecker(t.Rows[0].Cells, 2)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(t.Rows); i++ {
		cells := t.Rows[i].Cells
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

func ConvertTableWithHeaderToStructSlice(ctx context.Context,
	t *godog.Table,
	slicePtr interface{},
) error {
	if len(t.Rows) == 0 {
		return errors.New("table requires at least 1 row with the header")
	}
	if len(t.Rows) == 1 {
		// No data
		return nil
	}

	if reflect.TypeOf(slicePtr).Kind() != reflect.Ptr {
		return errors.New("expected a pointer to an slice of structs")
	}
	slicePtrValue := reflect.ValueOf(slicePtr)
	sliceValue := slicePtrValue.Elem()
	sliceElemType := sliceValue.Type().Elem()

	header := t.Rows[0].Cells
	for i := 1; i < len(t.Rows); i++ {
		elemValue := reflect.New(sliceElemType).Elem()
		for n, cell := range t.Rows[i].Cells {
			if err := assignValue(elemValue, header[n].Value, Value(ctx, cell.Value)); err != nil {
				return fmt.Errorf("failed setting element '%s' in struct of type '%s': %w",
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
//              | param | value     |
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
// Warning: still pending process values directly as arrays, i.e.:
//			| addresses | ["http://localhost:8080"] |
// use by now a CONF tag, i.e.:
//			| addresses | [CONF:elasticsearch.addresses] |

func ConvertTableWithoutHeaderToStruct(ctx context.Context, t *godog.Table, v interface{}) error {
	err := RemoveHeaders(t)
	if err != nil {
		return err
	}
	if len(t.Rows) == 0 {
		return nil
	}

	err = ColumnsChecker(t.Rows[0].Cells, 2)
	if err != nil {
		return err
	}

	ptrValue := reflect.ValueOf(v)
	value := ptrValue.Elem()
	for i := 0; i < len(t.Rows); i++ {
		cells := t.Rows[i].Cells
		propKey := cells[0].Value
		propValue := cells[1].Value
		if err := assignValue(value, propKey, Value(ctx, propValue)); err != nil {
			errStr := fmt.Sprintf("failed setting element '%s' in struct of type '%s': %s",
				propKey, value.Type(), err.Error())
			return errors.New(errStr)
		}
	}
	return nil
}

func assignValue(destination reflect.Value, name string, value interface{}) error {
	fieldValueStr := fmt.Sprintf("%v", value)
	if err := exctractField(&destination, name); err != nil {
		return err
	}
	return StrategyFormat[destination.Kind()].format(destination, fieldValueStr, value)
}
