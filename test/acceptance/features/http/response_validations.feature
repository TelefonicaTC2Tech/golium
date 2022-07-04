Feature: HTTP response validations


  @json
  Scenario: Validate a POST request defined by a json string in file
    When I send a "POST" request to "posts" with a JSON body that includes "example1"
    Then the "posts" response message should match with "example1" JSON message

  @json
  Scenario: Validate a POST request with path defined by a json string in file
    When I send a "POST" request to "posts" with path "1" with a JSON body that includes "example1"
    Then the "posts" response message should match with "empty" JSON message

  @json
  Scenario: Validate a POST request defined by a json string in file removing parameters
    When I send a "POST" request to "posts" with a JSON body that includes "example1" without
      | parameter |
      | boolean   |
    Then the "posts" response message should match with "example1" JSON message without
      | parameter |
      | boolean   |

  @json
  Scenario: Validate a POST request defined by a json string in file modifying parameters
    When I send a "POST" request to "posts" with a JSON body that includes "example1" modifying
      | parameter | value  |
      | boolean   | [TRUE] |
    Then the "posts" response message should match with "example1" JSON message modifying
      | parameter | value  |
      | boolean   | [TRUE] |

  @json
  Scenario: Validate response modifying nested parameters
    When I send a "GET" request to "users" endpoint with path "1"
    Then the "users" response message should match with "example1" JSON message modifying
      | parameter       | value    |
      | address.geo.lat | -37.3159 |
