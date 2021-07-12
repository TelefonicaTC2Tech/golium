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
	"github.com/sirupsen/logrus"
)

type Session struct {
	Client           *aws_s.Session
	CreatedDocuments []*CreatedDocument
}

type CreatedDocument struct {
	bucket string
	key    string
}

// newS3Session initiates a new aws session.
func (s *Session) newS3Session(ctx context.Context) error {
	logger := GetLogger()
	logger.LogMessage("Creating a new S3 session")
	var err error

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

	s.Client, err = aws_s.NewSession(s3Config)
	if err != nil {
		return fmt.Errorf("error creating s3 session. %v", err)
	}
	return nil
}

// uploadS3FileWithContent creates a new file in S3 with the content specified.
func (s *Session) uploadS3FileWithContent(ctx context.Context, bucket, key, message string) error {
	logger := GetLogger()
	logger.LogOperation("upload", bucket, key)
	uploader := s3manager.NewUploader(s.Client)
	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   strings.NewReader(message),
	})
	if err != nil {
		logger.LogMessage("unable to upload file")
		return fmt.Errorf("unable to upload %q to %q, %v", key, bucket, err)
	}

	s.CreatedDocuments = append(s.CreatedDocuments, &CreatedDocument{bucket: bucket, key: key})
	return nil
}

// createBucket creates a new bucket.
func (s *Session) createBucket(ctx context.Context, bucket string) error {
	logrus.Debugf("Creating a new bucket: %s", bucket)
	cparams := &s3.CreateBucketInput{
		Bucket: aws.String(bucket),
	}

	s3Client := s3.New(s.Client)
	_, err := s3Client.CreateBucket(cparams)
	if err != nil {
		logrus.Errorf("Error creating a new bucket: %s, err: %v", bucket, err)
	}
	return nil
}

// validateS3File checks the existence of a file in S3.
func (s *Session) validateS3File(ctx context.Context, bucket, key string) error {
	logger := GetLogger()
	logger.LogOperation("validate", bucket, key)
	s3svc := s3.New(s.Client)
	exists, err := s.s3KeyExists(s3svc, bucket, key)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("failed validating s3 file exits: no file exists in bucket '%s' with name '%s' ", bucket, key)
	}
	return nil
}

// validateS3FileWithContent checks the existence of a file in S3 with the content specified.
func (s *Session) validateS3FileWithContent(ctx context.Context, bucket, key, message string) error {
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

// deleteS3File deletes the file in S3.
func (s *Session) deleteS3File(ctx context.Context, bucket, key string) error {
	logger := GetLogger()
	logger.LogOperation("delete", bucket, key)
	s3svc := s3.New(s.Client)
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}
	_, err := s3svc.DeleteObject(input)
	if err != nil {
		return err
	}
	return nil
}

// s3KeyExists checks the existence of a key in a S3 bucket.
func (s *Session) s3KeyExists(s3svc *s3.S3, bucket string, key string) (bool, error) {
	_, err := s3svc.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case "NotFound": // s3.ErrCodeNoSuchKey does not work, aws is missing this error code so we hardwire a string
				return false, nil
			default:
				return false, err
			}
		}
		return false, err
	}
	return true, nil
}

// CleanUp cleans session by deleting all documents created in S3
func (s *Session) CleanUp(ctx context.Context) {
	for _, file := range s.CreatedDocuments {
		err := s.deleteS3File(ctx, file.bucket, file.key)
		if err != nil {
			logrus.Errorf("Failure on deletion of s3 file '%s' in bucket '%s', err %v", file.key, file.bucket, err)
		}
	}
}
