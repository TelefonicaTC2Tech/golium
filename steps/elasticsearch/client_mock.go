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
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/elastic/go-elasticsearch/v7/esutil"
)

var (
	NewClientError               error
	IndexError                   error
	ResIsError                   bool
	ResBodyDecodeParseError      error
	ResBodyDecodeGetIndexedError error
	SearchError                  error
	NewBulkIndexerError          error
	IndexerAddError              error
	IndexerCloseError            error
	IndexerStatsNumFailed        uint64
)

type ClientServiceFuncMock struct{}
type nopCloser struct {
	io.Reader
}

func (c ClientServiceFuncMock) NewClient(
	config elasticsearch.Config,
) (*elasticsearch.Client, error) {
	return nil, NewClientError
}

func (c ClientServiceFuncMock) Index(
	ctx context.Context,
	client *elasticsearch.Client,
	index, data string,
) (*esapi.Response, error) {
	return nil, IndexError
}
func (c ClientServiceFuncMock) Search(
	ctx context.Context,
	client *elasticsearch.Client,
	index string,
	body string,
) (*esapi.Response, error) {
	return &esapi.Response{
		Body: nopCloser{bytes.NewBufferString("test")},
	}, SearchError
}
func (c ClientServiceFuncMock) ResBodyClose(res *esapi.Response) {
	// Just to be runned on defer without failing
}

func (c ClientServiceFuncMock) ResIsError(res *esapi.Response) bool {
	return ResIsError
}
func (c ClientServiceFuncMock) ResBodyDecode(res *esapi.Response, v interface{}) error {
	switch dec := v.(type) {
	case *IndexedDocument:
		return ResBodyDecodeGetIndexedError
	case *Response:
		return ResBodyDecodeParseError
	default:
		return fmt.Errorf("dec should be *indexedDocument or *Response: %v", dec)
	}
}

func (c ClientServiceFuncMock) NewBulkIndexer(
	client *elasticsearch.Client,
) (esutil.BulkIndexer, error) {
	return nil, NewBulkIndexerError
}

func (c ClientServiceFuncMock) IndexerAdd(
	ctx context.Context,
	indexer esutil.BulkIndexer,
	document *IndexedDocument,
	correlator string,
) error {
	return IndexerAddError
}

func (c ClientServiceFuncMock) IndexerClose(
	ctx context.Context,
	indexer esutil.BulkIndexer,
) error {
	return IndexerCloseError
}

func (c ClientServiceFuncMock) IndexerStats(indexer esutil.BulkIndexer) esutil.BulkIndexerStats {
	stats := esutil.BulkIndexerStats{}
	stats.NumFailed = IndexerStatsNumFailed
	return stats
}

func (nopCloser) Close() error {
	return nil
}
