package pathutil

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJoin(t *testing.T) {
	tcs := []struct {
		elems []string
		out   string
	}{
		{[]string{"foo", "bar", "baz.hogi"}, "foo/bar/baz.hogi"},
		{[]string{"/foo", "bar", "baz.hogi"}, "/foo/bar/baz.hogi"},
		{[]string{"../foo", "bar", "baz.hogi"}, "/foo/bar/baz.hogi"},
		{[]string{"./foo", "bar", "baz.hogi"}, "foo/bar/baz.hogi"},
		{[]string{"foo", "/bar", "baz.hogi"}, "foo/bar/baz.hogi"},
		{[]string{"foo", "../bar", "baz.hogi"}, "foo/bar/baz.hogi"},
		{[]string{"foo", "./bar", "baz.hogi"}, "foo/bar/baz.hogi"},
		{[]string{"foo", "bar", "/baz.hogi"}, "foo/bar/baz.hogi"},
		{[]string{"foo", "bar", "../baz.hogi"}, "foo/bar/baz.hogi"},
		{[]string{"foo", "bar", "./baz.hogi"}, "foo/bar/baz.hogi"},
		{[]string{"/", "foo", "bar", "baz.hogi"}, "/foo/bar/baz.hogi"},
		{[]string{"../", "foo", "bar", "baz.hogi"}, "/foo/bar/baz.hogi"},
		{[]string{"./", "foo", "bar", "baz.hogi"}, "foo/bar/baz.hogi"},
	}

	for _, tc := range tcs {
		t.Run(strings.Join(tc.elems, ","), func(t *testing.T) {
			assert.Equal(t, tc.out, Join(tc.elems...))
		})
	}
}
