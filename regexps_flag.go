package gofactory

import (
	"fmt"
	"regexp"
	"strings"
)

type regexpsFlag struct {
	patternStrings []string
	patterns       []*regexp.Regexp
}

func (r *regexpsFlag) String() string {
	return strings.Join(r.patternStrings, ", ")
}

func (r *regexpsFlag) Set(pattern string) error {
	compiled, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("unable to compile factory pattern %s: %w", pattern, err)
	}

	r.patternStrings = append(r.patternStrings, pattern)
	r.patterns = append(r.patterns, compiled)

	return nil
}

func (r *regexpsFlag) Value() []*regexp.Regexp {
	return r.patterns
}
