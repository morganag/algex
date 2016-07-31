// Package factor defines basic factors: rationals and symbols.
package factor

import (
	"fmt"
	"math/big"
	"sort"
	"strings"
)

// Value captures a single factor. It is either a number or a symbol.
type Value struct {
	num *big.Rat

	pow int
	sym string
}

// String displays a single factor.
func (v Value) String() string {
	if v.num != nil {
		return v.num.RatString()
	}
	if v.sym != "" {
		if v.pow == 1 {
			return v.sym
		}
		return fmt.Sprintf("%s^%d", v.sym, v.pow)
	}
	return "<ERROR>"
}

// zero is a constant zero for comparisons.
var zero = big.NewRat(0, 1)

// one is a constant one for comparisons.
var one = big.NewRat(1, 1)

// minusOne is a constant -one for comparisons.
var minusOne = big.NewRat(-1, 1)

// R copies a rational value into a number value.
func R(n *big.Rat) Value {
	c := &big.Rat{}
	return Value{num: c.Set(n)}
}

// D converts two integers to a rational number value.
func D(num, den int64) Value {
	return Value{num: big.NewRat(num, den)}
}

// S converts a string into a symbol value.
func S(sym string) Value {
	return Value{sym: sym, pow: 1}
}

// Sp converts a string, power to a symbol value.
func Sp(sym string, pow int) Value {
	if pow == 0 {
		return D(1, 1)
	}
	return Value{sym: sym, pow: pow}
}

type ByAlpha []Value

func (a ByAlpha) Len() int      { return len(a) }
func (a ByAlpha) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByAlpha) Less(i, j int) bool {
	if a[i].sym < a[j].sym {
		return true
	}
	if a[i].sym > a[j].sym {
		return false
	}
	// Higher powers first (after simplify this is moot).
	return a[i].pow > a[j].pow
}

// Simplify condenses an unsorted array (product) of values into a
// simplified (ordered) form.
func Simplify(vs ...Value) []Value {
	if len(vs) == 0 {
		return nil
	}

	var syms []Value
	n := big.NewRat(1, 1)
	for _, v := range vs {
		if v.num != nil {
			if zero.Cmp(v.num) == 0 {
				return nil
			}
			n.Mul(n, v.num)
			continue
		}
		syms = append(syms, v)
	}
	sort.Sort(ByAlpha(syms))

	res := []Value{R(n)}
	for _, s := range syms {
		i := len(res) - 1
		last := res[i]
		if last.sym != s.sym {
			res = append(res, s)
			continue
		}
		last.pow += s.pow
		if last.pow == 0 {
			res = res[:i]
			continue
		}
		res = append(res[:i], last)
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

// Segment simplifies a set of factors and returns the numerical
// coefficient, the non-numeric array of factors and a string
// representation of this array of non-numeric factors.
func Segment(vs ...Value) (*big.Rat, []Value, string) {
	x := Simplify(vs...)
	return x[0].num, x[1:], Prod(x[1:]...)
}

// Replace replaces copies of b found in a with c. The number of times b
// appeared in a is returned as well as the replaced array of factors.
func Replace(a, b, c []Value) (int, []Value) {
	pn, pf, _ := Segment(b...)
	qf := Simplify(a...)
	r := pn.Inv(pn)
	n := 0
	for len(pf) > 0 {
		var nf []Value
		i := 0
		j := 0
	GIVEUP:
		for i < len(pf) && j < len(qf) {
			t := pf[i]
			for j < len(qf) {
				u := qf[j]
				j++
				if u.num != nil || t.sym != u.sym {
					nf = append(nf, u)
					continue
				}
				if t.pow*u.pow < 0 {
					// Same symbol, but we require that
					// the sign of the power is the same.
					break GIVEUP
				}
				np := u.pow - t.pow
				if np*t.pow < 0 {
					break GIVEUP
				}
				if np != 0 {
					nf = append(nf, Sp(t.sym, np))
				}
				i++
				break
			}
		}
		if i != len(pf) {
			break
		}
		// Whole match found.
		qf = Simplify(append(append(nf, qf[j:]...), append(c, R(r))...)...)
		n++
	}
	return n, qf
}
