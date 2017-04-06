package iae_test

import (
	"testing"

	"github.com/tweithoener/iae"
)

func simple(a int) error {
	return iae.CheckArg(a > 0, 1, a, ">0")
}

func multi(a int) error {
	return iae.Check().
		Arg(a > 1, 1, a, ">1").
		Arg(a > 2, 2, a, ">2").
		Err()
}

func Exported(a int) error {
	return iae.CheckArg(a > 0, 1, a, ">0")
}

func TestPrecond(t *testing.T) {
	tests := []struct {
		exp      iae.Mode
		unx      iae.Mode
		fun      func(int) error
		val      int
		expErr   bool
		expPanic bool
	}{
		{iae.OFF, iae.ERROR, simple, 1, false, false},
		{iae.ERROR, iae.ERROR, simple, 1, false, false},
		{iae.PANIC, iae.ERROR, simple, 1, false, false},

		{iae.OFF, iae.ERROR, simple, -1, true, false},
		{iae.ERROR, iae.ERROR, simple, -1, true, false},
		{iae.PANIC, iae.ERROR, simple, -1, true, false},

		{iae.OFF, iae.OFF, simple, -1, false, false},
		{iae.ERROR, iae.OFF, simple, -1, false, false},
		{iae.PANIC, iae.OFF, simple, -1, false, false},

		{iae.ERROR, iae.OFF, Exported, 1, false, false},
		{iae.ERROR, iae.ERROR, Exported, 1, false, false},
		{iae.ERROR, iae.PANIC, Exported, 1, false, false},

		{iae.ERROR, iae.OFF, Exported, -1, true, false},
		{iae.ERROR, iae.ERROR, Exported, -1, true, false},
		{iae.ERROR, iae.PANIC, Exported, -1, true, false},

		{iae.OFF, iae.OFF, Exported, -1, false, false},
		{iae.OFF, iae.ERROR, Exported, -1, false, false},
		{iae.OFF, iae.PANIC, Exported, -1, false, false},

		{iae.OFF, iae.ERROR, multi, 3, false, false},
		{iae.ERROR, iae.ERROR, multi, 3, false, false},
		{iae.PANIC, iae.ERROR, multi, 3, false, false},

		{iae.OFF, iae.ERROR, multi, 2, true, false},
		{iae.ERROR, iae.ERROR, multi, 2, true, false},
		{iae.PANIC, iae.ERROR, multi, 2, true, false},

		{iae.OFF, iae.ERROR, multi, 1, true, false},
		{iae.ERROR, iae.ERROR, multi, 1, true, false},
		{iae.PANIC, iae.ERROR, multi, 1, true, false},

		{iae.OFF, iae.OFF, multi, -1, false, false},
		{iae.ERROR, iae.OFF, multi, -1, false, false},
		{iae.PANIC, iae.OFF, multi, -1, false, false},

		{iae.OFF, iae.PANIC, simple, 1, false, false},
		{iae.ERROR, iae.PANIC, simple, 1, false, false},
		{iae.PANIC, iae.PANIC, simple, 1, false, false},

		{iae.OFF, iae.PANIC, simple, -1, false, true},
		{iae.ERROR, iae.PANIC, simple, -1, false, true},
		{iae.PANIC, iae.PANIC, simple, -1, false, true},

		{iae.OFF, iae.OFF, simple, -1, false, false},
		{iae.ERROR, iae.OFF, simple, -1, false, false},
		{iae.PANIC, iae.OFF, simple, -1, false, false},

		{iae.PANIC, iae.OFF, Exported, 1, false, false},
		{iae.PANIC, iae.ERROR, Exported, 1, false, false},
		{iae.PANIC, iae.PANIC, Exported, 1, false, false},

		{iae.PANIC, iae.OFF, Exported, -1, false, true},
		{iae.PANIC, iae.ERROR, Exported, -1, false, true},
		{iae.PANIC, iae.PANIC, Exported, -1, false, true},

		{iae.OFF, iae.PANIC, multi, 3, false, false},
		{iae.ERROR, iae.PANIC, multi, 3, false, false},
		{iae.PANIC, iae.PANIC, multi, 3, false, false},

		{iae.OFF, iae.PANIC, multi, 2, false, true},
		{iae.ERROR, iae.PANIC, multi, 2, false, true},
		{iae.PANIC, iae.PANIC, multi, 2, false, true},

		{iae.OFF, iae.PANIC, multi, 1, false, true},
		{iae.ERROR, iae.PANIC, multi, 1, false, true},
		{iae.PANIC, iae.PANIC, multi, 1, false, true},
	}
	for _, test := range tests {
		func() {
			defer func() {
				r := recover()
				if test.expPanic && r == nil {
					t.Error("precondition panic missing")
					t.Fail()
				}
				if !test.expPanic && r != nil {
					t.Errorf("unexpected precondition panic: %s", r)
					t.Fail()
				}

			}()
			iae.Exported = test.exp
			iae.NotExported = test.unx
			err := test.fun(test.val)
			if !test.expErr && err != nil {
				t.Errorf("unexpected precondition error: %s", err.Error())
				t.Fail()
			}
			if test.expErr && err == nil {
				t.Errorf("precondition error missing")
				t.Fail()
			}
		}()
	}
}

func BenchmarkErrorExported(b *testing.B) {
	iae.Exported = iae.ERROR
	for i := 0; i < b.N; i++ {
		err := Exported(-1)
		if err == nil {
			b.Error("precondition error missing")
			b.Fail()
		}
	}
}

func BenchmarkErrorMulti(b *testing.B) {
	iae.NotExported = iae.ERROR
	for i := 0; i < b.N; i++ {
		err := multi(2)
		if err == nil {
			b.Error("precondition error missing")
			b.Fail()
		}
	}
}
