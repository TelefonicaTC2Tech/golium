Feature: Rabbit client

  @rabbit
  Scenario: Publish and subscribe a text message
    Given the rabbit endpoint "amqp://guest:guest@localhost:5672/"
      And I subscribe to the rabbit topic "test-rabbitmq-text-topic"
     When I publish a message to the rabbit topic "test-rabbitmq-text-topic" with the text
          """
          This is a test message
          """
     Then I wait up to "3" seconds for a rabbit message
      And the rabbit message body has the text
          """
          This is a test message
          """
      And I unsubscribe from the rabbit topic "test-rabbitmq-text-topic"

  @rabbit
  Scenario: Publish and subscribe a JSON message
    Given the rabbit endpoint "amqp://guest:guest@localhost:5672/"
      And I subscribe to the rabbit topic "test-rabbitmq-json-topic"
     When I publish a message to the rabbit topic "test-rabbitmq-json-topic" with the JSON properties
          | id       | abc    |
          | name     | Golium |
     Then I wait up to "3" seconds for a rabbit message
      And the rabbit message body has the JSON properties
          | id       | abc    |
          | name     | Golium |
      And I unsubscribe from the rabbit topic "test-rabbitmq-json-topic"

  @rabbit
  Scenario: Publish and subscribe a text message with standard rabbitmq properties
    Given the rabbit endpoint "amqp://guest:guest@localhost:5672/"
      And I subscribe to the rabbit topic "test-rabbitmq-properties-topic"
      And I set standard rabbitmq properties
          | ContentType   | text/plain       |
          | CorrelationId | Unica-Correlator |
     When I publish a message to the rabbit topic "test-rabbitmq-properties-topic" with the text
          """
          """
     Then I wait up to "3" seconds for a rabbit message
      And the rabbit message has the standard rabbitmq properties
          | ContentType   | text/plain       |
          | CorrelationId | Unica-Correlator |
      And I unsubscribe from the rabbit topic "test-rabbitmq-properties-topic"

  @rabbit
  Scenario: Publish and subscribe a text message with rabbitmq headers
    Given the rabbit endpoint "amqp://guest:guest@localhost:5672/"
      And I subscribe to the rabbit topic "test-rabbitmq-headers-topic"
      And I set rabbitmq headers
          | ContentType     | text/plain       |
          | CorrelationId   | Unica-Correlator |
     When I publish a message to the rabbit topic "test-rabbitmq-headers-topic" with the text
          """
          """
     Then I wait up to "3" seconds for a rabbit message
      And the rabbit message has the rabbitmq headers
          | ContentType     | text/plain       |
          | CorrelationId   | Unica-Correlator |
      And I unsubscribe from the rabbit topic "test-rabbitmq-headers-topic"

  @rabbit
  Scenario: Subscribe and waits for no message
    Given the rabbit endpoint "amqp://guest:guest@localhost:5672/"
     When I subscribe to the rabbit topic "test-rabbitmq-empty-topic"
     Then I wait up to "3" seconds without receiving a rabbit message
      And I unsubscribe from the rabbit topic "test-rabbitmq-empty-topic"
