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
	"errors"
	"reflect"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/TelefonicaTC2Tech/golium"
)

// FUNCTIONS CALLED BY STEPS

func TestCheckMongoFieldDoesNotExistOrEmptyStep(t *testing.T) {
	// 1-Establish a Session instance, a context, a MongoDB collection, and a _id
	s := &Session{}
	ctx := golium.InitializeContext(context.Background())
	collectionName := "example"
	s.GenerateUUIDStoreItStep(ctx)

	// 2-Create a connection and insert a document into the "example" collection
	mongoConnection(s, t)
	s.DeleteAllDocumentscollectionNameStep(ctx, collectionName)
	s.CreateDocumentscollectionNameStep(ctx, 1, collectionName)

	// 3-Check if the field names do not exist or are empty
	t.Run("Checking that the 'fieldEmpty' field does not exist or is empty in the MongoDB collection", func(t *testing.T) {
		fieldSearched := "fieldEmpty"
		err := s.CheckMongoFieldDoesNotExistOrEmptyStep(ctx, collectionName, fieldSearched, s.idCollection+"_1")
		if err != nil {
			t.Error(err)
		}
	})
	t.Run("Checking that the 'fieldUnexist' field does not exist in the MongoDB collection", func(t *testing.T) {
		fieldSearched := "fieldUnexist"
		err := s.CheckMongoFieldDoesNotExistOrEmptyStep(ctx, collectionName, fieldSearched, s.idCollection+"_1")
		if err != nil {
			t.Error(err)
		}
	})

	// 4-Delete all documents and close connection
	s.DeleteAllDocumentscollectionNameStep(ctx, collectionName)
	s.MongoDisconnectionStep()
}

func TestCheckMongoFieldNameStep(t *testing.T) {
	// 1-Establish a Session instance, a context, a MongoDB collection, and a _id
	s := &Session{}
	ctx := golium.InitializeContext(context.Background())
	collectionName := "example"
	s.GenerateUUIDStoreItStep(ctx)

	// 2-Create a connection and insert a document into the "example" collection
	mongoConnection(s, t)
	s.DeleteAllDocumentscollectionNameStep(ctx, collectionName)
	s.CreateDocumentscollectionNameStep(ctx, 1, collectionName)

	// 3-Check whether or not the name of the searched field exists
	// If you want to check that it doesn't exist, you have to pass "not" as a parameter
	t.Run("Verify that the 'fieldString' field exists in the MongoDB collection", func(t *testing.T) {
		fieldSearched := "fieldString"
		err := s.CheckMongoFieldNameStep(ctx, collectionName, fieldSearched, "yes", s.idCollection+"_1")
		if err != nil {
			t.Error(err)
		}
	})
	t.Run("Verify that the 'fieldEmpty' field exists in the MongoDB collection", func(t *testing.T) {
		fieldSearched := "fieldEmpty"
		err := s.CheckMongoFieldNameStep(ctx, collectionName, fieldSearched, "does", s.idCollection+"_1")
		if err != nil {
			t.Error(err)
		}
	})
	t.Run("Verify that the 'fieldUnexist' field does not exist in the MongoDB collection", func(t *testing.T) {
		fieldSearched := "fieldUnexist"
		err := s.CheckMongoFieldNameStep(ctx, collectionName, fieldSearched, "not", s.idCollection+"_1")
		if err != nil {
			t.Error(err)
		}
	})

	// 4-Delete all documents and close connection
	s.DeleteAllDocumentscollectionNameStep(ctx, collectionName)
	s.MongoDisconnectionStep()
}

func TestCheckMongoValueIDStep(t *testing.T) {
	// 1-Establish a Session instance, a context, a MongoDB collection, and a _id
	s := &Session{}
	ctx := golium.InitializeContext(context.Background())
	collectionName := "example"
	s.GenerateUUIDStoreItStep(ctx)

	// 2-Create a connection and insert a document into the "example" collection
	mongoConnection(s, t)
	s.DeleteAllDocumentscollectionNameStep(ctx, collectionName)
	s.CreateDocumentscollectionNameStep(ctx, 1, collectionName)

	// 3-Check whether or not the _id sought exists
	// If you want to check that it doesn't exist, you have to pass "not" as a parameter
	t.Run("Verify that the '_id' exists in the MongoDB collection", func(t *testing.T) {
		err := s.CheckMongoValueIDStep(ctx, collectionName, s.idCollection+"_1", "yes")
		if err != nil {
			t.Error(err)
		}
	})
	t.Run("Verify that the '_id' does not exist in the MongoDB collection", func(t *testing.T) {
		err := s.CheckMongoValueIDStep(ctx, collectionName, "123456_id_unexist", "not")
		if err != nil {
			t.Error(err)
		}
	})

	// 4-Delete all documents and close connection
	s.DeleteAllDocumentscollectionNameStep(ctx, collectionName)
	s.MongoDisconnectionStep()
}

func TestCheckMongoValuesStep(t *testing.T) {
	// 1-Establish a Session instance, a context, a MongoDB collection, and a _id
	s := &Session{}
	ctx := golium.InitializeContext(context.Background())
	collectionName := "example"
	s.GenerateUUIDStoreItStep(ctx)

	// 2-Create a connection and insert a document into the "example" collection
	mongoConnection(s, t)
	s.DeleteAllDocumentscollectionNameStep(ctx, collectionName)
	s.CreateDocumentscollectionNameStep(ctx, 1, collectionName)

	// 3-Check whether or not the values you are looking for exist
	// If you want to check that it doesn't exist, you have to pass "not" as a parameter
	t.Run("Verify that the searched values exist in the MongoDB collection", func(t *testing.T) {
		// Data that will be searched in the collection and that will be found
		propTableYes := golium.NewTable([][]string{
			{"field", "value"},
			{"_id", "[CTXT:_ID]_1"},
			{"fieldString", "Example field string 1"},
			{"fieldInt", "[NUMBER:1]"},
			{"fieldFloat", "[NUMBER:3.14]"},
			{"fieldBool", "[TRUE]"},
			{"fieldSlice.#", "[NUMBER:3]"},
			{"fieldSlice.0", "itemSlice_1"},
			{"fieldSlice.1", "itemSlice20"},
			{"fieldSlice.2", "itemSlice30"},
			{"fieldEmpty", "[EMPTY]"},
			{"fieldMap.fieldString", "Example field in map string 1"},
			{"fieldMap.fieldInt", "[NUMBER:10]"},
			{"fieldMap.fieldFloat", "[NUMBER:1974.1976]"},
			{"fieldMap.fieldBool", "[FALSE]"},
			{"fieldMap.fieldSliceEmpty.#", "[NUMBER:0]"},
			{"fieldMap.fieldSliceEmpty.0", "[NULL]"},
			{"fieldMap.fieldSliceEmpty.1", ""},
			{"fieldMap.fieldSliceEmpty.2", "[EMPTY]"},
			{"fieldMap.fieldMap2.fieldString", "Example field in map map string 1"},
			{"fieldMap.fieldMap2.fieldInt", "[NUMBER:100]"},
			{"fieldMap.fieldMap2.fieldFloat", "[NUMBER:1974.1976]"},
			{"fieldMap.fieldMap2.fieldBool", "[FALSE]"},
			{"fieldMap.fieldMap2.fieldEmpty", "[NULL]"},
			{"fieldMap.fieldMap2.fieldEmptyText", "[EMPTY]"},
		})
		err := s.CheckMongoValuesStep(ctx, collectionName, s.idCollection+"_1", "yes", propTableYes)
		if err != nil {
			t.Error(err)
		}
	})
	t.Run("Verify that the searched values do not exist in the MongoDB collection", func(t *testing.T) {
		// Data that will be searched in the collection and will not be found
		propTableNot := golium.NewTable([][]string{
			{"field", "value"},
			{"_id", "123456_id_unexist"},
		})
		err := s.CheckMongoValuesStep(ctx, collectionName, s.idCollection+"_1", "not", propTableNot)
		if err != nil {
			t.Error(err)
		}
	})

	// 4-Delete all documents and close connection
	s.DeleteAllDocumentscollectionNameStep(ctx, collectionName)
	s.MongoDisconnectionStep()
}

func TestCheckNumberDocumentscollectionNameStep(t *testing.T) {
	// 1-Establish a Session instance, context, and MongoDB collection
	s := &Session{}
	ctx := golium.InitializeContext(context.Background())
	collectionName := "example"

	// 2-Create a connection and create 10 documents from the "example" collection
	mongoConnection(s, t)
	s.DeleteAllDocumentscollectionNameStep(ctx, collectionName)
	s.CreateDocumentscollectionNameStep(ctx, 10, collectionName)

	// 3-Count all the documents
	t.Run("Contar los documentos de una colección de MongoDB", func(t *testing.T) {
		err := s.CheckNumberDocumentscollectionNameStep(collectionName, 10)
		if err != nil {
			t.Error(err)
		}
	})

	// 4-Delete all documents and close connection
	s.DeleteAllDocumentscollectionNameStep(ctx, collectionName)
	s.MongoDisconnectionStep()
}

func TestCreateDocumentscollectionNameStep(t *testing.T) {
	// 1-Establish a Session instance, context, and MongoDB collection
	s := &Session{}
	ctx := golium.InitializeContext(context.Background())
	collectionName := "example"

	// 2-Create a connection and delete all documents in the "example" collection
	mongoConnection(s, t)
	s.DeleteAllDocumentscollectionNameStep(ctx, collectionName)
	err := s.CheckNumberDocumentscollectionNameStep(collectionName, 0)
	if err != nil {
		t.Error(err)
	}

	// 3-Insert a document
	err = s.CreateDocumentscollectionNameStep(ctx, 1, collectionName)
	if err != nil {
		t.Error(err)
	}

	// 4-Check that it has been inserted and then delete all documents
	t.Run("Verify that the document has been inserted into MongoDB", func(t *testing.T) {
		err := s.CheckNumberDocumentscollectionNameStep(collectionName, 1)
		if err != nil {
			t.Error(err)
		}
	})

	// 5-Delete all documents and close connection
	s.DeleteAllDocumentscollectionNameStep(ctx, collectionName)
	s.MongoDisconnectionStep()
}

func TestDeleteAllDocumentscollectionNameStep(t *testing.T) {
	// 1-Establish a Session instance, context, and MongoDB collection
	s := &Session{}
	ctx := golium.InitializeContext(context.Background())
	collectionName := "example"

	// 2-Create a connection and insert 10 documents into the "example" collection
	mongoConnection(s, t)
	s.CreateDocumentscollectionNameStep(ctx, 10, collectionName)

	// 3-Delete all documents
	t.Run("Deleting all documents in a MongoDB collection", func(t *testing.T) {
		err := s.DeleteAllDocumentscollectionNameStep(ctx, collectionName)
		if err != nil {
			t.Error(err)
		}
	})

	// 4-Verify that all documents have been deleted
	err := s.CheckNumberDocumentscollectionNameStep(collectionName, 0)
	if err != nil {
		t.Error(err)
	}
}

func TestDeleteDocumentscollectionNameStep(t *testing.T) {
	// 1-Establish a Session instance, context, and MongoDB collection
	ctx := InitializeContext(context.Background())
	s := GetSession(ctx)
	collectionName := "example"

	// 2-Create a connection and 3 documents from the "example" collection
	mongoConnection(s, t)
	s.DeleteAllDocumentscollectionNameStep(ctx, collectionName)
	s.CreateDocumentscollectionNameStep(ctx, 3, collectionName)

	// 3-Delete a document
	t.Run("Deleting a document in MongoDB", func(t *testing.T) {
		err := s.DeleteDocumentscollectionNameStep(ctx, collectionName, "fieldInt", "1")
		if err != nil {
			t.Error(err)
		}
	})

	// 4-Verify that it has been deleted, delete all documents and close connection
	err := s.CheckNumberDocumentscollectionNameStep(collectionName, 2)
	if err != nil {
		t.Error(err)
	}
	s.DeleteAllDocumentscollectionNameStep(ctx, collectionName)
	s.MongoDisconnectionStep()
}

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

func mongoConnection(s *Session, t *testing.T) {
	propTable := golium.NewTable([][]string{
		{"field", "value"},
		{"User", "mongoadmin"},
		{"Password", "mongoadmin"},
		{"Host", "localhost:27017"},
		{"AuthSource", "admin"},
		{"Database", "golium-demo"},
	})
	t.Run("Connecting to MongoDB", func(t *testing.T) {
		err := s.MongoConnectionStep(context.Background(), propTable)
		if err != nil {
			t.Error(err)
		}
	})
}

func TestMongoConnectionStep(t *testing.T) {
	s := &Session{}
	mongoConnection(s, t)
	s.MongoDisconnectionStep()
}

func TestMongoDisconnectionStep(t *testing.T) {
	// 1-Establish a Session instance, context, and MongoDB collection
	s := &Session{}
	ctx := golium.InitializeContext(context.Background())
	collectionName := "example"

	// 2-Create a connection and insert 3 collections
	mongoConnection(s, t)
	s.DeleteAllDocumentscollectionNameStep(ctx, collectionName)
	s.CreateDocumentscollectionNameStep(ctx, 3, collectionName)
	err := s.CheckNumberDocumentscollectionNameStep(collectionName, 3)
	if err != nil {
		t.Errorf("Error al contar los documentos de la colección '%s'", collectionName)
	}

	// 3-Close the connection to MongoDB
	t.Run("Close the connection to MongoDB", func(t *testing.T) {
		err := s.MongoDisconnectionStep()
		if err != nil {
			t.Error(err)
		}
	})

	// 4-Verify that the connection to MongoDB has been closed
	errString := s.CheckNumberDocumentscollectionNameStep(collectionName, 3).Error()
	if errString != "error: query error: 'client is disconnected'" {
		t.Errorf("Connection to MongoDB is still open")
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
		var value interface{} = nil
		filter := GetFilter(key, value)
		expectedFilter := bson.M{"name": nil}
		if !reflect.DeepEqual(filter, expectedFilter) {
			t.Errorf("Expected filter: %v, but got: %v", expectedFilter, filter)
		}
	})
}

func TestGetFilterConverted(t *testing.T) {
	t.Run("Convert to boolean (true)", func(t *testing.T) {
		field := "is_boolean"
		value := "true"
		filter := GetFilterConverted(field, value)
		expectedFilter := primitive.M{field: true}
		if filter[field] != expectedFilter[field] {
			t.Errorf("Expected filter: %v, but got: %v", expectedFilter, filter)
		}
	})
	t.Run("Convert to boolean (false)", func(t *testing.T) {
		field := "is_boolean"
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
		value := "[EMPTY]"
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

// SESSION FUNCTIONS

func TestCreateDocumentsCollection(t *testing.T) {
	// 1-Establish a Session instance and context
	s := &Session{}
	ctx := golium.InitializeContext(context.Background())

	// 2-Checks if the function can create an array with the correct number of interfaces
	// and if the values within the array are as expected
	t.Run("Create array with 3 interfaces", func(t *testing.T) {
		count := 3
		result := s.CreateDocumentsCollection(ctx, count)
		if len(result) != count {
			t.Errorf("Expected array length of %d, but got %d", count, len(result))
		}
		for i := 0; i < count; i++ {
			if result[i].(map[string]interface{})["fieldInt"] != i+1 {
				t.Errorf("Expected value at index %d to be %d, but got %v", i, i, result[i])
			}
		}
	})
	t.Run("Create empty array", func(t *testing.T) {
		count := 0
		result := s.CreateDocumentsCollection(ctx, count)
		if len(result) != count {
			t.Errorf("Expected empty array, but got an array of length %d", len(result))
		}
	})
}

func TestExistFieldCollection(t *testing.T) {
	s := &Session{}
	s.fieldsCollectionName = []string{"Alcorcon", "Madrid", "Barcelona"}

	t.Run("Verify that an element exists in fieldsCollectionName", func(t *testing.T) {
		if s.ExistFieldCollection("Alcorcon") != true {
			t.Errorf("'Alcorcon' exist in fieldsCollectionName")
		}
	})
	t.Run("Verify that an item does not exist in fieldsCollectionName", func(t *testing.T) {
		if s.ExistFieldCollection("Valladolid") != false {
			t.Errorf("'Valladolid' does not exist in fieldsCollectionName")
		}
	})
}

func TestGetDecodeDocument(t *testing.T) {
	// 1-Establish a Session instance, a context, a MongoDB collection, and a _id
	s := &Session{}
	ctx := golium.InitializeContext(context.Background())
	collectionName := "example"
	s.GenerateUUIDStoreItStep(ctx)
	fieldSearched := "_id"

	// 2-Create a connection and insert a document into the "example" collection
	mongoConnection(s, t)
	s.DeleteAllDocumentscollectionNameStep(ctx, collectionName)
	s.CreateDocumentscollectionNameStep(ctx, 1, collectionName)

	// 3-Create the s.singleResult and get the document decoded
	t.Run("Check GetDecodeDocument decodes the BSON document in the bsonDoc variable", func(t *testing.T) {
		searchedValue := s.idCollection + "_1"
		s.SetSingleResult(ctx, fieldSearched, interface{}(searchedValue))
		_, err := s.GetDecodeDocument(*s.singleResult)
		if err != nil {
			t.Error(err)
		}
	})
	t.Run("Check GetDecodeDocument does not decode the BSON document in the bsonDoc variable", func(t *testing.T) {
		searchedValue := s.idCollection + "_2"
		s.SetSingleResult(ctx, fieldSearched, interface{}(searchedValue))
		_, err := s.GetDecodeDocument(*s.singleResult)
		if err == nil {
			t.Error(err)
		}
	})

	// 4-Delete all documents and close connection
	s.DeleteAllDocumentscollectionNameStep(ctx, collectionName)
	s.MongoDisconnectionStep()
}

func TestSetCollection(t *testing.T) {
	// 1-Establish a Session instance, context, and MongoDB collection
	s := &Session{}
	ctx := golium.InitializeContext(context.Background())
	collectionName := "example"

	// 2-Create a connection and insert a document into the "example" collection
	mongoConnection(s, t)
	s.DeleteAllDocumentscollectionNameStep(ctx, collectionName)
	s.CreateDocumentscollectionNameStep(ctx, 1, collectionName)

	// 3-Verify that the collection is created
	t.Run("Verify that the collection is created", func(t *testing.T) {
		s.SetCollection(collectionName)
		if s.collection.Name() != collectionName {
			t.Errorf("s.collection is not %s, it is %s", collectionName, s.collection.Name())
		}
	})

	// 4-Delete all documents and close connection
	s.DeleteAllDocumentscollectionNameStep(ctx, collectionName)
	s.MongoDisconnectionStep()
}

func TestSetDataCollectionJSONBytes(t *testing.T) {
	// 1-Establish a Session instance
	s := &Session{}

	// 2- Check converts BSON to JSON
	t.Run("Check convert correct BSON object to JSON", func(t *testing.T) {
		doc := bson.D{
			{Key: "name", Value: "Juan Nadie"},
			{Key: "age", Value: 30},
		}
		err := s.SetDataCollectionJSONBytes(doc)
		if err != nil {
			t.Error(err)
		}
	})
	t.Run("Check convert incorrect BSON object to JSON", func(t *testing.T) {
		// A function in the value, which is not serializable in JSON
		doc := bson.D{
			{"name", func() {}},
		}
		err := s.SetDataCollectionJSONBytes(doc)
		if err == nil {
			t.Error(err)
		}
	})
}

func TestSetFieldsCollectionName(t *testing.T) {
	// 1-Establish a Session instance, a context, a MongoDB collection, and a _id
	s := &Session{}
	ctx := golium.InitializeContext(context.Background())
	collectionName := "example"
	s.GenerateUUIDStoreItStep(ctx)

	// 2-Create a connection and insert a document into the "example" collection
	mongoConnection(s, t)
	s.DeleteAllDocumentscollectionNameStep(ctx, collectionName)
	s.CreateDocumentscollectionNameStep(ctx, 1, collectionName)

	// 3-Verify that the s.fieldsCollectionName slice is empty
	if len(s.fieldsCollectionName) != 0 {
		t.Error("s.fieldsCollectionName is not empty: ", s.fieldsCollectionName)
	}

	// 4-Verify that a slice with field names is saved in s.fieldsCollectionName
	t.Run("Verify that a slice is saved", func(t *testing.T) {
		err := s.SetFieldsCollectionName(ctx, s.idCollection+"_1")
		if err != nil {
			t.Error(err)
		}
		if len(s.fieldsCollectionName) == 0 {
			t.Error("s.fieldsCollectionName is empty: ", s.fieldsCollectionName)
		}
	})
	t.Run("Verify that a slice is not saved if the searched value does not exist", func(t *testing.T) {
		err := s.SetFieldsCollectionName(ctx, s.idCollection+"_2")
		if err == nil {
			t.Error(err)
		}
	})

	// 5-Delete all documents and close connection
	s.DeleteAllDocumentscollectionNameStep(ctx, collectionName)
	s.MongoDisconnectionStep()
}

func TestSetSingleResult(t *testing.T) {
	// 1-Establish a Session instance, a context, a MongoDB collection, and a _id
	s := &Session{}
	ctx := golium.InitializeContext(context.Background())
	collectionName := "example"
	s.GenerateUUIDStoreItStep(ctx)
	fieldSearched := "_id"

	// 2-Create a connection and insert a document into the "example" collection
	mongoConnection(s, t)
	s.DeleteAllDocumentscollectionNameStep(ctx, collectionName)
	s.CreateDocumentscollectionNameStep(ctx, 1, collectionName)

	// 3-Check for s.singleResult
	t.Run("Verify that the s.singleResult exists", func(t *testing.T) {
		searchedValue := s.idCollection + "_1"
		err := s.SetSingleResult(ctx, fieldSearched, interface{}(searchedValue))
		if err != nil {
			t.Error(err)
		}
	})
	t.Run("Checking that the s.singleResult does not exist", func(t *testing.T) {
		searchedValue := s.idCollection + "_2"
		err := s.SetSingleResult(ctx, fieldSearched, interface{}(searchedValue))
		if err == nil {
			t.Error(err)
		}
	})

	// 4-Delete all documents and close connection
	s.DeleteAllDocumentscollectionNameStep(ctx, collectionName)
	s.MongoDisconnectionStep()
}

func TestValidateDataMongo(t *testing.T) {
	// 1-Establish a Session instance, a context, a MongoDB collection, and a _id
	s := &Session{}
	ctx := golium.InitializeContext(context.Background())
	collectionName := "example"
	s.GenerateUUIDStoreItStep(ctx)

	// 2-Create a connection and insert a document into the "example" collection
	mongoConnection(s, t)
	s.DeleteAllDocumentscollectionNameStep(ctx, collectionName)
	s.CreateDocumentscollectionNameStep(ctx, 1, collectionName)

	// 3-Check the existence of the searched values
	t.Run("Verify that the searched values exist in the MongoDB collection", func(t *testing.T) {
		// Data that will be searched in the collection and that will be found
		propTable := golium.NewTable([][]string{
			{"field", "value"},
			{"_id", s.idCollection + "_1"},
			{"fieldString", "Example field string 1"},
		})
		props, _ := golium.ConvertTableToMap(ctx, propTable)
		exist, errValidate := s.ValidateDataMongo(ctx, s.idCollection+"_1", props)
		if !exist || errValidate != nil {
			t.Error(errValidate)
		}
	})
	t.Run("Verify that the searched values do not exist in the MongoDB collection", func(t *testing.T) {
		// Data that will be searched in the collection and that will not be found
		propTable := golium.NewTable([][]string{
			{"field", "value"},
			{"_id", "123456_id_unexist"},
		})
		props, _ := golium.ConvertTableToMap(ctx, propTable)
		exist, errValidate := s.ValidateDataMongo(ctx, s.idCollection+"_1", props)
		if exist || errValidate == nil {
			t.Error(errValidate)
		}
	})

	// 4-Delete all documents and close connection
	s.DeleteAllDocumentscollectionNameStep(ctx, collectionName)
	s.MongoDisconnectionStep()
}

func TestVerifyExistAndMustExistValue(t *testing.T) {
	// 1-Establish a Session instance
	s := &Session{}

	//2-Creating test combinations
	t.Run("Checking that values exist and should exist", func(t *testing.T) {
		exist := true
		mustExist := true
		var errFunction error
		err := s.VerifyExistAndMustExistValue(exist, mustExist, errFunction)
		if err != nil {
			t.Error(err)
		}
	})
	t.Run("Checking that values exist and should not exist", func(t *testing.T) {
		exist := true
		mustExist := false
		var errFunction error
		err := s.VerifyExistAndMustExistValue(exist, mustExist, errFunction)
		if err == nil {
			t.Error(err)
		}
	})
	t.Run("Checking that values do not exist and should exist", func(t *testing.T) {
		exist := false
		mustExist := true
		var errFunction error
		err := s.VerifyExistAndMustExistValue(exist, mustExist, errFunction)
		if err == nil {
			t.Error(err)
		}
	})
	t.Run("Checking that values don't exist and shouldn't exist", func(t *testing.T) {
		exist := false
		mustExist := false
		var errFunction error
		err := s.VerifyExistAndMustExistValue(exist, mustExist, errFunction)
		if err != nil {
			t.Error(err)
		}
	})
	t.Run("Checking for a previous error", func(t *testing.T) {
		exist := false
		mustExist := true
		errFunction := errors.New("errorFunction is not nil")
		err := s.VerifyExistAndMustExistValue(exist, mustExist, errFunction)
		if err == nil {
			t.Error(err)
		}
	})
}

func TestVerifyExistID(t *testing.T) {
	// 1-Establish a Session instance, a context, a MongoDB collection, and a _id
	s := &Session{}
	ctx := golium.InitializeContext(context.Background())
	collectionName := "example"
	s.GenerateUUIDStoreItStep(ctx)

	// 2-Create a connection and insert a document into the "example" collection
	mongoConnection(s, t)
	s.DeleteAllDocumentscollectionNameStep(ctx, collectionName)
	s.CreateDocumentscollectionNameStep(ctx, 1, collectionName)

	// 3-Check if the s.idCollection exists
	t.Run("Verify that the value exists in the _id field", func(t *testing.T) {
		_, err := s.VerifyExistID(ctx, s.idCollection+"_1")
		if err != nil {
			t.Error(err)
		}
	})
	t.Run("Verify that the value does not exist in the _id field", func(t *testing.T) {
		_, err := s.VerifyExistID(ctx, s.idCollection+"_2")
		if err == nil {
			t.Error(err)
		}
	})

	// 4-Delete all documents and close connection
	s.DeleteAllDocumentscollectionNameStep(ctx, collectionName)
	s.MongoDisconnectionStep()
}
