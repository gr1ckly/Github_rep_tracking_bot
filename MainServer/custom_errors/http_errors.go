package custom_errors

import "strconv"

type StatusError struct {
	StatusCode int
	Status     string
}

func (se StatusError) Error() string {
	return se.Status + " " + strconv.Itoa(se.StatusCode)
}
