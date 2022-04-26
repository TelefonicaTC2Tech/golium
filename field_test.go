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
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/tidwall/gjson"
)

type Project struct {
	Name       string
	Commiters  []string
	Commits    int
	Archived   bool
	Stars      float64
	Branches   uint64
	Complexity complex128
}

type typeParser struct {
	kind        reflect.Kind
	destination reflect.Value
	name        string
	fieldValue  interface{}
	value       interface{}
}

type TestArg struct {
	name    string
	args    typeParser
	wantErr bool
}

var testArgs []TestArg

// typeParserBuilder to create test input values
func typeParserBuilder(ctx context.Context,
	kind reflect.Kind,
	name string,
	fieldValue interface{},
) typeParser {
	return typeParser{
		kind:        kind,
		destination: reflect.ValueOf(&Project{}).Elem(),
		name:        name,
		fieldValue:  fieldValue,
		value:       Value(ctx, fmt.Sprintf("%v", fieldValue)),
	}
}

func typeSlice(
	name string,
	fieldValue interface{},
) typeParser {
	return typeParser{
		kind:        reflect.Slice,
		destination: reflect.ValueOf(&Project{}).Elem(),
		name:        name,
		value:       fieldValue,
	}
}

func setup() {
	var ctx = context.Background()
	gjsonArg := gjson.Result{
		Index:   140,
		Indexes: []int{},
		Num:     float64(0),
		Raw:     "[ricardogarfe, jordipuigbou, sarmar11]",
		Str:     "[ricardogarfe, jordipuigbou, sarmar11]",
		Type:    3}
	var results []gjson.Result
	results = append(results, gjsonArg)
	testArgs = []TestArg{
		{
			name:    "Format to set array slice type value",
			args:    typeSlice("Commiters", results),
			wantErr: false,
		},
		{
			name:    "Format error when set array slice type value",
			args:    typeSlice("Commiters", "[results]"),
			wantErr: true,
		},
		{
			name:    "Format error when set a non bool type value",
			args:    typeParserBuilder(ctx, reflect.Bool, "Archived", "not a bool"),
			wantErr: true,
		},
		{
			name:    "Format to set float64 type value",
			args:    typeParserBuilder(ctx, reflect.Float64, "Stars", float64(34)),
			wantErr: false,
		},
		{
			name:    "Format error when set a non float64 type value",
			args:    typeParserBuilder(ctx, reflect.Float64, "Stars", "not a float64"),
			wantErr: true,
		},
		{
			name:    "Format to set uint64 type value",
			args:    typeParserBuilder(ctx, reflect.Uint64, "Branches", uint64(55)),
			wantErr: false,
		},
		{
			name:    "Format error when set a non uint64 type value",
			args:    typeParserBuilder(ctx, reflect.Uint64, "Branches", "not a uint64"),
			wantErr: true,
		},
		{
			name:    "Format to set complex64 type value",
			args:    typeParserBuilder(ctx, reflect.Complex64, "Complexity", complex64(100)),
			wantErr: false,
		},
		{
			name:    "Format error when set a non complex64 type value",
			args:    typeParserBuilder(ctx, reflect.Complex64, "Complexity", "not a complex"),
			wantErr: true,
		},
	}
}

func TestStrategyFormatTypes(t *testing.T) {
	setup()
	for _, tt := range testArgs {
		t.Run(tt.name, func(t *testing.T) {
			err := exctractField(&tt.args.destination, tt.args.name)
			if err != nil {
				t.Errorf("exctractField() error = %v, wantErr %v", err, tt.wantErr)
			}
			err = StrategyFormat[tt.args.kind].
				format(tt.args.destination, fmt.Sprint(tt.args.fieldValue), tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("format for kind %s error = %v, wantErr %v", tt.args.kind, err, tt.wantErr)
			}
		})
	}
}

func TestExtractField(t *testing.T) {
	var project = &Project{}
	type args struct {
		destination reflect.Value
		name        string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Validate existing field",
			args: args{
				destination: reflect.ValueOf(project).Elem(),
				name:        "Name"},
			wantErr: false,
		},
		{
			name: "Validate non existing field",
			args: args{
				destination: reflect.ValueOf(project).Elem(),
				name:        "Language"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := exctractField(&tt.args.destination, tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("exctractField() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
