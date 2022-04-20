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
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
)

// KindFormatter with format to apply conversion
type KindFormatter interface {
	format(destination reflect.Value, fieldValueStr string, value interface{}) error
}

// ConfigurePattern func to apply transformation into destination
type ConfigurePattern func(destination reflect.Value, fieldValueStr string, value interface{}) error

// FieldConversor
type FieldConversor struct {
	Pattern ConfigurePattern
}

// format applies defined pattern to format destination values
func (c FieldConversor) format(
	destination reflect.Value,
	fieldValueStr string,
	value interface{}) error {
	return c.Pattern(destination, fieldValueStr, value)
}

var StrategyFormat = map[reflect.Kind]*FieldConversor{
	reflect.Slice:      NewSliceConverter(),
	reflect.String:     NewStringConverter(),
	reflect.Bool:       NewBoolConverter(),
	reflect.Int:        NewInt64Converter(),
	reflect.Int8:       NewInt64Converter(),
	reflect.Int16:      NewInt64Converter(),
	reflect.Int32:      NewInt64Converter(),
	reflect.Int64:      NewInt64Converter(),
	reflect.Uint:       NewUInt64Converter(),
	reflect.Uint8:      NewInt64Converter(),
	reflect.Uint16:     NewUInt64Converter(),
	reflect.Uint32:     NewUInt64Converter(),
	reflect.Uint64:     NewUInt64Converter(),
	reflect.Float32:    NewFloat64Converter(),
	reflect.Float64:    NewFloat64Converter(),
	reflect.Complex64:  NewComplex64Converter(),
	reflect.Complex128: NewComplex64Converter(),
}

// NewSliceConverter Constructor
func NewSliceConverter() *FieldConversor {
	return &FieldConversor{Pattern: sliceType}
}

// sliceType conversion pattern to appy
func sliceType(destination reflect.Value, fieldValueStr string, value interface{}) error {
	array, ok := value.([]gjson.Result)
	if !ok {
		return fmt.Errorf(
			"failed parsing destination '%v', not a JSON array",
			value)
	}
	length := len(array)
	var fv reflect.Value
	if length > 0 {
		fv = makeSlice(array[0], length)
		for i, v := range array {
			setSliceValue(fv.Index(i), v)
		}
	}
	destination.Set(fv)
	return nil
}

// NewStringConverter Constructor
func NewStringConverter() *FieldConversor {
	return &FieldConversor{Pattern: stringType}
}

// stringType conversion pattern to appy
func stringType(destination reflect.Value, fieldValueStr string, value interface{}) error {
	destination.SetString(fieldValueStr)
	return nil
}

// NewBoolConverter Constructor
func NewBoolConverter() *FieldConversor {
	return &FieldConversor{Pattern: boolType}
}

// boolType conversion pattern to appy
func boolType(destination reflect.Value, fieldValueStr string, value interface{}) error {
	v, err := strconv.ParseBool(fieldValueStr)
	if err != nil {
		return fmt.Errorf("failed parsing to boolean the value '%s'",
			fieldValueStr)
	}
	destination.SetBool(v)
	return nil
}

// NewInt64Converter Constructor
func NewInt64Converter() *FieldConversor {
	return &FieldConversor{Pattern: int64Type}
}

// int64Type conversion pattern to appy
func int64Type(destination reflect.Value, fieldValueStr string, value interface{}) error {
	v, err := strconv.ParseInt(fieldValueStr, 10, 64)
	if err != nil {
		return fmt.Errorf("failed parsing to integer '%s' with destination '%v'",
			fieldValueStr,
			value)
	}
	destination.SetInt(v)
	return nil
}

// NewUInt64Converter Constructor
func NewUInt64Converter() *FieldConversor {
	return &FieldConversor{Pattern: uint64Type}
}

// uint64Type conversion pattern to appy
func uint64Type(destination reflect.Value, fieldValueStr string, value interface{}) error {
	v, err := strconv.ParseUint(fieldValueStr, 10, 64)
	if err != nil {
		return fmt.Errorf("failed parsing to unsigned integer '%s' with destination '%v'",
			fieldValueStr,
			value)
	}
	destination.SetUint(v)
	return nil
}

// NewFloat64Converter Constructor
func NewFloat64Converter() *FieldConversor {
	return &FieldConversor{Pattern: float64Type}
}

// float64Type conversion pattern to appy
func float64Type(destination reflect.Value, fieldValueStr string, value interface{}) error {
	v, err := strconv.ParseFloat(fieldValueStr, 64)
	if err != nil {
		return fmt.Errorf("failed parsing to float '%s' with destination '%v'",
			fieldValueStr,
			value)
	}
	destination.SetFloat(v)
	return nil
}

// NewComplex64Converter Constructor
func NewComplex64Converter() *FieldConversor {
	return &FieldConversor{Pattern: complex64Type}
}

// complex64Type conversion pattern to appy
func complex64Type(destination reflect.Value, fieldValueStr string, value interface{}) error {
	v, err := strconv.ParseComplex(fieldValueStr, 128)
	if err != nil {
		return fmt.Errorf("failed parsing to complex '%s' with destination '%v'",
			fieldValueStr,
			value)
	}
	destination.SetComplex(v)
	return nil
}

// validatedField to apply pattern
func validatedField(destination *reflect.Value, name string) error {
	if destination.Kind() == reflect.Ptr {
		*destination = reflect.Indirect(*destination)
	}
	if destination.Kind() != reflect.Struct {
		return fmt.Errorf("destination must be a struct")
	}
	*destination = destination.FieldByNameFunc(func(n string) bool {
		return strings.EqualFold(n, name)
	})
	if !destination.IsValid() {
		return fmt.Errorf("field '%s' is not valid", name)
	}
	if !destination.CanSet() {
		return fmt.Errorf("field '%s' cannot be set", name)
	}
	if destination.Kind() == reflect.Ptr {
		fv := reflect.New(destination.Type().Elem())
		destination.Set(fv)
		*destination = fv.Elem()
	}
	return nil
}

// makeSlice from element with selected length
func makeSlice(element gjson.Result, length int) reflect.Value {
	var rv reflect.Value
	switch element.Type {
	case gjson.False, gjson.True:
		var b bool
		rv = reflect.ValueOf(b)
	case gjson.Number:
		var i int
		rv = reflect.ValueOf(i)
	case gjson.String, gjson.JSON, gjson.Null:
		var s string
		rv = reflect.ValueOf(s)
	}
	return reflect.MakeSlice(reflect.SliceOf(rv.Type()), length, length)
}

// setSliceValue to field with value
func setSliceValue(field reflect.Value, value gjson.Result) {
	switch value.Type {
	case gjson.False, gjson.True:
		field.Set(reflect.ValueOf(value.Bool()))
	case gjson.Number:
		field.Set(reflect.ValueOf(value.Int()))
	case gjson.String, gjson.JSON, gjson.Null:
		field.Set(reflect.ValueOf(value.String()))
	}
}
