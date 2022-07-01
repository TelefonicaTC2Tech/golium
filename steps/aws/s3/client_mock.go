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

	"github.com/aws/aws-sdk-go-v2/aws"
	s3manager "github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	NewSessionError        error
	UploadError            error
	CreateBucketError      error
	DeleteBucketError      error
	GetBucketLocationError error
	HeadObjectError        error
	DownloadError          error
	DeleteObjectError      error
)

type ClientServiceFuncMock struct{}

func (c ClientServiceFuncMock) NewSession(cfgs aws.Config) *s3.Client {
	return nil
}

func (c ClientServiceFuncMock) NewUploader(s3Client *s3.Client) *s3manager.Uploader {
	return nil
}

func (c ClientServiceFuncMock) Upload(
	ctx context.Context,
	s3Client *s3.Client,
	uploader *s3manager.Uploader,
	bucket, key, message string,
) (*s3manager.UploadOutput, error) {
	return nil, UploadError
}

func (c ClientServiceFuncMock) CreateBucket(ctx context.Context,
	s3Client *s3.Client,
	input *s3.CreateBucketInput,
) (*s3.CreateBucketOutput, error) {
	return nil, CreateBucketError
}

func (c ClientServiceFuncMock) DeleteBucket(
	ctx context.Context,
	s3Client *s3.Client,
	input *s3.DeleteBucketInput,
) (*s3.DeleteBucketOutput, error) {
	return nil, DeleteBucketError
}

func (c ClientServiceFuncMock) GetBucketLocation(
	ctx context.Context,
	s3Client *s3.Client,
	input *s3.GetBucketLocationInput,
) (*s3.GetBucketLocationOutput, error) {
	return nil, GetBucketLocationError
}
func (c ClientServiceFuncMock) HeadObject(
	ctx context.Context,
	s3Client *s3.Client,
	input *s3.HeadObjectInput,
) (*s3.HeadObjectOutput, error) {
	return nil, HeadObjectError
}

func (c ClientServiceFuncMock) NewDownloader(s3Client *s3.Client) *s3manager.Downloader {
	return nil
}

func (c ClientServiceFuncMock) Download(
	ctx context.Context,
	s3Client *s3.Client,
	downloader *s3manager.Downloader,
	w io.WriterAt,
	input *s3.GetObjectInput,
) (int64, error) {
	return 0, DownloadError
}

func (c ClientServiceFuncMock) DeleteObject(
	ctx context.Context,
	s3Client *s3.Client,
	input *s3.DeleteObjectInput,
) (*s3.DeleteObjectOutput, error) {
	return nil, DeleteObjectError
}
