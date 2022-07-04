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

	"github.com/TelefonicaTC2Tech/golium"
)

type ServiceFunctions interface {
	ValueAsString(ctx context.Context, s string) string
	SendHTTPRequest(ctx context.Context, method string) error
}

type GoliumInterface struct{}

func NewGoliumInterface() *GoliumInterface {
	return &GoliumInterface{}
}

func (g GoliumInterface) ValueAsString(ctx context.Context, s string) string {
	return golium.ValueAsString(ctx, s)
}

func (g GoliumInterface) SendHTTPRequest(ctx context.Context, method string) error {
	return GetSession(ctx).SendHTTPRequest(ctx, method)
}
