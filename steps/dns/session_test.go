package dns

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"testing"

	"github.com/miekg/dns"
	"github.com/stretchr/testify/require"
)

const (
	fakeServer  string     = "fakeServer"
	udpServer   string     = "8.8.8.8:53"
	dohServer   string     = "https://dns.google:443/dns-query"
	dotServer   string     = "tls://dns.google:853"
	answerRType RecordType = Answer
)

func TestSendUDPQuery(t *testing.T) {
	tcs := []struct {
		name        string
		server      string
		options     bool
		qtype       uint16
		qdomain     string
		recursive   bool
		expectedErr error
	}{
		{
			name:        "Send UDP query ok",
			server:      udpServer,
			options:     false,
			qtype:       1,
			qdomain:     "any",
			recursive:   false,
			expectedErr: nil,
		},
		{
			name:        "Send UDP query with options",
			server:      udpServer,
			options:     true,
			qtype:       1,
			qdomain:     "any",
			recursive:   false,
			expectedErr: nil,
		},
		{
			name:      "Send UDP query fake server",
			server:    fakeServer,
			options:   false,
			qtype:     1,
			qdomain:   "any",
			recursive: false,
			expectedErr: fmt.Errorf("failed DNS query to '%s': %w", fakeServer,
				&net.OpError{Op: "dial", Net: "udp",
					Err: &net.AddrError{Err: "missing port in address", Addr: fakeServer},
				},
			),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ctx, s := getContextAndSession()
			createLogsDir()
			s.ConfigureServer(ctx, tc.server, "UDP")

			if tc.options {
				options := []dns.EDNS0{&dns.EDNS0_LOCAL{Code: 1, Data: []byte{1}}}
				s.ConfigureOptions(ctx, options)
			}

			err := s.SendUDPQuery(ctx, tc.qtype, tc.qdomain, tc.recursive)

			require.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestSendDoHQuery(t *testing.T) {
	tcs := []struct {
		name        string
		server      string
		transport   string
		method      string
		qtype       uint16
		qdomain     string
		recursive   bool
		expectedErr error
	}{
		{
			name:        "Send DoH query with GET ok",
			server:      dohServer,
			transport:   "Doh with GET",
			method:      "GET",
			qtype:       1,
			qdomain:     "www.telefonica.net",
			recursive:   true,
			expectedErr: nil,
		},
		{
			name:        "Send DoH query with POST ok",
			server:      dohServer,
			transport:   "Doh with POST",
			method:      "POST",
			qtype:       1,
			qdomain:     "www.telefonica.net",
			recursive:   true,
			expectedErr: nil,
		},
		{
			name:        "Send DoH query with unsupported method",
			server:      dohServer,
			transport:   "Doh with POST",
			method:      "UM",
			qtype:       1,
			qdomain:     "www.telefonica.net",
			recursive:   true,
			expectedErr: fmt.Errorf("unsupported method. %s", "UM"),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ctx, s := getContextAndSession()
			createLogsDir()
			s.ConfigureServer(ctx, tc.server, tc.transport)

			err := s.SendDoHQuery(ctx, tc.method, tc.qtype, tc.qdomain, tc.recursive)
			require.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestSendDoTQuery(t *testing.T) {
	tcs := []struct {
		name                   string
		server                 string
		qtype                  uint16
		qdomain                string
		recursive              bool
		expectedErr            error
		alternativeExpectedErr error
	}{
		{
			name:        "Send DoT query ok",
			server:      dotServer,
			qtype:       1,
			qdomain:     "www.telefonica.net.",
			recursive:   true,
			expectedErr: nil,
		},
		{
			name:      "Send DoT query fake server",
			server:    fakeServer,
			qtype:     1,
			qdomain:   "www.telefonica.net.",
			recursive: true,
			expectedErr: fmt.Errorf("cannot make the DNS request: %w",
				&net.OpError{Op: "dial", Net: "udp",
					Err: &net.DNSError{
						Err:         "Temporary failure in name resolution",
						Name:        fakeServer,
						IsTemporary: true,
						IsNotFound:  false,
					},
				}),
			alternativeExpectedErr: fmt.Errorf("cannot make the DNS request: %w",
				fmt.Errorf("dialing %q: %w", fakeServer+":53",
					fmt.Errorf("resolving hostname: %w", &net.DNSError{
						Err:         "server misbehaving",
						Name:        fakeServer,
						Server:      "127.0.0.53:53",
						IsTimeout:   false,
						IsTemporary: false,
						IsNotFound:  true,
					},
					),
				),
			),
		},
		{
			name:      "Send DoT query upstream failure",
			server:    "http://asdp-ert",
			qtype:     1,
			qdomain:   "www.telefonica.net.",
			recursive: true,
			expectedErr: fmt.Errorf("cannot create an upstream: %w",
				errors.New("unsupported url scheme: http")),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ctx, s := getContextAndSession()
			createLogsDir()
			s.ConfigureServer(ctx, tc.server, "DoT")

			err := s.SendDoTQuery(ctx, tc.qtype, tc.qdomain, tc.recursive)
			if tc.expectedErr == nil && err == nil {
				return
			}
			if tc.expectedErr.Error() == err.Error() {
				require.Equal(t, tc.expectedErr, err)
			} else if tc.alternativeExpectedErr != nil {
				require.Equal(t, tc.alternativeExpectedErr, err)
			}
		})
	}
}

func TestValidateResponseWithCode(t *testing.T) {
	tcs := []struct {
		name        string
		respCode    int
		expCodes    string
		expectedErr error
	}{
		{
			name:        "Validate Response with valid code",
			respCode:    1,
			expCodes:    "FORMERR",
			expectedErr: nil,
		},
		{
			name:        "Validate Response with invalid code",
			respCode:    24,
			expCodes:    "FORMERR",
			expectedErr: fmt.Errorf("invalid code '%d'", 24),
		},
		{
			name:     "Validate Response with different code",
			respCode: 2,
			expCodes: "NOERROR",
			expectedErr: fmt.Errorf(
				"expected DNS code '%s' but received '%s'", "NOERROR", "SERVFAIL"),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ctx, s := getContextAndSession()
			s.Response = &dns.Msg{MsgHdr: dns.MsgHdr{Rcode: tc.respCode}}

			err := s.ValidateResponseWithCode(ctx, tc.expCodes)

			require.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestValidateResponseWithOneOfCodes(t *testing.T) {
	tcs := []struct {
		name        string
		respCode    int
		expCodes    []string
		expectedErr error
	}{
		{
			name:        "Validate Response with valid codes",
			respCode:    1,
			expCodes:    []string{"FORMERR", "FORMERR"},
			expectedErr: nil,
		},
		{
			name:        "Validate Response with invalid codes",
			respCode:    24,
			expCodes:    []string{"FORMERR", "FORMERR"},
			expectedErr: fmt.Errorf("invalid code '%d'", 24),
		},
		{
			name:     "Validate Response with different codes",
			respCode: 2,
			expCodes: []string{"NOERROR", "FORMERR"},
			expectedErr: fmt.Errorf(
				"expected DNS code one of '%s' but received '%s'", "[NOERROR FORMERR]", "SERVFAIL"),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ctx, s := getContextAndSession()
			s.Response = &dns.Msg{MsgHdr: dns.MsgHdr{Rcode: tc.respCode}}

			err := s.ValidateResponseWithOneOfCodes(ctx, tc.expCodes)

			require.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestValidateResponseWithNumberOfRecords(t *testing.T) {
	tcs := []struct {
		name        string
		expRecords  int
		recordType  RecordType
		expectedErr error
	}{
		{
			name:        "Validate Response with valid number of records",
			expRecords:  1,
			recordType:  answerRType,
			expectedErr: nil,
		},
		{
			name:       "Validate Response with invalid number of records",
			expRecords: 0,
			recordType: answerRType,
			expectedErr: fmt.Errorf("expected '%d' records of type '%s' but received '%d'",
				0, answerRType, 1),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ctx, s := getContextAndSession()
			s.Response = &dns.Msg{Answer: []dns.RR{
				&dns.RR_Header{Name: "John", Class: dns.ClassINET, Rrtype: dns.TypeA},
			}}
			err := s.ValidateResponseWithNumberOfRecords(ctx, tc.expRecords, tc.recordType)

			require.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestValidateResponseWithRecords(t *testing.T) {
	rNname := "John"
	rType := dns.Type(dns.TypeA).String()
	rClass := dns.Class(dns.ClassINET).String()
	tcs := []struct {
		name        string
		recordType  RecordType
		expRecords  []Record
		expectedErr error
	}{
		{
			name:        "Validate Response with valid records",
			recordType:  answerRType,
			expRecords:  []Record{{Name: &rNname, Type: &rType, Class: &rClass}},
			expectedErr: nil,
		},
		{
			name:       "Validate Response with invalid records",
			recordType: "fake",
			expRecords: []Record{{Name: &rNname, Type: &rType, Class: &rClass}},
			expectedErr: fmt.Errorf("no '%s' record with '%+v' in",
				"fake", "{\"Name\":\"John\",\"Type\":\"A\",\"Class\":\"IN\",\"Data\":null,\"TTL\":null}"),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ctx, s := getContextAndSession()
			s.Response = &dns.Msg{Answer: []dns.RR{
				&dns.RR_Header{Name: "John", Class: dns.ClassINET, Rrtype: dns.TypeA},
			}}
			err := s.ValidateResponseWithRecords(ctx, tc.recordType, tc.expRecords)

			require.Equal(t, tc.expectedErr, err)
		})
	}
}

func createLogsDir() {
	newpath := filepath.Join(".", "logs")
	os.MkdirAll(newpath, os.ModePerm)
}

func getContextAndSession() (context.Context, *Session) {
	ctx := InitializeContext(context.Background())
	s := GetSession(ctx)
	return ctx, s
}
