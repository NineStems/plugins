package toml

import (
	"bytes"
)

func split(in []byte) lines {
	data := bytes.Split(in, []byte("\n"))
	lns := make(lines, len(data))
	for i, d := range data {
		lns[i] = bytes.TrimSpace(d)
	}
	return lns
}

type line []byte

func (l line) Header() ([]string, bool) {
	mark := (bytes.Contains(l, []byte("[")) || bytes.Contains(l, []byte("]"))) &&
		!bytes.Contains(l, []byte("="))
	if !mark {
		return nil, false
	}

	nl := bytes.TrimFunc(l, func(r rune) bool {
		return r == '[' || r == ']'
	})

	keys := bytes.Split(nl, []byte("."))
	if len(keys) == 0 {
		return nil, false
	}

	if len(keys) == 0 {
		return nil, false
	}

	headers := make([]string, len(keys))
	for i, key := range keys {
		headers[i] = string(key)
	}

	return headers, true
}

func (l line) Pair() (string, any) {
	if !bytes.Contains(l, []byte("=")) {
		return "", nil
	}

	data := bytes.Split(l, []byte("="))

	if len(data) != 2 {
		return "", nil
	}

	return string(bytes.TrimSpace(data[0])), convertValue(bytes.TrimSpace(data[1]))
}

func (l line) IsDelimLines() bool {
	return len(l) == 0
}

type lines []line

func (l lines) String() string {
	var buf bytes.Buffer
	for _, ln := range l {
		buf.Write(ln)
		if !bytes.Equal(ln, []byte("\n")) {
			buf.WriteByte('\n')
		}
	}
	return buf.String()
}
