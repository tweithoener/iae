package iae

import (
	"fmt"
	"runtime"
)

// Check starts a chain of calls to Arg(). After an arbitrary number of calls
// to Arg() the chain ends with a call to Err(). Each call to Arg() will check
// one argument. Err() will then return the first error that occurred during
// these checks.
// You can add a call to Dbg() at any point in the chain. This will tell the
// argument checking logic that all following checks are debug checks.
func Check() IAE {
	return &iae{
		err:     nil,
		execute: Release != OFF,
		debug:   false,
	}
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
	// This method should be called immediately after entering a method or
	// function to check arguments provided to the method or function. This
	// method is designed for chaining of multiple precondition checks. In a
	// chain of Arg() method calls Subsequent calls will return immediately if
	// a preceding call resulted in an error. You should place a call to the
	// Err() method at the end of the chain, to get the result of your
	// preconfition check. No call to Err() is required if you always panic on
	// failing argument checks.
	//
	// If you only have a single check use the CheckArg function of this package
	// instead.
	Arg(check bool, argument uint, value interface{}, condition string) IAE

	Dbg() IAE

	// Err returns the first error that occurred during previous calls to the
	// Arg() method or nil if no error occurred in any of the precondition
	// checks.
	Err() error
}

// Arg is documented in the IAE interface
func (d *iae) Arg(check bool, argument uint, value interface{}, condition string) IAE {
	if check {
		return d
	}
	if d.err != nil {
		return d
	}

	d.err = d.process(check, argument, value, condition)
	return d
}

func (d *iae) Dbg() IAE {
	d.debug = true
	d.execute = Debug != OFF
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
// everything looks good. The function will panic if the checks fails and the
// global Release variable is set to PANIC.
func CheckArg(check bool, argument uint, value interface{}, condition string) error {
	if check {
		return nil
	}
	if Release == OFF {
		return nil
	}
	d := &iae{
		err:     nil,
		execute: true,
		debug:   false,
	}
	return d.process(check, argument, value, condition)
}

// CheckDbgArg is equivalent to CheckArg with the difference that the argument
// check is considered a debug check. Thus the configuration for error reporting
// is taken from the global Debug variable not from the global Release variable.
func CheckDbgArg(check bool, argument uint, value interface{}, condition string) error {
	if check {
		return nil
	}
	if Debug == OFF {
		return nil
	}
	d := &iae{
		err:     nil,
		execute: true,
		debug:   true,
	}
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
// constants of this type. Together they define which argument checks should be
// performed (release checks or debug checks or both) and how errors should
// be reported (produce an error vs. panic.)
type Mode uint8

const (
	// OFF assigned to the Release or Debug variable of this package
	// means do not perform these checks.
	OFF Mode = iota

	// ERROR assigned to the Release or Debug variable of this package
	// means perform these checks and return an error in case of
	// a failing check.
	ERROR

	// PANIC assigned to the Release or Debug variable of this package
	// means perform these checks and panic in case of a failing
	// check.
	PANIC
)

// Release controls if release checks should be performed. Release checks are
// all those checks which are done using a CheckArg call and any chained check
// positioned before the first occurrence of Dbg() in the chain. If OFF is
// assigned no checks are performed. If ERROR is assigned to
// this variable checks will be performed and an error will be returned in case
// of a failing check. If PANIC is assigned to this variable checks will panic
// if they fail.
var Release = ERROR

// Debug controls if debug checks should be performed. Debug checks are
// all those checks which are done using a CheckDbgArg call and any chained check
// positioned after the first occurrence of Dbg() in the chain.
// If OFF is assigned no checks are performed. If ERROR
// is assigned to this variable checks will be performed and an error will be
// returned in case of a failing check. If PANIC is assigned to this variable
// checks will panic if they fail.
var Debug = PANIC

type iae struct {
	err     error // first error that occurred in a chain of argument checks.
	execute bool  // false if checks do not need to be executed
	debug   bool  // true when checks are debug checks
}

// process does the actual argument checking and error reporting
func (d *iae) process(check bool, argument uint, value interface{}, condition string) error {
	if check {
		return nil
	}
	if !d.execute {
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
	caller := runtime.FuncForPC(fpcs[1])
	if caller == nil {
		panic("can't get caller")
	}
	fileName, line := caller.FileLine(fpcs[1])
	err := &IAError{
		funcName:  funcName,
		fileName:  fileName,
		line:      line,
		argument:  argument,
		value:     value,
		condition: condition,
	}

	d.err = err
	d.execute = false

	if d.debug && Debug == PANIC {
		panic(err.Error())
	}
	if !d.debug && Release == PANIC {
		panic(err.Error())
	}

	return err
}
