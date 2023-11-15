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
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	PingError        error
	DisconnectError  error
	DatabaseDatabase *mongo.Database
	ContextCliFake   context.Context
	ClientCliFake    *mongo.Client
	NameCliFake      string
	ReadPrefCliFake  *readpref.ReadPref
)

type ClientServiceFuncMock struct{}

func (c ClientServiceFuncMock) Ping(ctx context.Context, rp *readpref.ReadPref,
	client *mongo.Client) error {
	ContextCliFake = ctx
	ReadPrefCliFake = rp
	ClientCliFake = client
	return PingError
}

func (c ClientServiceFuncMock) Disconnect(ctx context.Context, client *mongo.Client) error {
	ContextCliFake = ctx
	ClientCliFake = client
	return DisconnectError
}

func (c ClientServiceFuncMock) Database(name string, client *mongo.Client) *mongo.Database {
	NameCliFake = name
	ClientCliFake = client
	return DatabaseDatabase
}
