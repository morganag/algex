package matrix

import (
	"testing"

	"algex/factor"
	"algex/terms"
)

func TestNewMatrix(t *testing.T) {
	one, err := Identity(2)
	if err != nil {
		t.Fatalf("failed to make a 2x2 identity matrix!: %v", err)
	}
	if got, want := one.String(), "[[1, 0], [0, 1]]"; got != want {
		t.Errorf("one failed: got=%q, want=%q", got, want)
	}
}

func TestMul(t *testing.T) {
	a, err := Identity(2)
	if err != nil {
		t.Fatalf("failed to make 2x2 identity: %v", err)
	}
	b, err := Identity(2)
	if err != nil {
		t.Fatalf("failed to make 2x2 identity: %v", err)
	}
	b.Set(0, 1, terms.Mul(a.El(0, 0), terms.NewExp([]factor.Value{factor.Sp("x", 2)})))

	c, err := a.Mul(b)
	if err != nil {
		t.Fatalf("failed to multiply 2x2 matrices: %v", err)
	}
	if got, want := c.String(), b.String(); got != want {
		t.Errorf("matrix multiply %v*%v: got=%v, want=%v", a, b, got, want)
	}
	d, err := b.Mul(a)
	if err != nil {
		t.Fatalf("failed to multiply 2x2 matrices: %v", err)
	}
	if got, want := d.String(), b.String(); got != want {
		t.Errorf("matrix multiply %v*%v: got=%v, want=%v", a, b, got, want)
	}
}

func TestSum(t *testing.T) {
	a := terms.NewExp([]factor.Value{factor.D(-2, 3), factor.Sp("x", -1)})
	b := terms.NewExp([]factor.Value{factor.D(9, 4), factor.Sp("x", 2)})

	p, _ := NewMatrix(1, 2)
	p.Set(0, 0, a)
	q, _ := NewMatrix(1, 2)
	q.Set(0, 0, a)
	q.Set(0, 1, b)

	r, err := p.Sum(q, terms.Mul(a, a))
	if err != nil {
		t.Fatalf("can't sum things: %v", err)
	}
	if got, want := r.String(), "[[-2/3*x^-1-8/27*x^-3, 1]]"; got != want {
		t.Errorf("add: got=%q, want=%q", got, want)
	}
}
