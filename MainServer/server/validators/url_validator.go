package validators

import (
	"regexp"
	"sync"
)

type UrlValidator struct {
	regex        *regexp.Regexp
	checkerMutex sync.Mutex
	checker      Checker[string, bool]
}

func NewUrlValidator(pattern string, checker Checker[string, bool]) (*UrlValidator, error) {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	return &UrlValidator{regex: regex, checker: checker, checkerMutex: sync.Mutex{}}, nil
}

func (val *UrlValidator) Check(link string) bool {
	val.checkerMutex.Lock()
	defer val.checkerMutex.Unlock()
	if !val.checker.Check(link) {
		return false
	}
	return val.regex.MatchString(link)
}
