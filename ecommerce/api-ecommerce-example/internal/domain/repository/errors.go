package repository

type ErrNotFound struct {
	Message string
}

func (e *ErrNotFound) Error() string {
	return e.Message
}

type ErrConflict struct {
	Message string
}

func (e *ErrConflict) Error() string {
	return e.Message
}
