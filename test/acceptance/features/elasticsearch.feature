Feature: Elasticsearch client

  @elasticsearch
  Scenario: Create and search a document
    Given the elasticsearch server
          | param     | value                          |
          | addresses | [CONF:elasticsearch.addresses] |
     When I create the elasticsearch document with index "example" and the JSON properties
          | name   | example                             |
          | number | [NUMBER:1]                          |
          | bool   | [TRUE]                              |
          | date   | 2021-03-30T14:22:03.873813206+02:00 |
      And I wait for "1" seconds
     Then I search in the elasticsearch index "example" with the JSON body
          """
          {
               "query": {
                    "term": {
                         "name": "example"
                    }
               }
          }
          """
      And the search result must have the JSON properties
          | hits.total.value           | [NUMBER:1]                          |
          | hits.hits.#                | [NUMBER:1]                          |
          | hits.hits.0._source.name   | example                             |
          | hits.hits.0._source.number | [NUMBER:1]                          |
          | hits.hits.0._source.bool   | [TRUE]                              |
          | hits.hits.0._source.date   | 2021-03-30T14:22:03.873813206+02:00 |
