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

package dns

import (
	"context"
	"time"

	"github.com/miekg/dns"
)

// ClientContextKey defines a type to store the ClientContext in context.Context.
type ClientContextKey string

var clientContextKey ClientContextKey = "clientContextKey"

// ClientContext contains the context for the DNS client steps (e.g. to validate the response).
type ClientContext struct {
	// Server is the address to the DNS server, including the server port (e.g. 8.8.8.8:53).
	Server string
	// Query contains the DNS request message.
	Query *dns.Msg
	// Response contains the DNS response message.
	Response *dns.Msg
	// RTT is the response time.
	RTT time.Duration
}

// InitializeClientContext adds the ClientContext to the context.
// The new context is returned because context is immutable.
func InitializeClientContext(ctx context.Context) context.Context {
	clientContext := ClientContext{}
	return context.WithValue(ctx, clientContextKey, &clientContext)
}

// GetClientContext returns the ClientContext stored in context.
// Note that the context should be previously initialized with InitializeClientContext function.
func GetClientContext(ctx context.Context) *ClientContext {
	return ctx.Value(clientContextKey).(*ClientContext)
}

// ServerContextKey defines a type to store the ServerContext in context.Context.
type ServerContextKey string

var serverContextKey ServerContextKey = "serverContextKey"

// ServerContext contains the context for the DNS server steps.
// There is usually one single DNS server but it could be possible to start up multiple instances.
// Each DNS server instance is represented by ServerInstanceContext and added to the map using the
// name as key.
type ServerContext struct {
	Instances map[string]*ServerInstanceContext
}

// ServerInstanceContext contains the context for each DNS server.
type ServerInstanceContext struct {
	Name     string
	Protocol string
	Port     uint
}

// InitializeServerContext adds the ServerContext to the context.
// The new context is returned because context is immutable.
func InitializeServerContext(ctx context.Context) context.Context {
	serverContext := ServerContext{
		Instances: make(map[string]*ServerInstanceContext),
	}
	return context.WithValue(ctx, serverContextKey, &serverContext)
}

// GetServerContext returns the ServerContext stored in context.
// Note that the context should be previously initialized with InitializeServerContext function.
func GetServerContext(ctx context.Context) *ServerContext {
	return ctx.Value(serverContextKey).(*ServerContext)
}
