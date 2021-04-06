Feature: DOT client

  Background:
    Given the DNS server "[CONF:dot]"
    And a DNS timeout of "[CONF:timeout]" milliseconds

  @dot
  Scenario: DoT Query domain inspecting answer records
     When I send a DoT query of type "A" for "www.telefonica.com"
     Then the DNS response must have the code "NOERROR"
      And the DNS response must contain the following answer records
          | name                | class | type  | data          |
          | www.telefonica.com. | IN    | A     | 212.170.36.79 |

  @dot
  Scenario Outline: DoT Query domain with recursion
     When I send a DoT query of type "<type>" for "<domain>"
     Then the DNS response must have the code "<code>"
      And the DNS response must have "<answer>" answer record
      And the DNS response must have "<authority>" authority records
      And the DNS response must have "<additional>" additional records

    Examples: domain <domain> with type <type>
          | domain               | type | code     | answer | authority | additional |
          | www.telefonica.com.  | A    | NOERROR  | 1      | 0         | 0          |
          | www.telefonica.com.  | AAAA | NOERROR  | 1      | 0         | 0          |
          | www.telefonica.com.  | MX   | NOERROR  | 0      | 1         | 0          |
          | www.telefonica.com.  | NS   | NOERROR  | 0      | 1         | 0          |
          | w.invalid.dsfsd.     | A    | NXDOMAIN | 0      | 1         | 0          |
          | w.invalid.dsfsd.     | AAAA | NXDOMAIN | 0      | 1         | 0          |
          | w.invalid.dsfsd.     | MX   | NXDOMAIN | 0      | 1         | 0          |
          | w.invalid.dsfsd.     | NS   | NXDOMAIN | 0      | 1         | 0          |

  @dot
  Scenario Outline: DoT Query domain without recursion
    When I send a DoT query of type "<type>" for "<domain>" without recursion
    Then the DNS response must have one of the following codes: "<code>"

    Examples: 
         | domain               | type | code              |
         | www.telefonica.com.  | A    | NOERROR,SERVFAIL  |
         | www.telefonica.com.  | AAAA | NOERROR,SERVFAIL  |
         | www.telefonica.com.  | MX   | NOERROR,SERVFAIL  |
         | www.telefonica.com.  | NS   | NOERROR,SERVFAIL  |
         | w.invalid.dsfsd.     | A    | NXDOMAIN,SERVFAIL |
         | w.invalid.dsfsd.     | AAAA | NXDOMAIN,SERVFAIL |
         | w.invalid.dsfsd.     | MX   | NXDOMAIN,SERVFAIL |
         | w.invalid.dsfsd.     | NS   | NXDOMAIN,SERVFAIL |
