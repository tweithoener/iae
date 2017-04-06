package iae

import (
	"fmt"
	"runtime"
	"strings"
)

// Check starts a chain of calls to Arg(). After an arbitrary number of calls
// to Arg() the chain ends with a call to Err(). Each call to Arg() will check
// one argument. Err() will then return the first error that occurred during
// these checks.
func Check() IAE {
	return &iae{nil, true, false, false}
}

// IAE provides function argument checking functionality.
type IAE interface {
	// Arg checks preconditions for methods and functions. The first argument
	// will be your actual argument check which must result in a boolean value
	// (e.g. a>0 && a<10). The second argument is the argument number (i.e. the
	// position in the argument list of the function or method). The number will be
	// placed into the error message in case of a failing check. The value argument takes
	// the argument under test. The value will also be put into the error message
	// in case of a failing check. The final argument is a very brief human readable
	// description of the precondition test (e.g. "0<a<10" or "a between 0 and 10").
	//
	// Arg returns an IAError if the precondition check fails or nil if
	// everything looks good.
	//
	// This method should be called immediately after entering a method or
	// function to check arguments provided to the method or function. This
	// method is designed for chaining of multiple precondition checks. In a
	// chain of Arg() method calls Subsequent calls will return immediately if
	// a preceding call resulted in an error. You should place a call to the
	// Err() method at the end of the chain, to get the result of your
	// preconfition check. No call to Err() is required if you always panic on
	// failing argument checks.
	//
	// If you only have a single check use the Arg function of this package
	// instead.
	Arg(check bool, argument uint, value interface{}, condition string) IAE

	// Err returns the first error that occurred during previous calls to the
	// Arg() method or nil if no error occurred in any of the precondition
	// checks.
	Err() error
}

// Arg is documented in the IAE interface
func (d *iae) Arg(check bool, argument uint, value interface{}, condition string) IAE {
	if d.err != nil {
		return d
	}

	d.err = d.process(check, argument, value, condition)
	return d
}

// Err is documented in the IAE interface.
func (d *iae) Err() error {
	return d.err
}

// CheckArg checks preconditions for methods and functions. The first argument
// will be your actual argument check which must result in a boolean value
// (e.g. a>0 && a<10). The second argument is the argument number (i.e. the
// position in the argument list of the function or method). The number will be
// placed into the error message in case of problem. The value argument takes
// the argument under test. The value will also be put into the error message
// in case of an error. The final argument is a very brief human readable
// desription of the precondition test (e.g. "0<a<10" or "a between 0 and 10").
// CheckArg returns an IAError if the precondition check fails or nil if
// everything looks good.
func CheckArg(check bool, argument uint, value interface{}, condition string) error {
	if check {
		return nil
	}
	d := &iae{nil, true, false, false}
	return d.process(check, argument, value, condition)
}

// IAError is returned by the argument checks if an illegal argument was passed
// into a function.
type IAError struct {
	funcName  string
	fileName  string
	line      int
	argument  uint
	value     interface{}
	condition string
}

// Error returns a string representation of an IAError
func (pce *IAError) Error() string {
	return fmt.Sprintf("illegal argument error: argument %d of %s is '%v' but must be %s at %s:%d",
		pce.argument, pce.funcName, pce.value,
		pce.condition,
		pce.fileName, pce.line)
}

// Mode is a type used for variables which hold information about this
// package's mode of operation. There are two important variables and three
// constants of this type. Together they define when argument checks should be
// performed (on exported or unerported fuctions or both) and how errors should
// be reported (produce an error vs. panic.)
type Mode uint8

const (
	// OFF assigned to the Exported or NotExported variable of this package
	// means do not perform checks.
	OFF Mode = iota

	// ERROR assigned to the Exported or NotExported variable of this package
	// means perform checks for these functions and return an error in case of
	// a failing check.
	ERROR

	// PANIC assigned to the Exported or NotExported variable of this package
	// means perform checks for these functions and panic in case of a failing
	// check.
	PANIC
)

// Exported controls if checks should be performed for exported functions or
// mathods. If OFF is assigned no checks are performed. If ERROR is assigned to
// this variable checks will be performed and an error will be returned in case
// of a failing check. If PANIC is assigned to this variable checks will panic
// if they fail.
var Exported = ERROR

// NotExported controls if checks should be performed for unexported/private
// functions or mathods. If OFF is assigned no checks are performed. If ERROR
// is assigned to this variable checks will be performed and an error will be
// returned in case of a failing check. If PANIC is assigned to this variable
// checks will panic if they fail.
var NotExported = PANIC

type iae struct {
	// err is the first error that occurred in a chain of argument checks.
	err error
	// execute is false if checks do not need to be executed
	execute bool
	// exported is false if we know that the function is not exported
	exported bool
	// sure is set to true when we are sure that the value of exported is
	// correct.
	sure bool
}

// process does the actual argument checking and error reporting
func (d *iae) process(check bool, argument uint, value interface{}, condition string) (err error) {
	if check {
		return nil
	}

	if !d.execute {
		return nil
	}
	if Exported == OFF && NotExported == OFF {
		d.execute = false
		return nil
	}

	fpcs := make([]uintptr, 3)
	n := runtime.Callers(3, fpcs)
	if n == 0 {
		panic("can't get callers.")
	}

	callee := runtime.FuncForPC(fpcs[0])
	if callee == nil {
		panic("can't get callee")
	}

	funcName := callee.Name()
	if !d.sure && Exported != NotExported {
		i := strings.LastIndex(funcName, ".")
		first := funcName[i+1 : i+2]
		d.exported = first == strings.ToUpper(first)
		d.sure = true
	}

	if !d.exported && NotExported == OFF {
		return nil
	}
	if d.exported && Exported == OFF {
		return nil
	}

	caller := runtime.FuncForPC(fpcs[1])
	if caller == nil {
		panic("can't get caller")
	}
	fileName, line := caller.FileLine(fpcs[1])
	err = &IAError{
		funcName:  funcName,
		fileName:  fileName,
		line:      line,
		argument:  argument,
		value:     value,
		condition: condition,
	}

	d.err = err
	d.execute = false

	if !d.exported && NotExported == PANIC {
		panic(err.Error())
	}
	if d.exported && Exported == PANIC {
		panic(err.Error())
	}

	return
}
