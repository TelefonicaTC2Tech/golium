Feature: Rabbit client

  @rabbit
  Scenario: Publish and subscribe a text message
    Given the rabbit endpoint "amqp://guest:guest@localhost:5672/"
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
    Given the rabbit endpoint "amqp://guest:guest@localhost:5672/"
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

  @rabbit
  Scenario: Publish and subscribe a JSON message with list inside
    Given the rabbit endpoint "amqp://guest:guest@localhost:5672/"
      And I subscribe to the rabbit topic "test-topic"
     When I publish a message to the rabbit topic "test-topic" with the JSON
      """
      {
        "key1": "value1",   
        "listA": ["element1A", "element2A", "element3A"],
        "listB": ["element1B", "element2B", "element3B"],
        "nestedElement": {
             "nestA": "value1"
        }
      }
      """
     Then I wait up to "3" seconds for a rabbit message with the JSON
      """
      {
        "key1": "value1",   
        "listA": ["element1A", "element2A", "element3A"],
        "listB": ["element1B", "element2B", "element3B"],
        "nestedElement": {
             "nestA": "value1"
        }
      }
      """
     And I unsubscribe from the rabbit topic "test-topic"

  @rabbit
  Scenario: Publish and subscribe a JSON message with list inside and check JSON properties field by field
    Given the rabbit endpoint "amqp://guest:guest@localhost:5672/"
      And I subscribe to the rabbit topic "test-topic"
     When I publish a message to the rabbit topic "test-topic" with the JSON
      """
      {
        "key1": "value1",   
        "listA": ["element1A", "element2A", "element3A"],
        "listB": ["element1B", "element2B", "element3B"],
        "nestedElement": {
             "nestA": "value1"
        }
      }
      """
     Then I wait up to "3" seconds for a rabbit message with the JSON properties
          | key1                   | value1       |
          | listA.#                | [NUMBER:3]   |
          | listA.0                |  element1A   |
          | listA.1                |  element2A   |
          | listA.2                |  element3A   |
          | nestedElement.nestA    |  value1      |
     And I unsubscribe from the rabbit topic "test-topic"
