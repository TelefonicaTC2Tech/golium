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
	"io"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	aws_s "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type ClientFunctions interface {
	NewSession(cfgs *aws.Config) (*aws_s.Session, error)
	NewUploader(client *aws_s.Session) *s3manager.Uploader
	Upload(
		uploader *s3manager.Uploader,
		bucket, key, message string,
	) (*s3manager.UploadOutput, error)
	New(client *aws_s.Session) *s3.S3
	CreateBucket(
		s3Client *s3.S3,
		input *s3.CreateBucketInput,
	) (*s3.CreateBucketOutput, error)
	DeleteBucket(
		s3Client *s3.S3,
		input *s3.DeleteBucketInput,
	) (*s3.DeleteBucketOutput, error)
	GetBucketLocation(
		s3Client *s3.S3,
		input *s3.GetBucketLocationInput,
	) (*s3.GetBucketLocationOutput, error)
	HeadObject(
		s3Client *s3.S3,
		input *s3.HeadObjectInput,
	) (*s3.HeadObjectOutput, error)
	NewDownloader(client *aws_s.Session) *s3manager.Downloader
	Download(
		downloader *s3manager.Downloader,
		w io.WriterAt,
		input *s3.GetObjectInput,
	) (int64, error)
	DeleteObject(s3Client *s3.S3,
		input *s3.DeleteObjectInput,
	) (*s3.DeleteObjectOutput, error)
}

type ClientService struct{}

func NewS3ClientService() *ClientService {
	return &ClientService{}
}

func (c ClientService) NewSession(cfgs *aws.Config) (*aws_s.Session, error) {
	return aws_s.NewSession(cfgs)
}

func (c ClientService) NewUploader(client *aws_s.Session) *s3manager.Uploader {
	return s3manager.NewUploader(client)
}

func (c ClientService) Upload(
	uploader *s3manager.Uploader,
	bucket, key, message string,
) (*s3manager.UploadOutput, error) {
	return uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   strings.NewReader(message),
	})
}

func (c ClientService) New(client *aws_s.Session) *s3.S3 {
	return s3.New(client)
}

func (c ClientService) CreateBucket(
	s3Client *s3.S3,
	input *s3.CreateBucketInput,
) (*s3.CreateBucketOutput, error) {
	return s3Client.CreateBucket(input)
}

func (c ClientService) DeleteBucket(
	s3Client *s3.S3,
	input *s3.DeleteBucketInput,
) (*s3.DeleteBucketOutput, error) {
	return s3Client.DeleteBucket(input)
}
func (c ClientService) GetBucketLocation(
	s3Client *s3.S3,
	input *s3.GetBucketLocationInput,
) (*s3.GetBucketLocationOutput, error) {
	return s3Client.GetBucketLocation(input)
}

func (c ClientService) HeadObject(
	s3Client *s3.S3,
	input *s3.HeadObjectInput,
) (*s3.HeadObjectOutput, error) {
	return s3Client.HeadObject(input)
}

func (c ClientService) NewDownloader(client *aws_s.Session) *s3manager.Downloader {
	return s3manager.NewDownloader(client)
}

func (c ClientService) Download(
	downloader *s3manager.Downloader,
	w io.WriterAt,
	input *s3.GetObjectInput,
) (int64, error) {
	return downloader.Download(w, input)
}

func (c ClientService) DeleteObject(s3Client *s3.S3,
	input *s3.DeleteObjectInput,
) (*s3.DeleteObjectOutput, error) {
	return s3Client.DeleteObject(input)
}
