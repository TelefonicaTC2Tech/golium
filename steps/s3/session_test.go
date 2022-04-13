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
	"os"
	"testing"

	"github.com/TelefonicaTC2Tech/golium"
	"github.com/aws/aws-sdk-go/aws/awserr"
)

const (
	logsPath        = "./logs"
	environmentPath = "./environments"
	localConfFile   = `
minio: true
minioEndpoint: http://miniomock:9000
`
	testBucket  = "test_bucket"
	testKey     = "test_key"
	testMessage = "test_message"
)

type testMockedError struct {
	uploadErr            error
	createBucketErr      error
	deleteBucketErr      error
	getBucketLocationErr error
	headObjectErr        error
	downloadErr          error
	deleteObjectErr      error
}

type testS3 struct {
	name    string
	errors  *testMockedError
	wantErr bool
}

func TestNewS3Session(t *testing.T) {
	os.MkdirAll(logsPath, os.ModePerm)
	defer os.RemoveAll(logsPath)

	os.MkdirAll(environmentPath, os.ModePerm)
	defer os.RemoveAll(environmentPath)

	os.WriteFile("./environments/local.yml", []byte(localConfFile), os.ModePerm)
	tests := []struct {
		name            string
		newS3SessionErr error
		wantErr         bool
	}{
		{
			name:            "New session error",
			newS3SessionErr: fmt.Errorf("new s3 session error"),
			wantErr:         true,
		},
		{
			name:            "New session without error",
			newS3SessionErr: nil,
			wantErr:         false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{}
			s.S3ServiceClient = ClientServiceFuncMock{}
			goliumCtx := golium.InitializeContext(context.Background())

			NewSessionError = tt.newS3SessionErr

			if err := s.NewS3Session(goliumCtx); (err != nil) != tt.wantErr {
				t.Errorf("Session.NewS3Session() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUploadS3FileWithContent(t *testing.T) {
	os.MkdirAll(logsPath, os.ModePerm)
	defer os.RemoveAll(logsPath)
	tests := fillS3Tests("Upload", &testMockedError{
		uploadErr: fmt.Errorf("upload error"),
	})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := setS3SessionMockedClient(tt.errors)
			if err := s.UploadS3FileWithContent(
				context.Background(),
				testBucket,
				testKey,
				testMessage,
			); (err != nil) != tt.wantErr {
				t.Errorf("Session.UploadS3FileWithContent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateS3Bucket(t *testing.T) {
	tests := fillS3Tests("Create Bucket", &testMockedError{
		createBucketErr: fmt.Errorf("create bucket error"),
	})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := setS3SessionMockedClient(tt.errors)
			if err := s.CreateS3Bucket(context.Background(), testBucket); (err != nil) != tt.wantErr {
				t.Errorf("Session.CreateS3Bucket() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDeleteS3Bucket(t *testing.T) {
	tests := fillS3Tests("Delete Bucket", &testMockedError{
		deleteBucketErr: fmt.Errorf("delete bucket error"),
	})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := setS3SessionMockedClient(tt.errors)
			if err := s.DeleteS3Bucket(context.Background(), testBucket); (err != nil) != tt.wantErr {
				t.Errorf("Session.DeleteS3Bucket() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateS3BucketExists(t *testing.T) {
	tests := fillS3Tests("Get Bucket", &testMockedError{
		getBucketLocationErr: fmt.Errorf("get bucket location error"),
	})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := setS3SessionMockedClient(tt.errors)
			if err := s.ValidateS3BucketExists(
				context.Background(), testBucket); (err != nil) != tt.wantErr {
				t.Errorf(
					"Session.ValidateS3BucketExists() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateS3FileExists(t *testing.T) {
	tests := []testS3{
		{
			name: "Key exists error",
			errors: &testMockedError{
				headObjectErr: fmt.Errorf("head object error"),
			},
			wantErr: true,
		},
		{
			name: "Key exists without error",
			errors: &testMockedError{
				headObjectErr: nil,
			},
			wantErr: false,
		},
		{
			name: "Key not exists",
			errors: &testMockedError{
				headObjectErr: awserr.New("NotFound", "error", fmt.Errorf("error")),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := setS3SessionMockedClient(tt.errors)
			if err := s.ValidateS3FileExists(
				context.Background(), testBucket, testKey); (err != nil) != tt.wantErr {
				t.Errorf(
					"Session.ValidateS3FileExists() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func setS3SessionMockedClient(testErrors *testMockedError) *Session {
	s := &Session{}
	s.S3ServiceClient = ClientServiceFuncMock{}

	UploadError = testErrors.uploadErr
	CreateBucketError = testErrors.createBucketErr
	DeleteBucketError = testErrors.deleteBucketErr
	GetBucketLocationError = testErrors.getBucketLocationErr
	HeadObjectError = testErrors.headObjectErr
	DownloadError = testErrors.downloadErr
	DeleteObjectError = testErrors.deleteObjectErr

	return s
}

func fillS3Tests(name string, err *testMockedError) []testS3 {
	tests := []testS3{
		{
			name:    name + " error",
			errors:  err,
			wantErr: true,
		},
		{
			name: name + " without error",
			errors: &testMockedError{
				createBucketErr:      nil,
				uploadErr:            nil,
				deleteBucketErr:      nil,
				getBucketLocationErr: nil,
				headObjectErr:        nil,
				downloadErr:          nil,
				deleteObjectErr:      nil,
			},
			wantErr: false,
		},
	}
	return tests
}

func TestValidateS3FileExistsWithContent(t *testing.T) {
	os.MkdirAll(logsPath, os.ModePerm)
	defer os.RemoveAll(logsPath)
	tests := []struct {
		name    string
		message string
		errors  *testMockedError
		wantErr bool
	}{
		{
			name: "Download error",
			errors: &testMockedError{
				downloadErr: fmt.Errorf("download error"),
			},
			wantErr: true,
		},
		{
			name: "Expected vs actual mismatch",
			errors: &testMockedError{
				downloadErr: nil,
			},
			message: testMessage,
			wantErr: true,
		},
		{
			name: "Expected equal actual",
			errors: &testMockedError{
				downloadErr: nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := setS3SessionMockedClient(tt.errors)
			if err := s.ValidateS3FileExistsWithContent(
				context.Background(),
				testBucket,
				testKey,
				tt.message,
			); (err != nil) != tt.wantErr {
				t.Errorf(
					"Session.ValidateS3FileExistsWithContent() error = %v, wantErr %v",
					err, tt.wantErr)
			}
		})
	}
}

func TestDeleteS3File(t *testing.T) {
	os.MkdirAll(logsPath, os.ModePerm)
	defer os.RemoveAll(logsPath)
	tests := fillS3Tests("Delete object", &testMockedError{
		deleteObjectErr: fmt.Errorf("delete object error"),
	})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := setS3SessionMockedClient(tt.errors)
			if err := s.DeleteS3File(
				context.Background(),
				testBucket,
				testKey,
			); (err != nil) != tt.wantErr {
				t.Errorf(
					"Session.DeleteS3File() error = %v, wantErr %v",
					err, tt.wantErr)
			}
		})
	}
}

func TestCleanUp(t *testing.T) {
	os.MkdirAll(logsPath, os.ModePerm)
	defer os.RemoveAll(logsPath)

	tests := []testS3{
		{
			name: "delete object error",
			errors: &testMockedError{
				deleteObjectErr: fmt.Errorf("delete object error"),
			},
			wantErr: true,
		},
		{
			name: "delete bucket error",
			errors: &testMockedError{
				deleteBucketErr: fmt.Errorf("delete bucket error"),
			},
			wantErr: true,
		},
		{
			name: "delete object and delete bucket without error",
			errors: &testMockedError{
				deleteBucketErr: nil,
				deleteObjectErr: nil,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := setS3SessionMockedClient(tt.errors)
			testDocument := &CreatedDocument{
				bucket: testBucket,
				key:    testKey,
			}
			testBucket := &CreatedBucket{
				bucket: testBucket,
			}
			s.CreatedDocuments = []*CreatedDocument{
				testDocument,
			}
			s.CreatedBuckets = []*CreatedBucket{
				testBucket,
			}

			s.CleanUp(context.Background())
		})
	}
}
