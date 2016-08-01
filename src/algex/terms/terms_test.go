package terms

import (
	. "algex/factor"
	"testing"
)

func TestNewExp(t *testing.T) {
	vs := []struct {
		e [][]Value
		s string
	}{
		{
			e: [][]Value{
				{D(-3, 1)},
				{D(2, 1), S("a")},
				{D(-4, 1), Sp("b", -1)},
			},
			s: "-3+2*a-4*b^-1",
		},
		{
			e: [][]Value{
				{D(-3, 1)},
				{D(2, 1)},
				{D(2, 1), S("a")},
				{D(-4, 1), S("a")},
			},
			s: "-1-2*a",
		},
		{
			e: [][]Value{
				{D(-3, 1)},
				{D(3, 1)},
			},
			s: "0",
		},
		{
			e: [][]Value{
				{D(-3, 1), S("a"), S("b")},
				{D(2, 1), Sp("a", 2)},
				{D(2, 1), S("b"), S("a")},
				{D(1, 1), S("b"), S("a")},
			},
			s: "2*a^2",
		},
	}
	for i, v := range vs {
		e := NewExp(v.e...)
		if s := e.String(); s != v.s {
			t.Errorf("[%d] got=%q want=%q", i, s, v.s)
		}
	}
}

func TestAddSub(t *testing.T) {
	a := NewExp([]Value{Sp("a", 3), D(1, 3)},
		[]Value{D(2, 3), Sp("a", 3)},
		[]Value{Sp("a", -1)})
	b := NewExp([]Value{Sp("b", 5), D(1, 3)},
		[]Value{Sp("a", -1)})
	vs := []struct {
		e *Exp
		s string
	}{
		{e: a, s: "a^-1+a^3"},
		{e: Add(a, a), s: "2*a^-1+2*a^3"},
		{e: Sub(a, a), s: "0"},
		{e: Sub(a, Add(a, a)), s: "-a^-1-a^3"},
		{e: Sub(a, b), s: "a^3-1/3*b^5"},
	}
	for i, v := range vs {
		if s := v.e.String(); s != v.s {
			t.Errorf("[%d] got=%q want=%q", i, s, v.s)
		}
	}
}

func TestMul(t *testing.T) {
	a := NewExp([]Value{Sp("a", 3)}, []Value{Sp("b", 4)})
	b := NewExp([]Value{Sp("a", 3)}, []Value{Sp("b", 4), D(-1, 1)})
	c := NewExp([]Value{Sp("a", 6)}, []Value{Sp("b", 8)})

	vs := []struct {
		e *Exp
		s string
	}{
		{e: Mul(a, a), s: "2*a^3*b^4+a^6+b^8"},
		{e: Mul(b, b), s: "-2*a^3*b^4+a^6+b^8"},
		{e: Mul(a, b), s: "a^6-b^8"},
		{e: Mul(b, a), s: "a^6-b^8"},
		{e: Mul(b, a, c), s: "a^12-b^16"},
	}
	for i, v := range vs {
		if s := v.e.String(); s != v.s {
			t.Errorf("[%d] got=%q want=%q", i, s, v.s)
		}
	}
}

func TestSubstitute(t *testing.T) {
	vs := []struct {
		e, c *Exp
		b    []Value
		s    string
	}{
		{
			e: NewExp([]Value{S("a"), S("y")}),
			b: []Value{S("y")},
			c: NewExp([]Value{S("a")}, []Value{D(-1, 1), Sp("b", 2), Sp("a", -1)}),
			s: "a^2-b^2",
		},
		{
			e: NewExp([]Value{Sp("a", 2)}),
			b: []Value{S("a")},
			c: NewExp([]Value{S("b")}, []Value{D(-1, 1), S("c")}),
			s: "-2*b*c+b^2+c^2",
		},
		{
			e: NewExp([]Value{Sp("a", 2)}, []Value{Sp("b", 2)}),
			b: []Value{S("a")},
			c: NewExp(), // Zero.
			s: "b^2",
		},
	}
	for i, v := range vs {
		r := Substitute(v.e, v.b, v.c)
		if s := r.String(); s != v.s {
			t.Errorf("[%d] %q (%q -> %q) got=%q want=%q", i, v.e, Prod(v.b...), v.c, s, v.s)
		}
	}
}
