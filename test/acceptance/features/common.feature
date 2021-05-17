Feature: Common

  @common
  Scenario: Store a UUID in the context
    Given I store "[UUID]" in context "test.uuid"
      And I store "[SHA256:test.value]" in context "test.value"
    Given the HTTP endpoint "[CONF:url]/anything/[CTXT:test.uuid]/[CTXT:test.value]"
     When I send a HTTP "GET" request
     Then the HTTP status code must be "200"
      And the HTTP response body must have the JSON properties
          | method | GET                                                      |
          | url    | [CONF:url]/anything/[CTXT:test.uuid]/[SHA256:test.value] |
  
  @common
  Scenario: Base64 encoding
    Given I store "username" in context "test.base64.username"
      And I store "password" in context "test.base64.password"
     When I store "[BASE64:[CTXT:test.base64.username]:[CTXT:test.base64.password]]" in context "test.base64.auth"
     Then the value "[CTXT:test.base64.auth]" must be equal to "dXNlcm5hbWU6cGFzc3dvcmQ="

  @common
  Scenario: Wait
    Given I wait for "2" millis
    Given I wait for "1" seconds

  @common
  Scenario: Parse URL
    Given I parse the URL "https://www.elevenpaths.com:443/products-services/solutions?a=1&b=test" in context "url"
     Then the value "[CTXT:url.scheme]" must be equal to "https"
      And the value "[CTXT:url.host]" must be equal to "www.elevenpaths.com:443"
      And the value "[CTXT:url.hostname]" must be equal to "www.elevenpaths.com"
      And the value "[CTXT:url.path]" must be equal to "/products-services/solutions"
      And the value "[CTXT:url.rawquery]" must be equal to "a=1&b=test"
      And the value "[CTXT:url.query.a]" must be equal to "1"
      And the value "[CTXT:url.query.b]" must be equal to "test"

  @common
  Scenario: Store my local ip in context
    Given I store my local ip in context "context.ip"
