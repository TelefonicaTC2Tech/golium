Feature: DOH client

  Background:
    Given the DNS server "[CONF:doh]"
    And a DNS timeout of "[CONF:timeout]" milliseconds

  @doh
  Scenario Outline: DoH Query domain inspecting answer records
     When I send a DoH "<method>" request of type "A" for "www.telefonica.com"
     Then the DNS response must have the code "NOERROR"
      And the DNS response must contain the following answer records
          | name                | class | type  | data          |
          | www.telefonica.com. | IN    | A     | 212.170.36.79 |

    Examples: HTTP method: <method>
          | method |
          | GET    |
          | POST   |

  @doh
  Scenario Outline: DoH Query domain with recursion
     When I send a DoH "<method>" request of type "<type>" for "<domain>"
     Then the DNS response must have the code "<code>"
      And the DNS response must have "<answer>" answer record
      And the DNS response must have "<authority>" authority records
      And the DNS response must have "<additional>" additional records

    Examples: method <method> domain <domain> with type <type>
          | method | domain               | type | code     | answer | authority | additional |
          | GET    | www.telefonica.com.  | A    | NOERROR  | 1      | 0         | 0          |
          | GET    | www.telefonica.com.  | AAAA | NOERROR  | 1      | 0         | 0          |
          | GET    | www.telefonica.com.  | MX   | NOERROR  | 0      | 1         | 0          |
          | GET    | www.telefonica.com.  | NS   | NOERROR  | 0      | 1         | 0          |
          | GET    | w.invalid.dsfsd.     | A    | NXDOMAIN | 0      | 1         | 0          |
          | GET    | w.invalid.dsfsd.     | AAAA | NXDOMAIN | 0      | 1         | 0          |
          | GET    | w.invalid.dsfsd.     | MX   | NXDOMAIN | 0      | 1         | 0          |
          | GET    | w.invalid.dsfsd.     | NS   | NXDOMAIN | 0      | 1         | 0          |
          | POST   | www.telefonica.com.  | A    | NOERROR  | 1      | 0         | 0          |
          | POST   | www.telefonica.com.  | AAAA | NOERROR  | 1      | 0         | 0          |
          | POST   | www.telefonica.com.  | MX   | NOERROR  | 0      | 1         | 0          |
          | POST   | www.telefonica.com.  | NS   | NOERROR  | 0      | 1         | 0          |
          | POST   | w.invalid.dsfsd.     | A    | NXDOMAIN | 0      | 1         | 0          |
          | POST   | w.invalid.dsfsd.     | AAAA | NXDOMAIN | 0      | 1         | 0          |
          | POST   | w.invalid.dsfsd.     | MX   | NXDOMAIN | 0      | 1         | 0          |
          | POST   | w.invalid.dsfsd.     | NS   | NXDOMAIN | 0      | 1         | 0          |

  @doh
  Scenario Outline: DoH Query domain without recursion
    When I send a DoH "<method>" request of type "<type>" for "<domain>" without recursion
    Then the DNS response must have one of the following codes: "<code>"

    Examples: method <method> domain <domain> with type <type>
         | method | domain               | type | code              |
         | GET    | www.telefonica.com.  | A    | NOERROR,SERVFAIL  |
         | GET    | www.telefonica.com.  | AAAA | NOERROR,SERVFAIL  |
         | GET    | www.telefonica.com.  | MX   | NOERROR,SERVFAIL  |
         | GET    | www.telefonica.com.  | NS   | NOERROR,SERVFAIL  |
         | GET    | w.invalid.dsfsd.     | A    | NXDOMAIN,SERVFAIL |
         | GET    | w.invalid.dsfsd.     | AAAA | NXDOMAIN,SERVFAIL |
         | GET    | w.invalid.dsfsd.     | MX   | NXDOMAIN,SERVFAIL |
         | GET    | w.invalid.dsfsd.     | NS   | NXDOMAIN,SERVFAIL |
         | POST   | www.telefonica.com.  | A    | NOERROR,SERVFAIL  |
         | POST   | www.telefonica.com.  | AAAA | NOERROR,SERVFAIL  |
         | POST   | www.telefonica.com.  | MX   | NOERROR,SERVFAIL  |
         | POST   | www.telefonica.com.  | NS   | NOERROR,SERVFAIL  |
         | POST   | w.invalid.dsfsd.     | A    | NXDOMAIN,SERVFAIL |
         | POST   | w.invalid.dsfsd.     | AAAA | NXDOMAIN,SERVFAIL |
         | POST   | w.invalid.dsfsd.     | MX   | NXDOMAIN,SERVFAIL |
         | POST   | w.invalid.dsfsd.     | NS   | NXDOMAIN,SERVFAIL |
