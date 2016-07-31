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
