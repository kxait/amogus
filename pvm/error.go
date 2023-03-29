package pvm

import "fmt"

type PvmError struct {
	ErrorCode ErrorCode
}

func (e *PvmError) Error() string {
	return fmt.Sprintf("pvm3 returned error with code %d", e.ErrorCode)
}

func PvmErrorFromInt(code int) *PvmError {
	return &PvmError{ErrorCode: ErrorCode(code)}
}

func PvmErrorFromCInt(code _Ctype_int) *PvmError {
	return &PvmError{ErrorCode: ErrorCode(int(code))}
}
