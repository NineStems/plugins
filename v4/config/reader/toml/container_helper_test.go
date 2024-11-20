package toml

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_addValue(t *testing.T) {
	type args struct {
		headers []string
		cont    container
		in      line
		want    container
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "one level header added",
			args: args{
				headers: []string{"l1"},
				cont:    make(container),
				in:      []byte(`test=1`),
				want: container{
					"l1": &container{
						"test": value{int64(1)},
					},
				},
			},
		},
		{
			name: "two levels header added",
			args: args{
				headers: []string{"l1", "l2"},
				cont:    make(container),
				in:      []byte(`test=1`),
				want: container{
					"l1": &container{
						"l2": &container{
							"test": value{int64(1)},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processedValueLine(tt.args.headers, tt.args.cont, tt.args.in)
			require.Equal(t, tt.args.want, tt.args.cont)
		})
	}
}

func Test_convert(t *testing.T) {
	type args struct {
		in lines
	}
	tests := []struct {
		name string
		args args
		want containerI
	}{
		{
			name: "third level with two values",
			args: args{
				in: lines{
					line(`[db]`),
					line(`[db.postgres]`),
					line(`[db.postgres.book]`),
					line(`host = "localhost"`),
					line(`port = 5432`),
				},
			},
			want: container{
				"db": &container{
					"postgres": &container{
						"book": &container{
							"host": value{"localhost"},
							"port": value{int64(5432)},
						},
					},
				},
			},
		},
		{
			name: "several value on one and several on other",
			args: args{
				in: lines{
					line(`host="localhost"`),
					line(`[db]`),
					line(`name="name"`),
				},
			},
			want: container{
				"host": value{"localhost"},
				"db": &container{
					"name": value{"name"},
				},
			},
		},
		{
			name: "several different values with new line",
			args: args{
				in: lines{
					line(`[db]`),
					line(`name="name"`),
					line(``),
					line(`[mbs]`),
					line(`name="name"`),
				},
			},
			want: container{
				"db": &container{
					"name": value{"name"},
				},
				"mbs": &container{
					"name": value{"name"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := convert(tt.args.in)
			require.Equal(t, tt.want, got)
		})
	}
}

func Test_deleteContainer(t *testing.T) {
	type args struct {
		headers []string
		cont    container
	}
	tests := []struct {
		name string
		args struct {
			headers []string
			cont    container
		}
		want container
	}{
		{
			name: "one level deleting",
			args: args{
				headers: []string{"l1"},
				cont: container{
					"l1": &container{},
				},
			},
			want: container{},
		},
		{
			name: "second level deleting",
			args: args{
				headers: []string{"l1", "l2"},
				cont: container{
					"l1": &container{
						"l2": &container{},
					},
				},
			},
			want: container{
				"l1": &container{},
			},
		},
		{
			name: "third level deleting",
			args: args{
				headers: []string{"l1", "l2", "l3"},
				cont: container{
					"l1": &container{
						"l2": &container{
							"l3": &container{},
						},
					},
				},
			},
			want: container{
				"l1": &container{
					"l2": &container{},
				},
			},
		},

		{
			name: "delete branch",
			args: args{
				headers: []string{"l1", "l2"},
				cont: container{
					"l1": &container{
						"l2": &container{
							"l3": &container{},
						},
					},
				},
			},
			want: container{
				"l1": &container{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deleteContainer(tt.args.headers, tt.args.cont)
			require.Equal(t, tt.want, tt.args.cont)
		})
	}
}

func Test_builderMap(t *testing.T) {
	type args struct {
		in container
	}
	tests := []struct {
		name string
		args args
		want map[string]any
	}{
		{
			name: "happy path",
			args: args{
				container{
					"l1": &container{
						"head": value{"test"},
					},
				},
			},
			want: map[string]any{
				"l1": map[string]any{
					"head": "test",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, builderMap(tt.args.in))
		})
	}
}
