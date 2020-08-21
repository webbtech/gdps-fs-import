package pkgerrors

import (
	"errors"
	"testing"
)

func stdErrorFunc() error {
	return &StdError{Err: "My bogus error", Caller: "stdErrorFunc", Msg: "Error message"}
}

func TestStdErrorFunc(t *testing.T) {

	var err *StdError
	e := stdErrorFunc()

	if ok := errors.As(e, &err); ok {
		if err.Caller != "stdErrorFunc" {
			t.Errorf("got: %s, want: %s", err.Caller, "stdErrorFunc")
		}
	}

	if errors.Is(e, err) {
		if err.Msg != "Error message" {
			t.Errorf("got: %s, want: %s", err.Msg, "Error message")
		}
	}
}
