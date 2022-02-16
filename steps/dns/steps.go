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
	scenCtx.Step(`^the DNS server "([^"]*)"$`, func(svr string) {
		// DNS transport protocol is set to UDP by default
		transport := "UDP"
		session.ConfigureServer(ctx, golium.ValueAsString(ctx, svr), transport)
	})
	scenCtx.Step(`^the DNS server "([^"]*)" on "([^"]*)"$`, func(svr, transport string) {
		session.ConfigureServer(ctx, golium.ValueAsString(ctx, svr), golium.ValueAsString(ctx, transport))
	})
	scenCtx.Step(`^a DNS timeout of "([^"]*)" milliseconds$`, func(timeout string) error {
		to, err := golium.ValueAsInt(ctx, timeout)
		if err != nil {
			return fmt.Errorf("invalid timeout '%s': %w", timeout, err)
		}
		session.SetDNSResponseTimeout(ctx, to)
		return nil
	})
	scenCtx.Step(`^the DNS query options$`, func(t *godog.Table) error {
		options, err := parseOptionsTable(ctx, t)
		if err != nil {
			return fmt.Errorf("failed parsing DNS query options: %w", err)
		}
		session.ConfigureOptions(ctx, options)
		return nil
	})
	scenCtx.Step(`^I send a DNS query of type "([^"]*)" for "([^"]*)"(\s\bwithout recursion\b)?$`, func(qtype, qname, recursion string) error {
		recursive := recursion == ""
		qtype = golium.ValueAsString(ctx, qtype)
		qname = golium.ValueAsString(ctx, qname)
		qt, ok := QueryTypes[qtype]
		if !ok {
			return fmt.Errorf("invalid qtype '%s': permitted values '%s'", qtype, reflect.ValueOf(QueryTypes).MapKeys())
		}
		switch session.Transport {
		case "UDP":
			return session.SendUDPQuery(ctx, qt, qname, recursive)
		case "DoH with GET":
			return session.SendDoHQuery(ctx, "GET", qt, qname, recursive)
		case "DoH with POST":
			return session.SendDoHQuery(ctx, "POST", qt, qname, recursive)
		case "DoT":
			return session.SendDoTQuery(ctx, qt, qname, recursive)
		default:
			return fmt.Errorf("unsupported transport protocol. %s", session.Transport)
		}
	})
	scenCtx.Step(`^the DoH query parameters$`, func(t *godog.Table) error {
		params, err := golium.ConvertTableToMultiMap(ctx, t)
		if err != nil {
			return fmt.Errorf("failed processing query parameters from table: %w", err)
		}
		session.ConfigureDoHQueryParams(ctx, params)
		return nil
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
	scenCtx.Step(`the DNS response must have "(\d+)" ((\banswer\b)|(\bauthority\b)|(\badditional\b)) records?$`, func(number string, recordType string) error {
		n, err := golium.ValueAsInt(ctx, number)
		if err != nil {
			return fmt.Errorf("invalid number '%s': %w", number, err)
		}
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
		return nil, fmt.Errorf("failed mapping table with option struct: %w", err)
	}
	dnsOptions := make([]dns.EDNS0, len(options))
	for i, o := range options {
		data, err := hex.DecodeString(o.Data)
		if err != nil {
			// Do not raise an error if data is not hexadecimal. It converts the string to a byte slice.
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
		return fmt.Errorf("failed processing the table with the expected records: %w", err)
	}
	return session.ValidateResponseWithRecords(ctx, recordType, expectedRecords)
}
