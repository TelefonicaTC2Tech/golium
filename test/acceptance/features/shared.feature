Feature: HTTP client

  @shared
  Scenario: Send a GET request
    Given save the code "3003" from aggregated to shared session
    Then validate the code "3003" in shared session
