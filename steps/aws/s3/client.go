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
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	s3manager "github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type ClientFunctions interface {
	New(cfgs aws.Config) *s3.Client
	NewUploader(s3Client *s3.Client) *s3manager.Uploader
	Upload(
		ctx context.Context,
		s3Client *s3.Client,
		uploader *s3manager.Uploader,
		bucket, key, message string,
	) (*s3manager.UploadOutput, error)
	CreateBucket(ctx context.Context,
		s3Client *s3.Client,
		input *s3.CreateBucketInput,
	) (*s3.CreateBucketOutput, error)
	DeleteBucket(
		ctx context.Context,
		s3Client *s3.Client,
		input *s3.DeleteBucketInput,
	) (*s3.DeleteBucketOutput, error)
	GetBucketLocation(
		ctx context.Context,
		s3Client *s3.Client,
		input *s3.GetBucketLocationInput,
	) (*s3.GetBucketLocationOutput, error)
	HeadObject(
		ctx context.Context,
		s3Client *s3.Client,
		input *s3.HeadObjectInput,
	) (*s3.HeadObjectOutput, error)
	NewDownloader(s3Client *s3.Client) *s3manager.Downloader
	Download(
		ctx context.Context,
		s3Client *s3.Client,
		downloader *s3manager.Downloader,
		w io.WriterAt,
		input *s3.GetObjectInput,
	) (int64, error)
	DeleteObject(
		ctx context.Context,
		s3Client *s3.Client,
		input *s3.DeleteObjectInput,
	) (*s3.DeleteObjectOutput, error)
}

type ClientService struct{}

func NewS3ClientService() *ClientService {
	return &ClientService{}
}

func (c ClientService) New(cfgs aws.Config) *s3.Client {
	return s3.NewFromConfig(cfgs)
}
func (c ClientService) NewUploader(s3Client *s3.Client) *s3manager.Uploader {
	return s3manager.NewUploader(s3Client)
}

func (c ClientService) Upload(
	ctx context.Context,
	s3Client *s3.Client,
	uploader *s3manager.Uploader,
	bucket, key, message string,
) (*s3manager.UploadOutput, error) {
	return uploader.Upload(
		ctx,
		&s3.PutObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
			Body:   strings.NewReader(message),
		})
}

func (c ClientService) CreateBucket(
	ctx context.Context,
	s3Client *s3.Client,
	input *s3.CreateBucketInput,
) (*s3.CreateBucketOutput, error) {
	return s3Client.CreateBucket(ctx, input)
}

func (c ClientService) DeleteBucket(
	ctx context.Context,
	s3Client *s3.Client,
	input *s3.DeleteBucketInput,
) (*s3.DeleteBucketOutput, error) {
	return s3Client.DeleteBucket(ctx, input)
}

func (c ClientService) GetBucketLocation(
	ctx context.Context,
	s3Client *s3.Client,
	input *s3.GetBucketLocationInput,
) (*s3.GetBucketLocationOutput, error) {
	return s3Client.GetBucketLocation(ctx, input)
}

func (c ClientService) HeadObject(
	ctx context.Context,
	s3Client *s3.Client,
	input *s3.HeadObjectInput,
) (*s3.HeadObjectOutput, error) {
	return s3Client.HeadObject(ctx, input)
}
func (c ClientService) NewDownloader(s3Client *s3.Client) *s3manager.Downloader {
	return s3manager.NewDownloader(s3Client)
}

func (c ClientService) Download(
	ctx context.Context,
	s3Client *s3.Client,
	downloader *s3manager.Downloader,
	w io.WriterAt,
	input *s3.GetObjectInput,
) (int64, error) {
	return downloader.Download(ctx, w, input)
}

func (c ClientService) DeleteObject(
	ctx context.Context,
	s3Client *s3.Client,
	input *s3.DeleteObjectInput,
) (*s3.DeleteObjectOutput, error) {
	return s3Client.DeleteObject(ctx, input)
}
