Feature: HTTP client

  @shared
  Scenario: Send a GET request
    Given save the code "3003" on parent session
    Then validate the code "3003" on parent session
