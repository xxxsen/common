package errs

import (
	"bytes"
	"fmt"
)

var (
	ErrOK = New(0, "success")
)

// Deprecated: should not use this
type IError interface {
	error
	Code() int64
	Message() string
}

type Error struct {
	code   int64
	msg    string
	err    error
	extmsg []string
}

func (e *Error) Error() string {
	buf := bytes.NewBufferString(fmt.Sprintf("Error:[code:%d, msg:%s]", e.code, e.msg))
	if e.err != nil {
		buf.WriteString(fmt.Sprintf(", err:[%v]", e.err))
	}
	if len(e.extmsg) > 0 {
		buf.WriteString(fmt.Sprintf(", extmsg:%v", e.extmsg))
	}
	return buf.String()
}

func (e *Error) Code() int64 {
	return e.code
}

func (e *Error) Message() string {
	return e.msg
}

// Deprecated: should not use this
func New(code int64, fmtter string, args ...interface{}) *Error {
	return Wrap(
		code,
		fmt.Sprintf(fmtter, args...),
		nil,
	)
}

// Deprecated: should not use this
func Wrap(code int64, msg string, err error) *Error {
	return &Error{
		code: code,
		msg:  msg,
		err:  err,
	}
}

// Deprecated: should not use this
func (e *Error) WithDebugMsg(fmtter string, args ...interface{}) *Error {
	e.extmsg = append(e.extmsg, fmt.Sprintf(fmtter, args...))
	return e
}

// Deprecated: should not use this
func IsErrOK(err error) bool {
	ierr := FromError(err)
	if ierr == nil {
		return true
	}
	if ierr.Code() == 0 {
		return true
	}
	return false
}

// Deprecated: should not use this
func FromError(err error) IError {
	if err == nil {
		return nil
	}
	if e, ok := err.(IError); ok {
		return e
	}
	return Wrap(100000, "unknown error", err)
}
