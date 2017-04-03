package kvbuilder

import (
	"bytes"
	"reflect"
	"testing"
)

func TestBuilder_basic(t *testing.T) {
	var b Builder
	err := b.Add("foo=bar", "bar=baz", "baz=")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := map[string]interface{}{
		"foo": "bar",
		"bar": "baz",
		"baz": "",
	}
	actual := b.Map()
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: %#v", actual)
	}
}

func TestBuilder_escapedAt(t *testing.T) {
	var b Builder
	err := b.Add("foo=bar", "bar=\\@baz")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := map[string]interface{}{
		"foo": "bar",
		"bar": "@baz",
	}
	actual := b.Map()
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: %#v", actual)
	}
}

func TestBuilder_stdin(t *testing.T) {
	var b Builder
	b.Stdin = bytes.NewBufferString("baz")
	err := b.Add("foo=bar", "bar=-")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := map[string]interface{}{
		"foo": "bar",
		"bar": "baz",
	}
	actual := b.Map()
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: %#v", actual)
	}
}

func TestBuilder_stdinMap(t *testing.T) {
	var b Builder
	b.Stdin = bytes.NewBufferString(`{"foo": "bar"}`)
	err := b.Add("-", "bar=baz")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := map[string]interface{}{
		"foo": "bar",
		"bar": "baz",
	}
	actual := b.Map()
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: %#v", actual)
	}
}

func TestBuilder_stdinTwice(t *testing.T) {
	var b Builder
	b.Stdin = bytes.NewBufferString(`{"foo": "bar"}`)
	err := b.Add("-", "-")
	if err == nil {
		t.Fatal("should error")
	}
}

func TestBuilder_keyTwiceToCommaString(t *testing.T) {
	var b Builder
	err := b.Add("foo=bar", "foo=baz")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := map[string]interface{}{
		"foo": "bar,baz",
	}
	actual := b.Map()
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: %#v", actual)
	}
}
