Feature: S3 client

  Background:
    Given I generate a UUID and store it in context "Id"
    Given I store "[NOW::2006-01-02T15:04:05Z]" in context "now"
    Given I store "golium/[CTXT:Id]_[CTXT:now].txt" in context "key"
    Given I create a new S3 session
      And I create the S3 bucket "[CONF:bucket]"

  @s3
  Scenario: Upload file to S3 with content
     When I create a file in S3 bucket "[CONF:bucket]" with key "[CTXT:key]" and the content
      """
      0
      """
     Then the file "[CTXT:key]" exists in S3 bucket "[CONF:bucket]" with the content
      """
      0
      """
      And I delete the file in S3 bucket "[CONF:bucket]" with key "[CTXT:key]"
