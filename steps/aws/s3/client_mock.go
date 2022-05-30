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

	"github.com/aws/aws-sdk-go/aws"
	aws_s "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
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

func (c ClientServiceFuncMock) NewSession(cfgs *aws.Config) (*aws_s.Session, error) {
	return nil, NewSessionError
}

func (c ClientServiceFuncMock) NewUploader(client *aws_s.Session) *s3manager.Uploader {
	return nil
}

func (c ClientServiceFuncMock) Upload(
	uploader *s3manager.Uploader,
	bucket, key, message string,
) (*s3manager.UploadOutput, error) {
	return nil, UploadError
}

func (c ClientServiceFuncMock) New(client *aws_s.Session) *s3.S3 {
	return nil
}

func (c ClientServiceFuncMock) CreateBucket(
	s3Client *s3.S3,
	input *s3.CreateBucketInput,
) (*s3.CreateBucketOutput, error) {
	return nil, CreateBucketError
}

func (c ClientServiceFuncMock) DeleteBucket(
	s3Client *s3.S3,
	input *s3.DeleteBucketInput,
) (*s3.DeleteBucketOutput, error) {
	return nil, DeleteBucketError
}

func (c ClientServiceFuncMock) GetBucketLocation(
	s3Client *s3.S3,
	input *s3.GetBucketLocationInput,
) (*s3.GetBucketLocationOutput, error) {
	return nil, GetBucketLocationError
}
func (c ClientServiceFuncMock) HeadObject(
	s3Client *s3.S3,
	input *s3.HeadObjectInput,
) (*s3.HeadObjectOutput, error) {
	return nil, HeadObjectError
}

func (c ClientServiceFuncMock) NewDownloader(client *aws_s.Session) *s3manager.Downloader {
	return nil
}

func (c ClientServiceFuncMock) Download(
	downloader *s3manager.Downloader,
	w io.WriterAt,
	input *s3.GetObjectInput,
) (int64, error) {
	return 0, DownloadError
}

func (c ClientServiceFuncMock) DeleteObject(s3Client *s3.S3,
	input *s3.DeleteObjectInput,
) (*s3.DeleteObjectOutput, error) {
	return nil, DeleteObjectError
}
