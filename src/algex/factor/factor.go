// Package factor defines basic factors: rationals and symbols.
package factor

import (
	"math/big"
	"sort"
	"strings"
)

// Value captures a single factor. It is either a number or a symbol.
type Value struct {
	num *big.Rat
	sym string
}

// String displays a single factor.
func (v Value) String() string {
	if v.num != nil {
		return v.num.RatString()
	}
	if v.sym != "" {
		return v.sym
	}
	return "<ERROR>"
}

// zero is a constant zero for comparisons.
var zero = big.NewRat(0, 1)

// one is a constant one for comparisons.
var one = big.NewRat(1, 1)

// minusOne is a constant -one for comparisons.
var minusOne = big.NewRat(-1, 1)

// S converts a string into a symbol value.
func S(sym string) Value {
	return Value{sym: sym}
}

// R converts a rational value into a number value.
func R(n *big.Rat) Value {
	return Value{num: n}
}

// D converts two integers to a rational number value.
func D(num, den int64) Value {
	return Value{num: big.NewRat(num, den)}
}

// Simplify condenses an unsorted array (product) of values into a
// simplified (ordered) form.
func Simplify(vs ...Value) []Value {
	if len(vs) == 0 {
		return nil
	}

	var syms []string
	n := big.NewRat(1, 1)
	for _, v := range vs {
		if v.num != nil {
			if zero.Cmp(v.num) == 0 {
				return nil
			}
			n.Mul(n, v.num)
			continue
		}
		syms = append(syms, v.sym)
	}
	sort.Strings(syms)

	res := []Value{R(n)}
	for _, s := range syms {
		res = append(res, S(s))
	}
	return res
}

// Prod returns a string representing an product of values. This
// function does not attempt to simplify the array first.
func Prod(vs ...Value) string {
	if len(vs) == 0 {
		return "0"
	}
	var x []string
	prefix := ""
	for i, v := range vs {
		if v.num != nil && i == 0 && len(vs) != 1 {
			if one.Cmp(v.num) == 0 {
				continue
			}
			if minusOne.Cmp(v.num) == 0 {
				prefix = "-"
				continue
			}
		}
		x = append(x, v.String())
	}
	return prefix + strings.Join(x, "*")
}
