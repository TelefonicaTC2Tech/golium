package dns

import (
	"fmt"
	"testing"

	"github.com/miekg/dns"
	"github.com/stretchr/testify/require"
)

func createTestRecord() Record {
	rNname := "John"
	rType := dns.Type(dns.TypeMX).String()
	rClass := dns.Class(dns.ClassINET).String()
	record := Record{Name: &rNname, Type: &rType, Class: &rClass}
	return record
}

func TestMatches(t *testing.T) {
	tcs := []struct {
		name           string
		rr             dns.RR
		expectedResult bool
	}{
		{
			name:           "Matches RR",
			rr:             &dns.RR_Header{Name: "John", Class: dns.ClassINET, Rrtype: dns.TypeMX},
			expectedResult: true,
		},
		{
			name:           "Not Matches RR Name",
			rr:             &dns.RR_Header{Name: "Joh", Class: dns.ClassINET, Rrtype: dns.TypeMX},
			expectedResult: false,
		},
		{
			name:           "Not Matches RR Class",
			rr:             &dns.RR_Header{Name: "John", Class: dns.ClassCSNET, Rrtype: dns.TypeMX},
			expectedResult: false,
		},
		{
			name:           "Not Matches RR Type",
			rr:             &dns.RR_Header{Name: "John", Class: dns.ClassINET, Rrtype: dns.TypeANY},
			expectedResult: false,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			record := createTestRecord()

			rrMatch := record.Matches(tc.rr)

			require.Equal(t, tc.expectedResult, rrMatch)
		})
	}
}

func TestIsContained(t *testing.T) {
	tcs := []struct {
		name           string
		rrs            []dns.RR
		expectedResult bool
	}{
		{
			name: "Record is contained in []dns.RR",
			rrs: []dns.RR{
				&dns.RR_Header{Name: "John", Class: dns.ClassINET, Rrtype: dns.TypeMX},
			},
			expectedResult: true,
		},
		{
			name: "Record is not contained in []dns.RR",
			rrs: []dns.RR{
				&dns.RR_Header{Name: "Joh", Class: dns.ClassINET, Rrtype: dns.TypeMX},
			},
			expectedResult: false,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			record := createTestRecord()

			rrMatch := record.IsContained(tc.rrs)

			require.Equal(t, tc.expectedResult, rrMatch)
		})
	}
}

func TestGetValueFromRecord(t *testing.T) {
	tcs := []struct {
		name           string
		rr             dns.RR
		expectedResult string
		expectedErr    error
	}{
		{
			name: "Get Value from dns.TypeA",
			rr: &dns.A{
				Hdr: dns.RR_Header{Name: "John", Class: dns.ClassINET, Rrtype: dns.TypeA},
				A:   []byte("a"),
			},
			expectedResult: "?61",
			expectedErr:    nil,
		},
		{
			name: "Get Value from dns.TypeAAAA",
			rr: &dns.AAAA{
				Hdr:  dns.RR_Header{Name: "John", Class: dns.ClassINET, Rrtype: dns.TypeAAAA},
				AAAA: []byte("aaaa"),
			},
			expectedResult: "97.97.97.97",
			expectedErr:    nil,
		},
		{
			name: "Get Value from dns.TypeCNAME",
			rr: &dns.CNAME{
				Hdr:    dns.RR_Header{Name: "John", Class: dns.ClassINET, Rrtype: dns.TypeCNAME},
				Target: "cdomain-name",
			},
			expectedResult: "cdomain-name",
			expectedErr:    nil,
		},
		{
			name: "Get Value from dns.TypeMX",
			rr: &dns.MX{
				Hdr:        dns.RR_Header{Name: "John", Class: dns.ClassINET, Rrtype: dns.TypeMX},
				Preference: 1,
				Mx:         "cdomain-name",
			},
			expectedResult: "cdomain-name",
			expectedErr:    nil,
		},
		{
			name: "Get Value from dns.TypeNS",
			rr: &dns.NS{
				Hdr: dns.RR_Header{Name: "John", Class: dns.ClassINET, Rrtype: dns.TypeNS},
				Ns:  "cdomain-name",
			},
			expectedResult: "cdomain-name",
			expectedErr:    nil,
		},
		{
			name: "Get Value unsuported type",
			rr: &dns.NSAPPTR{
				Hdr: dns.RR_Header{Name: "John", Class: dns.ClassINET, Rrtype: dns.TypeNSAPPTR},
				Ptr: "domain-name",
			},
			expectedResult: "",
			expectedErr:    fmt.Errorf("unsupported record type '%s'", "NSAP-PTR"),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			value, err := getValueFromRecord(tc.rr)

			require.Equal(t, tc.expectedResult, value)
			require.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestGetRecordsForType(t *testing.T) {
	tcs := []struct {
		name           string
		msg            *dns.Msg
		recordType     RecordType
		expectedResult []dns.RR
	}{
		{
			name:           "Get Records for type Answer",
			msg:            &dns.Msg{Answer: []dns.RR{}},
			recordType:     Answer,
			expectedResult: []dns.RR{},
		},
		{
			name:           "Get Records for type Authority",
			msg:            &dns.Msg{Ns: []dns.RR{}},
			recordType:     Authority,
			expectedResult: []dns.RR{},
		},
		{
			name:           "Get Records for type Additional",
			msg:            &dns.Msg{Extra: []dns.RR{}},
			recordType:     Additional,
			expectedResult: []dns.RR{},
		},
		{
			name:           "Get Records for type any",
			msg:            &dns.Msg{Extra: []dns.RR{}},
			recordType:     "any",
			expectedResult: nil,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			resMsg := getRecordsForType(tc.msg, tc.recordType)

			require.Equal(t, tc.expectedResult, resMsg)
		})
	}
}
