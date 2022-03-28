package golium

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/cucumber/godog"
	"github.com/stretchr/testify/require"
)

func TestRemoveHeaders(t *testing.T) {
	tcs := []struct {
		name        string
		table       *godog.Table
		expectedErr error
	}{
		{
			name: "Remove header expected error",
			table: NewTable([][]string{
				{"John", "182"}, //  | John   | 182 |
			}),
			expectedErr: errors.New("table must have at least one header and one useful row"),
		},
		{
			name: "Remove header table ok",
			table: NewTable([][]string{
				{"name", "height"},
				{"John", "182"}, //  | John   | 182 |
			}),
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

func TestGetParamsFromTable(t *testing.T) {
	expectedResult := make(map[string]interface{})
	expectedResult["Name"] = "John"
	tcs := []struct {
		name           string
		table          *godog.Table
		expectedResult interface{}
		expectedErr    error
	}{
		{
			name: "Error from remove header",
			table: NewTable([][]string{
				{"Parameter", "Value"},
			}),
			expectedErr: fmt.Errorf("failed removing headers from table: %w",
				errors.New("table must have at least one header and one useful row")),
		},
		{
			name: "Converted ok",
			table: NewTable([][]string{
				{"Parameter", "Value"},
				{"Name", "John"},
			}),
			expectedErr:    nil,
			expectedResult: expectedResult,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			// Call the tested function
			ctx := context.Background()
			convertedTable, err := GetParamsFromTable(ctx,
				tc.table,
				func(context.Context, *godog.Table,
				) (interface{}, error) {
					return ConvertTableToMap(ctx, tc.table)
				})

			// Check expected behavior
			require.Equal(t, tc.expectedErr, err)
			require.Equal(t, tc.expectedResult, convertedTable)
		})
	}
}
