package toml

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go-micro.dev/v4/config/source"
)

func TestValues(t *testing.T) {
	emptyStr := ""
	raw := `foo="bar"
[baz]
	bar="cat"`
	testData := []struct {
		name     string
		csdata   []byte
		path     []string
		accepter interface{}
		value    interface{}
	}{
		{
			"simple scan",
			[]byte(raw),
			[]string{"foo"},
			emptyStr,
			"bar",
		},
		{
			"complex scan",
			[]byte(raw),
			[]string{"baz", "bar"},
			emptyStr,
			"cat",
		},
	}

	for _, test := range testData {
		t.Run(test.name, func(t *testing.T) {
			values, err := newValues(
				&source.ChangeSet{
					Data: test.csdata,
				},
			)
			require.NoError(t, err)
			v := values.Get(test.path...)
			err = v.Scan(&test.accepter)
			require.NoError(t, err)
			require.Equal(t, test.value, test.value)
		})
	}
}
