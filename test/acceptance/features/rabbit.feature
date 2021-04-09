Feature: Rabbit client

  @rabbit
  Scenario: Publish and subscribe a text message
    Given the rabbit endpoint "amqp://admin:password@localhost:5672/"
      And I subscribe to the rabbit topic "test-topic"
     When I publish a message to the rabbit topic "test-topic" with the text
          """
          This is a test message
          """
     Then I wait up to "3" seconds for a rabbit message with the text
          """
          This is a test message
          """
     And I unsubscribe from the rabbit topic "test-topic"

  @rabbit
  Scenario: Publish and subscribe a JSON message
    Given the rabbit endpoint "amqp://admin:password@localhost:5672/"
      And I subscribe to the rabbit topic "test-topic"
     When I publish a message to the rabbit topic "test-topic" with the JSON properties
          | id       | abc    |
          | name     | Golium |
     Then I wait up to "3" seconds for a rabbit message with the JSON properties
          | id       | abc    |
          | name     | Golium |
      And I wait up to "3" seconds without a rabbit message with the JSON properties
          | id       | abc        |
          | name     | unexpected |
     And I unsubscribe from the rabbit topic "test-topic"
