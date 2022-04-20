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
	"reflect"
	"testing"
)

type Project struct {
	Name      string
	Commiters []string
	Commits   int
}

func TestStringType(t *testing.T) {
	var project = &Project{}
	type args struct {
		destination   reflect.Value
		fieldValueStr string
		value         interface{}
	}

	ctx := context.Background()
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Set string type value",
			args: args{
				destination:   reflect.ValueOf(project).Elem(),
				fieldValueStr: "golium",
				value:         Value(ctx, "golium"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validatedField(&tt.args.destination, "Name")
			if (err != nil) != tt.wantErr {
				t.Errorf("validatedField() error = %v, wantErr %v", err, tt.wantErr)
			}
			err = stringType(tt.args.destination, tt.args.fieldValueStr, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("sliceType() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidatedField(t *testing.T) {
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
			if err := validatedField(&tt.args.destination, tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("validatedField() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
