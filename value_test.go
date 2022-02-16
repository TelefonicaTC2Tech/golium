package golium_test

import (
	"context"
	"testing"

	"github.com/TelefonicaTC2Tech/golium"
)

func TestStringTag(t *testing.T) {
	s := "test string tag"
	tag := golium.NewStringTag(s)
	v := tag.Value(context.Background())
	if v != s {
		t.Errorf("expected: %s, actual: %s", s, v)
	}
}

func TestSimpleNamedTag(t *testing.T) {
	tcs := map[string]interface{}{
		"[TRUE]":  true,
		"[FALSE]": false,
		"[NULL]":  nil,
		"[EMPTY]": "",
	}
	ctx := context.Background()
	for s, expectedValue := range tcs {
		v := golium.NewNamedTag(s).Value(ctx)
		if v != expectedValue {
			t.Errorf("expected: %s, actual: %s", expectedValue, v)
		}
	}
}

func TestComposedTag(t *testing.T) {
	tcs := map[string]string{
		"This is a test: [TRUE]":         "This is a test: true",
		"[TRUE]: This is a test":         "true: This is a test",
		"[TRUE]: This is a test.[FALSE]": "true: This is a test.false",
		"This [EMPTY]boolean is [TRUE].": "This boolean is true.",
	}
	ctx := context.Background()
	for s, expectedValue := range tcs {
		v := golium.NewComposedTag(s).Value(ctx)
		if v != expectedValue {
			t.Errorf("expected: %s, actual: %s", expectedValue, v)
		}
	}
}
