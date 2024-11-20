package toml

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"go-micro.dev/v4/config/reader"
	"go-micro.dev/v4/config/source"
)

type tomlValues struct {
	ch *source.ChangeSet
	sj containerI
	mu sync.RWMutex
}

type tomlValue struct {
	any
}

func newValues(ch *source.ChangeSet) (reader.Values, error) {
	data, _ := reader.ReplaceEnvVars(ch.Data)
	split(data)
	sj := convert(split(data))
	return &tomlValues{ch, sj, sync.RWMutex{}}, nil
}

func (j *tomlValues) Get(path ...string) reader.Value {
	j.mu.RLock()
	defer j.mu.RUnlock()

	return &tomlValue{j.sj.Path(path...).Value()}
}

func (j *tomlValues) Del(path ...string) {
	j.mu.Lock()
	defer j.mu.Unlock()

	j.sj.Delete(path...)
}

func (j *tomlValues) Set(val any, path ...string) {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.sj.Set(val, path...)
}

func (j *tomlValues) Bytes() []byte {
	j.mu.RLock()
	defer j.mu.RUnlock()

	b, _ := json.Marshal(j.ch)
	return b
}

func (j *tomlValues) Map() map[string]any {
	j.mu.RLock()
	defer j.mu.RUnlock()

	return j.sj.Map()
}

func (j *tomlValues) Scan(v interface{}) error {
	j.mu.RLock()
	defer j.mu.RUnlock()

	b, err := json.Marshal(j.sj.Map())
	if err != nil {
		return err
	}
	return json.Unmarshal(b, v)
}

func (j *tomlValues) String() string {
	j.mu.RLock()
	defer j.mu.RUnlock()

	return "toml"
}

func (j *tomlValue) Bool(def bool) bool {
	str, ok := j.any.(string)
	if !ok {
		return def
	}

	b, err := strconv.ParseBool(str)
	if err != nil {
		return def
	}

	return b
}

func (j *tomlValue) Int(def int) int {
	str, ok := j.any.(string)
	if !ok {
		return def
	}

	i, err := strconv.Atoi(str)
	if err != nil {
		return def
	}

	return i
}

func (j *tomlValue) String(def string) string {
	str, ok := j.any.(string)
	if !ok {
		return def
	}
	return str
}

func (j *tomlValue) Float64(def float64) float64 {
	str, ok := j.any.(string)
	if !ok {
		return def
	}

	f, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return def
	}

	return f
}

func (j *tomlValue) Duration(def time.Duration) time.Duration {
	str, ok := j.any.(string)
	if !ok {
		return def
	}

	v, err := time.ParseDuration(str)
	if err != nil {
		return def
	}

	return v
}

func (j *tomlValue) StringSlice(def []string) []string {
	sla, ok := j.any.([]any)
	if ok {
		sl := make([]string, 0, len(sla))
		for _, s := range sla {
			sl = append(sl, fmt.Sprintf("%v", s))
		}
		return sl
	}

	str, ok := j.any.(string)
	if !ok {
		return def
	}

	sl := strings.Split(str, ",")
	if len(sl) > 0 {
		return sl
	}

	return sl
}

func (j *tomlValue) StringMap(def map[string]string) map[string]string {
	m, ok := (j.any).(map[string]any)
	if !ok {
		return def
	}

	res := map[string]string{}

	for k, v := range m {
		res[k] = fmt.Sprintf("%v", v)
	}

	return res
}

func (j *tomlValue) Scan(v interface{}) error {
	b, err := json.Marshal(j.any)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, v)
}

func (j *tomlValue) Bytes() []byte {
	b, err := json.Marshal(j.any)
	if err != nil {
		return []byte{}
	}
	return b
}
