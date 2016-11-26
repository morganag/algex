// Package terms abstracts sums of products of factors.
package terms

import (
	"algex/factor"
	"math/big"
	"sort"
	"strings"
)

// term is a product of a coefficient and a set of non-numerical factors.
type term struct {
	coeff *big.Rat
	fact  []factor.Value
}

// Exp is a an expression or sum of terms.
type Exp struct {
	terms map[string]term
}

// NewExp creates a new expression.
func NewExp(ts ...[]factor.Value) *Exp {
	e := &Exp{
		terms: make(map[string]term),
	}
	for _, t := range ts {
		n, fs, s := factor.Segment(t...)
		if n == nil {
			continue
		}
		e.insert(n, fs, s)
	}
	return e
}

// String represents an expression of terms as a string.
func (e *Exp) String() string {
	if e == nil {
		return "0"
	} else if len(e.terms) == 0 {
		return "0"
	}
	var s []string
	for x := range e.terms {
		s = append(s, x)
	}
	// TODO: might want to prefer a non-ascii sorted expression.
	sort.Strings(s)
	for i, x := range s {
		f := e.terms[x]
		v := []factor.Value{factor.R(f.coeff)}
		t := factor.Prod(append(v, f.fact...)...)
		if i != 0 && t[0] != '-' {
			s[i] = "+" + t
		} else {
			s[i] = t
		}
	}
	return strings.Join(s, "")
}

// insert merges a coefficient, a product of factors to an expression
// indexed by s.
func (e *Exp) insert(n *big.Rat, fs []factor.Value, s string) {
	old, ok := e.terms[s]
	if !ok {
		e.terms[s] = term{
			coeff: n,
			fact:  fs,
		}
		return
	}
	// Combine with existing term.
	old.coeff = n.Add(n, e.terms[s].coeff)
	if old.coeff.Cmp(&big.Rat{}) == 0 {
		delete(e.terms, s)
		return
	}
	e.terms[s] = old
}

// Add adds together expressions. With only one argument, Add is a
// simple duplicate function.
func Add(as ...*Exp) *Exp {
	e := &Exp{
		terms: make(map[string]term),
	}
	for _, a := range as {
		for s, t := range a.terms {
			m := &big.Rat{}
			e.insert(m.Set(t.coeff), t.fact, s)
		}
	}
	return e
}

// Sub subtracts b from a into a new expression.
func Sub(a, b *Exp) *Exp {
	e := &Exp{
		terms: make(map[string]term),
	}
	for s, t := range a.terms {
		m := &big.Rat{}
		e.insert(m.Set(t.coeff), t.fact, s)
	}
	for s, t := range b.terms {
		m := big.NewRat(-1, 1)
		e.insert(m.Mul(m, t.coeff), t.fact, s)
	}
	return e
}

// Mod takes a numerical integer factor and eliminates obvious
// multiples of it from an expression. No attempt is made to
// simplify non-integer fractions.
func (e *Exp) Mod(x factor.Value) *Exp {
	if !x.IsNum() || !x.Num().IsInt() {
		return e
	}
	zero := &big.Int{}
	a := &Exp{terms: make(map[string]term)}
	for s, v := range e.terms {
		if !v.coeff.IsInt() {
			a.terms[s] = v
			continue
		}
		t := big.NewInt(1)
		u := big.NewInt(1)
		_, d := t.DivMod(v.coeff.Num(), x.Num().Num(), u)
		if d.Cmp(zero) == 0 {
			// Drop this term since it is a multiple of x.
			continue
		}
		r := &big.Rat{}
		r.SetInt(u)
		a.terms[s] = term{
			coeff: r,
			fact:  v.fact,
		}
	}
	return a
}

// Mul computes the product of a series of expressions.
func Mul(as ...*Exp) *Exp {
	var e *Exp
	for i, a := range as {
		if i == 0 {
			e = Add(a)
			continue
		}
		f := &Exp{
			terms: make(map[string]term),
		}
		for _, p := range a.terms {
			for _, q := range e.terms {
				x := []factor.Value{factor.R(p.coeff), factor.R(q.coeff)}
				n, fs, s := factor.Segment(append(x, append(p.fact, q.fact...)...)...)
				f.insert(n, fs, s)
			}
		}
		e = f
	}
	return e
}

// Substitute replaces each occurrence of b in an expression with the expression c.
func Substitute(e *Exp, b []factor.Value, c *Exp) *Exp {
	s := [][]factor.Value{}
	for _, t := range c.terms {
		s = append(s, append([]factor.Value{factor.R(t.coeff)}, t.fact...))
	}
	z := []factor.Value{factor.R(&big.Rat{})} // Zero.
	for {
		again := false
		f := &Exp{
			terms: make(map[string]term),
		}
		for _, x := range e.terms {
			a := append([]factor.Value{factor.R(x.coeff)}, x.fact...)
			hit, y := factor.Replace(a, b, z, 1)
			if hit == 0 {
				n, fs, tag := factor.Segment(y...)
				f.insert(n, fs, tag)
				// If nothing substituted, then only insert once.
				continue
			}
			if len(s) == 0 {
				// If we are substituting 0 then we won't need anything.
				continue
			}
			again = true
			for _, t := range s {
				_, y := factor.Replace(a, b, t, 1)
				n, fs, tag := factor.Segment(y...)
				f.insert(n, fs, tag)
			}
		}
		e = f
		if !again {
			break
		}
	}
	return e
}
