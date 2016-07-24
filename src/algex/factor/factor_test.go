package factor

import (
	"testing"
)

func TestString(t *testing.T) {
	vs := []struct {
		v Value
		s string
	}{
		{v: D(1, 1), s: "1"},
		{v: D(0, 1), s: "0"},
		{v: D(-3, 1), s: "-3"},
		{v: S("x"), s: "x"},
	}
	for i, x := range vs {
		if s := x.v.String(); s != x.s {
			t.Errorf("[%d] got=%q want=%q", i, s, x.s)
		}
	}
}

func TestSimplify(t *testing.T) {
	vs := []struct {
		v []Value
		s string
	}{
		{v: Simplify(), s: "0"},
		{v: Simplify(S("x"), S("y"), D(1, 3), S("a")), s: "1/3*a*x*y"},
		{v: Simplify(D(3, 1), S("y"), D(1, 3), S("a")), s: "a*y"},
		{v: Simplify(D(3, 1), D(-1, 6), S("a"), D(2, 1)), s: "-a"},
	}
	for i, x := range vs {
		if s := Prod(x.v...); s != x.s {
			t.Errorf("[%d] got=%q want=%q", i, s, x.s)
		}
	}
}
