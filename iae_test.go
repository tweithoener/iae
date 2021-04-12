package iae_test

import (
	"testing"

	"github.com/tweithoener/iae"
)

func multi(a int) error {
	return iae.Check().Dbg().
		Arg(a > 1, 1, a, ">1").
		Arg(a > 2, 2, a, ">2").
		Err()
}

func Exported(a int) error {
	return iae.Check().Arg(a > 0, 1, a, ">0").Err()
}

func Mixed(a int) error {
	return iae.Check().
		Arg(a > 0, 1, a, ">0").
		Dbg().
		Arg(a > 1, 1, a, ">1").Err()
}

func TestPrecond(t *testing.T) {
	tests := []struct {
		rel      iae.Mode
		dbg      iae.Mode
		fun      func(int) error
		val      int
		expErr   bool
		expPanic bool
	}{
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

		{iae.OFF, iae.OFF, Mixed, -1, false, false},
		{iae.ERROR, iae.OFF, Mixed, -1, true, false},
		{iae.ERROR, iae.OFF, Mixed, 0, true, false},
		{iae.ERROR, iae.OFF, Mixed, 1, false, false},
		{iae.ERROR, iae.ERROR, Mixed, 1, true, false},
		{iae.ERROR, iae.ERROR, Mixed, 2, false, false},

		{iae.OFF, iae.OFF, Mixed, -1, false, false},
		{iae.PANIC, iae.OFF, Mixed, -1, false, true},
		{iae.PANIC, iae.OFF, Mixed, 0, false, true},
		{iae.PANIC, iae.OFF, Mixed, 1, false, false},
		{iae.PANIC, iae.PANIC, Mixed, 1, false, true},
		{iae.PANIC, iae.PANIC, Mixed, 2, false, false},
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
			iae.Release = test.rel
			iae.Debug = test.dbg
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
	iae.Release = iae.ERROR
	for i := 0; i < b.N; i++ {
		err := Exported(-1)
		if err == nil {
			b.Error("precondition error missing")
			b.Fail()
		}
	}
}

func BenchmarkErrorMulti(b *testing.B) {
	iae.Debug = iae.ERROR
	for i := 0; i < b.N; i++ {
		err := multi(2)
		if err == nil {
			b.Error("precondition error missing")
			b.Fail()
		}
	}
}
