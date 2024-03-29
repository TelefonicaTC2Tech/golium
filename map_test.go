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
	"testing"
)

const values = `
{
	"name":{
		"first": "John",
		"last": "Doe"
	},
	"age":47,
	"commiter": true,
	"developer": false,
	"projects": [
		"golium",
		"toolium"
	]
}`

func TestGet(t *testing.T) {
	type args struct {
		path string
	}

	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			name: "Get string value from json map",
			args: args{path: "name.first"},
			want: "John",
		},
		{
			name: "Get number value from json map",
			args: args{path: "age"},
			want: float64(47),
		},
		{
			name: "Get nil value from json map",
			args: args{path: "none"},
			want: nil,
		},
		{
			name: "Get true value from json map",
			args: args{path: "commiter"},
			want: true,
		},
		{
			name: "Get false value from json map",
			args: args{path: "developer"},
			want: false,
		},
		{
			name: "Get nil value from json map",
			args: args{path: "nil"},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMapFromJSONBytes([]byte(values))
			if got := m.Get(tt.args.path); got != tt.want {
				t.Errorf("gjsonMap.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}
