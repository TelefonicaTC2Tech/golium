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
	"encoding/json"
	"fmt"

	"github.com/miekg/dns"
)

// RecordType
// enumerates the possible types of DNS records: answer, authority and additional records.
type RecordType string

const (
	// Answer record type for DNS response messages
	Answer RecordType = "answer"
	// Authority record type for DNS response messages
	Authority = "authority"
	// Additional record type for DNS response messages
	Additional = "additional"
)

// QueryTypes is a map of the DNS query type with the correspondence type in dns package.
var QueryTypes = map[string]uint16{
	"A":     dns.TypeA,
	"AAAA":  dns.TypeAAAA,
	"CNAME": dns.TypeCNAME,
	"MX":    dns.TypeMX,
	"NS":    dns.TypeNS,
}

// Record is an abstraction of a DNS records (based on dns.RR).
// It aims to checks the most relevant fields of a DNS records.
type Record struct {
	Name  *string
	Type  *string
	Class *string
	Data  *string
	TTL   *uint
}

func (r Record) String() string {
	b, err := json.Marshal(r)
	if err != nil {
		return "Error marshalling the record"
	}
	return string(b)
}

// Matches checks if the record contains the same information in the fields of the
// actual DNS record (of type dns.RR).
func (r *Record) Matches(rr dns.RR) bool {
	h := rr.Header()
	if r.Name != nil && *r.Name != h.Name {
		return false
	}
	if r.Class != nil && *r.Class != dns.Class(h.Class).String() {
		return false
	}
	if r.Type != nil && *r.Type != dns.Type(h.Rrtype).String() {
		return false
	}
	if r.Data != nil {
		if data, _ := getValueFromRecord(rr); data != *r.Data {
			return false
		}
	}
	return true
}

// IsContained checks if the record matches with any of the actual DNS records of the slice.
func (r *Record) IsContained(rrs []dns.RR) bool {
	for _, rr := range rrs {
		if r.Matches(rr) {
			return true
		}
	}
	return false
}

func getValueFromRecord(rr dns.RR) (string, error) {
	var value string
	switch rr.Header().Rrtype {
	case dns.TypeA:
		a := rr.(*dns.A).A
		value = a.String()
	case dns.TypeAAAA:
		aaaa := rr.(*dns.AAAA).AAAA
		value = aaaa.String()
	case dns.TypeCNAME:
		value = rr.(*dns.CNAME).Target
	case dns.TypeMX:
		value = rr.(*dns.MX).Mx
	case dns.TypeNS:
		value = rr.(*dns.NS).Ns
	default:
		return "", fmt.Errorf(
			"unsupported record type '%s'", dns.Type(rr.Header().Rrtype).String())
	}
	return value, nil
}

func getRecordsForType(msg *dns.Msg, recordType RecordType) []dns.RR {
	switch recordType {
	case Answer:
		return msg.Answer
	case Authority:
		return msg.Ns
	case Additional:
		return msg.Extra
	default:
		return nil
	}
}
