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

// Foo takes an int argument, negates it and
// stores it in the receiver. The argument must be greater than 10 but not 12.
func (b *B) Foo(a int) error {
	// The preconditions of this function: See how we chain the check
	// methods. The final Err() returns the first error that occurred.
	// If the receiver is nil we panic in any case
	if err := iae.Check().
		Recv(b != nil, b, "not nil").
		Panic().
		Arg(a > 10, 1, a, ">10").
		Arg(a != 12, 1, a, "!=12").
		Err(); err != nil {
		return err
	}

	// The logic of this function:
	b.B = -1 * a

	return nil
}

func Example() {
	// Make sure we get errors not panic
	iae.Release = iae.ERROR

	b := &B{}
	// Let's call Foo() with different arguments
	if err := b.Foo(5); err != nil {
		fmt.Println("5 is not OK:", strings.Split(err.Error(), " at")[0], "...")
	}

	if err := b.Foo(11); err != nil {
		fmt.Println("11 is fine it must not cause an error:", err)
	}

	defer func() {
		r := recover()
		if r != nil {
			msg := fmt.Sprintf("%v", r)
			fmt.Println("panic:", strings.Split(msg, " at")[0], "...")
		}
	}()

	b = nil
	b.Foo(5)

	// Output:
	// 5 is not OK: illegal argument error: argument 1 of git.weithoener.net/repo/iae_test.(*B).Foo is '5' but must be >10 ...
	// panic: illegal argument error: receiver of git.weithoener.net/repo/iae_test.(*B).Foo is '<nil>' but must be not nil ...
}
