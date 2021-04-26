Feature: Rabbit client

  @rabbit
  Scenario: Publish and subscribe a text message
    Given the rabbit endpoint "amqp://guest:guest@localhost:5672/"
      And I subscribe to the rabbit topic "test-rabbit-text-properties-topic"
     When I publish a message to the rabbit topic "test-rabbit-text-properties-topic" with the text
          """
          This is a test message
          """
     Then I wait up to "3" seconds for a rabbit message with the text
          """
          This is a test message
          """
      And I unsubscribe from the rabbit topic "test-rabbit-text-properties-topic"

  @rabbit
  Scenario: Publish and subscribe a JSON message
    Given the rabbit endpoint "amqp://guest:guest@localhost:5672/"
      And I subscribe to the rabbit topic "test-rabbit-json-properties-topic"
      And I set standard rabbit properties
          | ContentType   | application/json     |
          | CorrelationId | [CTXT:CorrelationId] |
     When I publish a message to the rabbit topic "test-rabbit-json-properties-topic" with the JSON properties
          | id       | abc    |
          | name     | Golium |
     Then I wait up to "3" seconds for a rabbit message with the JSON properties
          | id       | abc    |
          | name     | Golium |
      And I unsubscribe from the rabbit topic "test-rabbit-json-properties-topic"

  @rabbit
  Scenario: Publish and subscribe a JSON message. Use standard properties
    Given I generate a UUID and store it in context "CorrelationId"
    Given the rabbit endpoint "amqp://guest:guest@localhost:5672/"
      And I subscribe to the rabbit topic "test-rabbit-json-properties-[CTXT:CorrelationId]"
      And I set standard rabbit properties
          | ContentType   | application/json     |
          | CorrelationId | [CTXT:CorrelationId] |
     When I publish a message to the rabbit topic "test-rabbit-json-properties-[CTXT:CorrelationId]" with the JSON properties
          | id       | abc    |
          | name     | Golium |
     Then I wait up to "3" seconds for a rabbit message with the standard properties
          | ContentType   | application/json     |
          | CorrelationId | [CTXT:CorrelationId] |
      And the rabbit message body has the JSON properties
          | id       | abc    |
          | name     | Golium |
      And I unsubscribe from the rabbit topic "test-rabbit-json-properties-[CTXT:CorrelationId]"

  @rabbit
  Scenario: Publish and subscribe a JSON message with rabbit headers and standard properties
    Given the rabbit endpoint "amqp://guest:guest@localhost:5672/"
      And I subscribe to the rabbit topic "test-rabbit-headers-properties-topic"
      And I set rabbit headers
          | Header1   | value1 |
          | Header2   | value2 |
     When I publish a message to the rabbit topic "test-rabbit-headers-properties-topic" with the JSON properties
          | id       | abc    |
          | name     | Golium |
     Then I wait up to "3" seconds for a rabbit message with the JSON properties
          | id       | abc    |
          | name     | Golium |
      And the rabbit message has the rabbit headers
          | Header1   | value1 |
          | Header2   | value2 |
      And I unsubscribe from the rabbit topic "test-rabbit-headers-properties-topic"

  @rabbit
  Scenario: Publish and subscribe a text message with rabbit headers
    Given the rabbit endpoint "amqp://guest:guest@localhost:5672/"
      And I subscribe to the rabbit topic "test-rabbit-headers-topic"
      And I set rabbit headers
          | Header1   | value1 |
          | Header2   | value2 |
     When I publish a message to the rabbit topic "test-rabbit-headers-topic" with the text
          """
          """
     Then I wait up to "3" seconds for a rabbit message with the text
          """
          """
      And the rabbit message has the rabbit headers
          | Header1   | value1 |
          | Header2   | value2 |
      And I unsubscribe from the rabbit topic "test-rabbit-headers-topic"

  @rabbit
  Scenario: Subscribe and waits for no message
    Given I generate a UUID and store it in context "CorrelationId"
    Given the rabbit endpoint "amqp://guest:guest@localhost:5672/"
     When I subscribe to the rabbit topic "test-rabbit-empty-topic"
     Then I wait up to "3" seconds without a rabbit message with the standard properties
          | CorrelationId | [CTXT:CorrelationId] |
      And I unsubscribe from the rabbit topic "test-rabbit-empty-topic"
