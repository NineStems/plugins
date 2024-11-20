package toml

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go-micro.dev/v4/config/source"
)

func TestReader(t *testing.T) {
	data := []byte(`foo="bar"
baz.bar="cat"`)

	testData := []struct {
		name  string
		path  []string
		value string
	}{
		{
			"simple get",
			[]string{"foo"},
			"bar",
		},
		{
			"complex get",
			[]string{"baz", "bar"},
			"cat",
		},
	}

	r := NewReader()

	c, err := r.Merge(&source.ChangeSet{Data: data}, &source.ChangeSet{})
	if err != nil {
		t.Fatal(err)
	}

	values, err := r.Values(c)
	if err != nil {
		t.Fatal(err)
	}

	for _, test := range testData {
		t.Run(test.name, func(t *testing.T) {
			v := values.Get(test.path...).String("")
			require.Equal(t, test.value, v, test.path)
		})
	}
}
