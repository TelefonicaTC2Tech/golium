Feature: Common

  Scenario: Send a GET request
    Given I generate a UUID and store it in context "test.uuid"
      And I store "[SHA256:test.value]" in context "test.value"
    Given the HTTP endpoint "[CONF:url]/anything/[CTXT:test.uuid]/[CTXT:test.value]"
     When I send a HTTP "GET" request
     Then the HTTP status code must be "200"
      And the HTTP response body must have the JSON properties
          | method | GET                                                      |
          | url    | [CONF:url]/anything/[CTXT:test.uuid]/[SHA256:test.value] |

