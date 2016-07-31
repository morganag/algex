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
