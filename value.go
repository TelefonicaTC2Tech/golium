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
	"crypto/sha256"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// ValueAsString invokes Value and convert the return value to string.
func ValueAsString(ctx context.Context, s string) string {
	return fmt.Sprintf("%v", Value(ctx, s))
}

// ValueAsInt invokes Value and convert the return value to int.
func ValueAsInt(ctx context.Context, s string) (int, error) {
	v := Value(ctx, s)
	if n, ok := v.(float64); ok {
		return int(n), nil
	}
	return strconv.Atoi(s)
}

// Value converts a value as a string to consider some golium patterns.
// Supported patterns:
// - Booleans: [TRUE] or [FALSE]
// - Null value: [NULL]
// - Empty value: [EMPTY]
// - Number: [NUMBER:1234] or [NUMBER:1234.67]
// - Configuration parameters: [CONF:test.parameter]
// - Context values: [CTXT:test.context]
// - SHA256: [SHA256:text.to.be.hashed]
// - Time: [NOW:+24h:unix] with the format: [NOW:{duration}:{format}]
//   The value {duration} can be empty (there is no change from now timestamp) or a format valid for
//   time.ParseDuration function. Currently, it supports the following units: "ns", "us", "ms", "s", "m", "h".
//   The format can be "unix" or a layout valid for time.Format function.
//   It is possible to use [NOW]. In that case, it returns an int64 with the now timestamp in unix format.
//
// Most cases, the return value is a string except for the following cases:
// - [TRUE] and [FALSE] return a bool type.
// - [NUMBER:1234] returns a float64 if s only contains this tag and there is no surrounding text.
// - [NOW:{duration}:{format}] returns an int64 when {format} is "unix".
func Value(ctx context.Context, s string) interface{} {
	composedTag := NewComposedTag(s)
	return composedTag.Value(ctx)
}

var simpleTagFuncs = map[string]func(ctx context.Context) (interface{}, error){
	"TRUE":  func(ctx context.Context) (interface{}, error) { return true, nil },
	"FALSE": func(ctx context.Context) (interface{}, error) { return false, nil },
	"EMPTY": func(ctx context.Context) (interface{}, error) { return "", nil },
	"NOW":   func(ctx context.Context) (interface{}, error) { return time.Now().Unix(), nil },
	"NULL":  func(ctx context.Context) (interface{}, error) { return nil, nil },
	"UUID": func(ctx context.Context) (interface{}, error) {
		guid, err := uuid.NewUUID()
		if err != nil {
			return "", err
		}
		return guid.String(), nil
	},
}

var valuedTagFuncs = map[string]func(ctx context.Context, s string) (interface{}, error){
	"CONF": func(ctx context.Context, s string) (interface{}, error) {
		m := GetEnvironment()
		return m.Get(s), nil
	},
	"CTXT": func(ctx context.Context, s string) (interface{}, error) {
		return GetContext(ctx).Get(s), nil
	},
	"SHA256": func(ctx context.Context, s string) (interface{}, error) {
		return fmt.Sprintf("%x", sha256.Sum256([]byte(s))), nil
	},
	"NUMBER": func(ctx context.Context, s string) (interface{}, error) {
		return strconv.ParseFloat(s, 64)
	},
	"NOW": func(ctx context.Context, s string) (interface{}, error) {
		return processNow(s)
	},
}

// processNow processes tag "NOW" with the format [NOW:{duration}:{format}].
// So, tagName has the format: {duration}:{format}
func processNow(s string) (interface{}, error) {
	parts := strings.SplitN(s, ":", 2)
	if len(parts) != 2 {
		return nil, errors.New("invalid NOW tag")
	}
	duration := parts[0]
	format := parts[1]
	now := time.Now()
	if duration != "" {
		d, err := time.ParseDuration(duration)
		if err != nil {
			return nil, fmt.Errorf("invalid duration in NOW tag: %w", err)
		}
		now = now.Add(d)
	}
	switch format {
	case "unix":
		return now.Unix(), nil
	default:
		return now.Format(format), nil
	}
}

// Tag interface to calculate the value of a tag.
// A golium tag is a text surrounded by brackets that can be evaluated into a value.
// For example: [CONF:property]
type Tag interface {
	Value(ctx context.Context) interface{}
}

// StringTag represents a implicit tag composed of a text.
// This tag is used to compose a string with a tag to generate a new string.
type StringTag struct {
	s string
}

// NewStringTag creates a Tag that evaluated to the string without any modification.
func NewStringTag(s string) Tag {
	return &StringTag{s: s}
}

func (s StringTag) Value(ctx context.Context) interface{} {
	return s.s
}

// NamedTag is a Tag that can be evaluated with a tag function depending on the name of the tag.
type NamedTag struct {
	s string
}

// NewNamedTag creates a NamedTag.
func NewNamedTag(s string) Tag {
	return &NamedTag{s: s}
}

func (t NamedTag) Value(ctx context.Context) interface{} {
	if v, err := t.valueWithError(ctx); err == nil {
		return v
	}
	return t.s
}

func (t NamedTag) valueWithError(ctx context.Context) (interface{}, error) {
	tag := t.s[1 : len(t.s)-1]
	parts := strings.SplitN(tag, ":", 2)
	tagName := parts[0]
	if len(parts) == 2 {
		tagValue := parts[1]
		return t.processValuedTag(ctx, tagName, tagValue)
	}
	return t.processSimpleTag(ctx, tagName)
}

func (t NamedTag) processSimpleTag(ctx context.Context, tagName string) (interface{}, error) {
	if f, ok := simpleTagFuncs[tagName]; ok {
		return f(ctx)
	}
	return nil, fmt.Errorf("invalid tag %s", tagName)
}

func (t NamedTag) processValuedTag(ctx context.Context, tagName, tagValue string) (interface{}, error) {
	if f, ok := valuedTagFuncs[tagName]; ok {
		composedTag := NewComposedTag(tagValue)
		composedTagValue := composedTag.Value(ctx)
		composedTagValueString := fmt.Sprintf("%v", composedTagValue)
		return f(ctx, composedTagValueString)
	}
	return nil, fmt.Errorf("invalid tag %s", tagName)
}

type separator struct {
	opener bool
	pos    int
}

// ComposedTag is a composition of tags, including StringTags, NamedTags and other ComposedTags
// to provide an evaluation.
type ComposedTag struct {
	s string
}

// NewComposedTag creates a ComposedTag.
func NewComposedTag(s string) Tag {
	return &ComposedTag{s: s}
}

func (t ComposedTag) findSeparators() (separators []separator) {
	for i, c := range t.s {
		if c == '[' {
			sep := separator{opener: true, pos: i}
			separators = append(separators, sep)
		} else if c == ']' {
			sep := separator{opener: false, pos: i}
			separators = append(separators, sep)
		}
	}
	return
}

func (t ComposedTag) buildTags(separators []separator) []Tag {
	tags := []Tag{}
	if len(separators) < 2 {
		return tags
	}
	lastCloser := -1
	for i := 0; i < len(separators)-1; i++ {
		if !separators[i].opener {
			// Discard it because we must start with an opener
			continue
		}
		distance := 1
		for j := i + 1; j < len(separators); j++ {
			if separators[j].opener {
				distance++
			} else {
				distance--
			}
			if distance != 0 {
				continue
			}
			opener := separators[i].pos
			closer := separators[j].pos
			// Add a tag text if there is a text prefix
			if lastCloser+1 < opener {
				tag := NewStringTag(t.s[lastCloser+1 : opener])
				tags = append(tags, tag)
			}
			// Found end of tag
			tag := NewNamedTag(t.s[opener : closer+1])
			tags = append(tags, tag)
			i = j
			lastCloser = closer
			break
		}
	}
	// Add a tag text if there is a text suffix
	if lastCloser+1 < len(t.s) {
		tag := NewStringTag(t.s[lastCloser+1:])
		tags = append(tags, tag)
	}
	return tags
}

func (t ComposedTag) Value(ctx context.Context) interface{} {
	tags := t.buildTags(t.findSeparators())
	if len(tags) == 0 {
		return t.s
	}
	if len(tags) == 1 {
		return tags[0].Value(ctx)
	}
	// If multiple tags, it returns a string with the concatenation of each tag value
	var v strings.Builder
	for _, tag := range tags {
		v.WriteString(fmt.Sprintf("%v", tag.Value(ctx)))
	}
	return v.String()
}
