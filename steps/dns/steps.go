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
	"encoding/hex"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/Telefonica/golium"
	"github.com/cucumber/godog"
	"github.com/miekg/dns"
)

// Steps type is responsible to initialize the DNS client steps in godog framework.
type Steps struct {
}

// InitializeSteps adds client DNS steps to the scenario context.
// It implements StepInitializer interface.
// It returns a new context (context is immutable) with the DNS Context.
func (s Steps) InitializeSteps(ctx context.Context, scenCtx *godog.ScenarioContext) context.Context {
	// Initialize the DNS session in the context
	ctx = InitializeContext(ctx)
	session := GetSession(ctx)

	// Initialize the steps
	scenCtx.Step(`^the DNS server "([^"]*)"$`, func(svr string) error {
		return session.ConfigureServer(ctx, golium.ValueAsString(ctx, svr))
	})
	scenCtx.Step(`^a DNS timeout of "([^"]*)" milliseconds$`, func(time string) error {
		timeout, err := strconv.Atoi(golium.ValueAsString(ctx, time))
		if err != nil {
			return fmt.Errorf("Error casting timeout parameter. %s", err)
		}
		return session.SetDNSResponseTimeout(ctx, timeout)
	})
	scenCtx.Step(`^the DNS query options$`, func(t *godog.Table) error {
		options, err := parseOptionsTable(ctx, t)
		if err != nil {
			return fmt.Errorf("Error parsing DNS query options. %s", err)
		}
		return session.ConfigureOptions(ctx, options)
	})
	scenCtx.Step(`^I send a DNS query of type "([^"]*)" for "([^"]*)"(\s\bwithout recursion\b)?$`, func(qtype, qname, recursion string) error {
		recursive := recursion == ""
		qtype = golium.ValueAsString(ctx, qtype)
		qname = golium.ValueAsString(ctx, qname)
		qt, ok := QueryTypes[qtype]
		if !ok {
			return fmt.Errorf("Invalid qtype: %s. Permitted values: %s", qtype, reflect.ValueOf(QueryTypes).MapKeys())
		}
		return session.SendQuery(ctx, qt, qname, recursive)
	})
	scenCtx.Step(`the DNS response must have the code "([^"]*)"$`, func(code string) error {
		return session.ValidateResponseWithCode(ctx, golium.ValueAsString(ctx, code))
	})
	scenCtx.Step(`the DNS response must have one of the following codes: "([^"]*)"$`, func(list string) error {
		codes := strings.Split(list, ",")
		for i := range codes {
			codes[i] = strings.TrimSpace(golium.ValueAsString(ctx, codes[i]))
		}
		return session.ValidateResponseWithOneOfCodes(ctx, codes)
	})
	scenCtx.Step(`the DNS response must have "(\d+)" ((\banswer\b)|(\bauthority\b)|(\badditional\b)) records?$`, func(n int, recordType string) error {
		recordType = golium.ValueAsString(ctx, recordType)
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

func parseOptionsTable(ctx context.Context, t *godog.Table) ([]dns.EDNS0, error) {
	type option struct {
		Code uint16
		Data string // Hexadecimal string
	}
	var options []option
	if err := golium.ConvertTableWithHeaderToStructSlice(ctx, t, &options); err != nil {
		return nil, fmt.Errorf("Error mapping table with option struct. %s", err)
	}
	dnsOptions := make([]dns.EDNS0, len(options))
	for i, o := range options {
		data, err := hex.DecodeString(o.Data)
		if err != nil {
			//return nil, fmt.Errorf("Error converting to byte array: %s. %s", o.Data, err)
			data = []byte(o.Data)
		}
		dnsOptions[i] = &dns.EDNS0_LOCAL{
			Code: o.Code,
			Data: data,
		}
	}
	return dnsOptions, nil
}

func validateResponseWithRecords(ctx context.Context, session *Session, recordType RecordType, t *godog.Table) error {
	expectedRecords := []Record{}
	if err := golium.ConvertTableWithHeaderToStructSlice(ctx, t, &expectedRecords); err != nil {
		return fmt.Errorf("Error processing the table with the expected records. %s", err)
	}
	return session.ValidateResponseWithRecords(ctx, recordType, expectedRecords)
}
