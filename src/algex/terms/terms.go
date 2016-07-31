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

// String represents an expression of terms as a string.
func (e *Exp) String() string {
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

// NewExp creates a new expression.
func NewExp(ts ...[]factor.Value) *Exp {
	e := &Exp{
		terms: make(map[string]term),
	}
	for _, t := range ts {
		n, fs, s := factor.Segment(t...)
		old, ok := e.terms[s]
		if !ok {
			e.terms[s] = term{
				coeff: n,
				fact:  fs,
			}
			continue
		}
		// Combine with existing term.
		old.coeff = n.Add(n, e.terms[s].coeff)
		e.terms[s] = old
	}
	return e
}
