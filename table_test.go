package golium

import (
	"errors"
	"testing"

	"github.com/cucumber/godog"
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
