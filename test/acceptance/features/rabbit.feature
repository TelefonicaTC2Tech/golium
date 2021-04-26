Feature: Rabbit client

  @rabbit
  Scenario: Publish and subscribe a text message
    Given I generate a UUID and store it in context "CorrelationId"
    Given the rabbit endpoint "amqp://guest:guest@localhost:5672/"
      And I subscribe to the rabbit topic "test-rabbitmq-text-properties-topic"
      And I set standard rabbitmq properties
          | ContentType   | text/plain           |
          | CorrelationId | [CTXT:CorrelationId] |
     When I publish a message to the rabbit topic "test-rabbitmq-text-properties-topic" with the text
          """
          This is a test message
          """
     Then I wait up to "3" seconds for a rabbit message with the standard properties
          | ContentType   | text/plain           |
          | CorrelationId | [CTXT:CorrelationId] |
      And the rabbit message body has the text
          """
          This is a test message
          """
      And I unsubscribe from the rabbit topic "test-rabbitmq-text-properties-topic"

  @rabbit
  Scenario: Publish and subscribe a JSON message
    Given I generate a UUID and store it in context "CorrelationId"
    Given the rabbit endpoint "amqp://guest:guest@localhost:5672/"
      And I subscribe to the rabbit topic "test-rabbitmq-json-properties-topic"
      And I set standard rabbitmq properties
          | ContentType   | application/json     |
          | CorrelationId | [CTXT:CorrelationId] |
     When I publish a message to the rabbit topic "test-rabbitmq-json-properties-topic" with the JSON properties
          | id       | abc    |
          | name     | Golium |
     Then I wait up to "3" seconds for a rabbit message with the standard properties
          | ContentType   | application/json     |
          | CorrelationId | [CTXT:CorrelationId] |
      And the rabbit message body has the JSON properties
          | id       | abc    |
          | name     | Golium |
      And I unsubscribe from the rabbit topic "test-rabbitmq-json-properties-topic"

  @rabbit
  Scenario: Publish and subscribe a JSON message with rabbitmq headers and standard properties
    Given I generate a UUID and store it in context "CorrelationId"
    Given the rabbit endpoint "amqp://guest:guest@localhost:5672/"
      And I subscribe to the rabbit topic "test-rabbitmq-headers-properties-topic"
      And I set rabbitmq headers
          | Header1   | value1 |
          | Header2   | value2 |
      And I set standard rabbitmq properties
          | CorrelationId   | [CTXT:CorrelationId] |
     When I publish a message to the rabbit topic "test-rabbitmq-headers-properties-topic" with the JSON properties
          | id       | abc    |
          | name     | Golium |
     Then I wait up to "3" seconds for a rabbit message with the standard properties
          | CorrelationId | [CTXT:CorrelationId] |
      And the rabbit message body has the JSON properties
          | id       | abc    |
          | name     | Golium |
      And the rabbit message has the rabbitmq headers
          | Header1   | value1 |
          | Header2   | value2 |

      And I unsubscribe from the rabbit topic "test-rabbitmq-headers-properties-topic"

  @rabbit
  Scenario: Publish and subscribe a text message with rabbitmq headers
    Given I generate a UUID and store it in context "CorrelationId"
    Given the rabbit endpoint "amqp://guest:guest@localhost:5672/"
      And I subscribe to the rabbit topic "test-rabbitmq-headers-topic"
      And I set rabbitmq headers
          | Header1   | value1 |
          | Header2   | value2 |
      And I set standard rabbitmq properties
          | CorrelationId   | [CTXT:CorrelationId] |
     When I publish a message to the rabbit topic "test-rabbitmq-headers-topic" with the text
          """
          """
     Then I wait up to "3" seconds for a rabbit message with the standard properties
          | CorrelationId | [CTXT:CorrelationId] |
      And the rabbit message has the rabbitmq headers
          | Header1   | value1 |
          | Header2   | value2 |
      And I unsubscribe from the rabbit topic "test-rabbitmq-headers-topic"

  @rabbit
  Scenario: Subscribe and waits for no message
    Given I generate a UUID and store it in context "CorrelationId"
    Given the rabbit endpoint "amqp://guest:guest@localhost:5672/"
     When I subscribe to the rabbit topic "test-rabbitmq-empty-topic"
     Then I wait up to "3" seconds without a rabbit message with the standard properties
          | CorrelationId | [CTXT:CorrelationId] |
      And I unsubscribe from the rabbit topic "test-rabbitmq-empty-topic"
