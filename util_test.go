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

import "testing"

func TestContainsString(t *testing.T) {
	list := []string{"attribute1", "attribute2", "attribute3"}

	tests := []struct {
		name           string
		values         []string
		expectedValue  string
		valueNonString int
		expectedResult bool
	}{
		{
			name:           "testing with a correct value",
			values:         list,
			expectedValue:  "attribute1",
			expectedResult: true,
		},
		{
			name:           "testing with a incorrect value",
			values:         list,
			expectedValue:  "failValue",
			expectedResult: false,
		},
		{
			name:           "testing with a value that is not a string",
			values:         list,
			valueNonString: 3,
			expectedResult: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ContainsString(tt.expectedValue, tt.values); got != tt.expectedResult {
				t.Errorf("ContainsString() = %v, expectedResult %v", got, tt.expectedResult)
			}
		})
	}
}
