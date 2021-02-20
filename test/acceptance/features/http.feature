Feature: HTTP client

  Scenario: Send a GET request
    Given the HTTP endpoint "[CONF:url]/anything"
      And the HTTP path "/test1"
      And the HTTP query parameters
          | exists | true |
          | sort   | name |
      And the HTTP request headers
          | Authorization | Bearer access-token |
     When I send a HTTP "GET" request
     Then the HTTP status code must be "200"
      And the HTTP response body must comply with the JSON schema "test-schema"
      And the HTTP response body must have the JSON properties
          | args.exists           | true                                            |
          | args.sort             | name                                            |
          | headers.Authorization | Bearer access-token                             |
          | headers.Host          | [CONF:host]                                     |
          | method                | GET                                             |
          | url                   | [CONF:url]/anything/test1?exists=true&sort=name |

  Scenario: Send a POST request
    Given the HTTP endpoint "[CONF:url]/anything"
      And the HTTP path "/test2"
      And the HTTP request headers
          | Authorization | Bearer access-token |
      And the JSON properties in the HTTP request body
          | name     | golium  |
          | active   | [TRUE]  |
          | inactive | [FALSE] |
          | empty    | [EMPTY] |
          | null     | [NULL]  |
     When I send a HTTP "POST" request
     Then the HTTP status code must be "200"
      And the HTTP response body must comply with the JSON schema "test-schema"
      And the HTTP response body must have the JSON properties
          | headers.Authorization | Bearer access-token       |
          | headers.Host          | [CONF:host]               |
          | method                | POST                      |
          | url                   | [CONF:url]/anything/test2 |
          | json.name             | golium                    |
          | json.active           | [TRUE]                    |
          | json.inactive         | [FALSE]                   |
          | json.empty            | [EMPTY]                   |
          | json.null             | [NULL]                    |
