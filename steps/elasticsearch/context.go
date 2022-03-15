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

package elasticsearch

import (
	"context"
)

// ContextKey defines a type to store the elasticsearch session in context.Context.
type ContextKey string

const contextKey ContextKey = "elasticsearchSession"

// InitializeContext adds the elasticsearch session to the context.
// The new context is returned because context is immutable.
func InitializeContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, contextKey, &Session{})
}

// GetSession returns the elasticsearch session stored in context.
// Note that the context should be previously initialized with InitializeContext function.
func GetSession(ctx context.Context) *Session {
	return ctx.Value(contextKey).(*Session)
}
