package toml

import (
	"errors"
	"time"

	"github.com/go-micro/plugins/v4/config/encoder/toml"
	"github.com/imdario/mergo"
	"go-micro.dev/v4/config/encoder"
	"go-micro.dev/v4/config/reader"
	"go-micro.dev/v4/config/source"
)

type tomlReader struct {
	opts reader.Options
	toml encoder.Encoder
}

func (j *tomlReader) Merge(changes ...*source.ChangeSet) (*source.ChangeSet, error) {
	var merged map[string]interface{}

	for _, m := range changes {
		if m == nil {
			continue
		}

		if len(m.Data) == 0 {
			continue
		}

		codec, ok := j.opts.Encoding[m.Format]
		if !ok {
			codec = j.toml
		}

		var data map[string]interface{}
		if err := codec.Decode(m.Data, &data); err != nil {
			return nil, err
		}
		if err := mergo.Map(&merged, data, mergo.WithOverride); err != nil {
			return nil, err
		}
	}

	b, err := j.toml.Encode(merged)
	if err != nil {
		return nil, err
	}

	cs := &source.ChangeSet{
		Timestamp: time.Now(),
		Data:      b,
		Source:    "toml",
		Format:    j.toml.String(),
	}
	cs.Checksum = cs.Sum()

	return cs, nil
}

func (j *tomlReader) Values(ch *source.ChangeSet) (reader.Values, error) {
	if ch == nil {
		return nil, errors.New("changeset is nil")
	}
	if ch.Format != "toml" {
		return nil, errors.New("unsupported format")
	}
	return newValues(ch)
}

func (j *tomlReader) String() string {
	return "toml"
}

func NewOptions(opts ...reader.Option) reader.Options {
	options := reader.Options{
		Encoding: map[string]encoder.Encoder{
			"toml": toml.NewEncoder(),
		},
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}

// NewReader creates a toml reader.
func NewReader(opts ...reader.Option) reader.Reader {
	options := NewOptions(opts...)
	return &tomlReader{
		toml: toml.NewEncoder(),
		opts: options,
	}
}
