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
)

type ClientFunctions interface {

	Ping(ctx context.Context, client *mongo.Client) error
	
}
type ClientService struct{}

func NewMongoClientService() *ClientService {
	return &ClientService{}
}

//Ping check that there is a connection to MongoDB
func (c ClientService) Ping(ctx context.Context, client *mongo.Client) error {
	return client.Ping(context.Background(), nil)
}
