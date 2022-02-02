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

package shared

import (
	"context"
	"fmt"
)

// Session contains the information of shared session
type Session struct {
	StatusCode int
}

// SaveStatusCode save the status code.
func (s *Session) SaveStatusCode(ctx context.Context, code int) error {
	s.StatusCode = code
	return nil
}

// ValidateStatusCode validates the status code from the base response.
func (s *Session) ValidateSharedStatusCode(ctx context.Context, expectedCode int) error {
	if expectedCode != s.StatusCode {
		return fmt.Errorf("status code mismatch: expected '%d', actual '%d'", expectedCode, s.StatusCode)
	}
	return nil
}
