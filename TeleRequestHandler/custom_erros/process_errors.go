package custom_erros

type ProcessError struct {
	Message string
}

func (pe ProcessError) Error() string {
	return pe.Message
}
