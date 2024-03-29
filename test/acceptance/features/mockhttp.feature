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
      And the HTTP response must contain the headers
          | param        | value            |
          | Content-Type | application/json |
      And the HTTP response body must have the JSON properties
          | param | value              |
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
  Scenario: Mock permanent request with multiending path
    Given I mock the HTTP request at "[CONF:httpMockUrl]" with the JSON
      """
      {
        "permanent": true,
        "request": {
          "method": "DELETE",
          "path": "/test-me-path/<*>"
        },
        "response": {
          "status": 204,
          "body": ""
        }
      }
      """
    Given the HTTP endpoint "[CONF:httpMockUrl]/test-me-path/a0b1c2"
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
  Scenario: Mock one-shot request with multiending path
    Given I mock the HTTP request at "[CONF:httpMockUrl]" with the JSON
      """
      {
        "permanent": false,
        "request": {
          "method": "DELETE",
          "path": "/test/<*>"
        },
        "response": {
          "status": 204,
          "body": ""
        }
      }
      """
    Given the HTTP endpoint "[CONF:httpMockUrl]/test/a0b1c2"
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
      And the HTTP response must contain the headers
          | param        | value            |
          | Content-Type | application/json |
      And the HTTP response body must have the JSON properties
          | param | value              |
          | value | test mock response |

  @mockhttp
  Scenario: Mock request with full configuration and complex tag nesting
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
          "body": "{ \"values\": [ { \"type\": \"blacklist\", \"path\": \"/[CTXT:id]-[CTXT:id]\" } ] }"
        }
      }
      """
    Given the HTTP endpoint "[CONF:httpMockUrl]/test/[CTXT:id]"
     When I send a HTTP "POST" request
     Then the HTTP status code must be "201"
      And the HTTP response must contain the headers
          | param        | value            |
          | Content-Type | application/json |
      And the HTTP response body must have the JSON properties
          | param         | value                |
          | values.0.path | /[CTXT:id]-[CTXT:id] |

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

  @mockhttp
  Scenario: Mock request with text/plain content
    Given I store "[UUID]" in context "id"
      And I mock the HTTP request at "[CONF:httpMockUrl]" with the JSON
      """
      {
        "request": {
          "method": "GET",
          "path": "/test/text-plain/[CTXT:id]"
        },
        "response": {
          "status": 200,
          "headers": {
            "Content-Type": ["text/plain"]
          },
          "body": "Just a plain text format"
        }
      }
      """
    Given the HTTP endpoint "[CONF:httpMockUrl]/test/text-plain/[CTXT:id]"
     When I send a HTTP "GET" request
     Then the HTTP status code must be "200"
      And the HTTP response must contain the headers
          | param        | value      |
          | Content-Type | text/plain |
      And the HTTP response body must be the text
          """
          Just a plain text format
          """
