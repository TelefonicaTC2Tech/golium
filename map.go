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
	"github.com/tidwall/gjson"
)

// Map is an interface to get elements from a data structure with a dot notation.
type Map interface {
	Get(path string) interface{}
}

// gjsonMap provides methods to get elements from a JSON document with a dot notation.
type gjsonMap struct {
	gmap gjson.Result
}

// NewMapFromJSONBytes creates a Map from a slice of bytes of a JSON document.
func NewMapFromJSONBytes(buf []byte) Map {
	return &gjsonMap{
		gmap: gjson.ParseBytes(buf),
	}
}

// Get an element from the map by a path with dot notation.
func (m *gjsonMap) Get(path string) interface{} {
	result := m.gmap.Get(path)
	switch result.Type {
	case gjson.String:
		return result.String()
	case gjson.Null:
		return nil
	case gjson.True:
		return true
	case gjson.False:
		return false
	case gjson.Number:
		return result.Float()
	default:
		if result.IsArray() {
			return result.Array()
		}
		return result.String()
	}
}
