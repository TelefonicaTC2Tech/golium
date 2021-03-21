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
