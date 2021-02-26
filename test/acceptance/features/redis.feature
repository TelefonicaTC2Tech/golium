Feature: Redis client

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
