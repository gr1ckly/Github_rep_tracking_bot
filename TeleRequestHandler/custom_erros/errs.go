package custom_erros

import (
	"Common"
	"strconv"
)

type ServerError struct {
	StatusCode int
	Status     string
	Url        string
	ErrDTO     Common.ErrorDTO
}

func (se ServerError) Error() string {
	return "Server error: " + se.Status + " " + strconv.Itoa(se.StatusCode) + " on " + se.Url + " error: " + se.ErrDTO.Error
}
