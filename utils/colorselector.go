package utils

type ColorSelector interface {
	Select(field int) int
}

type ColorSelector16 struct {
}

func (s *ColorSelector16) Select(field int) int {
	return (field % 6) + 1
}
