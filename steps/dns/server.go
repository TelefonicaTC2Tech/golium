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
	"fmt"
	"strconv"
	"strings"

	"github.com/cucumber/godog"
	"github.com/miekg/dns"
)

// ServerSteps type is responsible to initialize the DNS server steps in godog framework.
type ServerSteps struct {
}

// InitializeSteps adds DNS server steps to the scenario context.
// It implements StepInitializer interface.
// It returns a new context (context is immutable) with the ServerContext.
func (ss ServerSteps) InitializeSteps(ctx context.Context, scenCtx *godog.ScenarioContext) context.Context {
	// Add the ServerContext to context
	ctx = InitializeServerContext(ctx)
	// Update the ServerContext
	//serverContext := GetServerContext(ctx)

	// Initialize the steps
	scenCtx.Step(`^I start up a DNS server$`, func(serverConfig *godog.Table) error {
		return StartUpServerStep(ctx, serverConfig)
	})
	// scenCtx.Step(`^I configure on DNS server "" the following resolutions$`, func(name string, resolutions *godog.Table) error {
	// 	return ConfigureResolutionsServerStep(ctx, name, resolutions)
	// })

	return ctx
}

// StartUpServerStep starts up a DNS server using the configuration from the config table.
func StartUpServerStep(ctx context.Context, serverConfig *godog.Table) error {
	if len(serverConfig.Rows) == 0 {
		return fmt.Errorf("No row found for server config")
	}
	if len(serverConfig.Rows[0].Cells) != 2 {
		return fmt.Errorf("Server config must have 2 columns: key and value")
	}
	serverInstanceContext := &ServerInstanceContext{
		Name:     "",
		Protocol: "udp",
	}
	for _, row := range serverConfig.Rows {
		switch row.Cells[0].Value {
		case "name":
			serverInstanceContext.Name = row.Cells[1].Value
		case "port":
			u64, err := strconv.ParseUint(row.Cells[1].Value, 10, 32)
			if err != nil {
				return fmt.Errorf("DNS server port must be an unsigned integer")
			}
			serverInstanceContext.Port = uint(u64)
		case "protocol":
			protocol := strings.ToLower(row.Cells[1].Value)
			if protocol != "udp" {
				return fmt.Errorf("DNS server protocol must be udp")
			}
		}
	}
	// Port is mandatory
	if serverInstanceContext.Port == 0 {
		return fmt.Errorf("DNS server port is mandatory")
	}
	// Store the instance context in the server context using the name as the key
	serverCtx := GetServerContext(ctx)
	serverCtx.Instances[serverInstanceContext.Name] = serverInstanceContext
	// Start up the server
	server := &dns.Server{
		Addr:    fmt.Sprintf(":%d", serverInstanceContext.Port),
		Net:     serverInstanceContext.Protocol,
		Handler: &handler{instance: serverInstanceContext},
	}
	go server.ListenAndServe()
	return nil
}

// func ConfigureResolutionsServerStep(ctx context.Context, name string, resolutions *godog.Table) error {

// }

type handler struct {
	instance *ServerInstanceContext
}

func (h *handler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	fmt.Printf("in handler with %+v", r)
}
