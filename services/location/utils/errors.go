package utils

// NotImplementedError indicates that the function has not been implemented yet.
type NotImplementedError string

func (e NotImplementedError) Error() string {
	return string(e) + " is not implemented"
}
