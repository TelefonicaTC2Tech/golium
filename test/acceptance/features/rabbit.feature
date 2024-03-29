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
  Scenario: Publish, subscribe and consume two text messages
    Given the rabbit endpoint "[CONF:rabbitmq]"
    And I subscribe to the rabbit topic "test-rabbit-text-properties-topic"
    When I publish a message to the rabbit topic "test-rabbit-text-properties-topic" with the text
          """
          This is a test message
          """
    And  I publish a message to the rabbit topic "test-rabbit-text-properties-topic" with the text
          """
          This is a second test message
          """
    Then I wait up to "3" seconds for a rabbit message with the text
          """
          This is a test message
          """
    And I wait up to "3" seconds for a rabbit message with the text
          """
          This is a second test message
          """

  @rabbit
  Scenario: Publish and subscribe a JSON message
    Given the rabbit endpoint "[CONF:rabbitmq]"
    And I subscribe to the rabbit topic "test-rabbit-json-properties-topic"
    And I set standard rabbit properties
      | param         | value                |
      | ContentType   | application/json     |
      | CorrelationId | [CTXT:CorrelationId] |
    When I publish a message to the rabbit topic "test-rabbit-json-properties-topic" with the JSON properties
      | param | value  |
      | id    | abc    |
      | name  | Golium |
    Then I wait up to "3" seconds for a rabbit message with the JSON properties
      | param | value  |
      | id    | abc    |
      | name  | Golium |

  @rabbit
  Scenario: Publish, subscribe and consume two JSON messages
    Given the rabbit endpoint "[CONF:rabbitmq]"
    And I subscribe to the rabbit topic "test-rabbit-json-properties-topic"
    And I set standard rabbit properties
      | param         | value                |
      | ContentType   | application/json     |
      | CorrelationId | [CTXT:CorrelationId] |
    When I publish a message to the rabbit topic "test-rabbit-json-properties-topic" with the JSON properties
      | param | value  |
      | id    | abc    |
      | name  | Golium |
    And I publish a message to the rabbit topic "test-rabbit-json-properties-topic" with the JSON properties
      | param | value  |
      | id    | def    |
      | name  | Golium |
    Then I wait up to "3" seconds for a rabbit message with the JSON properties
      | param | value  |
      | id    | abc    |
      | name  | Golium |
    And I wait up to "3" seconds without a rabbit message with the JSON properties
      | param | value  |
      | id    | abc    |
      | name  | Golium |
    And I wait up to "3" seconds for a rabbit message with the JSON properties
      | param | value  |
      | id    | def    |
      | name  | Golium |
    And the rabbit message body has the JSON properties
      | param | value  |
      | id    | def    |
      | name  | Golium |
    And I wait up to "3" seconds without a rabbit message with the JSON properties
      | param | value  |
      | id    | def    |
      | name  | Golium |

  @rabbit
  Scenario: Publish and subscribe a JSON message. Use standard properties
    Given I generate a UUID and store it in context "CorrelationId"
    And the rabbit endpoint "[CONF:rabbitmq]"
    And I subscribe to the rabbit topic "test-rabbit-json-properties-[CTXT:CorrelationId]"
    And I set standard rabbit properties
      | param         | value                |
      | ContentType   | application/json     |
      | CorrelationId | [CTXT:CorrelationId] |
    When I publish a message to the rabbit topic "test-rabbit-json-properties-[CTXT:CorrelationId]" with the JSON properties
      | param | value  |
      | id    | abc    |
      | name  | Golium |
    Then I wait up to "3" seconds for a rabbit message with the standard properties
      | param         | value                |
      | ContentType   | application/json     |
      | CorrelationId | [CTXT:CorrelationId] |
    And the rabbit message body has the JSON properties
      | param | value  |
      | id    | abc    |
      | name  | Golium |

  @rabbit
  Scenario: Publish, subscribe and consume three JSON messages. Use standard properties
    Given I generate a UUID and store it in context "CorrelationId"
    And the rabbit endpoint "[CONF:rabbitmq]"
    And I subscribe to the rabbit topic "test-rabbit-json-properties-[CTXT:CorrelationId]"
    And I set standard rabbit properties
      | param         | value                |
      | ContentType   | application/json     |
      | CorrelationId | [CTXT:CorrelationId] |
    When I publish a message to the rabbit topic "test-rabbit-json-properties-[CTXT:CorrelationId]" with the JSON properties
      | param | value   |
      | id0   | abc0    |
      | name0 | Golium0 |
    And I publish a message to the rabbit topic "test-rabbit-json-properties-[CTXT:CorrelationId]" with the JSON properties
      | param | value   |
      | id1   | abc1    |
      | name1 | Golium1 |
    And I publish a message to the rabbit topic "test-rabbit-json-properties-[CTXT:CorrelationId]" with the JSON properties
      | param | value   |
      | id2   | abc2    |
      | name2 | Golium2 |
    Then I wait up to "5" seconds for exactly "3" rabbit messages with the standard properties
      | param         | value                |
      | ContentType   | application/json     |
      | CorrelationId | [CTXT:CorrelationId] |
    And the body of the rabbit message in position "0" has the JSON properties
      | param | value   |
      | id0   | abc0    |
      | name0 | Golium0 |
    And the body of the rabbit message in position "1" has the JSON properties
      | param | value   |
      | id1   | abc1    |
      | name1 | Golium1 |
    And the body of the rabbit message in position "2" has the JSON properties
      | param | value   |
      | id2   | abc2    |
      | name2 | Golium2 |
    And I wait up to "3" seconds without a rabbit message with the standard properties
      | param         | value                |
      | ContentType   | application/json     |
      | CorrelationId | [CTXT:CorrelationId] |

  @rabbit
  Scenario: Publish, subscribe and consume three JSON messages with different CorrelationId. Use standard properties
    Given I generate a UUID and store it in context "CorrelationId"
    And I generate a UUID and store it in context "CorrelationId2"
    And the rabbit endpoint "[CONF:rabbitmq]"
    And I subscribe to the rabbit topic "test-rabbit-json-properties-[CTXT:CorrelationId]"
    And I set standard rabbit properties
      | param         | value                |
      | ContentType   | application/json     |
      | CorrelationId | [CTXT:CorrelationId] |
    When I publish a message to the rabbit topic "test-rabbit-json-properties-[CTXT:CorrelationId]" with the JSON properties
      | param | value   |
      | id0   | abc0    |
      | name0 | Golium0 |
    And I publish a message to the rabbit topic "test-rabbit-json-properties-[CTXT:CorrelationId]" with the JSON properties
      | param | value   |
      | id1   | abc1    |
      | name1 | Golium1 |
    And I set standard rabbit properties
      | param         | value                 |
      | ContentType   | application/json      |
      | CorrelationId | [CTXT:CorrelationId2] |
    And I publish a message to the rabbit topic "test-rabbit-json-properties-[CTXT:CorrelationId]" with the JSON properties
      | param | value   |
      | id2   | abc2    |
      | name2 | Golium2 |
    Then I wait up to "5" seconds for exactly "2" rabbit messages with the standard properties
      | param         | value                |
      | ContentType   | application/json     |
      | CorrelationId | [CTXT:CorrelationId] |
    And the body of the rabbit message in position "0" has the JSON properties
      | param | value   |
      | id0   | abc0    |
      | name0 | Golium0 |
    And the body of the rabbit message in position "1" has the JSON properties
      | param | value   |
      | id1   | abc1    |
      | name1 | Golium1 |
    And I wait up to "3" seconds without a rabbit message with the standard properties
      | param         | value                |
      | ContentType   | application/json     |
      | CorrelationId | [CTXT:CorrelationId] |
    And I wait up to "3" seconds for a rabbit message with the standard properties
      | param         | value                 |
      | ContentType   | application/json      |
      | CorrelationId | [CTXT:CorrelationId2] |
    And the rabbit message body has the JSON properties
      | param | value   |
      | id2   | abc2    |
      | name2 | Golium2 |
    And I wait up to "3" seconds without a rabbit message with the standard properties
      | param         | value                 |
      | ContentType   | application/json      |
      | CorrelationId | [CTXT:CorrelationId2] |

  @rabbit
  Scenario: Publish and subscribe a JSON message with rabbit headers and standard properties
    Given the rabbit endpoint "[CONF:rabbitmq]"
    And I subscribe to the rabbit topic "test-rabbit-headers-properties-topic"
    And I set rabbit headers
      | param   | value  |
      | Header1 | value1 |
      | Header2 | value2 |
    When I publish a message to the rabbit topic "test-rabbit-headers-properties-topic" with the JSON properties
      | param | value  |
      | id    | abc    |
      | name  | Golium |
    Then I wait up to "3" seconds for a rabbit message with the JSON properties
      | param | value  |
      | id    | abc    |
      | name  | Golium |
    And the rabbit message has the rabbit headers
      | param   | value  |
      | Header1 | value1 |
      | Header2 | value2 |

  @rabbit
  Scenario: Publish and subscribe a text message with rabbit headers
    Given the rabbit endpoint "[CONF:rabbitmq]"
    And I subscribe to the rabbit topic "test-rabbit-headers-topic"
    And I set rabbit headers
      | param   | value  |
      | Header1 | value1 |
      | Header2 | value2 |
    When I publish a message to the rabbit topic "test-rabbit-headers-topic" with the text
          """
          """
    Then I wait up to "3" seconds for a rabbit message with the text
          """
          """
    And the rabbit message has the rabbit headers
      | param   | value  |
      | Header1 | value1 |
      | Header2 | value2 |

  @rabbit
  Scenario: Subscribe and waits for no message
    Given I generate a UUID and store it in context "CorrelationId"
    And the rabbit endpoint "[CONF:rabbitmq]"
    When I subscribe to the rabbit topic "test-rabbit-empty-topic"
    Then I wait up to "3" seconds without a rabbit message with the standard properties
      | param         | value                |
      | CorrelationId | [CTXT:CorrelationId] |
