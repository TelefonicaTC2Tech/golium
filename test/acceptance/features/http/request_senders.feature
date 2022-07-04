Feature: HTTP Request Senders


  @http @request_senders
  Scenario: Send a GET request without parameters
    When I send a "GET" request to "posts" endpoint
    Then the HTTP status code must be "200"

  @http @request_senders
  Scenario: Send a GET request without parameters and backslash
    When I send a "GET" request to "posts" endpoint without last backslash
    Then the HTTP status code must be "200"

  @http @request_senders
  Scenario: Send a GET request with valid API-KEY
    When I send a "GET" request to "posts" endpoint with "valid" API-KEY
    Then the HTTP status code must be "200"

  @http @request_senders
  Scenario: Send a GET request without credentials
    When I send a "GET" request to "posts" endpoint without credentials
    Then the HTTP status code must be "200"

  @http @request_senders
  Scenario: Send a GET request with query parameters
    When I send a "GET" request to "posts" endpoint with query params
      | field  | value |
      | userId | 1     |
      | id     | 8     |
    Then the HTTP status code must be "200"

  @http @request_senders
  Scenario Outline: Send a GET request with single filter
    When I send a "GET" request to "posts" endpoint with "<filters>" filters
    Then the HTTP status code must be "200"
    And the HTTP response body must have the JSON properties
      | property | value             |
      | #        | <filtered_values> |
    Examples:
      | filters  | filtered_values |
      | userId=1 | [NUMBER:10]     |
      | id=8     | [NUMBER:1]      |

  @http @request_senders
  Scenario Outline: Send a GET request with multiple filter
    When I send a "GET" request to "posts" endpoint with "<filters>" filters
    Then the HTTP status code must be "200"
    And the HTTP response body must have the JSON properties
      | property | value             |
      | #        | <filtered_values> |
    Examples:
      | filters       | filtered_values |
      | id=8&userId=1 | [NUMBER:1]      |

  @http @request_senders
  Scenario: Send a GET request with multiple filter
    When I send a "GET" request to "posts" endpoint with path "8"
    Then the HTTP status code must be "200"
    And the HTTP response body must have the JSON properties
      | property | value      |
      | id       | [NUMBER:8] |
      | userId   | [NUMBER:1] |

  @http @request_senders @json
  Scenario: Send a POST request defined by a json string in file
    When I send a "POST" request to "posts" with a JSON body that includes "example1"
    Then the HTTP status code must be "201"

  @http @request_senders @json
  Scenario: Send a POST request defined by a json string in file removing parameters
    When I send a "POST" request to "posts" with a JSON body that includes "example1" without
      | parameter |
      | boolean   |
    Then the HTTP status code must be "201"

  @http @request_senders @json
  Scenario: Send a POST request defined by a json string in file modifying parameters
    When I send a "POST" request to "posts" with a JSON body that includes "example1" modifying
      | parameter | value  |
      | boolean   | [TRUE] |
    Then the HTTP status code must be "201"

  @http @request_senders @json
  Scenario: Send a POST request with pre configured header
    Given the HTTP request headers
      | Parameter     | Value                  |
      | Authorization | Bearer 1234567890AEIOU |
    When I send a "POST" request to "posts" with a JSON body that includes "example1"
    Then the HTTP status code must be "201"    