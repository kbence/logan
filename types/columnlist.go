package types

import (
	"bytes"
	"hash/crc64"
	"sort"
)

type ColumnList map[int]string

var crc64Table = crc64.MakeTable(crc64.ECMA)

func (l ColumnList) Equals(other ColumnList) bool {
	if len(l) != len(other) {
		return false
	}

	for key, val := range l {
		otherVal, found := other[key]

		if !found || val != otherVal {
			return false
		}
	}

	return true
}

func (l ColumnList) Join() string {
	buffer := bytes.NewBufferString("")

	for _, k := range l.SortedKeys() {
		buffer.WriteString(l[k])
		buffer.WriteByte(0)
	}

	return buffer.String()
}

func (l ColumnList) Keys() []int {
	keys := []int{}

	for k := range l {
		keys = append(keys, k)
	}

	return keys
}

func (l ColumnList) SortedKeys() []int {
	keys := l.Keys()
	sort.Ints(keys)
	return keys
}
