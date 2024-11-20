package toml

import (
	"reflect"
	"testing"
)

func Test_line_Header(t *testing.T) {
	tests := []struct {
		name string
		l    line
		want []string
		mark bool
	}{
		{
			name: "only level 1",
			l:    []byte("[l1]"),
			want: []string{"l1"},
			mark: true,
		},
		{
			name: "take last from second level",
			l:    []byte("[l1.l2]"),
			want: []string{"l1", "l2"},
			mark: true,
		},
		{
			name: "take last from third level",
			l:    []byte("[l1.l2.l3]"),
			want: []string{"l1", "l2", "l3"},
			mark: true,
		},
		{
			name: "not a header with array with",
			l:    []byte("numbers=[ 0.1, 0.2, 0.5, 1, 2, 5 ]"),
			want: nil,
			mark: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.l.Header()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Header() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.mark {
				t.Errorf("Header() got1 = %v, want %v", got1, tt.mark)
			}
		})
	}
}
