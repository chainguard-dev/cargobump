package cargobump

import (
	"reflect"
	"testing"

	"github.com/chainguard-dev/cargobump/pkg/types"
)

func TestParsePackageList(t *testing.T) {
	t.Run("happy-path", func(t *testing.T) {
		got, err := parsePackageList("foo@1.2.3 bar@4.5.6")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		want := map[string]*types.Package{
			"foo": {Name: "foo", Version: "1.2.3"},
			"bar": {Name: "bar", Version: "4.5.6"},
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("mismatch\nwant: %#v\ngot : %#v", want, got)
		}
	})

	t.Run("invalid-format", func(t *testing.T) {
		_, err := parsePackageList("oops-no-at-sign")
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
	})

	t.Run("duplicate-package", func(t *testing.T) {
		_, err := parsePackageList("foo@1.0.0 foo@2.0.0")
		if err == nil {
			t.Fatalf("expected duplicate-package error, got nil")
		}
		want := `duplicate package foo@2.0.0 found, already defined as foo@1.0.0`
		if err.Error() != want {
			t.Fatalf("unexpected error message:\nwant: %q\ngot : %q", want, err.Error())
		}
	})
}
