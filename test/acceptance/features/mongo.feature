# Copyright (c) Telefónica Cybersecurity & Cloud Tech S.L.
# SDET Team <sdetteam_tcct@telefonica.com>

Feature: MongoDB client
      Examples of MongoDB access and data processing. It is possible to:
  - Creation of documents.
  - Deletion of documents.
  - Checking the number of documents.
  - Verification of the existence or not of the _id in the documents.
  - Checking the existence or not of certain data in documents (texts, numbers, booleans, lists, objects).
  - Checking the existence or not of fields in the collection.
  - Checking null, empty and non-existent values.


  Background:
    # Prepare data: Collection "example" empty and create "_id" field
    Given I connect to MongoDB
      | field      | value                  |
      | User       | [CONF:mongoUsername]   |
      | Password   | [CONF:mongoPassword]   |
      | Host       | [CONF:mongoHost]       |
      | AuthSource | [CONF:mongoAuthSource] |
      | Database   | [CONF:mongoDatabase]   |
    And I delete all documents from the MongoDB "example" collection
    And I check that the number of documents in collection "example" is "0"
    And I generate a UUID and store it


  @mongodb
  Scenario: Create a collection, check it and delete it.
    Given I connect to MongoDB
    When I create "2" documents in the MongoDB "test-colection" collection
    Then I check that the number of documents in collection "test-colection" is "2"
    When I delete documents from the MongoDB "test-colection" collection whose "fieldString" field is "Example field string 1" value
    Then I check that the number of documents in collection "test-colection" is "1"
    When I delete documents from the MongoDB "test-colection" collection whose "fieldString" field is "Example field string 2" value
    Then I check that the number of documents in collection "test-colection" is "0"


  @mongodb
  Scenario: Creating documents in the "Example" Collection
    Given I create "2" documents in the MongoDB "example" collection
    Then I check that the number of documents in collection "example" is "2"


  @mongodb
  Scenario: Delete documents from the "example" collection
    Given I create "10" documents in the MongoDB "example" collection
    When I check that the number of documents in collection "example" is "10"
    # Delete a document by searching for a string
    When I delete documents from the MongoDB "example" collection whose "fieldString" field is "Example field string 1" value
    Then I check that the number of documents in collection "example" is "9"
    # Delete a document by searching for an int
    When I delete documents from the MongoDB "example" collection whose "fieldInt" field is "2" value
    Then I check that the number of documents in collection "example" is "8"
    # Delete a document by searching for an item in a slice
    When I delete documents from the MongoDB "example" collection whose "fieldSlice.0" field is "itemSlice_3" value
    Then I check that the number of documents in collection "example" is "7"
    # Delete a document by searching for a string element in a map
    When I delete documents from the MongoDB "example" collection whose "fieldMap.fieldString" field is "Example field in map string 4" value
    Then I check that the number of documents in collection "example" is "6"
    # Delete a document by searching for an int element in a map
    When I delete documents from the MongoDB "example" collection whose "fieldMap.fieldInt" field is "50" value
    Then I check that the number of documents in collection "example" is "5"
    # Delete a document by searching for a string element on a map within another map
    When I delete documents from the MongoDB "example" collection whose "fieldMap.fieldMap2.fieldString" field is "Example field in map map string 6" value
    Then I check that the number of documents in collection "example" is "4"
    # Delete a document by searching for an int element on a map within another map
    When I delete documents from the MongoDB "example" collection whose "fieldMap.fieldMap2.fieldInt" field is "700" value
    Then I check that the number of documents in collection "example" is "3"
    # Delete a document by searching for an empty item on a map within another map
    When I delete documents from the MongoDB "example" collection whose "fieldMap.fieldMap2.fieldEmptyText" field is "" value
    Then I check that the number of documents in collection "example" is "0"


  @mongodb
  Scenario: Check the number of documents in the "example" collection
    Given I create "5" documents in the MongoDB "example" collection
    Then I check that the number of documents in collection "example" is "5"
    When I delete documents from the MongoDB "example" collection whose "fieldInt" field is "1" value
    Then I check that the number of documents in collection "example" is "4"


  @mongodb
  Scenario: Check that the _id "[CTXT:_ID]_1" exists in the "example" collection
    Given I create "1" documents in the MongoDB "example" collection
    Then  I check that in the MongoDB "example" collection, "_id" field "does" exist for the "[CTXT:_ID]_1" _id


  @mongodb
  Scenario: Check that the "[CTXT:_ID]_0" _id does not exist in the "example" collection
    Given I create "1" documents in the MongoDB "example" collection
    When I check that the number of documents in collection "example" is "1"
    Then I check that in the MongoDB "example" collection, "_id" field "does not" exist for the "[CTXT:_ID]_0" _id


  @mongodb
  Scenario: Check that there is some data in the "example" collection through a table in the feature
    Given I create "1" documents in the MongoDB "example" collection
    Then I check that these values of the MongoDB "example" collection with "[CTXT:_ID]_1" _id "do" exist
      | field | value |
      # General values
      | _id         | [CTXT:_ID]_1           |
      | fieldString | Example field string 1 |
      | fieldInt    | [NUMBER:1]             |
      | fieldFloat  | [NUMBER:3.14]          |
      | fieldBool   | [TRUE]                 |
      # Slice with data
      | fieldSlice.# | [NUMBER:3]  |
      | fieldSlice.0 | itemSlice_1 |
      | fieldSlice.1 | itemSlice20 |
      | fieldSlice.2 | itemSlice30 |
      | fieldEmpty   | [EMPTY]     |
      # Map with data
      | fieldMap.fieldString | Example field in map string 1 |
      | fieldMap.fieldInt    | [NUMBER:10]                   |
      | fieldMap.fieldFloat  | [NUMBER:1974.1976]            |
      | fieldMap.fieldBool   | [FALSE]                       |
      # Slice empty in map. It is possible to compare empty items
      | fieldMap.fieldSliceEmpty.# | [NUMBER:0] |
      | fieldMap.fieldSliceEmpty.0 | [NULL]     |
      | fieldMap.fieldSliceEmpty.1 |            |
      | fieldMap.fieldSliceEmpty.2 | [EMPTY]    |
      # Map in map with similar data
      | fieldMap.fieldMap2.fieldString    | Example field in map map string 1 |
      | fieldMap.fieldMap2.fieldInt       | [NUMBER:100]                      |
      | fieldMap.fieldMap2.fieldFloat     | [NUMBER:1974.1976]                |
      | fieldMap.fieldMap2.fieldBool      | [FALSE]                           |
      | fieldMap.fieldMap2.fieldEmpty     | [NULL]                            |
      | fieldMap.fieldMap2.fieldEmptyText | [EMPTY]                           |


  @mongodb
  Scenario: Check that there is no data in the "example" collection through a table in the feature
    Given I create "1" documents in the MongoDB "example" collection
    And I check that the number of documents in collection "example" is "1"
    Then I check that these values of the MongoDB "example" collection with "[CTXT:_ID]_1" _id "do not" exist
      | field       | value                |
      | fieldString | Value does not exist |


  @mongodb
  Scenario: Check that the "fieldString" field exists in the "example" collection
    Given I create "1" documents in the MongoDB "example" collection
    Then I check that in the MongoDB "example" collection, "fieldString" field "does" exist for the "[CTXT:_ID]_1" _id


  @mongodb
  Scenario: Check that the "campoInexistente" field does not exist in the "example" collection
    Given I create "1" documents in the MongoDB "example" collection
    And I check that the number of documents in collection "example" is "1"
    Then I check that in the MongoDB "example" collection, "campoInexistente" field "does not" exist for the "[CTXT:_ID]_1" _id
    And I check that in the MongoDB "example" collection, "campoInexistente" field does not exist or is empty for the "[CTXT:_ID]_1" _id


  @mongodb
  Scenario: Check that the value of the "fiedlEmpty" field in the "example" collection does not exist or is null
    Given I create "1" documents in the MongoDB "example" collection
    When  I check that in the MongoDB "example" collection, "_id" field "does" exist for the "[CTXT:_ID]_1" _id
    Then I check that in the MongoDB "example" collection, "fieldEmpty" field does not exist or is empty for the "[CTXT:_ID]_1" _id
