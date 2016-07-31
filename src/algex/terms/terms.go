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
		e.insert(n, fs, s)
	}
	return e
}

// String represents an expression of terms as a string.
func (e *Exp) String() string {
	if len(e.terms) == 0 {
		return "0"
	}
	var s []string
	for x := range e.terms {
		s = append(s, x)
	}
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
