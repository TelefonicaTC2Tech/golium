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
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Telefonica/golium"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/tidwall/sjson"
)

// Session contains the information of a elasticsearch session.
type Session struct {
	Client       *elasticsearch.Client
	SearchResult golium.Map
	Correlator   string
}

// ConfigureClient creates a elasticsearch connection based on the URI.
func (s *Session) ConfigureConnection(ctx context.Context, config elasticsearch.Config) error {
	var err error
	if s.Client, err = elasticsearch.NewClient(config); err != nil {
		return errors.Wrap(err, "failed configuring elasticsearch client")
	}
	s.Correlator = uuid.NewString()
	return nil
}

// CreatesDocument creates a document in index with given JSON properties.
func (s *Session) CreatesDocument(ctx context.Context, index string, props map[string]interface{}) error {
	logger := GetLogger()
	var err error
	data := ""
	for key, value := range props {
		data, err = sjson.Set(data, key, value)
		if err != nil {
			return errors.Wrapf(err, "failed setting property '%s' with value '%s' in the request body", key, value)
		}
	}
	res, err := s.Client.Index(
		index,
		strings.NewReader(data),
		s.Client.Index.WithContext(ctx),
	)
	if err != nil {
		return errors.Wrapf(err, "failed creating index '%s' with body '%s", index, data)
	}
	defer res.Body.Close()
	if res.IsError() {
		return s.parseErrorResponse(res)
	}
	logger.LogCreateIndex(res, data, index, s.Correlator)
	return nil
}

// SearchDocument searchs in elasticsearch with given index and JSON body and saves the result in the application context.
func (s *Session) SearchDocument(ctx context.Context, index string, body string) error {
	logger := GetLogger()
	res, err := s.Client.Search(
		s.Client.Search.WithIndex(index),
		s.Client.Search.WithBody(strings.NewReader(body)),
	)
	if err != nil {
		return errors.Wrapf(err, "failed searching index '%s' with body '%s'", index, body)
	}
	defer res.Body.Close()
	if res.IsError() {
		return errors.Wrapf(s.parseErrorResponse(res), "failed in searching response with index '%s' and body '%s", index, body)
	}
	logger.LogSearchIndex(res, body, index, s.Correlator)
	buff := new(bytes.Buffer)
	if _, err := buff.ReadFrom(res.Body); err != nil {
		return errors.Wrap(err, "failed decoding search result body")
	}
	s.SearchResult = golium.NewMapFromJSONBytes(buff.Bytes())
	return nil
}

// ValidateDocumentJSONProperties validates that the search result in the application context has the given properties.
func (s *Session) ValidateDocumentJSONProperties(ctx context.Context, props map[string]interface{}) error {
	for key, expectedValue := range props {
		value := s.SearchResult.Get(key)
		if value != expectedValue {
			return fmt.Errorf("mismatch of json property '%s': expected '%s', actual '%s'", key, expectedValue, value)
		}
	}
	return nil
}

func (s *Session) parseErrorResponse(res *esapi.Response) error {
	resDAO := &Response{}
	if err := json.NewDecoder(res.Body).Decode(&resDAO); err != nil {
		return errors.Wrapf(err, "failed parsing the response body: %s", err)
	}
	// Print the response status and error information.
	return fmt.Errorf("error response [%s] %s: %s",
		res.Status(),
		resDAO.Error.Type,
		resDAO.Error.Reason,
	)
}
