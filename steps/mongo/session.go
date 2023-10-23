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
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/TelefonicaTC2Tech/golium"
	"github.com/cucumber/godog"
	"github.com/google/uuid"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Session contains the information of a MongoDB session.
type Session struct {

	// Save MongoDB login credentials: username, password, and AuthSource
	credentials options.Credential

	// Saves the host from accessing MongoDB
	host string

	// Save the MongoDB database
	database string

	// Set the host and credentials for the connection
	clientOptions *options.ClientOptions

	// Create the mongo client with the MongoDB connection
	client *mongo.Client

	// Points to a record in a collection
	singleResult *mongo.SingleResult

	// Saves the collection to be used
	collection *mongo.Collection

	// Save the name of the collectionName fields
	fieldsCollectionName []string

	// Save collection items as JSON
	dataCollectionJSONBytes []byte

	// Access MongoDB features
	MongoClientService ClientFunctions
}

// FUNCTIONS CALLED BY STEPS

// CheckMongoFieldExistOrEmptyStep check that a field does not exist or it does exist and is empty
func (s *Session) CheckMongoFieldExistOrEmptyStep(
	ctx context.Context, fieldSearched, collectionName, idCollection string,
) error {
	// 1-Setting the Collection Name in the Session
	s.SetCollection(collectionName)

	// 2-Set fields collection name in Session: s.fieldsCollectionName
	err := s.SetFieldsCollectionName(ctx, idCollection)
	if err != nil {
		return err
	}

	// 3-If fieldSearched exists, check that the value of the fieldSearched is null
	if s.ExistFieldCollection(fieldSearched) {
		err = s.SetSingleResult(ctx, fieldSearched, nil)
	}

	return err
}

// CheckMongoFieldNameStep check if the name of the field searched of the user collection is correct
func (s *Session) CheckMongoFieldNameStep(
	ctx context.Context, collectionName, fieldSearched, exist, idCollection string,
) error {
	// 1-Set collection name and fields collection name in Session
	s.SetCollection(collectionName)
	s.SetFieldsCollectionName(ctx, idCollection)

	// 2-Get boolean if field exist and must exist in collection
	existField := s.ExistFieldCollection(fieldSearched)
	mustExistField := VerifyMustExist(exist)

	// 3-Verify exist and must exist
	return s.VerifyExistAndMustExistValue(existField, mustExistField, nil)
}

// CheckMongoValueIDStep checks if the past idCollection exists in the collection
func (s *Session) CheckMongoValueIDStep(
	ctx context.Context, collectionName, idCollection, exist string,
) error {
	// 1-Set collection name and fields collection name in Session
	s.SetCollection(collectionName)

	// 2-Get boolean if _id exist and must exist in collection
	existID, err := s.VerifyExistID(ctx, idCollection)
	mustExistID := VerifyMustExist(exist)

	// 3-Verify exist and must exist
	return s.VerifyExistAndMustExistValue(existID, mustExistID, err)
}

// CheckMongoValuesStep checks the value of the MongoDB fields in the specified collection
func (s *Session) CheckMongoValuesStep(
	ctx context.Context, collectionName, idCollection, exist string, t *godog.Table,
) error {
	// 1-Set collection name and fields collection name in Session
	s.SetCollection(collectionName)

	// 2-Get value of specified table
	props, err := golium.ConvertTableToMap(ctx, t)
	if err != nil {
		return fmt.Errorf("ERROR: failed processing the table for validating the body: '%w'", err)
	}

	// 3-Get boolean if data exist and must exist in collection
	existValue, err := s.ValidateDataMongo(ctx, idCollection, props)
	mustExistValue := VerifyMustExist(exist)

	// 4-Verify exist and must exist
	return s.VerifyExistAndMustExistValue(existValue, mustExistValue, err)
}

// MongoConnectionStep establishes a connection in MongoDB.
// The connection, client and database data are saved in s.clientOptions, s.client, and s.database
func (s *Session) MongoConnectionStep(ctx context.Context) error {
	// 1-Set credentials and host in session
	s.SetCredentials(ctx)
	s.SetHost(ctx)

	// 2-Set clientOptions in session
	s.clientOptions = options.Client().SetHosts(strings.Split(s.host, ",")).SetAuth(s.credentials)

	// 3-Connect to the MongoDB server and set client in session
	var err error
	s.client, err = mongo.Connect(context.Background(), s.clientOptions)
	if err != nil {
		return fmt.Errorf("error: problems with the client options or with the context. '%s'", err)
	}

	// 4-Check the connection to the MongoDB server
	err = s.MongoClientService.Ping(s.client)
	if err != nil {
		return fmt.Errorf("error: problems with connection to MongoDB. '%s'", err)
	}

	// 5-Set the database in session
	s.SetDatabase(ctx)

	return nil
}

// CreateDocumentcollectionNameStep creates a number of documents in the specified collection
func (s *Session) CreateDocumentscollectionNameStep(
	ctx context.Context, num int, collectionName string,
) error {
	// 1-collection in which the insertion will be made, if it does not exist it is created
	s.SetCollection(collectionName)

	// 2-The documents to be inserted are created
	allDocuments := s.CreateDocumentsCollection(ctx, num)

	// 3-Insert the documents into the "collectionName" collection of the database
	_, err := s.collection.InsertMany(context.TODO(), allDocuments)
	if err != nil {
		return err
	}

	return nil
}

// DeleteDocumentscollectionNameStep delete a document from the MongoDB collection
func (s *Session) DeleteDocumentscollectionNameStep(
	ctx context.Context, collectionName, field, value string,
) error {
	// 1-The collection in which the deletion is to be made is established.
	s.SetCollection(collectionName)

	// 2-Performs the deletion of documents that match the filter after the data type is converted.
	// Only it is possible filter by string, int, float, or boolean values, also in slices and maps.
	_, err := s.collection.DeleteMany(ctx, GetFilterConverted(field, value))
	if err != nil {
		return err
	}

	return nil
}

// CheckNumberDocumentscollectionNameStep verify the number of documents in collection
func (s *Session) CheckNumberDocumentscollectionNameStep(collectionName string, num int) error {
	// 1-The collection from which the documents are to be counted is established
	s.SetCollection(collectionName)

	// 2-Make a query to get all the documents in the collection
	cursor, err := s.collection.Find(context.Background(), bson.D{})
	if err != nil {
		return fmt.Errorf("error: query error: '%s'", err)
	}

	// 3-Iterate through the documents and count the ones that are there
	count := 0
	for cursor.Next(context.Background()) {
		count++
	}

	// 4-Check the result
	if count != num {
		return fmt.Errorf("error: the number of documents is '%d' and should be '%d'", count, num)
	}
	return nil
}

// GENERIC FUNCTIONS

// ContainsElement check if an item exists in a slice
func ContainsElement(expectedElement, sliceElements interface{}) bool {
	// 1-Create a reflection object from the slice
	sliceValue := reflect.ValueOf(sliceElements)

	// 2-It is verified that the object of reflection is of the "slice" type
	if sliceValue.Kind() != reflect.Slice {
		return false
	}

	for i := 0; i < sliceValue.Len(); i++ {
		// 3-The slice element converted into an interface is compared with the searched element.
		// The comparison is made in value and type
		if reflect.DeepEqual(sliceValue.Index(i).Interface(), expectedElement) {
			return true
		}
	}

	return false
}

// GetFilter returns a filter to search for a record from a field
func GetFilter(key string, value interface{}) primitive.M {
	if value == nil {
		return bson.M{key: nil}
	}
	return bson.M{key: value}
}

// GetFilterConverted returns a filter with field and data type required for delete a record
func GetFilterConverted(field, value string) primitive.M {
	/* In Golang:
	- the values "1", "true", "t", "T", "TRUE", "True" are interpreted as true
	- the values "0", "false", "f", "F", "FALSE", "False" are interpreted as false.
	To force only a few values in a cell to be considered Boolean, this slice is created.
	*/
	boolSlice := []string{"true", "TRUE", "True", "false", "FALSE", "False"}

	// Try to convert to bool
	if ContainsElement(value, boolSlice) {
		convertedValue, err := strconv.ParseBool(value)
		if err == nil {
			return bson.M{field: convertedValue}
		}
	}

	// Try to convert to int
	if convertedValue, err := strconv.Atoi(value); err == nil {
		return bson.M{field: convertedValue}
	}

	// Try converting to float64
	if convertedValue, err := strconv.ParseFloat(value, 64); err == nil {
		return bson.M{field: convertedValue}
	}

	// If the passed value is "[EMPTY]" or "[NULL]", it evaluates to nil
	if value == "[EMPTY]" || value == "[NULL]" {
		return bson.M{field: nil}
	}

	// If none of the above is true, the value remains as the original type (string)
	return bson.M{field: value}
}

// GetOptionsSearchAllFields creates an "Options" to search all fields in the collection
func GetOptionsSearchAllFields() *options.FindOneOptions {
	return options.FindOne().SetProjection(bson.M{})
}

// VerifyMustExist returns a boolean indicating whether the element should exist
func VerifyMustExist(exist string) bool {
	return !strings.Contains(strings.ToLower(exist), "not")
}

// SESSION FUNCTIONS

// CreateDocumentsCollection creates num documents in a slice and inserts them into a collection
func (s *Session) CreateDocumentsCollection(ctx context.Context, num int) []interface{} {
	// 1-Initialize the document slice
	allDocuments := []interface{}{}

	// 2-Obtaining the _id of the context, if it does not exist it is created
	id := golium.GetContext(ctx).Get("_ID")
	if id == nil {
		var err error
		id, err = uuid.NewRandom()
		if err != nil {
			return nil
		}
		// Transforming id into an interface
		id = id.(uuid.UUID).String()
	}
	// 3-Creating Documents and Inserting Into the Slice
	for i := 1; i <= num; i++ {
		// Defines the document to be inserted.
		// The _id will be the same across all + _ + iteration number
		document := map[string]interface{}{
			"_id":         id.(string) + "_" + strconv.Itoa(i),
			"fieldString": "Example field string " + strconv.Itoa(i),
			"fieldInt":    i,
			"fieldFloat":  3.14,
			"fieldBool":   true,
			"fieldSlice":  []string{"itemSlice_" + strconv.Itoa(i), "itemSlice20", "itemSlice30"},
			"fieldEmpty":  nil,
			"fieldMap": map[string]interface{}{
				"fieldString":     "Example field in map string " + strconv.Itoa(i),
				"fieldInt":        i * 10,
				"fieldFloat":      1974.1976,
				"fieldBool":       false,
				"fieldSliceEmpty": []string{},
				"fieldMap2": map[string]interface{}{
					"fieldString":    "Example field in map map string " + strconv.Itoa(i),
					"fieldInt":       i * 100,
					"fieldFloat":     1974.1976,
					"fieldBool":      false,
					"fieldEmpty":     nil,
					"fieldEmptyText": "",
				},
			},
		}
		allDocuments = append(allDocuments, document)
	}
	return allDocuments
}

// SetCredentials set an Options element. Credential with MongoDB access credentials
func (s *Session) SetCredentials(ctx context.Context) {
	credentials := options.Credential{
		Username:   golium.ValueAsString(ctx, "[CONF:mongoUsername]"),
		Password:   golium.ValueAsString(ctx, "[CONF:mongoPassword]"),
		AuthSource: golium.ValueAsString(ctx, "[CONF:mongoAuthSource]"),
	}
	s.credentials = credentials
}

// SetHost set a string with the host
func (s *Session) SetHost(ctx context.Context) {
	s.host = golium.ValueAsString(ctx, "[CONF:mongoHost]")
}

// SetDatabase set a string with the database
func (s *Session) SetDatabase(ctx context.Context) {
	s.database = golium.ValueAsString(ctx, "[CONF:mongoDatabase]")
}

// SetSingleResult set a Single Result from a Filter Search (GetFilter(...) function)
func (s *Session) SetSingleResult(
	ctx context.Context, fieldSearched string, value interface{},
) error {
	s.singleResult = s.collection.FindOne(
		ctx, GetFilter(fieldSearched, value), GetOptionsSearchAllFields())
	if s.singleResult.Err() != nil {
		return fmt.Errorf("error: the searched '%s' field does not have the '%s' value "+
			"in the '%s' collection", fieldSearched, value, s.collection.Name())
	}
	return nil
}

// SetCollection sets the collection. If the collection does not exist, no error is returned.
// Collections are created dynamically when you insert a document
func (s *Session) SetCollection(collectionName string) {
	s.collection = s.client.Database(s.database).Collection(collectionName)
}

// ExistFieldCollection evaluate whether or not the searched field exists
func (s *Session) ExistFieldCollection(fieldSearched string) bool {
	existField := false
	for _, element := range s.fieldsCollectionName {
		if element == fieldSearched {
			existField = true
		}
	}
	return existField
}

// GetDecodeDocument decodes the BSON document in the bsonDoc variable
func (s *Session) GetDecodeDocument(singleResult mongo.SingleResult) (bson.D, error) {
	var bsonDoc bson.D
	if err := singleResult.Decode(&bsonDoc); err != nil {
		err = fmt.Errorf("error: the decoding of the BSON has been erroneous: '%s'", err)
		return nil, err
	}
	return bsonDoc, nil
}

// SetDataCollectionJSONBytes convert BSON object to JSON
func (s *Session) SetDataCollectionJSONBytes(bsonDoc bson.D) error {
	var err error
	s.dataCollectionJSONBytes, err = bson.MarshalExtJSON(bsonDoc, false, false)
	if err != nil {
		return fmt.Errorf("error: the conversion from BSON to JSON has been erroneous: '%s'", err)
	}
	return nil
}

// SetFieldsCollectionName save a slice with the names of the fields in s.fieldsCollectionName
func (s *Session) SetFieldsCollectionName(ctx context.Context, idCollection string) error {
	// Make a query to find past _id's document
	var document bson.M
	err := s.collection.FindOne(ctx, GetFilter("_id", idCollection)).Decode(&document)
	if err == mongo.ErrNoDocuments {
		return fmt.Errorf("error: no documents matching the filter were found")
	} else if err != nil {
		return fmt.Errorf("error: '%s'", err)
	} else {
		// s.fieldsCollectionName is flushed, and the names of the fields in the document are added
		s.fieldsCollectionName = s.fieldsCollectionName[:0]
		for fieldName := range document {
			s.fieldsCollectionName = append(s.fieldsCollectionName, fieldName)
		}
	}
	return nil
}

// VerifyExistID returns a boolean indicating whether the _id searched exists in the collection
func (s *Session) VerifyExistID(ctx context.Context, idCollection string) (bool, error) {
	// Perform the search and get the result as singleResult
	err := s.SetSingleResult(ctx, "_id", idCollection)
	if err != nil {
		return false, fmt.Errorf("error: searched _id '%s' does not exist in the '%s' collection",
			idCollection, s.collection.Name())
	}
	return true, nil
}

// VerifyExistAndMustExistValue check if the values exist and should exist
func (s *Session) VerifyExistAndMustExistValue(exist, mustExist bool, err error) error {
	// If Exist and should NOT exist or NOT exist and should exist return error
	// If Exist and shoud exist OR not exist and should not exist return nil
	if exist == mustExist {
		return nil
	}
	if err != nil {
		return err
	}
	return fmt.Errorf("error: the value DOES NOT EXIST and SHOULD, or EXIST and SHOULD NOT, "+
		"in '%s' collection", s.collection.Name())
}

// ValidateDataMongo verifies that the feature table data exists in the MongoDB collection
func (s *Session) ValidateDataMongo(
	ctx context.Context, idCollection string, props map[string]interface{},
) (bool, error) {
	// 1-Sets a document to s.singleResult from a filter search
	s.SetSingleResult(ctx, "_id", idCollection)

	// 2-Decodes the singleResult document into a BSON (bsonDoc)
	bsonDoc, err := s.GetDecodeDocument(*s.singleResult)
	if err != nil {
		return false, err
	}

	// 3-Convert the bsonDoc object to JSON
	err = s.SetDataCollectionJSONBytes(bsonDoc)
	if err != nil {
		return false, err
	}

	// 4-Set the list of fields in the document to fieldsCollectionName
	s.SetFieldsCollectionName(ctx, idCollection)
	if err != nil {
		return false, err
	}

	// 5-Navigate through the feature table and check the data
	m := golium.NewMapFromJSONBytes(s.dataCollectionJSONBytes)
	for key, expectedValue := range props {
		value := m.Get(key)
		// Verify that the name of the clean field (fieldTableFeature) exists
		// in the list of fields in the collection (s.fieldsCollectionName)
		fieldTableFeature := strings.Split(key, ".")[0]
		if ContainsElement(fieldTableFeature, s.fieldsCollectionName) {
			if value == nil && expectedValue == "" {
				// The value of the field in MongoDB is null and expectedValue is [EMPTY]
				continue
			} else if value != expectedValue {
				return false, fmt.Errorf("error: mismatch of mongo field '%s': expected '%s',"+
					"actual '%s'", key, expectedValue, value)
			}
		} else {
			return false, fmt.Errorf("error: the field '%s': does not exist in '%s' collection",
				key, s.collection.Name())
		}
	}
	return true, nil
}

// MongoDisconnection closes the connection to MongoDB if it exists
func (s *Session) MongoDisconnection() error {
	if s.client != nil {
		err := s.client.Disconnect(context.Background())
		if err != nil {
			return fmt.Errorf("error: problem in MongoDB disconnection: '%s'", err)
		}
	}
	return nil
}
