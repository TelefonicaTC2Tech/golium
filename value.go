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
	"fmt"
	"strconv"
	"strings"
)

// ValueAsString invokes Value and convert the return value to string.
func ValueAsString(ctx context.Context, s string) string {
	return fmt.Sprintf("%s", Value(ctx, s))
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
//
// Most cases, the return value is a string except for the following cases:
// - [TRUE] and [FALSE] return a bool type.
// - [NUMBER:1234] returns a float64 if s only contains this tag and there is no surrounding text.
func Value(ctx context.Context, s string) interface{} {
	switch s {
	case "[TRUE]":
		return true
	case "[FALSE]":
		return false
	case "[NULL]":
		return nil
	case "[EMPTY]":
		return ""
	default:
		orig := s
		s = processTag(s, "CONF", func(tagName string) string {
			m := GetEnvironment()
			return fmt.Sprintf("%s", m.Get(tagName))
		})
		s = processTag(s, "CTXT", func(tagName string) string {
			return fmt.Sprintf("%s", GetContext(ctx).Get(tagName))
		})
		s = processTag(s, "SHA256", func(tagName string) string {
			return fmt.Sprintf("%x", sha256.Sum256([]byte(tagName)))
		})
		s = processTag(s, "NUMBER", func(tagName string) string {
			return tagName
		})
		// If there is only a NUMBER tag, without any surrounding text, then return a float number
		if orig == fmt.Sprintf("[NUMBER:%s]", s) {
			if v, err := strconv.ParseFloat(s, 64); err == nil {
				return v
			}
		}
		return s
	}
}

func processTag(s string, tag string, getTagValue func(string) string) string {
	tagNames := getTagNames(s, tag)
	for _, tagName := range tagNames {
		token := fmt.Sprintf("[%s:%s]", tag, tagName)
		s = strings.ReplaceAll(s, token, getTagValue(tagName))
	}
	return s
}

func getTagNames(s string, tag string) []string {
	tagNames := []string{}
	tokens := strings.Split(s, fmt.Sprintf("[%s:", tag))
	if len(tokens) > 1 {
		for _, token := range tokens[1:] {
			n := strings.Index(token, "]")
			if n < 1 {
				continue
			}
			tagName := token[:n]
			tagNames = append(tagNames, tagName)
		}
	}
	return tagNames
}
