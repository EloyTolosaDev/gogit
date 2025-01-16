package errors

type InitError struct {
}

func (err InitError) Error() string {
	return ""
}

type CommitError struct {
}

func (err CommitError) Error() string {
	return ""
}

type ConfigError struct {
}

func (ce ConfigError) Error() string {
	return ""
}
