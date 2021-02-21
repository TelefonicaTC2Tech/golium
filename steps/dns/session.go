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
	"time"

	"github.com/miekg/dns"
)

// Session contains the information related to a DNS query and response.
type Session struct {
	// Server is the address to the DNS server, including the server port (e.g. 8.8.8.8:53).
	Server string
	// Query contains the DNS request message.
	Query *dns.Msg
	// Response contains the DNS response message.
	Response *dns.Msg
	// RTT is the response time.
	RTT time.Duration
}

// ConfigureServer configures the DNS server location.
func (s *Session) ConfigureServer(ctx context.Context, svr string) error {
	s.Server = svr
	return nil
}

// SendQuery sends a DNS query to resolve a domain.
func (s *Session) SendQuery(ctx context.Context, qtype uint16, qdomain string, recursive bool) error {
	c := dns.Client{}
	m := dns.Msg{}
	m.SetQuestion(dns.Fqdn(qdomain), qtype)
	m.RecursionDesired = recursive
	s.Query = &m
	r, rtt, err := c.ExchangeContext(ctx, &m, s.Server)
	if err != nil {
		return fmt.Errorf("Error in DNS query to %s. %s", s.Server, err)
	}
	s.Response = r
	s.RTT = rtt
	return nil
}

// ValidateResponseWithCode validates the code of the DNS response.
func (s *Session) ValidateResponseWithCode(ctx context.Context, expectedCode string) error {
	responseCode, ok := dns.RcodeToString[s.Response.Rcode]
	if !ok {
		return fmt.Errorf("Invalid code: %d", s.Response.Rcode)
	}
	if responseCode != expectedCode {
		return fmt.Errorf("Expected DNS code: %s but received: %s", expectedCode, responseCode)
	}
	return nil
}

// ValidateResponseWithOneOfCodes validates the code of the DNS response againt a list of valid codes.
func (s *Session) ValidateResponseWithOneOfCodes(ctx context.Context, expectedCodes []string) error {
	responseCode, ok := dns.RcodeToString[s.Response.Rcode]
	if !ok {
		return fmt.Errorf("Invalid code: %d", s.Response.Rcode)
	}
	for _, code := range expectedCodes {
		if responseCode == code {
			return nil
		}
	}
	return fmt.Errorf("Expected DNS code one of: %s but received: %s", expectedCodes, responseCode)
}

// ValidateResponseWithNumberOfRecords validates the amount of records in a DNS response for one of the
// record types: answer, authority, additional.
func (s *Session) ValidateResponseWithNumberOfRecords(ctx context.Context, expectedRecords int, recordType RecordType) error {
	records := len(getRecordsForType(s.Response, recordType))
	if records != expectedRecords {
		return fmt.Errorf("Expected %d records of type: %s but received %d", expectedRecords, recordType, records)
	}
	return nil
}

// ValidateResponseWithRecords validates that the response contains the following records for one of the record
// types: answer, authority, additional.
func (s *Session) ValidateResponseWithRecords(ctx context.Context, recordType RecordType, expectedRecords []Record) error {
	records := getRecordsForType(s.Response, recordType)
	for _, expectedRecord := range expectedRecords {
		if !expectedRecord.IsContained(records) {
			return fmt.Errorf("No %s record with: %+v", recordType, expectedRecord)
		}
	}
	return nil
}
