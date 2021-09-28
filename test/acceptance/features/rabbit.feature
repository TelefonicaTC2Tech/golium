Feature: Rabbit client

  @rabbit
  Scenario: Publish and subscribe a text message
    Given the rabbit endpoint "[CONF:rabbitmq]"
      And I subscribe to the rabbit topic "test-rabbit-text-properties-topic"
     When I publish a message to the rabbit topic "test-rabbit-text-properties-topic" with the text
          """
          This is a test message
          """
     Then I wait up to "3" seconds for a rabbit message with the text
          """
          This is a test message
          """

  @rabbit
  Scenario: Publish and subscribe a JSON message
    Given the rabbit endpoint "[CONF:rabbitmq]"
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

  @rabbit1
  Scenario: Publish and subscribe a JSON message. Use standard properties
    Given I generate a UUID and store it in context "CorrelationId"
    Given the rabbit endpoint "[CONF:rabbitmq]"
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

  @rabbit
  Scenario: Publish and subscribe three JSON message. Use standard properties
    Given I generate a UUID and store it in context "CorrelationId"
    Given the rabbit endpoint "[CONF:rabbitmq]"
      And I subscribe to the rabbit topic "test-rabbit-json-properties-[CTXT:CorrelationId]"
      And I set standard rabbit properties
          | ContentType   | application/json     |
          | CorrelationId | [CTXT:CorrelationId] |
     When I publish a message to the rabbit topic "test-rabbit-json-properties-[CTXT:CorrelationId]" with the JSON properties
          | id0       | abc0    |
          | name0     | Golium0 |
     When I publish a message to the rabbit topic "test-rabbit-json-properties-[CTXT:CorrelationId]" with the JSON properties
          | id1       | abc1    |
          | name1     | Golium1 |
     When I publish a message to the rabbit topic "test-rabbit-json-properties-[CTXT:CorrelationId]" with the JSON properties
          | id2       | abc2    |
          | name2     | Golium2 |
     Then I wait up to "5" seconds for "3" rabbit messages with the standard properties
          | ContentType   | application/json     |
          | CorrelationId | [CTXT:CorrelationId] |
      And the body of the rabbit message in position "0" has the JSON properties
          | id0      | abc0    |
          | name0    | Golium0 |
      And the body of the rabbit message in position "1" has the JSON properties
          | id1      | abc1    |
          | name1    | Golium1 |
      And the body of the rabbit message in position "2" has the JSON properties
          | id2      | abc2    |
          | name2    | Golium2 |

  @rabbit
  Scenario: Publish and subscribe a JSON message with rabbit headers and standard properties
    Given the rabbit endpoint "[CONF:rabbitmq]"
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

  @rabbit
  Scenario: Publish and subscribe a text message with rabbit headers
    Given the rabbit endpoint "[CONF:rabbitmq]"
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

  @rabbit
  Scenario: Subscribe and waits for no message
    Given I generate a UUID and store it in context "CorrelationId"
    Given the rabbit endpoint "[CONF:rabbitmq]"
     When I subscribe to the rabbit topic "test-rabbit-empty-topic"
     Then I wait up to "3" seconds without a rabbit message with the standard properties
          | CorrelationId | [CTXT:CorrelationId] |
