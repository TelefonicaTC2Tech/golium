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

	"github.com/TelefonicaTC2Tech/golium"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/tidwall/sjson"
)

// Session contains the information of a elasticsearch session.
type Session struct {
	Client           *elasticsearch.Client
	SearchResult     golium.Map
	Correlator       string
	indexedDocuments []*IndexedDocument
	ESServiceClient  ClientFunctions
}

// ConfigureClient creates a elasticsearch connection based on the URI.
func (s *Session) ConfigureClient(ctx context.Context, config elasticsearch.Config) error {
	var err error
	if s.Client, err = s.ESServiceClient.NewClient(config); err != nil {
		return errors.Wrap(err, "failed configuring elasticsearch client")
	}
	s.Correlator = uuid.NewString()
	return nil
}

// NewDocument creates a new document, with the given JSON and for the given index.
func (s *Session) NewDocument(
	ctx context.Context,
	index string,
	props map[string]interface{},
) error {
	logger := GetLogger()
	var err error
	data := ""
	for key, value := range props {
		data, err = sjson.Set(data, key, value)
		if err != nil {
			return errors.Wrapf(
				err, "failed setting property '%s' with value '%s' in the request body",
				key, value)
		}
	}
	res, err := s.ESServiceClient.Index(ctx, s.Client, index, data)
	if err != nil {
		return errors.Wrapf(err, "failed creating index '%s' with body '%s", index, data)
	}
	defer s.ESServiceClient.ResBodyClose(res)
	if s.ESServiceClient.ResIsError(res) {
		return s.parseErrorResponse(res)
	}
	if document, err := s.getResponseIndexedDocument(res); err == nil {
		s.indexedDocuments = append(s.indexedDocuments, document)
	} else {
		logger.LogError(errors.Wrap(err, "failed storing indexed document"), s.Correlator)
	}
	logger.LogCreateIndex(res, data, index, s.Correlator)
	return nil
}

// SearchDocument
// searches in elasticsearch with given index and JSON body
// and saves the result in the application context.
func (s *Session) SearchDocument(
	ctx context.Context,
	index string,
	body string,
) error {
	logger := GetLogger()
	res, err := s.ESServiceClient.Search(ctx, s.Client, index, body)
	if err != nil {
		return errors.Wrapf(
			err, "failed searching index '%s' with body '%s'", index, body)
	}
	defer s.ESServiceClient.ResBodyClose(res)
	if s.ESServiceClient.ResIsError(res) {
		return errors.Wrapf(
			s.parseErrorResponse(res), "failed in searching response with index '%s' and body '%s",
			index, body)
	}
	logger.LogSearchIndex(res, body, index, s.Correlator)
	buff := new(bytes.Buffer)
	if _, err := buff.ReadFrom(res.Body); err != nil {
		return errors.Wrap(err, "failed decoding search result body")
	}
	s.SearchResult = golium.NewMapFromJSONBytes(buff.Bytes())
	return nil
}

// ValidateDocumentJSONProperties validates that the search result in the application context
// has the given properties.
func (s *Session) ValidateDocumentJSONProperties(
	ctx context.Context,
	props map[string]interface{},
) error {
	for key, expectedValue := range props {
		value := s.SearchResult.Get(key)
		if value != expectedValue {
			return fmt.Errorf(
				"mismatch of json property '%s': expected '%s', actual '%s'",
				key, expectedValue, value)
		}
	}
	return nil
}

// CleanUp cleans session by deleting all indexed documents in Elasticsearch
func (s *Session) CleanUp(ctx context.Context) {
	logger := GetLogger()
	indexer, err := s.ESServiceClient.NewBulkIndexer(s.Client)
	if err != nil {
		logger.LogError(errors.Wrap(err, "failed creating indexer to clean up indexes"), s.Correlator)
		return
	}
	for _, document := range s.indexedDocuments {
		if err := s.ESServiceClient.IndexerAdd(ctx, indexer, document, s.Correlator); err != nil {
			logger.LogError(
				errors.Wrapf(
					err, "failed adding indexed item to clean up indexes for document '%+v",
					document), s.Correlator)
			return
		}
	}
	if err := s.ESServiceClient.IndexerClose(ctx, indexer); err != nil {
		logger.LogError(errors.Wrap(err, "failed closing indexer to clean up indexes"), s.Correlator)
		return
	}
	stats := s.ESServiceClient.IndexerStats(indexer)
	if stats.NumFailed > 0 {
		logger.LogError(errors.Errorf(
			"failed cleaning up indexes: indexer stats %+v", stats), s.Correlator)
		return
	}
}

func (s *Session) parseErrorResponse(res *esapi.Response) error {
	resDAO := &Response{}
	if err := s.ESServiceClient.ResBodyDecode(res, resDAO); err != nil {
		return errors.Wrapf(err, "failed parsing the response body: %s", err)
	}
	// Print the response status and error information.
	return fmt.Errorf("error response [%s] %s: %s",
		res.Status(),
		resDAO.Error.Type,
		resDAO.Error.Reason,
	)
}

func (s *Session) getResponseIndexedDocument(res *esapi.Response) (*IndexedDocument, error) {
	resDAO := &IndexedDocument{}
	if err := s.ESServiceClient.ResBodyDecode(res, resDAO); err != nil {
		return nil, errors.Wrapf(err, "failed parsing the index response body: %s", err)
	}
	return resDAO, nil
}
