package factory

import "strings"

type stringsFlag []string

func (s stringsFlag) String() string {
	return strings.Join(s, ", ")
}

func (s *stringsFlag) Set(value string) error {
	*s = append(*s, value)

	return nil
}

func (s stringsFlag) Value() []string {
	res := make([]string, 0, len(s))

	for _, str := range s {
		res = append(res, strings.TrimSpace(str))
	}

	return res
}
