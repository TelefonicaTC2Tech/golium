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

package aggregate

import (
	"context"

	"github.com/Telefonica/golium/steps/http"
)

// AggregateSession contains the information of shared session
type AggregateSession struct {
	session *http.Session
}

// SaveStatusCode saves code in shared session.
func (a *AggregateSession) SaveStatusCode(ctx context.Context, code int) error {
	return a.session.SaveStatusCode(ctx, code)
}
