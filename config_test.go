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

package golium

import (
	"os"
	"reflect"
	"testing"

	"github.com/tidwall/gjson"
)

const (
	logsPath        = "./logs"
	environmentPath = "./environments"
	localConfFile   = `
# Local	
minio: true
minioEndpoint: http://miniomock:9000
`
	privateLocalConfFile = `
# Private	
minio: true
minioEndpoint: http://miniomock:9000
`
)

func Test_initEnvironment(t *testing.T) {
	os.MkdirAll(logsPath, os.ModePerm)
	defer os.RemoveAll(logsPath)
	os.MkdirAll(environmentPath, os.ModePerm)
	defer os.RemoveAll(environmentPath)

	configMap := &gjsonMap{
		gmap: gjson.ParseBytes([]byte(`{"minio":true,"minioEndpoint":"http://miniomock:9000"}`)),
	}

	tests := []struct {
		name               string
		privateConfFileErr bool
		want               Map
	}{
		{
			name:               "Load private environment configuration from file error",
			privateConfFileErr: true,
			want:               configMap,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.WriteFile("./environments/local.yml", []byte(localConfFile), os.ModePerm)

			if !tt.privateConfFileErr {
				os.WriteFile("./environments/local-private.yml", []byte(privateLocalConfFile), os.ModePerm)
			}
			if got := initEnvironment(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("initEnvironment() = %v, want %v", got, tt.want)
			}
		})
	}
}
