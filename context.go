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
)

// ContextKey defines a type to store the Context in context.Context.
type ContextKey string

const contextKey ContextKey = "contextKey"

// Context contains the context required for common utilities.
// It contains a map[string]interface{} to store global values and find them with [CTXT:xxx] tag.
type Context struct {
	m map[string]interface{}
}

// InitializeContext adds the Context to the context.
// The new context is returned because context is immutable.
func InitializeContext(ctx context.Context) context.Context {
	commonCtx := Context{
		m: make(map[string]interface{}),
	}
	return context.WithValue(ctx, contextKey, &commonCtx)
}

// GetContext returns the Context stored in context.
// Note that the context should be previously initialized with InitializeContext function.
func GetContext(ctx context.Context) *Context {
	return ctx.Value(contextKey).(*Context)
}

// Get returns an element from Context.Ctx.
// If the value does not exist, it returns nil.
func (c *Context) Get(key string) interface{} {
	return c.m[key]
}

// Put writes an element in the context.
func (c *Context) Put(key string, value interface{}) {
	c.m[key] = value
}
