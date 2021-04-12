// Package iae provides functionality for method/function argument checks. This
// tiny, self-contained, robust library has the following aims:
//
// 1) Less typing: A short notation with only one if statement for an arbitrary
// number of argument checks and no need to type the same boring error messages
// over and over again.
//
// 2) Consistency: Get consistent syntax for checks and consistent error types and
// error messages if a check fails.
//
// 3) Flexibility: Want a panic during development and no checks at all in
// production? Want always errors? Set the Debug and Release variables of this
// Package respectively. Consider using go build tags to set these variables
// depending on your build type.
//
// Have a look at the example to see the functionality in action. It's really easy
// to use.

//
// A note: The author thinks it is vital to check function arguments. This
// is especially true for exported functions, You may call this 'preconditions'. But
// do not ask for postconditions and invariants. The author is not going to add
// support for these test to this package. Consistent error checking is a must for
// solid go programs. This renders postconditions useless. You should not write a
// long function and then at the very end detect that there must have been an
// error somewhere. And if you would like to have invariants, no problem: it's
// just a method (with a special name) called at the end of every method (defer
// invariants()).
package iae
