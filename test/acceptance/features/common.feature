Feature: Common

  Scenario: Send a GET request
    Given I generate a UUID and store it in context "test.uuid"
    Given the HTTP endpoint "[CONF:url]/anything/[CTXT:test.uuid]"
     When I send a HTTP "GET" request
     Then the HTTP status code must be "200"
      And the HTTP response body must have the JSON properties
          | method | GET                                  |
          | url    | [CONF:url]/anything/[CTXT:test.uuid] |

