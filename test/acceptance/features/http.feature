Feature: HTTP client

  @http
  Scenario: Send a GET request
    Given the HTTP endpoint "[CONF:http_bin_url]/anything"
      And the HTTP path "/test-query"
      And the HTTP query parameters
          | param  | value |
          | exists | true  |
          | sort   | name  |
      And the HTTP request headers
          | param         | value               |
          | Authorization | Bearer access-token |
     When I send a HTTP "GET" request
     Then the HTTP status code must be "200"
      And the HTTP response must contain the headers
          | param        | value            |
          | Content-Type | application/json |
      And the HTTP response body must comply with the JSON schema "test-schema"
      And the HTTP response body must have the JSON properties
          | param                 | value                                                |
          | headers.Authorization | Bearer access-token                                  |
          | headers.Content-Type  | [NULL]                                               |
          | headers.Host          | [CONF:host]                                          |
          | args.exists           | true                                                 |
          | args.sort             | name                                                 |
          | method                | GET                                                  |
          | url                   | [CONF:http_bin_url]/anything/test-query?exists=true&sort=name |

  @http
  Scenario: Send a POST request with a JSON body using a table of properties
    Given the HTTP endpoint "[CONF:http_bin_url]/anything"
      And the HTTP path "/test-json"
      And the HTTP request headers
          | param         | value               |
          | Authorization | Bearer access-token |
      And the JSON properties in the HTTP request body
          | param    | value          |
          | name     | golium         |
          | active   | [TRUE]         |
          | inactive | [FALSE]        |
          | empty    | [EMPTY]        |
          | null     | [NULL]         |
          | integer  | [NUMBER:1234]  |
          | float    | [NUMBER:-34.6] |
     When I send a HTTP "POST" request
     Then the HTTP status code must be "200"
      And the HTTP response must contain the headers
          | param        | value            |
          | Content-Type | application/json |
      And the HTTP response body must comply with the JSON schema "test-schema"
      And the HTTP response body must have the JSON properties
          | param                 | value                         |
          | headers.Authorization | Bearer access-token           |
          | headers.Content-Type  | application/json              |
          | headers.Host          | [CONF:host]                   |
          | method                | POST                          |
          | url                   | [CONF:http_bin_url]/anything/test-json |
          | json.name             | golium                        |
          | json.active           | [TRUE]                        |
          | json.inactive         | [FALSE]                       |
          | json.empty            | [EMPTY]                       |
          | json.null             | [NULL]                        |
          | json.integer          | [NUMBER:1234]                 |
          | json.float            | [NUMBER:-34.6]                |

  @http
  Scenario: Send a POST request with a x-www-form-urlencoded body using a table of properties
    Given the HTTP endpoint "[CONF:http_bin_url]/anything"
      And the HTTP path "/test-urlencoded"
      And the HTTP request headers
          | param         | value               |
          | Authorization | Bearer access-token |
      And the HTTP request body with the URL encoded properties
          | param   | value          |
          | name    | test           |
          | surname | golium         |
          | boolean | [TRUE]         |
          | float   | [NUMBER:-34.6] |
     When I send a HTTP "POST" request
     Then the HTTP status code must be "200"
      And the HTTP response must contain the headers
          | param        | value            |
          | Content-Type | application/json |
      And the HTTP response body must comply with the JSON schema "test-schema"
      And the HTTP response body must have the JSON properties
          | param                 | value                               |
          | headers.Authorization | Bearer access-token                 |
          | headers.Content-Type  | application/x-www-form-urlencoded   |
          | headers.Host          | [CONF:host]                         |
          | method                | POST                                |
          | url                   | [CONF:http_bin_url]/anything/test-urlencoded |
          | form.name             | test                                |
          | form.surname          | golium                              |
          | form.boolean          | true                                |
          | form.float            | -34.6                               |

  @http
  Scenario: Send a POST request defined by a json string using the context storage
    Given the HTTP endpoint "[CONF:http_bin_url]/anything"
      And the HTTP path "/test-content"
      And the HTTP request headers
          | param        | value            |
          | Content-Type | application/json |
      And the HTTP request body with the JSON
      """
      {
        "empty": "[EMPTY]",
        "boolean": [FALSE],
        "list": [
          { "attribute": "attribute0", "value": "value0"},
          { "attribute": "attribute1", "value": "value1"},
          { "attribute": "attribute2", "value": "value2"}
        ]
      }
      """
     When I send a HTTP "POST" request
      And the HTTP status code must be "200"
      And the HTTP response body must comply with the JSON schema "test-schema"
      And the HTTP response body must have the JSON properties
          | param                 | value      |
          | json.empty            | [EMPTY]    |
          | json.boolean          | [FALSE]    |
          | json.list.#           | [NUMBER:3] |
          | json.list.0.attribute | attribute0 |
          | json.list.0.value     | value0     |
          | json.list.1.attribute | attribute1 |
          | json.list.1.value     | value1     |
          | json.list.2.attribute | attribute2 |
          | json.list.2.value     | value2     |
      And I store the element "json.list.0.attribute" from the JSON HTTP response body in context "attr"
      And I store the element "json.list.0.value" from the JSON HTTP response body in context "val"
     Then the HTTP endpoint "[CONF:http_bin_url]/anything"
      And the HTTP path "/test3-2"
      And the HTTP request headers
          | param        | value            |
          | Content-Type | application/json |
      And the JSON properties in the HTTP request body
          | param     | value       |
          | attribute | [CTXT:attr] |
          | value     | [CTXT:val]  |
      And I send a HTTP "POST" request
      And the HTTP status code must be "200"
      And the HTTP response body must comply with the JSON schema "test-schema"
      And the HTTP response body must have the JSON properties
          | param          | value      |
          | json.attribute | attribute0 |
          | json.value     | value0     |
  
  @http
  Scenario: Send a POST request defined by a json string in file and check the response
    Given the HTTP endpoint "[CONF:http_bin_url]/anything"
      And the HTTP path "/test-content"
      And the HTTP request headers
          | param        | value            |
          | Content-Type | application/json |
      And the HTTP request body with the JSON "example1" from "http" file
     When I send a HTTP "POST" request
      And the HTTP status code must be "200"
      And the HTTP response "json" must match with the JSON "example1" from "http" file

  @http
  Scenario: Send a POST request defined by a json string in file removing parameters and check the response
    Given the HTTP endpoint "[CONF:http_bin_url]/anything"
      And the HTTP path "/test-content"
      And the HTTP request headers
          | param        | value            |
          | Content-Type | application/json |
      And the HTTP request body with the JSON "example1" from "http" file without
          | param   |
          | boolean |
     When I send a HTTP "POST" request
      And the HTTP status code must be "200"
      And the HTTP response "json" must match with the JSON "example1" from "http" file without
          | param   |
          | boolean |

  @http
  Scenario: Follow redirection
    Given the HTTP endpoint "http://www.github.com"
     When I send a HTTP "GET" request
      And the HTTP status code must be "200"

  @http
  Scenario: Follow no redirection
    Given the HTTP endpoint "http://www.github.com"
      And the HTTP client does not follow any redirection
     When I send a HTTP "GET" request
      And the HTTP status code must be "301"
      And the HTTP response must contain the headers
          | param    | value                        |
          | Location | https://www.github.com/ |
      And I store the header "Location" from the HTTP response in context "header.location"

  @http
  Scenario: Set HTTP request host
    Given the HTTP endpoint "[CONF:http_bin_url]/headers"
      And the HTTP request headers
          | param | value       |
          | Host  | example.com |
     When I send a HTTP "GET" request
      And the HTTP status code must be "200"
      And the HTTP response body must have the JSON properties
          | param        | value       |
          | headers.Host | example.com |

  @http
  Scenario: Validate Not found path with trailing slash
    Given the HTTP endpoint "[CONF:http_bin_url]/image/jpeg/"
     When I send a HTTP "POST" request
     Then the HTTP status code must be "404"

  @http
    Scenario: Send a GET request and check if a specific header is not in response headers
      Given the HTTP endpoint "[CONF:http_bin_url]/anything"
        And the HTTP path "/test-query"
        And the HTTP query parameters
            | param  | value |
            | exists | true  |
            | sort   | name  |
        And the HTTP request headers
            | param         | value               |
            | Authorization | Bearer access-token |
      When I send a HTTP "GET" request
      Then the HTTP status code must be "200"
        And the HTTP response must not contain the headers
            | param               |
            | non-existent-header |

  @http
    Scenario: Send a GET request to a site with self-signed cert and not verify it
      Given the HTTP endpoint "[CONF:urlSelfSigned]"
        And the HTTP client does not verify https cert
       When I send a HTTP "GET" request
       Then the HTTP status code must be "200"
