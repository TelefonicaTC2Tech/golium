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
	"bytes"
	"context"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/AdguardTeam/dnsproxy/upstream"
	"github.com/google/uuid"
	"github.com/miekg/dns"
)

// Session contains the information related to a DNS query and response.
type Session struct {
	// Server is the address to the DNS server, including the server port (e.g. 8.8.8.8:53).
	Server string
	// Transport is the network protocol used to send the queries (valid values: UDP, DoT, Doh with GET, DoH with POST)
	Transport string
	// DNS query options (EDNS0)
	Options []dns.EDNS0
	// Query parameters
	DoHQueryParams map[string][]string
	// Query contains the DNS request message.
	Query *dns.Msg
	// Response contains the DNS response message.
	Response *dns.Msg
	// RTT is the response time.
	RTT time.Duration
	// Timeout is the maximum time to wait for a response. Expressed in Milliseconds
	Timeout time.Duration
}

// ConfigureServer configures the DNS server location and the transport protocol.
func (s *Session) ConfigureServer(ctx context.Context, svr string, transport string) {
	s.Server = svr
	s.Transport = transport
}

// SetDNSResponseTimeout configures a DNS response timeout.
func (s *Session) SetDNSResponseTimeout(ctx context.Context, timeout int) {
	s.Timeout = time.Duration(timeout) * time.Millisecond
}

// ConfigureOptions adds EDNS0 options to be included in the DNS query.
func (s *Session) ConfigureOptions(ctx context.Context, options []dns.EDNS0) {
	s.Options = append(s.Options, options...)
}

// ConfigureDoHQueryParams stores a table of query parameters in the application context.
func (s *Session) ConfigureDoHQueryParams(ctx context.Context, params map[string][]string) {
	s.DoHQueryParams = params
}

// SendUDPQuery sends a DNS query to resolve a domain.
func (s *Session) SendUDPQuery(ctx context.Context, qtype uint16, qdomain string, recursive bool) error {
	logger := GetLogger()
	corr := uuid.New().String()
	c := dns.Client{Timeout: s.Timeout}
	m := &dns.Msg{}
	m.SetQuestion(dns.Fqdn(qdomain), qtype)
	m.RecursionDesired = recursive
	// Add EDNS0 options (if registered in s.Options)
	if len(s.Options) > 0 {
		opt := &dns.OPT{
			Hdr: dns.RR_Header{
				Name:   ".",
				Rrtype: dns.TypeOPT,
			},
			Option: s.Options,
		}
		m.Extra = append(m.Extra, opt)
	}
	s.Query = m
	logger.LogRequest(m, corr)
	r, rtt, err := c.ExchangeContext(ctx, m, s.Server)
	if err != nil {
		return fmt.Errorf("failed DNS query to '%s': %w", s.Server, err)
	}
	logger.LogResponse(m, corr)
	s.Response = r
	s.RTT = rtt
	return nil
}

// SendDoHQuery sends a DoH query to resolve a domain.
func (s *Session) SendDoHQuery(ctx context.Context, method string, qtype uint16, qdomain string, recursive bool) error {
	logger := GetLogger()
	corr := uuid.New().String()
	//Set DNS query
	m := &dns.Msg{}
	m.SetQuestion(dns.Fqdn(qdomain), qtype)
	m.RecursionDesired = recursive
	s.Query = m
	logger.LogRequest(m, corr)
	// Pack the DNS query to convert to a DNS wireformat
	data, err := s.Query.Pack()
	if err != nil {
		return err
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Timeout:   s.Timeout,
		Transport: tr,
	}
	var request *http.Request

	switch method {
	case "GET":
		dq := base64.RawURLEncoding.EncodeToString(data)
		request, err = http.NewRequest("GET", fmt.Sprintf("%s?dns=%s", s.Server, dq), nil)
		if err != nil {
			return err
		}
	case "POST":
		u, err := url.Parse(s.Server)
		if err != nil {
			return err
		}
		params := url.Values(s.DoHQueryParams)
		u.RawQuery = params.Encode()
		request, err = http.NewRequest("POST", u.String(), bytes.NewReader(data))
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported method. %s", method)
	}

	request.Header.Set("Content-Type", "application/dns-message")
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("error sending request. %s", err)
	}
	// Check Content-Type
	if response.Header.Get("Content-Type") != "application/dns-message" {
		return fmt.Errorf("error in Content-Type Header. Value: %s, Expected: %s", response.Header.Get("Content-Type"), "application/dns-message")
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("error reading the response body. %s", err)
	}
	// Get the response body and Unpack it to convert from a DNS wireformat
	dnsResp := new(dns.Msg)
	err = dnsResp.Unpack(body)
	if err != nil {
		return fmt.Errorf("error unpacking body. %s", err)
	}
	logger.LogResponse(dnsResp, corr)
	// Set response in dns session struct
	s.Response = dnsResp
	return nil
}

// SendDoTQuery sends a DoT query to resolve a domain.
func (s *Session) SendDoTQuery(ctx context.Context, qtype uint16, qdomain string, recursive bool) error {
	logger := GetLogger()
	corr := uuid.New().String()
	//Set DNS query
	opts := upstream.Options{
		Timeout:            s.Timeout,
		InsecureSkipVerify: true,
	}
	u, err := upstream.AddressToUpstream(s.Server, &opts)
	if err != nil {
		logger.log.Fatalf("Cannot create an upstream: %s", err)
	}
	m := &dns.Msg{}
	m.Id = dns.Id()
	m.SetQuestion(dns.Fqdn(qdomain), qtype)
	m.RecursionDesired = recursive
	s.Query = m
	logger.LogRequest(m, corr)
	dnsResp, err := u.Exchange(m)
	if err != nil {
		logger.log.Fatalf("Cannot make the DNS request: %s", err)
	}
	logger.LogResponse(m, corr)
	// Set response in dns session struct
	s.Response = dnsResp
	return nil
}

// ValidateResponseWithCode validates the code of the DNS response.
func (s *Session) ValidateResponseWithCode(ctx context.Context, expectedCode string) error {
	responseCode, ok := dns.RcodeToString[s.Response.Rcode]
	if !ok {
		return fmt.Errorf("invalid code '%d'", s.Response.Rcode)
	}
	if responseCode != expectedCode {
		return fmt.Errorf("expected DNS code '%s' but received '%s'", expectedCode, responseCode)
	}
	return nil
}

// ValidateResponseWithOneOfCodes validates the code of the DNS response againt a list of valid codes.
func (s *Session) ValidateResponseWithOneOfCodes(ctx context.Context, expectedCodes []string) error {
	responseCode, ok := dns.RcodeToString[s.Response.Rcode]
	if !ok {
		return fmt.Errorf("invalid code '%d'", s.Response.Rcode)
	}
	for _, code := range expectedCodes {
		if responseCode == code {
			return nil
		}
	}
	return fmt.Errorf("expected DNS code one of '%s' but received '%s'", expectedCodes, responseCode)
}

// ValidateResponseWithNumberOfRecords validates the amount of records in a DNS response for one of the
// record types: answer, authority, additional.
func (s *Session) ValidateResponseWithNumberOfRecords(ctx context.Context, expectedRecords int, recordType RecordType) error {
	records := len(getRecordsForType(s.Response, recordType))
	if records != expectedRecords {
		return fmt.Errorf("expected '%d' records of type '%s' but received '%d'", expectedRecords, recordType, records)
	}
	return nil
}

// ValidateResponseWithRecords validates that the response contains the following records for one of the record
// types: answer, authority, additional.
func (s *Session) ValidateResponseWithRecords(ctx context.Context, recordType RecordType, expectedRecords []Record) error {
	records := getRecordsForType(s.Response, recordType)
	for _, expectedRecord := range expectedRecords {
		if !expectedRecord.IsContained(records) {
			return fmt.Errorf("no '%s' record with '%+v'", recordType, expectedRecord)
		}
	}
	return nil
}
