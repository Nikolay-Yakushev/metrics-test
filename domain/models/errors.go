package models

import (
	"errors"
	"runtime"
	"runtime/debug"
)
type Traceble interface{
	Cause () error
	Stack ()
}

type stackTrace struct{
	trace     []byte
	cause     error
	file      string
	line      int
}

type RuntimeError struct{
	strace *stackTrace
	cause  error
	msg    string
	code   int
}


func (re *RuntimeError) Error() string {
	return re.msg
}

func (re *RuntimeError) Unwrap() error {
	return re.strace.cause
}

func New(msg string, code int, cause error) error {
	var (
		re RuntimeError
		buf []byte
	)

	if ok := errors.As(cause, &re); ok{
		return cause
	}
	
	_, file , line, _ := runtime.Caller(1)
	buf = debug.Stack()
	//TODO how to check stack
	// is already in stack
	re = RuntimeError{
		msg:  msg,
		code: code,
		strace: &stackTrace{
			cause: cause,
			file: file,
			line: line,
			trace: buf,
		},
	}
	return &re
}


var (
	NotFoundErr = errors.New("not found")
	NotAutorizedErr = errors.New("not authorized")
	InternalErr = errors.New("server error")
	ForbiddenErr = errors.New("access forbiden")

)