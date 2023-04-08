package pvm

// #cgo CFLAGS: -g -Wall
// #cgo LDFLAGS: -lpvm3
/*
	#include "pvm_wrapper.h"
	#include<stdlib.h>
*/
import "C"
import (
	"reflect"
	"strings"
	"unsafe"
)

func Mytid() (int, error) {
	result := int(C.pvm_mytid())

	dbgln("pvm_mytid() = %d", result)

	if result < 0 {
		return 0, PvmErrorFromInt(result)
	}

	return result, nil
}

func Parent() (int, error) {
	result := int(C.pvm_parent())

	dbgln("pvm_parent() = %d", result)

	if result < 0 {
		return 0, PvmErrorFromInt(result)
	}

	return result, nil
}

func CatchoutStdout() error {
	if result := C.pvm_catchout_stdout(); result < 0 {
		return PvmErrorFromCInt(result)
	}

	dbgln("pvm_catchout(stdout)")

	return nil
}

type Spawn_result struct {
	Numt int
	TIds []int
}

// if error code > 0, then some tasks failed to spawn - check TIds for error results
func Spawn(task string, args []string, flag SpawnOptions, where string, ntask int) (*Spawn_result, error) {
	task_cstr := C.CString(task)
	defer C.free(unsafe.Pointer(task_cstr))

	args_cstr_arr := stringSliceToCStringArray(args)
	defer func() {
		for _, c := range args_cstr_arr {
			C.free(unsafe.Pointer(c))
		}
	}()

	tIds_cint_ptr := (*C.int)(C.malloc(C.sizeof_ulong * C.ulong(ntask)))
	defer C.free(unsafe.Pointer(tIds_cint_ptr))

	where_cstr := C.CString(where)
	defer C.free(unsafe.Pointer(where_cstr))

	numt_cint := C.pvm_spawn(task_cstr, &args_cstr_arr[0], C.int(flag), where_cstr, C.int(ntask), tIds_cint_ptr)

	if int(numt_cint) < 0 {
		return nil, PvmErrorFromCInt(numt_cint)
	}

	tIds := cArrayToSlice(tIds_cint_ptr, ntask)

	dbgln("pvm_spawn(%s, %s, %d, %s, %d) = %d, %s", task, strings.Join(args, ", "), flag, where, ntask, numt_cint, tIds)

	if int(numt_cint) < ntask {
		return &Spawn_result{
			Numt: int(numt_cint),
			TIds: tIds,
		}, PvmErrorFromCInt(numt_cint)
	}

	return &Spawn_result{
		Numt: int(numt_cint),
		TIds: tIds,
	}, nil
}

func Perror(msg string) error {
	msg_cstr := C.CString(msg)
	defer C.free(unsafe.Pointer(msg_cstr))

	if info := C.pvm_perror(msg_cstr); info != 0 {
		return PvmErrorFromCInt(info)
	}

	dbgln("pvm_perror(%s)", msg)

	return nil
}

func Exit() error {
	if info := C.pvm_exit(); info != 0 {
		return PvmErrorFromCInt(info)
	}

	dbgln("pvm_exit()")

	return nil
}

func Initsend(encoding DataPackingStyle) (int, error) {
	bufid := C.pvm_initsend(C.int(encoding))

	if bufid < 0 {
		return 0, PvmErrorFromCInt(bufid)
	}

	dbgln("pvm_initsend(%d) = %d", encoding, bufid)

	return int(bufid), nil
}

func Kill(tId int) error {
	info := C.pvm_kill(C.int(tId))

	if info < 0 {
		return PvmErrorFromCInt(info)
	}

	dbgln("pvm_kill(%d) = %d", tId, info)

	return nil
}

/* PVM_PACKF STATIC BINDINGS */

func PackfString(fmt string, arg string) (int, error) {
	fmt_cstr := C.CString(fmt)
	defer C.free(unsafe.Pointer(fmt_cstr))

	arg_cstr := C.CString(arg)
	defer C.free(unsafe.Pointer(arg_cstr))

	info := C.pvm_packf_string(fmt_cstr, arg_cstr)
	if info < 0 {
		return 0, PvmErrorFromCInt(info)
	}

	dbgln("pvm_packf(%s, %s)", fmt, arg)

	return int(info), nil
}

/* PVM_UNPACKF STATIC BINDINGS */

func UnpackfString(fmt string, buflen int) (string, error) {
	fmt_cstr := C.CString(fmt)
	defer C.free(unsafe.Pointer(fmt_cstr))

	arg_cstr := (*C.char)(C.malloc(C.sizeof_char * C.ulong(buflen)))
	defer C.free(unsafe.Pointer(arg_cstr))

	info := C.pvm_unpackf_string(fmt_cstr, arg_cstr)

	if info < 0 {
		return "", PvmErrorFromCInt(info)
	}

	dbgln("pvm_unpackf(%s, %d) = %s", fmt, buflen, C.GoString(arg_cstr))

	return C.GoString(arg_cstr), nil
}

func Send(tid int, msgtag int) error {
	if info := C.pvm_send(C.int(tid), C.int(msgtag)); info < 0 {
		return PvmErrorFromCInt(info)
	}

	dbgln("pvm_send(%d, %d)", tid, msgtag)

	return nil
}

func Recv(tid int, msgtag int) (int, error) {
	info := C.pvm_recv(C.int(tid), C.int(msgtag))

	if info < 0 {
		return 0, PvmErrorFromCInt(info)
	}

	dbgln("pvm_send(%d, %d) = %d", tid, msgtag, info)

	return int(info), nil
}

func Nrecv(tid int, msgtag int) (int, error) {
	info := C.pvm_nrecv(C.int(tid), C.int(msgtag))

	if info < 0 {
		return 0, PvmErrorFromCInt(info)
	}

	dbgln("pvm_nrecv(%d, %d) = %d", tid, msgtag, info)

	return int(info), nil
}

type BufinfoResult struct {
	Bytes  int
	MsgTag int
	TId    int
}

func Bufinfo(bufid int) (*BufinfoResult, error) {
	bytes_cint_ptr := (*C.int)(C.malloc(C.sizeof_ulong))
	defer C.free(unsafe.Pointer(bytes_cint_ptr))

	msgtag_cint_ptr := (*C.int)(C.malloc(C.sizeof_ulong))
	defer C.free(unsafe.Pointer(msgtag_cint_ptr))

	tid_cint_ptr := (*C.int)(C.malloc(C.sizeof_ulong))
	defer C.free(unsafe.Pointer(tid_cint_ptr))

	if info := C.pvm_bufinfo(C.int(bufid), bytes_cint_ptr, msgtag_cint_ptr, tid_cint_ptr); info < 0 {
		return nil, PvmErrorFromCInt(info)
	}

	return &BufinfoResult{
		Bytes:  int(*bytes_cint_ptr),
		MsgTag: int(*msgtag_cint_ptr),
		TId:    int(*tid_cint_ptr),
	}, nil
}

func stringSliceToCStringArray(strings []string) []*C.char {
	var result []*C.char
	for _, c := range strings {
		result = append(result, C.CString(c))
	}

	return result
}

func cArrayToSlice(array *C.int, len int) []int {
	var result []int
	var list []C.int
	sliceHeader := (*reflect.SliceHeader)((unsafe.Pointer(&list)))
	sliceHeader.Cap = len
	sliceHeader.Len = len
	sliceHeader.Data = uintptr(unsafe.Pointer(array))

	for _, c := range list {
		result = append(result, int(c))
	}

	return result
}
