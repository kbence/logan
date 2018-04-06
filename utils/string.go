package utils

import (
	"bytes"
	"fmt"
)

func StringListsEqual(a, b []string) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func FormatStringArray(strings []string) string {
	buf := bytes.NewBufferString("[")

	for i := range strings {
		if i > 0 {
			buf.WriteString(", ")
		}

		buf.WriteString(fmt.Sprintf("\"%s\"", strings[i]))
	}

	buf.WriteRune(']')

	return buf.String()
}
