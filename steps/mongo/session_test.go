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

package mongo

import (
	"context"
	"fmt"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
)

func TestPing(t *testing.T) {
	tests := []struct {
		name    string
		pingErr error
		wantErr bool
	}{
		{
			name:    "Ping error",
			pingErr: fmt.Errorf("ping error"),
			wantErr: true,
		},
		{
			name:    "Ping done",
			pingErr: nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{}
			s.MongoClientService = ClientServiceFuncMock{}
			PingError = tt.pingErr
			if err := s.MongoClientService.Ping(context.Background(), nil, s.client); (err != nil) != tt.wantErr {
				t.Errorf("Session.Ping() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDatabase(t *testing.T) {
	tests := []struct {
		name     string
		database *mongo.Database
		wantErr  bool
	}{
		{
			name:     "Database error",
			database: nil,
			wantErr:  true,
		},
		{
			name:     "Database done",
			database: &mongo.Database{},
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{}
			s.MongoClientService = ClientServiceFuncMock{}
			DatabaseDatabase = tt.database
			if d := s.MongoClientService.Database("test", s.client); (d != nil) == tt.wantErr {
				t.Errorf("Session.Database() error, wantErr %v", tt.wantErr)
			}
		})
	}
}

func TestDisconnect(t *testing.T) {
	tests := []struct {
		name          string
		disconnectErr error
		wantErr       bool
	}{
		{
			name:          "Disconnect error",
			disconnectErr: fmt.Errorf("Disconnect error"),
			wantErr:       true,
		},
		{
			name:          "Disconnect done",
			disconnectErr: nil,
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{}
			s.MongoClientService = ClientServiceFuncMock{}
			DisconnectError = tt.disconnectErr
			if err := s.MongoClientService.Disconnect(context.Background(), s.client); (err != nil) != tt.wantErr {
				t.Errorf("Session.Disconnect() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCollectionCollection(t *testing.T) {
	tests := []struct {
		name       string
		collection *mongo.Collection
		wantErr    bool
	}{
		{
			name:       "Collection error",
			collection: nil,
			wantErr:    true,
		},
		{
			name:       "Collection done",
			collection: &mongo.Collection{},
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{}
			s.MongoCollectionService = CollectionServiceFuncMock{}
			CollectionCollection = tt.collection
			if d := s.MongoCollectionService.Collection("test", &mongo.Database{}); (d != nil) == tt.wantErr {
				t.Errorf("Session.Collection() error, wantErr %v", tt.wantErr)
			}
		})
	}
}

func TestCollectionName(t *testing.T) {
	tests := []struct {
		name       string
		nameString string
		wantErr    bool
	}{
		{
			name:       "Collection name error",
			nameString: "",
			wantErr:    true,
		},
		{
			name:       "Collection name done",
			nameString: "test-collection",
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{}
			s.collection = &mongo.Collection{}
			s.MongoCollectionService = CollectionServiceFuncMock{}
			NameString = tt.nameString
			if n := s.MongoCollectionService.Name(s.collection); (len(n) == 0) != tt.wantErr {
				t.Errorf("Session.Name() error, wantErr %v", tt.wantErr)
			}
		})
	}
}

func TestCollectionFind(t *testing.T) {
	tests := []struct {
		name       string
		findCursor *mongo.Cursor
		findError  error
		wantErr    bool
	}{
		{
			name:       "Collection find error",
			findCursor: nil,
			findError:  fmt.Errorf("find error"),
			wantErr:    true,
		},
		{
			name:       "Collection find done",
			findCursor: &mongo.Cursor{},
			findError:  nil,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{}
			s.collection = &mongo.Collection{}
			s.MongoCollectionService = CollectionServiceFuncMock{}
			FindCursor = tt.findCursor
			FindError = tt.findError
			_, err := s.MongoCollectionService.Find(context.Background(), nil, s.collection)
			if (err != nil) != tt.wantErr {
				t.Errorf("Session.Find() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCollectionFindOne(t *testing.T) {
	tests := []struct {
		name                string
		findOneSingleResult *mongo.SingleResult
		wantErr             bool
	}{
		{
			name:                "Collection find one error",
			findOneSingleResult: nil,
			wantErr:             true,
		},
		{
			name:                "Collection find one done",
			findOneSingleResult: &mongo.SingleResult{},
			wantErr:             false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{}
			s.collection = &mongo.Collection{}
			s.MongoCollectionService = CollectionServiceFuncMock{}
			FindOneSingleResult = tt.findOneSingleResult
			if f := s.MongoCollectionService.FindOne(context.Background(), nil, s.collection); (f != nil) == tt.wantErr {
				t.Errorf("Session.FindOne() error, wantErr %v", tt.wantErr)
			}
		})
	}
}

func TestCollectionInsertMany(t *testing.T) {
	tests := []struct {
		name             string
		insertManyResult *mongo.InsertManyResult
		insertManyError  error
		wantErr          bool
	}{
		{
			name:             "Collection insert many error",
			insertManyResult: nil,
			insertManyError:  fmt.Errorf("insert many error"),
			wantErr:          true,
		},
		{
			name:             "Collection insert many done",
			insertManyResult: &mongo.InsertManyResult{},
			insertManyError:  nil,
			wantErr:          false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{}
			s.collection = &mongo.Collection{}
			s.MongoCollectionService = CollectionServiceFuncMock{}
			InsertManyInsertManyResult = tt.insertManyResult
			InsertManyError = tt.insertManyError
			_, err := s.MongoCollectionService.InsertMany(context.Background(), nil, s.collection)
			if (err != nil) != tt.wantErr {
				t.Errorf("Session.InsertMany() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCollectionDeletetMany(t *testing.T) {
	tests := []struct {
		name              string
		deleteResult      *mongo.DeleteResult
		deleteResultError error
		wantErr           bool
	}{
		{
			name:              "Collection delete many error",
			deleteResult:      nil,
			deleteResultError: fmt.Errorf("delete many error"),
			wantErr:           true,
		},
		{
			name:              "Collection delete many done",
			deleteResult:      &mongo.DeleteResult{},
			deleteResultError: nil,
			wantErr:           false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{}
			s.collection = &mongo.Collection{}
			s.MongoCollectionService = CollectionServiceFuncMock{}
			DeleteResultDeleteResult = tt.deleteResult
			DeleteResultError = tt.deleteResultError
			_, err := s.MongoCollectionService.DeleteMany(context.Background(), nil, nil, s.collection)
			if (err != nil) != tt.wantErr {
				t.Errorf("Session.DeleteMany() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
