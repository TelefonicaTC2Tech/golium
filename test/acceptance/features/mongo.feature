Feature: MongoDB client

  @mongodb
  Scenario: Create a collection, check it and delete it.
    Given I connect to MongoDB
    When I create "2" documents in the MongoDB "test-colection" collection
    Then I check that the number of documents in collection "test-colection" is "2"
    When I delete documents from the MongoDB "test-colection" collection whose "fieldString" field is "Example field string 1" value
    Then I check that the number of documents in collection "test-colection" is "1"
    When I delete documents from the MongoDB "test-colection" collection whose "fieldString" field is "Example field string 2" value
    Then I check that the number of documents in collection "test-colection" is "0"
