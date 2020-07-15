package cerror

import (
	"errors"
	"fmt"
	"github.com/m2c/kiplestar/commons"
	"runtime"
	"strings"
)

type CommonsError struct {
	Message    string
	StatusCode commons.ResponseCode
	rawErr     error
	stackPC    []uintptr
}

func (e *CommonsError) Error() string {
	return e.Message
}

// RawErr the origin err
func (e *CommonsError) RawErr() error {
	return e.rawErr
}

// CallStack get function call stack
func (e *CommonsError) CallStack() string {
	frames := runtime.CallersFrames(e.stackPC)
	var (
		f      runtime.Frame
		more   bool
		result string
		index  int
	)
	for {
		f, more = frames.Next()
		if index = strings.Index(f.File, "src"); index != -1 {
			f.File = string(f.File[index+4:])
		}
		result = fmt.Sprintf("%s%s\n\t%s:%d\n", result, f.Function, f.File, f.Line)
		if !more {
			break
		}
	}
	return result
}

/**
* Package error code
 */
func ConstructionErr(err error, code commons.ResponseCode, fmtAndArgs ...interface{}) *CommonsError {
	msg := fmtErrMsg(fmtAndArgs...)
	if err == nil {
		err = errors.New(msg)
	}
	if e, ok := err.(*CommonsError); ok {
		if msg != "" {
			e.Message = msg
		}
		if code != 0 {
			e.StatusCode = code
		}
		return e
	}

	pcs := make([]uintptr, 32)
	// skip the first 3 invocations
	count := runtime.Callers(3, pcs)
	e := &CommonsError{
		StatusCode: code,
		Message:    msg,
		rawErr:     err,
		stackPC:    pcs[:count],
	}
	if e.Message == "" {
		e.Message = err.Error()
	}
	return e
}

/**
 * format error message
 */
func fmtErrMsg(msgs ...interface{}) string {
	if len(msgs) > 1 {
		return fmt.Sprintf(msgs[0].(string), msgs[1:]...)
	}
	if len(msgs) == 1 {
		if v, ok := msgs[0].(string); ok {
			return v
		}
		if v, ok := msgs[0].(error); ok {
			return v.Error()
		}
	}
	return ""
}
func ServiceWrapErr(err error, code commons.ResponseCode, fmtAndArgs ...interface{}) *CommonsError {
	return ConstructionErr(err, code, fmtAndArgs...)
}

/**
 * console log the business log
 */
func ServiceError(obj interface{}) {
	panic(obj)
}
