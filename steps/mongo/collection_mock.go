// Copyright (c) Telefónica Cybersecurity & Cloud Tech S.L.
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

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	CollectionCollection       *mongo.Collection
	NameString                 string
	FindOneSingleResult        *mongo.SingleResult
	FindCursor                 *mongo.Cursor
	FindError                  error
	InsertManyInsertManyResult *mongo.InsertManyResult
	InsertManyError            error
	DeleteResultDeleteResult   *mongo.DeleteResult
	DeleteResultError          error
	ContextColFake             context.Context
	CollectionColFake          *mongo.Collection
	FindOptsColFake            []*options.FindOptions
	FindOneOptsColFake         []*options.FindOneOptions
	DeleteOptsColFake          []*options.DeleteOptions
	DatabaseColFake            *mongo.Database
	NameColFake                string
	FilterColFake              interface{}
	DocumentsColFake           interface{}
)

type CollectionServiceFuncMock struct{}

func (c CollectionServiceFuncMock) Collection(name string,
	database *mongo.Database) *mongo.Collection {
	NameColFake = name
	DatabaseColFake = database
	return CollectionCollection
}

func (c CollectionServiceFuncMock) Name(collection *mongo.Collection) string {
	CollectionColFake = collection
	return NameString
}

func (c CollectionServiceFuncMock) FindOne(ctx context.Context, filter interface{},
	collection *mongo.Collection, opts ...*options.FindOneOptions) *mongo.SingleResult {
	ContextColFake = ctx
	FilterColFake = filter
	CollectionColFake = collection
	FindOneOptsColFake = opts
	return FindOneSingleResult
}

func (c CollectionServiceFuncMock) Find(ctx context.Context, filter interface{},
	collection *mongo.Collection, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	ContextColFake = ctx
	FilterColFake = filter
	CollectionColFake = collection
	FindOptsColFake = opts
	return FindCursor, FindError
}

func (c CollectionServiceFuncMock) InsertMany(ctx context.Context, documents []interface{},
	collection *mongo.Collection) (*mongo.InsertManyResult, error) {
	ContextColFake = ctx
	DocumentsColFake = documents
	CollectionColFake = collection
	return InsertManyInsertManyResult, InsertManyError
}

func (c CollectionServiceFuncMock) DeleteMany(ctx context.Context, filter interface{},
	collection *mongo.Collection, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	ContextColFake = ctx
	FilterColFake = filter
	CollectionColFake = collection
	DeleteOptsColFake = opts
	return DeleteResultDeleteResult, DeleteResultError
}
