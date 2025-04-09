package Common

import (
	"strconv"
)

type StatusError struct {
	StatusCode int
	Status     string
	Url        string
}

func (se StatusError) Error() string {
	return se.Status + " " + strconv.Itoa(se.StatusCode) + " on " + se.Url
}
