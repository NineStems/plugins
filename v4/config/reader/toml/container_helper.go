package toml

func convert(in lines) containerI {
	result := make(container, len(in))
	if len(in) == 0 {
		return result
	}

	headers := make([]string, 0, 1)
	for _, ln := range in {
		if ln.IsDelimLines() {
			continue
		}

		h, ok := ln.Header()
		if ok {
			headers = h
			continue
		}
		k, v := ln.Pair()
		addPath(headers, result, k, v)
	}

	return result
}

func processedValueLine(headers []string, cont container, in line) {
	k, v := in.Pair()
	addPath(headers, cont, k, v)
}

func calculatePath(in []string) (string, []string) {
	key := in[len(in)-1]
	var path []string
	if len(in) >= 2 {
		path = in[:len(in)-1]
	}

	return key, path
}

func pathContainer(path []string, cont container) container {
	current := cont
	for _, header := range path {
		ci, ok := current[header]
		if !ok {
			ci = &container{}
		}
		current[header] = ci
		cc := ci.(*container)
		current = *cc
	}
	return current
}

func addPath(path []string, cont container, key string, val any) {
	current := pathContainer(path, cont)
	current[key] = value{val}
}

func deleteContainer(in []string, cont container) {
	if len(in) == 0 {
		return
	}

	key, path := calculatePath(in)
	current := pathContainer(path, cont)
	delete(current, key)
}

func builderMap(in container) map[string]any {
	result := make(map[string]any, len(in))
	for key, val := range in {
		if val.Type() == containerType {
			cc := val.(*container)
			result[key] = builderMap(*cc)
			continue
		}
		result[key] = val.Value()
	}
	return result
}
