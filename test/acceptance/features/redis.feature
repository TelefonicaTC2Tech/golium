Feature: Redis client

  @redis
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

  @redis
  Scenario: Set and get a mapped message
    Given the redis endpoint
          | addr     | localhost:6379 |
          | db       | 0              |
     When I set the redis key "golium:key:mapped" with hash properties
          | golium.number  | [NUMBER:4] |
          | golium.string  | test       |
          | golium.bool    | [TRUE]     |
     Then the redis key "golium:key:mapped" must have hash properties
          | golium.number  | 4    |
          | golium.string  | test |
          | golium.bool    | 1    |

  @redis
  Scenario: Set and get a mapped message with TTL
    Given the redis endpoint
          | addr     | localhost:6379 |
          | db       | 0              |
      And the redis TTL of "500" millis
     When I set the redis key "golium:key:ttl:mapped" with hash properties
          | golium.number  | [NUMBER:4] |
          | golium.string  | test       |
          | golium.bool    | [TRUE]     |
     Then the redis key "golium:key:ttl:mapped" must have hash properties
          | golium.number  | 4    |
          | golium.string  | test |
          | golium.bool    | 1    |
     When I wait for "600" millis
     Then the redis key "golium:key:ttl:mapped" must not exist

  @redis
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

  @redis
  Scenario: Set and get a JSON message with TTL
    Given the redis endpoint
          | addr     | localhost:6379 |
          | db       | 0              |
      And the redis TTL of "500" millis
     When I set the redis key "golium:key:ttl:json" with the text
          """
          {
               "golium": {
                    "number": 4,
                    "string": "test",
                    "bool": true
               }
          }
          """
     Then the redis key "golium:key:ttl:json" must have the JSON properties
          | golium.number  | [NUMBER:4] |
          | golium.string  | test       |
          | golium.bool    | [TRUE]     |
     When I wait for "600" millis
     Then the redis key "golium:key:ttl:json" must not exist

  @redis
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

  @redis
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
      And I wait up to "3" seconds without a redis message with the JSON properties
          | id       | abc        |
          | name     | unexpected |
     And I unsubscribe from the redis topic "test-topic"
