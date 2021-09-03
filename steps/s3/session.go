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

	"github.com/Telefonica/golium"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	aws_s "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type Session struct {
	Client           *aws_s.Session
	CreatedBuckets   []*CreatedBucket
	CreatedDocuments []*CreatedDocument
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

	s3Config := &aws.Config{
		Region: aws.String("eu-west-1"),
	}
	// Check minio
	minio := golium.Value(ctx, "[CONF:minio]").(bool)
	if minio {
		s3Config = &aws.Config{
			Credentials:      credentials.NewStaticCredentials(golium.Value(ctx, "[CONF:minioAwsAccessKeyId]").(string), golium.Value(ctx, "[CONF:minioAwsSecretAccessKey]").(string), ""),
			Endpoint:         aws.String(golium.Value(ctx, "[CONF:minioEndpoint]").(string)),
			Region:           aws.String(golium.Value(ctx, "[CONF:minioAwsRegion]").(string)),
			DisableSSL:       aws.Bool(true),
			S3ForcePathStyle: aws.Bool(true),
		}
	}
	var err error
	if s.Client, err = aws_s.NewSession(s3Config); err != nil {
		return fmt.Errorf("error creating s3 session. %v", err)
	}

	return nil
}

// UploadS3FileWithContent creates a new file in S3 with the content specified.
func (s *Session) UploadS3FileWithContent(ctx context.Context, bucket, key, message string) error {
	logger := GetLogger()
	logger.LogOperation("upload", bucket, key)
	uploader := s3manager.NewUploader(s.Client)
	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   strings.NewReader(message),
	})
	if err != nil {
		return fmt.Errorf("unable to upload %q to %q, %v", key, bucket, err)
	}

	s.CreatedDocuments = append(s.CreatedDocuments, &CreatedDocument{bucket: bucket, key: key})
	return nil
}

// CreateS3Bucket creates a new bucket.
func (s *Session) CreateS3Bucket(ctx context.Context, bucket string) error {
	logger := GetLogger()
	logger.LogMessage(fmt.Sprintf("creating a new bucket: %s", bucket))
	cparams := &s3.CreateBucketInput{
		Bucket: aws.String(bucket),
	}

	s3Client := s3.New(s.Client)

	if _, err := s3Client.CreateBucket(cparams); err != nil {
		return fmt.Errorf("error creating a new bucket: %s, err: %v", bucket, err)
	}

	s.CreatedBuckets = append(s.CreatedBuckets, &CreatedBucket{bucket: bucket})
	return nil
}

// DeleteS3Bucket deletes the bucket in S3.
func (s *Session) DeleteS3Bucket(ctx context.Context, bucket string) error {
	logger := GetLogger()
	logger.LogMessage(fmt.Sprintf("deleting bucket: %s", bucket))
	cparams := &s3.DeleteBucketInput{
		Bucket: aws.String(bucket),
	}

	s3Client := s3.New(s.Client)

	if _, err := s3Client.DeleteBucket(cparams); err != nil {
		return fmt.Errorf("error deleting bucket: %s, err: %v", bucket, err)
	}

	return nil
}

// ValidateS3BucketExists verifies the existence of a bucket.
func (s *Session) ValidateS3BucketExists(ctx context.Context, bucket string) error {
	logger := GetLogger()
	logger.LogMessage(fmt.Sprintf("validating the existence of bucket: %s", bucket))
	// GetBucketLocation is used to validate whether the bucket exists
	s3svc := s3.New(s.Client)
	input := &s3.GetBucketLocationInput{
		Bucket: aws.String(bucket),
	}
	if _, err := s3svc.GetBucketLocation(input); err != nil {
		return fmt.Errorf("bucket: '%s' does not exist", bucket)
	}
	return nil
}

// ValidateS3FileExists checks the existence of a file in S3.
func (s *Session) ValidateS3FileExists(ctx context.Context, bucket, key string) error {
	logger := GetLogger()
	logger.LogOperation("validate", bucket, key)
	s3svc := s3.New(s.Client)
	exists, err := s.s3KeyExists(s3svc, bucket, key)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("failed validating s3 file exists: file '%s' does not exist in bucket '%s', err %v ", key, bucket, err)
	}
	return nil
}

// ValidateS3FileWithContent checks the existence of a file in S3 with the content specified.
func (s *Session) ValidateS3FileExistsWithContent(ctx context.Context, bucket, key, message string) error {
	expected := strings.TrimSpace(message)
	logger := GetLogger()
	logger.LogOperation("validate", bucket, key)
	downloader := s3manager.NewDownloader(s.Client)
	buf := aws.NewWriteAtBuffer([]byte{})
	_, err := downloader.Download(buf, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("unable to upload %q to %q, %v", key, bucket, err)
	}
	actual := strings.TrimSpace(string(buf.Bytes()))
	if expected != actual {
		return fmt.Errorf("failed validating s3 bucket '%s' file '%s' content: expected:\n%s\n\nactual:\n%s", bucket, key, expected, actual)
	}
	return nil
}

// DeleteS3File deletes the file in S3.
func (s *Session) DeleteS3File(ctx context.Context, bucket, key string) error {
	logger := GetLogger()
	logger.LogOperation("delete", bucket, key)
	s3svc := s3.New(s.Client)
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}
	if _, err := s3svc.DeleteObject(input); err != nil {
		return fmt.Errorf("error deleting file '%s' in s3 bucket '%s', err: %v", key, bucket, err)
	}
	return nil
}

// s3KeyExists checks the existence of a key in a S3 bucket.
func (s *Session) s3KeyExists(s3svc *s3.S3, bucket string, key string) (bool, error) {
	_, err := s3svc.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	if err == nil {
		return true, nil
	}
	var aerr awserr.Error
	var ok bool
	if aerr, ok = err.(awserr.Error); !ok {
		return false, err
	}
	if aerr.Code() == "NotFound" {
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
			logger.LogMessage(fmt.Sprintf("failure on deletion of s3 file '%s' in bucket '%s', err %v", file.key, file.bucket, err))
		}
	}
	// Remove buckets
	for _, file := range s.CreatedBuckets {
		if err := s.DeleteS3Bucket(ctx, file.bucket); err != nil {
			logger.LogMessage(fmt.Sprintf("failure on deletion of s3 bucket '%s', err %v", file.bucket, err))
		}
	}
}
