package utils

import "os"

type Exception struct {
	code  int
	trace []string
	Err   error
}

func WrapError(err error, code int) *Exception {
	trace := []string{}
	return &Exception{
		code:  code,
		trace: trace,
		Err:   err,
	}
}

func (e *Exception) Error() string {
	return e.Err.Error()
}

func (e *Exception) Trace() []string {
	return e.trace
}
func (e *Exception) Code() int {
	return e.code
}

func (e *Exception) SetCode(code int) *Exception {
	e.code = code
	return e
}
func (e *Exception) Explain(trace string) *Exception {
	e.trace = append(e.trace, trace)
	return e
}

func (e *Exception) Terminate() {
	trace := e.Trace()
	message := ""
	for i := len(trace) - 1; i >= 0; i-- {
		message += trace[i] + "\n"
	}
	message += e.Error() + "\n"
	print(message)
	os.Exit(e.Code())
}
