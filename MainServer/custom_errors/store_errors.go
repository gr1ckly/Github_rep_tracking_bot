package custom_errors

type NoValuesError struct {
	err error
}

func NewNoValuesError(err error) NoValuesError {
	return NoValuesError{err}
}

func (ne NoValuesError) Error() string {
	return "No values error: " + ne.err.Error()
}

func (ne NoValuesError) Unwrap() error {
	return ne.err
}

type AlreadyExists struct {
	err error
}

func NewAlreadyExistsError(err error) AlreadyExists {
	return AlreadyExists{err}
}

func (ae AlreadyExists) Error() string {
	return "No values error: " + ae.err.Error()
}

func (ae AlreadyExists) Unwrap() error {
	return ae.err
}
