Feature: HTTP Mock server

  @mockhttp
  Scenario: Mock request with simple configuration
    Given I store "[UUID]" in context "id"
      And I mock the HTTP request at "[CONF:httpMockUrl]" for path "/test/[CTXT:id]" with status "200" and JSON body
      """
      {
        "value": "test mock response"
      }
      """
    Given the HTTP endpoint "[CONF:httpMockUrl]/test/[CTXT:id]"
     When I send a HTTP "GET" request
     Then the HTTP status code must be "200"
      And the HTTP request headers
          | Content-Type | application/json |
      And the HTTP response body must have the JSON properties
          | value | test mock response |

  @mockhttp
  Scenario: Mock permanent request
    Given I store "[UUID]" in context "id"
      And I mock the HTTP request at "[CONF:httpMockUrl]" with the JSON
      """
      {
        "permanent": true,
        "request": {
          "method": "DELETE",
          "path": "/test/[CTXT:id]"
        },
        "response": {
          "status": 204,
          "body": ""
        }
      }
      """
    Given the HTTP endpoint "[CONF:httpMockUrl]/test/[CTXT:id]"
     When I send a HTTP "DELETE" request
     Then the HTTP status code must be "204"
      And the HTTP response body must be empty
     When I send a HTTP "DELETE" request
     Then the HTTP status code must be "204"
      And the HTTP response body must be empty

  @mockhttp
  Scenario: Mock one-shot request
    Given I store "[UUID]" in context "id"
      And I mock the HTTP request at "[CONF:httpMockUrl]" with the JSON
      """
      {
        "permanent": false,
        "request": {
          "method": "DELETE",
          "path": "/test/[CTXT:id]"
        },
        "response": {
          "status": 204,
          "body": ""
        }
      }
      """
    Given the HTTP endpoint "[CONF:httpMockUrl]/test/[CTXT:id]"
     When I send a HTTP "DELETE" request
     Then the HTTP status code must be "204"
      And the HTTP response body must be empty
     When I send a HTTP "DELETE" request
     Then the HTTP status code must be "404"

  @mockhttp
  Scenario: Mock request with full configuration
    Given I store "[UUID]" in context "id"
      And I mock the HTTP request at "[CONF:httpMockUrl]" with the JSON
      """
      {
        "request": {
          "method": "POST",
          "path": "/test/[CTXT:id]"
        },
        "response": {
          "status": 201,
          "headers": {
            "Content-Type": ["application/json"]
          },
          "body": "{\"value\": \"test mock response\"}"
        }
      }
      """
    Given the HTTP endpoint "[CONF:httpMockUrl]/test/[CTXT:id]"
     When I send a HTTP "POST" request
     Then the HTTP status code must be "201"
      And the HTTP request headers
          | Content-Type | application/json |
      And the HTTP response body must have the JSON properties
          | value | test mock response |

  @mockhttp
  Scenario: Mock request with timeout
    Given I store "[UUID]" in context "id"
      And I mock the HTTP request at "[CONF:httpMockUrl]" with the JSON
      """
      {
        "request": {
          "method": "GET",
          "path": "/test/[CTXT:id]"
        },
        "latency": 1000
      }
      """
    Given the HTTP endpoint "[CONF:httpMockUrl]/test/[CTXT:id]"
      And an HTTP timeout of "300" milliseconds
     When I send a HTTP "GET" request
     Then the HTTP response timed out

  @mockhttp
  Scenario: Send unmatching request to mock server
    Given the HTTP endpoint "[CONF:httpMockUrl]/test/[UUID]"
     When I send a HTTP "GET" request
     Then the HTTP status code must be "404"
