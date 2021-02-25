Feature: DNS client

  @wip
  Scenario: Query domain inspecting answer records
    Given the DNS server "[CONF:dns]"
      And the DNS query options
          | code | data               |
          | 1001 | [SHA256:test-6001] |
          | 1002 | [SHA256:test-1002] |
     When I send a DNS query of type "A" for "www.elevenpaths.com"
     Then the DNS response must have the code "NOERROR"
      And the DNS response must contain the following answer records
          | name                                                   | class | type  | data                                                   |
          | www.elevenpaths.com.                                   | IN    | CNAME | legacy-web-pro-1067844838.eu-west-1.elb.amazonaws.com. |
          | legacy-web-pro-1067844838.eu-west-1.elb.amazonaws.com. | IN    | A     | 52.16.121.202                                          |
          | legacy-web-pro-1067844838.eu-west-1.elb.amazonaws.com. | IN    | A     | 52.30.171.81                                           |
          | legacy-web-pro-1067844838.eu-west-1.elb.amazonaws.com. | IN    | A     | 54.194.165.93                                          |

  Scenario: Query domain with eDNS0 options
    Given the DNS server "[CONF:dns]"
     When I send a DNS query of type "A" for "www.elevenpaths.com"
     Then the DNS response must have the code "NOERROR"
      And the DNS response must contain the following answer records
          | name                                                   | class | type  | data                                                   |
          | www.elevenpaths.com.                                   | IN    | CNAME | legacy-web-pro-1067844838.eu-west-1.elb.amazonaws.com. |
          | legacy-web-pro-1067844838.eu-west-1.elb.amazonaws.com. | IN    | A     | 52.16.121.202                                          |
          | legacy-web-pro-1067844838.eu-west-1.elb.amazonaws.com. | IN    | A     | 52.30.171.81                                           |
          | legacy-web-pro-1067844838.eu-west-1.elb.amazonaws.com. | IN    | A     | 54.194.165.93                                          |

  Scenario Outline: Query domain with recursion
    Given the DNS server "[CONF:dns]"
     When I send a DNS query of type "<type>" for "<domain>"
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

  Scenario Outline: Query domain without recursion
    Given the DNS server "[CONF:dns]"
    When I send a DNS query of type "<type>" for "<domain>" without recursion
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
