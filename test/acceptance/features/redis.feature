Feature: Redis client

  Scenario: Set and get a text message
    Given the redis endpoint
          | addr     | localhost:6379 |
          | db       | 0              |
     When I set the redis key "golium:key:text" with the text
          """
          This is a test value
          """
     Then the redis key "golium:key:text" must have the text
          """
          This is a test value
          """

  Scenario: Set and get a JSON message
    Given the redis endpoint
          | addr     | localhost:6379 |
          | db       | 0              |
     When I set the redis key "golium:key:json" with the JSON properties
          | golium.number  | [NUMBER:4] |
          | golium.string  | test       |
          | golium.bool    | [TRUE]     |
     Then the redis key "golium:key:json" must have the JSON properties
          | golium.number  | [NUMBER:4] |
          | golium.string  | test       |
          | golium.bool    | [TRUE]     |

  Scenario: Set and get a JSON message with TTL
    Given the redis endpoint
          | addr     | localhost:6379 |
          | db       | 0              |
      And the redis TTL of "500" millis
     When I set the redis key "golium:key:ttl" with the text
          """
          {
               "golium": {
                    "number": 4,
                    "string": "test",
                    "bool": true
               }
          }
          """
     Then the redis key "golium:key:ttl" must have the JSON properties
          | golium.number  | [NUMBER:4] |
          | golium.string  | test       |
          | golium.bool    | [TRUE]     |
     When I wait for "600" millis
     Then the redis key "golium:key:ttl" must be empty

  Scenario: Publish and subscribe a text message
    Given the redis endpoint
          | addr     | localhost:6379 |
          | db       | 0              |
      And I subscribe to the redis topic "test-topic"
     When I publish a message to the redis topic "test-topic" with the text
          """
          This is a test message
          """
     Then I wait up to "3" seconds for a redis message with the text
          """
          This is a test message
          """
     And I unsubscribe from the redis topic "test-topic"

  Scenario: Publish and subscribe a JSON message
    Given the redis endpoint
          | addr     | localhost:6379 |
          | db       | 0              |
      And I subscribe to the redis topic "test-topic"
     When I publish a message to the redis topic "test-topic" with the JSON properties
          | id       | abc    |
          | name     | Golium |
     Then I wait up to "3" seconds for a redis message with the JSON properties
          | id       | abc    |
          | name     | Golium |
     And I unsubscribe from the redis topic "test-topic"
