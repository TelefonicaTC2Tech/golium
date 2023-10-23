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

	"github.com/TelefonicaTC2Tech/golium"
	"github.com/cucumber/godog"
	 "github.com/cucumber/messages-go/v16"		

	"go.mongodb.org/mongo-driver/mongo"

)

var client *mongo.Client
var collection *mongo.Collection

type MongoSteps struct {
}

// InitializeSteps initializes all the steps
func (us MongoSteps) InitializeSteps(ctx context.Context, scenCtx *godog.ScenarioContext) context.Context {
	// Initialize the HTTP session in the context
	ctx = InitializeContext(ctx)
	session := GetSession(ctx)

	scenCtx.Step(`^I connect to MongoDB$`, func() error {
		return session.MongoConnectionStep(ctx)
	})
	scenCtx.Step(`^I check that these values of the MongoDB "([^"]*)" collection with "([^"]*)" _id "([^"]*)" exist$`, func(collectionName string, idCollection string, exist string, t *godog.Table) error {
		return session.CheckMongoValuesStep(ctx, golium.ValueAsString(ctx, collectionName), golium.ValueAsString(ctx, idCollection), golium.ValueAsString(ctx, exist), t)
	})
	scenCtx.Step(`^I check that the MongoDB "([^"]*)" collection with "([^"]*)" _id "([^"]*)" exist$`, func(collectionName string, idCollection string, exist string) error {
		return session.CheckMongoValueIDStep(ctx, golium.ValueAsString(ctx, collectionName), golium.ValueAsString(ctx, idCollection), golium.ValueAsString(ctx, exist))
	})
	scenCtx.Step(`^I check that in the MongoDB "([^"]*)" collection, "([^"]*)" field "([^"]*)" exist for the "([^"]*)" _id$`, func(collectionName string, fieldSearched string, exist string, idCollection string) error {
		return session.CheckMongoFieldNameStep(ctx, golium.ValueAsString(ctx, collectionName), golium.ValueAsString(ctx, fieldSearched), golium.ValueAsString(ctx, exist), golium.ValueAsString(ctx, idCollection))
	})
	scenCtx.Step(`^I check that the "([^"]*)" field in the MongoDB "([^"]*)" collection does not exist or is empty for the "([^"]*)" _id$`, func(fieldSearched string, collectionName string, idCollection string) error {
		return session.CheckMongoFieldExistOrEmptyStep(ctx, golium.ValueAsString(ctx, fieldSearched), golium.ValueAsString(ctx, collectionName), golium.ValueAsString(ctx, idCollection))
	})
	scenCtx.Step(`^I create "(\d+)" documents in the MongoDB "([^"]*)" collection$`, func(num int, collectionName string) error {
		return session.CreateDocumentscollectionNameStep(ctx, num, golium.ValueAsString(ctx, collectionName))
	})
	scenCtx.Step(`^I delete documents from the MongoDB "([^"]*)" collection whose "([^"]*)" field is "([^"]*)" value$`, func(collectionName string, field string, value string) error {
		return session.DeleteDocumentscollectionNameStep(ctx, golium.ValueAsString(ctx, collectionName), golium.ValueAsString(ctx, field), value)
	})	
	scenCtx.Step(`^I check that the number of documents in collection "([^"]*)" is "(\d+)"$`, func(collectionName string, num int) error {
		return session.CheckNumberDocumentscollectionNameStep(ctx, golium.ValueAsString(ctx, collectionName), num)
	})	

	scenCtx.AfterScenario(func(sc *messages.Pickle, err error) {
		session.MongoDisconnection(ctx)
	})

	return ctx
}
