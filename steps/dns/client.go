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
	"github.com/miekg/dns"
)

// ClientSteps type is responsible to initialize the DNS client steps in godog framework.
type ClientSteps struct {
	Server string
}

// InitializeSteps adds client DNS steps to the scenario context.
// It implements StepInitializer interface.
// It returns a new context (context is immutable) with the ClientContext.
func (cs ClientSteps) InitializeSteps(ctx context.Context, scenCtx *godog.ScenarioContext) context.Context {
	// Add the ClientContext to context
	ctx = InitializeClientContext(ctx)
	// Update the ClientContext with the server information (if configured)
	clientContext := GetClientContext(ctx)
	clientContext.Server = cs.Server

	// Initialize the steps
	scenCtx.Step(`^I configure the DNS server at "([^"]*)"$`, func(svr string) error {
		return ConfigureServerClientStep(ctx, svr)
	})
	scenCtx.Step(`^I send a DNS query of type "([^"]*)" for "([^"]*)"(\s\bwithout recursion\b)?$`, func(qtype, qname, recursion string) error {
		recursive := recursion == ""
		qt, ok := QueryTypes[qtype]
		if !ok {
			return fmt.Errorf("Invalid qtype: %s. Permitted values: %s", qtype, reflect.ValueOf(QueryTypes).MapKeys())
		}
		return SendQueryClientStep(ctx, qt, qname, recursive)
	})
	scenCtx.Step(`the DNS response must have the code "([^"]*)"$`, func(code string) error {
		return ValidateResponseWithCodeClientStep(ctx, code)
	})
	scenCtx.Step(`the DNS response must have one of the following codes: "([^"]*)"$`, func(list string) error {
		codes := strings.Split(list, ",")
		for i := range codes {
			codes[i] = strings.TrimSpace(codes[i])
		}
		return ValidateResponseWithOneOfCodesClientStep(ctx, codes)
	})
	scenCtx.Step(`the DNS response must have "(\d+)" ((\banswer\b)|(\bauthority\b)|(\badditional\b)) records?$`, func(n int, recordType string) error {
		return ValidateResponseWithNumberOfRecordsClientStep(ctx, n, RecordType(recordType))
	})
	scenCtx.Step(`the DNS response must contain the following answer records?$`, func(expectedRecords *godog.Table) error {
		return ValidateResponseWithRecordsClientStep(ctx, Answer, expectedRecords)
	})
	scenCtx.Step(`the DNS response must contain the following authority records?$`, func(expectedRecords *godog.Table) error {
		return ValidateResponseWithRecordsClientStep(ctx, Authority, expectedRecords)
	})
	scenCtx.Step(`the DNS response must contain the following additional records?$`, func(expectedRecords *godog.Table) error {
		return ValidateResponseWithRecordsClientStep(ctx, Additional, expectedRecords)
	})
	return ctx
}

// ConfigureServerClientStep configures the DNS server location.
// This step may override the DNS server initially configured.
func ConfigureServerClientStep(ctx context.Context, svr string) error {
	clientCtx := GetClientContext(ctx)
	clientCtx.Server = svr
	return nil
}

// SendQueryClientStep sends a DNS query to resolve a domain.
func SendQueryClientStep(ctx context.Context, qtype uint16, qdomain string, recursive bool) error {
	clientCtx := GetClientContext(ctx)
	c := dns.Client{}
	m := dns.Msg{}
	m.SetQuestion(dns.Fqdn(qdomain), qtype)
	m.RecursionDesired = recursive
	clientCtx.Query = &m
	r, rtt, err := c.ExchangeContext(ctx, &m, clientCtx.Server)
	if err != nil {
		return fmt.Errorf("Error in DNS query to %s. %s", clientCtx.Server, err)
	}
	clientCtx.Response = r
	clientCtx.RTT = rtt
	return nil
}

// ValidateResponseWithCodeClientStep validates the code of the DNS response.
func ValidateResponseWithCodeClientStep(ctx context.Context, expectedCode string) error {
	clientCtx := GetClientContext(ctx)
	responseCode, ok := dns.RcodeToString[clientCtx.Response.Rcode]
	if !ok {
		return fmt.Errorf("Invalid code: %d", clientCtx.Response.Rcode)
	}
	if responseCode != expectedCode {
		return fmt.Errorf("Expected DNS code: %s but received: %s", expectedCode, responseCode)
	}
	return nil
}

// ValidateResponseWithOneOfCodesClientStep validates the code of the DNS response againt a list of valid codes.
func ValidateResponseWithOneOfCodesClientStep(ctx context.Context, expectedCodes []string) error {
	clientCtx := GetClientContext(ctx)
	responseCode, ok := dns.RcodeToString[clientCtx.Response.Rcode]
	if !ok {
		return fmt.Errorf("Invalid code: %d", clientCtx.Response.Rcode)
	}
	for _, code := range expectedCodes {
		if responseCode == code {
			return nil
		}
	}
	return fmt.Errorf("Expected DNS code one of: %s but received: %s", expectedCodes, responseCode)
}

// ValidateResponseWithNumberOfRecordsClientStep validates the amount of records in a DNS response for one of the
// record types: answer, authority, additional.
func ValidateResponseWithNumberOfRecordsClientStep(ctx context.Context, expectedRecords int, recordType RecordType) error {
	clientCtx := GetClientContext(ctx)
	records := len(getRecordsForType(clientCtx.Response, recordType))
	if records != expectedRecords {
		return fmt.Errorf("Expected %d records of type: %s but received %d", expectedRecords, recordType, records)
	}
	return nil
}

// ValidateResponseWithRecordsClientStep validates that the response contains the following records for one of the record
// types: answer, authority, additional.
// The expected information of the record is available in the godog table.
func ValidateResponseWithRecordsClientStep(ctx context.Context, recordType RecordType, t *godog.Table) error {
	clientCtx := GetClientContext(ctx)
	records := getRecordsForType(clientCtx.Response, recordType)
	expectedRecords := []Record{}
	if err := golium.ConvertTableWithHeaderToStructSlice(ctx, t, &expectedRecords); err != nil {
		return err
	}
	for _, expectedRecord := range expectedRecords {
		if !expectedRecord.IsContained(records) {
			return fmt.Errorf("No %s record with: %+v", recordType, expectedRecord)
		}
	}
	return nil
}

// func ValidateResponseWithRecordsClientStep(ctx context.Context, recordType RecordType, expectedRecords *godog.Table) error {
// 	clientCtx := GetClientContext(ctx)
// 	records := getRecordsForType(clientCtx.Response, recordType)
// 	header := expectedRecords.Rows[0].Cells
// 	recordInfo := make(map[string](string))
// 	for i := 1; i < len(expectedRecords.Rows); i++ {
// 		for n, cell := range expectedRecords.Rows[i].Cells {
// 			recordInfo[header[n].Value] = cell.Value
// 		}
// 		if err := validateResponseWithRecord(recordType, records, recordInfo); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }
