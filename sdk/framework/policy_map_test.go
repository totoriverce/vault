package framework

import (
	"context"
	"reflect"
	"testing"

	"github.com/hashicorp/vault/sdk/logical"
)

func TestPolicyMap(t *testing.T) {
	p := &PolicyMap{}
	p.PathMap.Name = "foo"
	s := new(logical.InmemStorage)

	ctx := context.Background()

	p.Put(ctx, s, "foo", map[string]any{"value": "bar"})
	p.Put(ctx, s, "bar", map[string]any{"value": "foo,baz "})

	// Read via API
	actual, err := p.Policies(ctx, s, "foo", "bar")
	if err != nil {
		t.Fatalf("bad: %#v", err)
	}

	expected := []string{"bar", "baz", "foo"}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: %#v", actual)
	}
}
