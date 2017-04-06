// Package iae provides functionality for method/function argument checks. This
// tiny, self-contained, robust library has the following aims:
//
// 1) Less typing: A short notation with only one if statement for an arbitrary
// number of argument checks and no need to type the same boring error messages
// over and over again.
//
// 2) Consistency: Get consistent syntax for checks and consistent error types
// and error messages if a check fails.
//
// 3) Flexibility: Want a panic during development and no checks at all in
// production? Want always errors? Exported and unexported functions treated
// differently. No problem. The behavior is controlled by two globel variables
// in this package (Exported and NotExported). Consider using golang build tags
// to configure these variables as needed in your debug or release builds
//
// Have a look at the example to see the functionality in action. It's really
// easy to use.
//
// A note: The author intends to add a code generator to this project at some
// point in the future. The code generator will add argument checks based on
// some special comments.
//
// Another note: The author thinks it is vital to check function arguments.
// This is especially true for exported functions, You may call this
// predonditions. But do not ask for postconditions and invariants. The author
// is not going to add support for these test to this package. Consistent error
// checking is a must for solid go programs. This renders postconditions
// useless. You should not write a long function and then at the very end
// detect that there must have been an error somewhere. And if you would like
// to have invariants, no problem: it's just a method (with a special name)
// called at the end of every method (defer invariants()).
package iae
