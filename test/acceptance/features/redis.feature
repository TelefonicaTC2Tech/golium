Feature: Redis client

  Background:
      Given the redis endpoint
          | param | value                 |
          | addr  | [CONF:redis.endpoint] |
          | db    | 0                     |

  @redis
  Scenario: Set and get a text message
     Given I generate a UUID and store it in context "key"
     When I set the redis key "golium:key:text:[CTXT:key]" with the text
          """
          This is a test value with id: [CTXT:key]
          """
     Then the redis key "golium:key:text:[CTXT:key]" must have the text
          """
          This is a test value with id: [CTXT:key]
          """

  @redis
  Scenario: Set and get a mapped message
    Given I generate a UUID and store it in context "key"
     When I set the redis key "golium:key:mapped:[CTXT:key]" with hash properties
          | param         | value      |
          | golium.number | [NUMBER:4] |
          | golium.string | test       |
          | golium.bool   | [TRUE]     |
          | golium.id     | [CTXT:key] |
     Then the redis key "golium:key:mapped:[CTXT:key]" must have hash properties
          | param         | value      |
          | golium.number | 4          |
          | golium.string | test       |
          | golium.bool   | 1          |
          | golium.id     | [CTXT:key] |

  @redis
  Scenario: Set and get a mapped message with TTL
    Given I generate a UUID and store it in context "key"
      And the redis TTL of "500" millis
     When I set the redis key "golium:key:ttl:mapped:[CTXT:key]" with hash properties
          | param         | value      |
          | golium.number | [NUMBER:4] |
          | golium.string | test       |
          | golium.bool   | [TRUE]     |
          | golium.id     | [CTXT:key] |
     Then the redis key "golium:key:ttl:mapped:[CTXT:key]" must have hash properties
          | param         | value      |
          | golium.number | 4          |
          | golium.string | test       |
          | golium.bool   | 1          |
          | golium.id     | [CTXT:key] |
     When I wait for "600" millis
     Then the redis key "golium:key:ttl:mapped:[CTXT:key]" must not exist

  @redis
  Scenario: Set and get a JSON message
    Given I generate a UUID and store it in context "key"
     When I set the redis key "golium:key:json:[CTXT:key]" with the JSON properties
          | param         | value      |
          | golium.number | [NUMBER:4] |
          | golium.string | test       |
          | golium.bool   | [TRUE]     |
          | golium.id     | [CTXT:key] |
     Then the redis key "golium:key:json:[CTXT:key]" must have the JSON properties
          | param         | value      |
          | golium.number | [NUMBER:4] |
          | golium.string | test       |
          | golium.bool   | [TRUE]     |
          | golium.id     | [CTXT:key] |

  @redis
  Scenario: Set and get a JSON message with TTL
    Given I generate a UUID and store it in context "key"
      And the redis TTL of "500" millis
     When I set the redis key "golium:key:ttl:json:[CTXT:key]" with the text
          """
          {
               "golium": {
                    "number": 4,
                    "string": "test",
                    "bool": true,
                    "id": "[CTXT:key]"
               }
          }
          """
     Then the redis key "golium:key:ttl:json:[CTXT:key]" must have the JSON properties
          | param         | value      |
          | golium.number | [NUMBER:4] |
          | golium.string | test       |
          | golium.bool   | [TRUE]     |
          | golium.id     | [CTXT:key] |
     When I wait for "600" millis
     Then the redis key "golium:key:ttl:json:[CTXT:key]" must not exist

  @redis
  Scenario: Publish and subscribe a text message
    Given I subscribe to the redis topic "test-topic"
     When I publish a message to the redis topic "test-topic" with the text
          """
          This is a test message with id: [CTXT:key]
          """
     Then I wait up to "3" seconds for a redis message with the text
          """
          This is a test message with id: [CTXT:key]
          """
     And I unsubscribe from the redis topic "test-topic"

  @redis
  Scenario: Publish and subscribe a JSON message
    Given I subscribe to the redis topic "test-topic"
     When I publish a message to the redis topic "test-topic" with the JSON properties
          | param | value  |
          | id    | abc    |
          | name  | Golium |
     Then I wait up to "3" seconds for a redis message with the JSON properties
          | param | value  |
          | id    | abc    |
          | name  | Golium |
      And I wait up to "3" seconds without a redis message with the JSON properties
          | param | value      |
          | id    | abc        |
          | name  | unexpected |
     And I unsubscribe from the redis topic "test-topic"

  @redis
  Scenario: Select database, set and get a text message
     Given I generate a UUID and store it in context "key"
      When I select the redis database "1"
       And I set the redis key "golium:key:text:[CTXT:key]" with the text
           """
           This is a test value with id: [CTXT:key]
           """
      Then the redis key "golium:key:text:[CTXT:key]" must have the text
           """
           This is a test value with id: [CTXT:key]
           """

  @redis
  Scenario: Select database, set key, select previous database and key must not exists
     Given I generate a UUID and store it in context "key"
      When I select the redis database "1"
       And I set the redis key "golium:key:text:[CTXT:key]" with the text
           """
           This is a test value with id: [CTXT:key]
           """
       And I select the redis database "0"
      Then the redis key "golium:key:text:[CTXT:key]" must not exists
