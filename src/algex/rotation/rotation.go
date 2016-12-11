// Package rotation generates matrices for 3D rotations.
//
// This package prefixes angles with 's', 'c' or 't' for sine, cosine
// and tangent respectively.
package rotation

import (
	"algex/factor"
	"algex/matrix"
	"algex/terms"
)

// A matrix for rotating anticlockwise around the X-axis.
func RX(theta string) *matrix.Matrix {
	m, _ := matrix.NewMatrix(3, 3)
	c := []factor.Value{factor.S("c" + theta)}
	s := []factor.Value{factor.S("s" + theta)}
	mS := append([]factor.Value{factor.D(-1, 1)}, s...)
	one := []factor.Value{factor.D(1, 1)}
	m.Set(0, 0, terms.NewExp(one))
	m.Set(1, 1, terms.NewExp(c))
	m.Set(1, 2, terms.NewExp(mS))
	m.Set(2, 1, terms.NewExp(s))
	m.Set(2, 2, terms.NewExp(c))
	return m
}

// A matrix for rotating anticlockwise around the Y-axis.
func RY(theta string) *matrix.Matrix {
	m, _ := matrix.NewMatrix(3, 3)
	c := []factor.Value{factor.S("c" + theta)}
	s := []factor.Value{factor.S("s" + theta)}
	mS := append([]factor.Value{factor.D(-1, 1)}, s...)
	one := []factor.Value{factor.D(1, 1)}
	m.Set(0, 0, terms.NewExp(c))
	m.Set(0, 2, terms.NewExp(s))
	m.Set(1, 1, terms.NewExp(one))
	m.Set(2, 0, terms.NewExp(mS))
	m.Set(2, 2, terms.NewExp(c))
	return m
}

// A matrix for rotating anticlockwise around the Z-axis.
func RZ(theta string) *matrix.Matrix {
	m, _ := matrix.NewMatrix(3, 3)
	c := []factor.Value{factor.S("c" + theta)}
	s := []factor.Value{factor.S("s" + theta)}
	mS := append([]factor.Value{factor.D(-1, 1)}, s...)
	one := []factor.Value{factor.D(1, 1)}
	m.Set(0, 0, terms.NewExp(c))
	m.Set(0, 1, terms.NewExp(mS))
	m.Set(1, 0, terms.NewExp(s))
	m.Set(1, 1, terms.NewExp(c))
	m.Set(2, 2, terms.NewExp(one))
	return m
}
