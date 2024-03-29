Feature: DNS client

  @dns
  Scenario: Query domain inspecting answer records
    Given the DNS server "[CONF:dns]"
      And a DNS timeout of "[CONF:timeout]" milliseconds
     When I send a DNS query of type "A" for "www.telefonica.net"
     Then the DNS response must have the code "NOERROR"
      And the DNS response must contain the following answer records
          | name                | class | type  | data          |
          | www.telefonica.net. | IN    | A     | 213.4.130.95  |

  @dns
  Scenario: Query domain with eDNS0 options
    Given the DNS server "[CONF:dns]"
      And a DNS timeout of "[CONF:timeout]" milliseconds
      And the DNS query options
          | code | data               |
          | 1001 | [SHA256:test-1001] |
          | 1002 | [SHA256:test-1002] |
     When I send a DNS query of type "A" for "www.telefonica.net"
     Then the DNS response must have the code "NOERROR"
      And the DNS response must contain the following answer records
          | name                | class | type  | data          |
          | www.telefonica.net. | IN    | A     | 213.4.130.95  |

  @dns
  Scenario Outline: Query domain with recursion
    Given the DNS server "[CONF:dns]"
      And a DNS timeout of "[CONF:timeout]" milliseconds
     When I send a DNS query of type "<type>" for "<domain>"
     Then the DNS response must have the code "<code>"
      And the DNS response must have "<answer>" answer record
      And the DNS response must have "<authority>" authority records
      And the DNS response must have "<additional>" additional records

    Examples: domain <domain> with type <type>
          | domain               | type | code     | answer | authority | additional |
          | www.telefonica.net.  | A    | NOERROR  | 1      | 0         | 0          |
          | www.telefonica.net.  | AAAA | NOERROR  | 0      | 1         | 0          |
          | www.telefonica.net.  | MX   | NOERROR  | 0      | 1         | 0          |
          | www.telefonica.net.  | NS   | NOERROR  | 0      | 1         | 0          |
          | w.invalid.dsfsd.     | A    | NXDOMAIN | 0      | 1         | 0          |
          | w.invalid.dsfsd.     | AAAA | NXDOMAIN | 0      | 1         | 0          |
          | w.invalid.dsfsd.     | MX   | NXDOMAIN | 0      | 1         | 0          |
          | w.invalid.dsfsd.     | NS   | NXDOMAIN | 0      | 1         | 0          |

  @dns
  Scenario Outline: Query domain without recursion
    Given the DNS server "[CONF:dns]"
    When I send a DNS query of type "<type>" for "<domain>" without recursion
    Then the DNS response must have one of the following codes: "<code>"

    Examples: 
         | domain               | type | code              |
         | www.telefonica.net.  | A    | NOERROR,SERVFAIL  |
         | www.telefonica.net.  | AAAA | NOERROR,SERVFAIL  |
         | www.telefonica.net.  | MX   | NOERROR,SERVFAIL  |
         | www.telefonica.net.  | NS   | NOERROR,SERVFAIL  |
         | w.invalid.dsfsd.     | A    | NXDOMAIN,SERVFAIL |
         | w.invalid.dsfsd.     | AAAA | NXDOMAIN,SERVFAIL |
         | w.invalid.dsfsd.     | MX   | NXDOMAIN,SERVFAIL |
         | w.invalid.dsfsd.     | NS   | NXDOMAIN,SERVFAIL |
