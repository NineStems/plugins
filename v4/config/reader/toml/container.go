package toml

// containerI describe method working with tree.
type containerI interface {
	Path(path ...string) containerI
	Value() any

	Set(val any, path ...string)
	Delete(path ...string)
	Map() map[string]any

	Type() string
}

const (
	containerType = "container"
	valueType     = "value"
)

// container alias for map of map/values.
type container map[string]containerI

// Path return container by his path.
func (c container) Path(path ...string) containerI {
	current := c
	var result containerI = current
	for idx, header := range path {
		ci, ok := current[header]
		if !ok {
			return nil
		}

		if ci.Type() == containerType {
			cc := ci.(*container)
			current = *cc
			result = ci

			continue
		}

		if ci.Type() == valueType && idx != len(path)-1 {
			return value{}
		}

		if ci.Type() == valueType && idx == len(path)-1 {
			result = ci
			continue
		}
	}
	return result
}

func (c container) Value() any {
	return c
}

func (c container) Set(val any, path ...string) {
	key, path := calculatePath(path)
	addPath(path, c, key, val)
}

func (c container) Delete(path ...string) {
	deleteContainer(path, c)
}

func (c container) Map() map[string]any {
	return builderMap(c)
}

func (c container) Type() string {
	return containerType
}

type value struct {
	any
}

func (v value) Path(_ ...string) containerI {
	return v
}

func (v value) Value() any {
	return v.any
}

func (v value) Set(_ any, _ ...string) {}

func (v value) Delete(_ ...string) {}

func (v value) Map() map[string]any {
	m, ok := (v.any).(map[string]any)
	if !ok {
		return nil
	}
	return m
}

func (v value) Type() string {
	return valueType
}
