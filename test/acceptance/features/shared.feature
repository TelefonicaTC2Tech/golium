Feature: HTTP client

  @shared
  Scenario: Send a GET request
    Given save the code "3003" from aggregate to shared session
    Then validate the code "3003" on shared session
