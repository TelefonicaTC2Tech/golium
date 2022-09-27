// Copyright 2021 Telefonica Cybersecurity & Cloud Tech SL
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package s3steps

import (
	"context"
	"fmt"
	"strings"

	"github.com/TelefonicaTC2Tech/golium"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	s3manager "github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	awserr "github.com/aws/smithy-go"
)

const (
	nilSessionMessage = "nil session: may forget step 'I create a new S3 session'"
)

type Session struct {
	Client           *s3.Client
	CreatedBuckets   []*CreatedBucket
	CreatedDocuments []*CreatedDocument
	S3ServiceClient  ClientFunctions
}

type CreatedBucket struct {
	bucket string
}

type CreatedDocument struct {
	bucket string
	key    string
}

// NewS3Session initiates a new aws session.
func (s *Session) NewS3Session(ctx context.Context) error {
	logger := GetLogger()
	logger.LogMessage("Creating a new S3 session")

	s3Config := aws.Config{}

	// Check if minio and adapt s3 session properly
	if golium.GetEnvironment().Get("minio") != nil {
		if minio := golium.Value(ctx, "[CONF:minio]").(bool); minio {
			var err error
			s3Config, err = awsconfig.LoadDefaultConfig(ctx,
				awsconfig.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
					func(service, region string, options ...interface{}) (aws.Endpoint, error) {
						return aws.Endpoint{
							URL:               golium.Value(ctx, "[CONF:minioEndpoint]").(string),
							HostnameImmutable: true,
						}, nil
					})),
			)
			if err != nil {
				return fmt.Errorf("error setting aws config: %v", err)
			}
		}
	} else {
		s3Config, err = awsConfig.LoadDefaultConfig(ctx)
	}

	s.Client = s.S3ServiceClient.New(s3Config)

	return nil
}

// UploadS3FileWithContent creates a new file in S3 with the content specified.
func (s *Session) UploadS3FileWithContent(ctx context.Context, bucket, key, message string) error {
	if s.Client == nil {
		return fmt.Errorf("failed uploading S3 file: " + nilSessionMessage)
	}
	logger := GetLogger()
	logger.LogOperation("upload", bucket, key)
	uploader := s.S3ServiceClient.NewUploader(s.Client)
	_, err := s.S3ServiceClient.Upload(ctx, s.Client, uploader, bucket, key, message)
	if err != nil {
		return fmt.Errorf("unable to upload %q to %q, %v", key, bucket, err)
	}

	s.CreatedDocuments = append(s.CreatedDocuments, &CreatedDocument{bucket: bucket, key: key})
	return nil
}

// CreateS3Bucket creates a new bucket.
func (s *Session) CreateS3Bucket(ctx context.Context, bucket string) error {
	if s.Client == nil {
		return fmt.Errorf("failed creating S3 bucket: " + nilSessionMessage)
	}
	logger := GetLogger()
	logger.LogMessage(fmt.Sprintf("creating a new bucket: %s", bucket))
	cparams := &s3.CreateBucketInput{
		Bucket: aws.String(bucket),
	}

	if _, err := s.S3ServiceClient.CreateBucket(ctx, s.Client, cparams); err != nil {
		return fmt.Errorf("error creating a new bucket: %s, err: %v", bucket, err)
	}

	s.CreatedBuckets = append(s.CreatedBuckets, &CreatedBucket{bucket: bucket})
	return nil
}

// DeleteS3Bucket deletes the bucket in S3.
func (s *Session) DeleteS3Bucket(ctx context.Context, bucket string) error {
	if s.Client == nil {
		return fmt.Errorf("failed deleting S3 bucket: " + nilSessionMessage)
	}
	logger := GetLogger()
	logger.LogMessage(fmt.Sprintf("deleting bucket: %s", bucket))
	cparams := &s3.DeleteBucketInput{
		Bucket: aws.String(bucket),
	}

	if _, err := s.S3ServiceClient.DeleteBucket(ctx, s.Client, cparams); err != nil {
		return fmt.Errorf("error deleting bucket: %s, err: %v", bucket, err)
	}

	return nil
}

// ValidateS3BucketExists verifies the existence of a bucket.
func (s *Session) ValidateS3BucketExists(ctx context.Context, bucket string) error {
	if s.Client == nil {
		return fmt.Errorf("failed validating S3 bucket: " + nilSessionMessage)
	}
	logger := GetLogger()
	logger.LogMessage(fmt.Sprintf("validating the existence of bucket: %s", bucket))
	// GetBucketLocation is used to validate whether the bucket exists
	input := &s3.GetBucketLocationInput{
		Bucket: aws.String(bucket),
	}
	if _, err := s.S3ServiceClient.GetBucketLocation(ctx, s.Client, input); err != nil {
		return fmt.Errorf("bucket: '%s' does not exist", bucket)
	}
	return nil
}

// ValidateS3FileExists checks the existence of a file in S3.
func (s *Session) ValidateS3FileExists(ctx context.Context, bucket, key string) error {
	if s.Client == nil {
		return fmt.Errorf("failed validating S3 file: " + nilSessionMessage)
	}
	logger := GetLogger()
	logger.LogOperation("validate", bucket, key)
	exists, err := s.s3KeyExists(ctx, bucket, key)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf(
			"failed validating s3 file exists: file '%s' does not exist in bucket '%s', err %v ",
			key, bucket, err)
	}
	return nil
}

// ValidateS3FileWithContent checks the existence of a file in S3 with the content specified.
func (s *Session) ValidateS3FileExistsWithContent(
	ctx context.Context,
	bucket, key, message string,
) error {
	if s.Client == nil {
		return fmt.Errorf("failed validating S3 file with content: " + nilSessionMessage)
	}
	expected := strings.TrimSpace(message)
	logger := GetLogger()
	logger.LogOperation("validate", bucket, key)
	downloader := s.S3ServiceClient.NewDownloader(s.Client)
	buf := s3manager.NewWriteAtBuffer([]byte{})
	_, err := s.S3ServiceClient.Download(ctx, s.Client, downloader, buf, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("unable to upload %q to %q, %v", key, bucket, err)
	}

	actual := strings.TrimSpace(string(buf.Bytes()))
	if expected != actual {
		return fmt.Errorf(
			"failed validating s3 bucket '%s' file '%s' content: expected:\n%s\n\nactual:\n%s",
			bucket, key, expected, actual)
	}
	return nil
}

// DeleteS3File deletes the file in S3.
func (s *Session) DeleteS3File(ctx context.Context, bucket, key string) error {
	if s.Client == nil {
		return fmt.Errorf("failed deleting S3 file : " + nilSessionMessage)
	}
	logger := GetLogger()
	logger.LogOperation("delete", bucket, key)
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}
	if _, err := s.S3ServiceClient.DeleteObject(ctx, s.Client, input); err != nil {
		return fmt.Errorf("error deleting file '%s' in s3 bucket '%s', err: %v", key, bucket, err)
	}
	return nil
}

// s3KeyExists checks the existence of a key in a S3 bucket.
func (s *Session) s3KeyExists(ctx context.Context, bucket, key string) (bool, error) {
	_, err := s.S3ServiceClient.HeadObject(ctx, s.Client, &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	if err == nil {
		return true, nil
	}
	var aerr awserr.APIError
	var ok bool
	if aerr, ok = err.(awserr.APIError); !ok {
		return false, err
	}
	if aerr.ErrorCode() == "NotFound" {
		return false, nil
	}
	return false, err
}

// CleanUp cleans session by deleting all documents created in S3
func (s *Session) CleanUp(ctx context.Context) {
	logger := GetLogger()
	// Remove keys
	for _, file := range s.CreatedDocuments {
		if err := s.DeleteS3File(ctx, file.bucket, file.key); err != nil {
			logger.LogMessage(
				fmt.Sprintf(
					"failure on deletion of s3 file '%s' in bucket '%s', err %v",
					file.key, file.bucket, err))
		}
	}
	// Remove buckets
	for _, file := range s.CreatedBuckets {
		if err := s.DeleteS3Bucket(ctx, file.bucket); err != nil {
			logger.LogMessage(
				fmt.Sprintf("failure on deletion of s3 bucket '%s', err %v",
					file.bucket, err,
				),
			)
		}
	}
}
