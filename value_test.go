package golium_test

import (
	"context"
	"os"
	"testing"

	"github.com/TelefonicaTC2Tech/golium"
)

const (
	environmentPath = "./environments"
	localConfFile   = `
minio: true
minioEndpoint: http://miniomock:9000
`
	testSHA256 = "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08"
	testBASE64 = "dGVzdA=="
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

func TestValuedTag(t *testing.T) {
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
	os.MkdirAll(environmentPath, os.ModePerm)
	defer os.RemoveAll(environmentPath)

	os.WriteFile("./environments/local.yml", []byte(localConfFile), os.ModePerm)
	tcs := map[string]interface{}{
		"[CONF:minio]":  true,
		"[NUMBER:1]":    float64(1),
		"[SHA256:test]": testSHA256,
		"[BASE64:test]": testBASE64,
		"[CTXT:test]":   "contextTest",
	}

	for s, expectedValue := range tcs {
		ctx := golium.InitializeContext(context.Background())
		if s == "[CTXT:test]" {
			golium.GetContext(ctx).Put(golium.ValueAsString(ctx, "test"), "contextTest")
		}
		v := golium.NewComposedTag(s).Value(ctx)
		if v != expectedValue {
			t.Errorf("expected: %s, actual: %s", expectedValue, v)
		}
	}
}

func TestTagWithoutChecks(t *testing.T) {
	tcs := []string{
		"[NOW]",
		"[UUID]",
		"[NOW:1s]",
		"[NOW:-1:]",
		"[NOW:300ms:]",
		"[NOW::]",
		"[NOW::unix]",
	}
	ctx := context.Background()
	for _, s := range tcs {
		_ = golium.NewComposedTag(s).Value(ctx)
	}
}
