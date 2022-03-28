package golium

import (
	"context"
	"errors"
	"testing"

	"github.com/cucumber/godog"
)

var (
	table1x2 = [][]string{{"John", "182"}}
	table2x1 = [][]string{{"name"}, {"John"}}
	table2x2 = [][]string{{"name", "height"}, {"John", "182"}}
	table2x3 = [][]string{{"name", "height", "age"}, {"John", "182", "32"}}
)

func TestRemoveHeaders(t *testing.T) {
	tcs := []struct {
		name        string
		table       *godog.Table
		expectedErr error
	}{
		{
			name:  "Remove header expected error",
			table: NewTable(table1x2),
			expectedErr: errors.New(
				"cannot remove header: table must have at least one header and one useful row"),
		},
		{
			name:        "Remove header table ok",
			table:       NewTable(table2x2),
			expectedErr: nil,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			// Call the tested function
			err := RemoveHeaders(tc.table)

			// Check expected behavior
			if err != nil && err.Error() != tc.expectedErr.Error() {
				t.Errorf("Expected error: %s, Error received: %s", tc.expectedErr, err)
			}

			if tc.expectedErr == nil && err != nil {
				t.Errorf("Expected error: %s, Error received: %s", tc.expectedErr, err)
			}
		})
	}
}

func TestColumnsChecker(t *testing.T) {
	tcs := []struct {
		name        string
		table       *godog.Table
		n           int
		expectedErr error
	}{
		{
			name:  "Check columns less than n",
			table: NewTable(table2x1),
			n:     2,
			expectedErr: errors.New(
				"table must have 2 columns"),
		},
		{
			name:  "Check columns more than n",
			table: NewTable(table2x3),
			n:     2,
			expectedErr: errors.New(
				"table must have 2 columns"),
		},
		{
			name:        "Check columns ok",
			table:       NewTable(table2x2),
			n:           2,
			expectedErr: nil,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			// Call the tested function
			err := ColumnsChecker(tc.table.Rows[0].Cells, tc.n)

			// Check expected behavior
			if err != nil && err.Error() != tc.expectedErr.Error() {
				t.Errorf("Expected error: %s, Error received: %s", tc.expectedErr, err)
			}

			if tc.expectedErr == nil && err != nil {
				t.Errorf("Expected error: %s, Error received: %s", tc.expectedErr, err)
			}
		})
	}
}

func TestConvertTableToMap(t *testing.T) {
	tcs := []struct {
		name        string
		table       *godog.Table
		expectedErr error
	}{
		{
			name:        "Convert table to map ok",
			table:       NewTable(table2x2),
			expectedErr: nil,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			// Call the tested function
			_, err := ConvertTableToMap(ctx, tc.table)

			// Check expected behavior
			if err != nil && err.Error() != tc.expectedErr.Error() {
				t.Errorf("Expected error: %s, Error received: %s", tc.expectedErr, err)
			}

			if tc.expectedErr == nil && err != nil {
				t.Errorf("Expected error: %s, Error received: %s", tc.expectedErr, err)
			}
		})
	}
}

func TestConvertTableColumnToArray(t *testing.T) {
	tcs := []struct {
		name        string
		table       *godog.Table
		expectedErr error
	}{
		{
			name:        "Convert table column to array ok",
			table:       NewTable(table2x1),
			expectedErr: nil,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			// Call the tested function
			_, err := ConvertTableColumnToArray(ctx, tc.table)

			// Check expected behavior
			if err != nil && err.Error() != tc.expectedErr.Error() {
				t.Errorf("Expected error: %s, Error received: %s", tc.expectedErr, err)
			}

			if tc.expectedErr == nil && err != nil {
				t.Errorf("Expected error: %s, Error received: %s", tc.expectedErr, err)
			}
		})
	}
}

func TestConvertTableToMultiMap(t *testing.T) {
	tcs := []struct {
		name        string
		table       *godog.Table
		expectedErr error
	}{
		{
			name:        "Convert table to multi map ok",
			table:       NewTable(table2x2),
			expectedErr: nil,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			// Call the tested function
			_, err := ConvertTableToMultiMap(ctx, tc.table)

			// Check expected behavior
			if err != nil && err.Error() != tc.expectedErr.Error() {
				t.Errorf("Expected error: %s, Error received: %s", tc.expectedErr, err)
			}

			if tc.expectedErr == nil && err != nil {
				t.Errorf("Expected error: %s, Error received: %s", tc.expectedErr, err)
			}
		})
	}
}
