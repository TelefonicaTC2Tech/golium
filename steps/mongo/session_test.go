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
	"reflect"
	"testing"

	"github.com/TelefonicaTC2Tech/golium"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var ISBOOLEAN = "is_boolean"

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
			if err := s.MongoClientService.Ping(context.Background(),
				nil, s.client); (err != nil) != tt.wantErr {
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
			if err := s.MongoClientService.Disconnect(context.Background(),
				s.client); (err != nil) != tt.wantErr {
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
			if d := s.MongoCollectionService.Collection("test",
				&mongo.Database{}); (d != nil) == tt.wantErr {
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
			if n := s.MongoCollectionService.Name(s.collection); (n == "") != tt.wantErr {
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
			s.MongoCollectionService = CollectionServiceFuncMock{}
			FindCursor = tt.findCursor
			FindError = tt.findError
			c, _ := s.MongoCollectionService.Find(context.Background(), nil, &mongo.Collection{})
			if (c != nil) == tt.wantErr {
				t.Errorf("Session.Find() error, wantErr %v", tt.wantErr)
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
			if f := s.MongoCollectionService.FindOne(context.Background(),
				nil, s.collection); (f != nil) == tt.wantErr {
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
			d, _ := s.MongoCollectionService.DeleteMany(context.Background(), nil, s.collection)
			if (d == nil) != tt.wantErr {
				t.Errorf("Session.DeleteMany() error, wantErr %v", tt.wantErr)
			}
		})
	}
}

// FUNCTIONS CALLED BY STEPS

func TestGenerateUUIDStoreItStep(t *testing.T) {
	// 1-Establish a Session instance and context
	s := &Session{}
	ctx := golium.InitializeContext(context.Background())

	// 2-Generating a Random UUID
	t.Run("Generating a Random UUID", func(t *testing.T) {
		if s.GenerateUUIDStoreItStep(ctx) != nil {
			t.Errorf("Error generating a random UUID")
		}
	})

	// 3-Verifying that a UUID has been created
	if len(s.idCollection) != 36 {
		t.Errorf("s.idCollection was not created")
	}
}

// GENERIC FUNCTIONS

func TestContainsElements(t *testing.T) {
	t.Run("Element exists in slice", func(t *testing.T) {
		slice := []int{1, 2, 3, 4, 5}
		element := 3
		result := ContainsElements(element, slice)
		if !result {
			t.Errorf("Expected %v to be in the slice, but it wasn't.", element)
		}
	})
	t.Run("Element does not exist in slice", func(t *testing.T) {
		slice := []string{"Alcorcon", "Madrid", "Barcelona"}
		element := "Valladolid"
		result := ContainsElements(element, slice)
		if result {
			t.Errorf("Expected %v not to be in the slice, but it was.", element)
		}
	})
}

func TestGetFilter(t *testing.T) {
	t.Run("Create filter with non-nil value", func(t *testing.T) {
		key := "age"
		value := 30
		filter := GetFilter(key, value)
		expectedFilter := bson.M{"age": 30}
		if !reflect.DeepEqual(filter, expectedFilter) {
			t.Errorf("Expected filter: %v, but got: %v", expectedFilter, filter)
		}
	})
	t.Run("Create filter with nil value", func(t *testing.T) {
		key := "name"
		filter := GetFilter(key, nil)
		expectedFilter := bson.M{"name": nil}
		if !reflect.DeepEqual(filter, expectedFilter) {
			t.Errorf("Expected filter: %v, but got: %v", expectedFilter, filter)
		}
	})
}

func TestGetFilterConverted(t *testing.T) {
	t.Run("Convert to boolean (true)", func(t *testing.T) {
		field := ISBOOLEAN
		value := "true"
		filter := GetFilterConverted(field, value)
		expectedFilter := primitive.M{field: true}
		if filter[field] != expectedFilter[field] {
			t.Errorf("Expected filter: %v, but got: %v", expectedFilter, filter)
		}
	})
	t.Run("Convert to boolean (false)", func(t *testing.T) {
		field := ISBOOLEAN
		value := "false"
		filter := GetFilterConverted(field, value)
		expectedFilter := primitive.M{field: false}
		if filter[field] != expectedFilter[field] {
			t.Errorf("Expected filter: %v, but got: %v", expectedFilter, filter)
		}
	})
	t.Run("Convert to integer", func(t *testing.T) {
		field := "quantity number"
		value := "42"
		filter := GetFilterConverted(field, value)
		expectedFilter := primitive.M{field: 42}
		if filter[field] != expectedFilter[field] {
			t.Errorf("Expected filter: %v, but got: %v", expectedFilter, filter)
		}
	})
	t.Run("Convert to float64", func(t *testing.T) {
		field := "quantity decimal number"
		value := "99.99"
		filter := GetFilterConverted(field, value)
		expectedFilter := primitive.M{field: 99.99}
		if filter[field] != expectedFilter[field] {
			t.Errorf("Expected filter: %v, but got: %v", expectedFilter, filter)
		}
	})
	t.Run("Convert to nil (empty)", func(t *testing.T) {
		field := "nil or empty"
		value := EMPTY
		filter := GetFilterConverted(field, value)
		expectedFilter := primitive.M{field: nil}
		if filter[field] != expectedFilter[field] {
			t.Errorf("Expected filter: %v, but got: %v", expectedFilter, filter)
		}
	})
	t.Run("No conversion (string)", func(t *testing.T) {
		field := "Nadie"
		value := "Juan"
		filter := GetFilterConverted(field, value)
		expectedFilter := primitive.M{field: value}
		if filter[field] != expectedFilter[field] {
			t.Errorf("Expected filter: %v, but got: %v", expectedFilter, filter)
		}
	})
}

func TestGetOptionsSearchAllFields(t *testing.T) {
	t.Run("Check creates an 'Options' to search all fields in the collection", func(t *testing.T) {
		options := GetOptionsSearchAllFields()
		if options.Projection == nil {
			t.Errorf("Expected Projection to be nil, but got: %v", options.Projection)
		}
	})
}

func TestVerifyMustExist(t *testing.T) {
	t.Run("exist is true", func(t *testing.T) {
		if VerifyMustExist("does exist") != true {
			t.Errorf("Expected true, but got false")
		}
	})
	t.Run("exist is false", func(t *testing.T) {
		if VerifyMustExist("does not exist") != false {
			t.Errorf("Expected false, but got true")
		}
	})
}
