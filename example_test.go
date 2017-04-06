package iae_test

import (
	"fmt"
	"strings"

	"git.weithoener.net/repo/iae"
)

// B wraps an int value
type B struct {
	B int
}

// Foo takes and int argument, negates it and wraps it in a B struct. The
// argument must be greater than 10.
func Foo(a int) (b *B, err error) {
	// A single, simple precondition
	if err = iae.CheckArg(a > 10, 1, a, ">10"); err != nil {
		return
	}

	// The logic of this function:
	b = &B{
		B: -1 * a,
	}

	return
}

// Bar is a more complex version of Foo. It takes an int argumenti, negates it
// and wraps it into a B struct. The argument must be greater than 10 but not
// 12.
func Bar(a int) (b *B, err error) {
	// First define the preconditions of this function. See how we chain the
	// check method calls. The final Err() returns the first error that
	// occurred.
	if err = iae.Check().
		Arg(a > 10, 1, a, ">10").
		Arg(a != 12, 1, a, "!=12").Err(); err != nil {
		return
	}

	// The logic of this function:
	b = &B{
		B: -1 * a,
	}

	return
}

func Example() {
	// Make sure we get errors not panic
	iae.Release = iae.ERROR

	// Let's call Foo() with different arguments
	if _, err := Foo(5); err != nil {
		fmt.Println("5 is not OK:", strings.Split(err.Error(), " at")[0], "...")
	}

	if _, err := Foo(11); err != nil {
		fmt.Println("11 is fine it must not cause an error:", err)
	}

	// Once again but now we are calling the more complex Bar() function:
	if _, err := Bar(5); err != nil {
		fmt.Println("5 is still not OK:", strings.Split(err.Error(), " at")[0], "...")
	}

	if _, err := Bar(11); err != nil {
		fmt.Println("11 is fine it must not cause an error:", err)
	}

	// Output:
	// 5 is not OK: illegal argument error: argument 1 of git.weithoener.net/repo/iae_test.Foo is '5' but must be >10 ...
	// 5 is still not OK: illegal argument error: argument 1 of git.weithoener.net/repo/iae_test.Bar is '5' but must be >10 ...

}
