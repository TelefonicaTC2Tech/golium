Feature: S3 client

  Background:
    Given I generate a UUID and store it in context "Id"
    Given I store "[NOW::2006-01-02T15:04:05Z]" in context "now"
    Given I store "golium/[CTXT:Id]_[CTXT:now].txt" in context "key"
    Given I create a new S3 session
      And I create the S3 bucket "[CONF:s3Bucket]"

  @s3
  Scenario: Upload file to S3 with content
     Given the S3 bucket "[CONF:s3Bucket]" exists
     When I create a file in S3 bucket "[CONF:s3Bucket]" with key "[CTXT:key]" and the content
      """
      Document content line 1
      Document content line 2
      """
     Then the file "[CTXT:key]" exists in S3 bucket "[CONF:s3Bucket]" with the content
      """
      Document content line 1
      Document content line 2
      """
      And I delete the file in S3 bucket "[CONF:s3Bucket]" with key "[CTXT:key]"
      And I delete the S3 bucket "[CONF:s3Bucket]"
