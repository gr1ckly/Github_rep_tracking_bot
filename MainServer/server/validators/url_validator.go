package validators

import (
	"regexp"
)

type UrlValidator struct {
	regex   *regexp.Regexp
	checker Checker[string, bool]
}

func NewUrlValidator(pattern string, checker Checker[string, bool]) (*UrlValidator, error) {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	return &UrlValidator{regex: regex, checker: checker}, nil
}

func (val *UrlValidator) Check(link string) bool {
	return val.regex.MatchString(link) && val.checker.Check(link)
}
