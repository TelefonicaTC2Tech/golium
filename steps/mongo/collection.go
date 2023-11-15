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

type CollectionFunctions interface {
	Collection(name string, database *mongo.Database) *mongo.Collection
	Name(collection *mongo.Collection) string
	Find(ctx context.Context, filter interface{}, collection *mongo.Collection, opts ...*options.FindOptions) (*mongo.Cursor, error)
	FindOne(ctx context.Context, filter interface{}, collection *mongo.Collection, opt ...*options.FindOneOptions) *mongo.SingleResult
	InsertMany(ctx context.Context, documents []interface{}, collection *mongo.Collection) (*mongo.InsertManyResult, error)
	DeleteMany(ctx context.Context, filter interface{}, opt *options.DeleteOptions, collection *mongo.Collection) (*mongo.DeleteResult, error)
}

type CollectionService struct{}

func NewMongoCollectionService() *CollectionService {
	return &CollectionService{}
}

func (c CollectionService) Collection(name string, database *mongo.Database) *mongo.Collection {
	return database.Collection(name)
}

func (c CollectionService) Name(collection *mongo.Collection) string {
	return collection.Name()
}

func (c CollectionService) FindOne(ctx context.Context, filter interface{}, collection *mongo.Collection, opts ...*options.FindOneOptions) *mongo.SingleResult {
	return collection.FindOne(ctx, filter, opts...)
}

func (c CollectionService) Find(ctx context.Context, filter interface{}, collection *mongo.Collection, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	return collection.Find(ctx, filter, opts...)
}

func (c CollectionService) InsertMany(ctx context.Context, documents []interface{}, collection *mongo.Collection) (*mongo.InsertManyResult, error) {
	return collection.InsertMany(ctx, documents)
}

func (c CollectionService) DeleteMany(ctx context.Context, filter interface{}, opt *options.DeleteOptions, collection *mongo.Collection) (*mongo.DeleteResult, error) {
	return collection.DeleteMany(ctx, filter, opt)
}
