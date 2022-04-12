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
	"encoding/json"
	"fmt"
	"strings"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/elastic/go-elasticsearch/v7/esutil"
	"github.com/pkg/errors"
)

type ClientFunctions interface {
	NewClient(
		config elasticsearch.Config,
	) (*elasticsearch.Client, error)
	Index(
		ctx context.Context,
		client *elasticsearch.Client,
		index, data string,
	) (*esapi.Response, error)
	Search(
		ctx context.Context,
		client *elasticsearch.Client,
		index string,
		body string,
	) (*esapi.Response, error)
	ResBodyClose(res *esapi.Response)
	ResIsError(res *esapi.Response) bool
	ResBodyDecode(res *esapi.Response, v interface{}) error
	NewBulkIndexer(
		client *elasticsearch.Client,
	) (esutil.BulkIndexer, error)
	IndexerAdd(ctx context.Context,
		indexer esutil.BulkIndexer,
		document *IndexedDocument,
		correlator string,
	) error
	IndexerClose(ctx context.Context, indexer esutil.BulkIndexer) error
	IndexerStats(indexer esutil.BulkIndexer) esutil.BulkIndexerStats
}

func NewElasticsearchClientService() *ClientService {
	return &ClientService{}
}

type ClientService struct{}

func (c ClientService) NewClient(
	config elasticsearch.Config,
) (*elasticsearch.Client, error) {
	return elasticsearch.NewClient(config)
}

func (c ClientService) Index(
	ctx context.Context,
	client *elasticsearch.Client,
	index, data string,
) (*esapi.Response, error) {
	return client.Index(index,
		strings.NewReader(data),
		client.Index.WithContext(ctx))
}

func (c ClientService) Search(
	ctx context.Context,
	client *elasticsearch.Client,
	index string,
	body string,
) (*esapi.Response, error) {
	return client.Search(
		client.Search.WithIndex(index),
		client.Search.WithBody(strings.NewReader(body)),
	)
}

func (c ClientService) ResBodyClose(res *esapi.Response) {
	res.Body.Close()
}
func (c ClientService) ResIsError(res *esapi.Response) bool {
	return res.IsError()
}
func (c ClientService) ResBodyDecode(res *esapi.Response, v interface{}) error {
	return json.NewDecoder(res.Body).Decode(v)
}

func (c ClientService) NewBulkIndexer(
	client *elasticsearch.Client,
) (esutil.BulkIndexer, error) {
	return esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Client: client,
	})
}
func (c ClientService) IndexerAdd(
	ctx context.Context,
	indexer esutil.BulkIndexer,
	document *IndexedDocument,
	correlator string,
) error {
	return indexer.Add(ctx,
		esutil.BulkIndexerItem{
			Action:     "delete",
			Index:      document.Index,
			DocumentID: document.ID,
			OnFailure: func(
				ctx context.Context,
				item esutil.BulkIndexerItem,
				res esutil.BulkIndexerResponseItem,
				err error,
			) {
				errMsg := fmt.Sprintf(
					"failed deleting document to clean up indexes for document %+v",
					document)
				if err != nil {
					logger.LogError(errors.Wrap(err, errMsg), correlator)
				} else {
					logger.LogError(
						errors.Wrap(
							errors.Errorf(
								"elasticsearch error. %s: %s",
								res.Error.Type, res.Error.Reason), errMsg), correlator)
				}
			},
		})
}

func (c ClientService) IndexerClose(ctx context.Context, indexer esutil.BulkIndexer) error {
	return indexer.Close(ctx)
}

func (c ClientService) IndexerStats(indexer esutil.BulkIndexer) esutil.BulkIndexerStats {
	return indexer.Stats()
}
