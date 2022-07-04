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

package http

import (
	"context"
	"fmt"
)

var (
	FakeResponse   string
	ValuesAsString map[string]string
)

type GoliumInterfaceMock struct{}

func (g GoliumInterfaceMock) ValueAsString(ctx context.Context, s string) string {
	return ValuesAsString[s]
}

func (g GoliumInterfaceMock) SendHTTPRequest(ctx context.Context, method string) error {
	switch FakeResponse {
	case "error":
		return fmt.Errorf("error with the HTTP request. %v", "fake_error")

	default:
		return GetSession(ctx).SendHTTPRequest(ctx, method)
	}
}
