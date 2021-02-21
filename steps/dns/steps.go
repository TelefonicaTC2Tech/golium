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
	"reflect"
	"strings"

	"github.com/Telefonica/golium"
	"github.com/cucumber/godog"
)

// Steps type is responsible to initialize the DNS client steps in godog framework.
type Steps struct {
}

// InitializeSteps adds client DNS steps to the scenario context.
// It implements StepInitializer interface.
// It returns a new context (context is immutable) with the ClientContext.
func (s Steps) InitializeSteps(ctx context.Context, scenCtx *godog.ScenarioContext) context.Context {
	// Initialize the DNS session in the context
	ctx = InitializeContext(ctx)
	session := GetSession(ctx)

	// Initialize the steps
	scenCtx.Step(`^I configure the DNS server at "([^"]*)"$`, func(svr string) error {
		return session.ConfigureServer(ctx, svr)
	})
	scenCtx.Step(`^I send a DNS query of type "([^"]*)" for "([^"]*)"(\s\bwithout recursion\b)?$`, func(qtype, qname, recursion string) error {
		recursive := recursion == ""
		qt, ok := QueryTypes[qtype]
		if !ok {
			return fmt.Errorf("Invalid qtype: %s. Permitted values: %s", qtype, reflect.ValueOf(QueryTypes).MapKeys())
		}
		return session.SendQuery(ctx, qt, qname, recursive)
	})
	scenCtx.Step(`the DNS response must have the code "([^"]*)"$`, func(code string) error {
		return session.ValidateResponseWithCode(ctx, code)
	})
	scenCtx.Step(`the DNS response must have one of the following codes: "([^"]*)"$`, func(list string) error {
		codes := strings.Split(list, ",")
		for i := range codes {
			codes[i] = strings.TrimSpace(codes[i])
		}
		return session.ValidateResponseWithOneOfCodes(ctx, codes)
	})
	scenCtx.Step(`the DNS response must have "(\d+)" ((\banswer\b)|(\bauthority\b)|(\badditional\b)) records?$`, func(n int, recordType string) error {
		return session.ValidateResponseWithNumberOfRecords(ctx, n, RecordType(recordType))
	})
	scenCtx.Step(`the DNS response must contain the following answer records?$`, func(t *godog.Table) error {
		return validateResponseWithRecords(ctx, session, Answer, t)
	})
	scenCtx.Step(`the DNS response must contain the following authority records?$`, func(t *godog.Table) error {
		return validateResponseWithRecords(ctx, session, Authority, t)
	})
	scenCtx.Step(`the DNS response must contain the following additional records?$`, func(t *godog.Table) error {
		return validateResponseWithRecords(ctx, session, Additional, t)
	})
	return ctx
}

func validateResponseWithRecords(ctx context.Context, session *Session, recordType RecordType, t *godog.Table) error {
	expectedRecords := []Record{}
	if err := golium.ConvertTableWithHeaderToStructSlice(ctx, t, &expectedRecords); err != nil {
		return fmt.Errorf("Error processing the table with the expected records. %s", err)
	}
	return session.ValidateResponseWithRecords(ctx, RecordType(recordType), expectedRecords)
}
