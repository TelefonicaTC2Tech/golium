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

package elasticsearch

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/TelefonicaTC2Tech/golium"
	"github.com/elastic/go-elasticsearch/v7"
)

const (
	logsPath   = "./logs"
	correlator = "correlator"
)

func TestConfigureClient(t *testing.T) {
	tests := []struct {
		name         string
		newClientErr error
		config       elasticsearch.Config
		wantErr      bool
	}{
		{
			name:         "New client error",
			newClientErr: fmt.Errorf("new client error"),
			wantErr:      true,
		},
		{
			name:         "New client without errors",
			newClientErr: nil,
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{}
			s.ESServiceClient = ClientServiceFuncMock{}
			NewClientError = tt.newClientErr
			if err := s.ConfigureClient(context.Background(), tt.config); (err != nil) != tt.wantErr {
				t.Errorf("Session.ConfigureClient() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewDocument(t *testing.T) {
	os.MkdirAll(logsPath, os.ModePerm)
	defer os.RemoveAll(logsPath)
	testProps := make(map[string]interface{})
	testProps["name"] = "testExample"
	type args struct {
		index string
		props map[string]interface{}
	}
	tests := []struct {
		name                    string
		args                    args
		indexErr                error
		resIsErr                bool
		wantErr                 bool
		resBodyDecParseErr      error
		resBodyDecGetIndexedErr error
	}{
		{
			name: "Index error",
			args: args{
				index: "testIndex",
				props: testProps,
			},
			indexErr: fmt.Errorf("index error"),
			wantErr:  true,
		},
		{
			name: "Response error",
			args: args{
				index: "testIndex",
				props: testProps,
			},
			indexErr:           nil,
			resIsErr:           true,
			resBodyDecParseErr: fmt.Errorf("response body decode error"),
			wantErr:            true,
		},
		{
			name: "Append indexed document",
			args: args{
				index: "testIndex",
				props: testProps,
			},
			indexErr:           nil,
			resIsErr:           false,
			resBodyDecParseErr: nil,
			wantErr:            false,
		},
		{
			name: "Append indexed document error",
			args: args{
				index: "testIndex",
				props: testProps,
			},
			indexErr:                nil,
			resIsErr:                false,
			resBodyDecParseErr:      nil,
			resBodyDecGetIndexedErr: fmt.Errorf("document indexed error"),
			wantErr:                 false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{}
			s.ESServiceClient = ClientServiceFuncMock{}
			IndexError = tt.indexErr
			ResIsError = tt.resIsErr
			ResBodyDecodeParseError = tt.resBodyDecParseErr
			ResBodyDecodeGetIndexedError = tt.resBodyDecGetIndexedErr
			if err := s.NewDocument(
				context.Background(), tt.args.index, tt.args.props); (err != nil) != tt.wantErr {
				t.Errorf("Session.NewDocument() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSearchDocument(t *testing.T) {
	os.MkdirAll(logsPath, os.ModePerm)
	defer os.RemoveAll(logsPath)
	type args struct {
		index string
		body  string
	}
	tests := []struct {
		name      string
		args      args
		searchErr error
		resIsErr  bool
		wantErr   bool
	}{
		{
			name:      "Search error",
			searchErr: fmt.Errorf("search error"),
			wantErr:   true,
		},
		{
			name:      "Response error",
			searchErr: nil,
			resIsErr:  true,
			wantErr:   true,
		},
		{
			name:      "Read without error",
			searchErr: nil,
			resIsErr:  false,
			args: args{
				body:  "body",
				index: "index",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{}
			s.Correlator = correlator
			s.ESServiceClient = ClientServiceFuncMock{}
			SearchError = tt.searchErr
			ResIsError = tt.resIsErr
			if err := s.SearchDocument(
				context.Background(), tt.args.index, tt.args.body); (err != nil) != tt.wantErr {
				t.Errorf("Session.SearchDocument() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateDocumentJSONProperties(t *testing.T) {
	testProps := make(map[string]interface{})

	tests := []struct {
		name      string
		testKey   string
		testValue string
		wantErr   bool
	}{
		{
			name:      "Match value",
			testKey:   "key",
			testValue: "value",
			wantErr:   false,
		},
		{
			name:      "Mismatch value",
			testKey:   "key",
			testValue: "value2",
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{}
			testProps[tt.testKey] = tt.testValue
			s.SearchResult = golium.NewMapFromJSONBytes([]byte(`{"key":"value"}`))
			if err := s.ValidateDocumentJSONProperties(
				context.Background(), testProps); (err != nil) != tt.wantErr {
				t.Errorf("Session.ValidateDocumentJSONProperties() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCleanUp(t *testing.T) {
	os.MkdirAll(logsPath, os.ModePerm)
	defer os.RemoveAll(logsPath)
	tests := []struct {
		name            string
		bulkIndexErr    error
		indexerAddErr   error
		indexerCloseErr error
		statsNumFailed  uint64
	}{
		{
			name:         "New BulkIndexer error",
			bulkIndexErr: fmt.Errorf("bulk indexer error"),
		},
		{
			name:          "Indexer Add error",
			bulkIndexErr:  nil,
			indexerAddErr: fmt.Errorf("indexer add error"),
		},
		{
			name:            "Indexer Close error",
			bulkIndexErr:    nil,
			indexerAddErr:   nil,
			indexerCloseErr: fmt.Errorf("indexer close error"),
		},
		{
			name:            "Stats with failed",
			bulkIndexErr:    nil,
			indexerAddErr:   nil,
			indexerCloseErr: nil,
			statsNumFailed:  1,
		},
		{
			name:            "Stats without failed",
			bulkIndexErr:    nil,
			indexerAddErr:   nil,
			indexerCloseErr: nil,
			statsNumFailed:  0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{}
			s.ESServiceClient = ClientServiceFuncMock{}
			s.Correlator = correlator
			NewBulkIndexerError = tt.bulkIndexErr
			IndexerAddError = tt.indexerAddErr
			IndexerCloseError = tt.indexerCloseErr
			IndexerStatsNumFailed = tt.statsNumFailed
			s.indexedDocuments = []*IndexedDocument{
				{
					Index: "index",
					ID:    "1",
				},
			}
			s.CleanUp(context.Background())
		})
	}
}
