// Package matrix manages matrices of expressions.
package matrix

import (
	"fmt"
	"strings"

	"algex/factor"
	"algex/terms"
)

type Matrix struct {
	// row count and col count
	rows, cols int
	// The matrix elements arranged, [r=0,c=0], [0,1], [0,2] ...
	data []*terms.Exp
}

// NewMatrix creates a rows x cols matrix.
func NewMatrix(rows, cols int) (*Matrix, error) {
	if rows <= 0 || cols <= 0 {
		return nil, fmt.Errorf("need positive dimensions, not %dx%d", rows, cols)
	}
	m := &Matrix{
		rows: rows,
		cols: cols,
		data: make([]*terms.Exp, rows*cols),
	}
	return m, nil
}

// String serializes a matrix for displaying.
func (m *Matrix) String() string {
	var rs []string
	for r := 0; r < m.rows; r++ {
		var cs []string
		for c := 0; c < m.cols; c++ {
			cs = append(cs, m.data[c+m.cols*r].String())
		}
		rs = append(rs, "["+strings.Join(cs, ", ")+"]")
	}
	return "[" + strings.Join(rs, ", ") + "]"
}

// Set sets the value of a matrix element.
func (m *Matrix) Set(row, col int, e *terms.Exp) error {
	if row < 0 || col < 0 || row >= m.rows || col >= m.cols {
		return fmt.Errorf("bad cell: [%d,%d] in %dx%d matrix", row, col, m.rows, m.cols)
	}
	m.data[col+m.cols*row] = e
	return nil
}

// El returns the row,col element of the matrix.
func (m *Matrix) El(row, col int) *terms.Exp {
	return m.data[col+m.cols*row]
}

// Identity returns a square identity matrix of dimension n.
func Identity(n int) (*Matrix, error) {
	if n <= 0 {
		return nil, fmt.Errorf("invalid identity matrix of dimension n=%d", n)
	}
	m, _ := NewMatrix(n, n)
	for i := 0; i < n; i++ {
		m.Set(i, i, terms.NewExp([]factor.Value{factor.D(1, 1)}))
	}
	return m, nil
}

// Mul multiplies m x n with conventional matrix multiplication.
func (m *Matrix) Mul(n *Matrix) (*Matrix, error) {
	if m.cols != n.rows {
		return nil, fmt.Errorf("a cols(%d) != b rows(%d)", m.cols, n.rows)
	}
	a, err := NewMatrix(m.rows, n.cols)
	if err != nil {
		return nil, err
	}
	for r := 0; r < a.rows; r++ {
		for c := 0; c < a.cols; c++ {
			var e []*terms.Exp
			for i := 0; i < m.cols; i++ {
				x, y := m.El(r, i), n.El(i, c)
				if x != nil && y != nil {
					e = append(e, terms.Mul(x, y))
				}
			}
			a.Set(r, c, terms.Add(e...))
		}
	}
	return a, nil
}
