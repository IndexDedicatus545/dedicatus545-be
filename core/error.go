package core

import "strconv"

const (
	ErrInvalidBook = Error(iota)
	ErrInvalidISBN
)

type Error uint

//todo подумать над лучшим представлением ошибок (после того как станет понятно как их будут использовать)
type errorWithMsg struct {
	code   Error
	msg    string
	clause error
}

func (e errorWithMsg) Error() string {
	return e.msg + ": " + e.clause.Error()
}

func (e errorWithMsg) Unwrap() error {
	return e.clause
}

func newInvalidISBNError(invalidISBN string) error {
	return errorWithMsg{code: ErrInvalidISBN, msg: "invalid isbn: [" + strconv.Quote(invalidISBN) + "]"}
}

func newInvalidBookError(msg string) error {
	return errorWithMsg{code: ErrInvalidBook, msg: "invalid book: " + msg}
}

func newInvalidBookErrorWithClause(msg string, err error) error {
	return errorWithMsg{code: ErrInvalidBook, msg: "invalid book: " + msg, clause: err}
}
