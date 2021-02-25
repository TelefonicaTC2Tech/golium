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
// - Configuration parameters: [CONF:test.parameter]
// - Context values: [CTXT:test.context]
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
