package CFF

import (
	"errors"
	"fmt"
	"gf2n"
)

//d-CFF((dk+1)q, q^(k+1)) w-CFF(m, 2^n)
//assume that q=2^a, a<64
// w = d
// m = (dk+1)q
// n = a(k+1)
// k: highest degree of polynomial f over F_q
// X: F_{dk+1} x F_q
// B_f: {(x, f(x))}
//\cal{B}: {B_f} where f's degree is equal or less than k
type CFF struct {
	k     uint64
	d     uint64
	field *gf2n.GF2nField
}

// (x, f(x)) x\in F_q
type CFFBSetElement struct {
	argument *gf2n.GF2nElement
	image    *gf2n.GF2nElement
}

//generate a CFF instance
func NewCFF(a uint8, k, d, irre uint64) (*CFF, error) {
	if k < 1 || 1<<a < d*k+1 {
		return nil, errors.New("Invalid parameters")
	}

	field, err := gf2n.NewGF2nField(a, irre)
	if err != nil {
		return nil, err
	}
	return &CFF{k, d, field}, nil
}

//b: binary string with a specific length
func (c *CFF) FindBSet(b string) ([]*CFFBSetElement, error) {
	n := c.field.ExtDegree()
	if uint64(len(b)) != uint64(n)*(c.k+1) {
		return nil, errors.New("The input string has a wrong length")
	}

	//select f using b, generate k+1 n-bit coefficients
	coeff := make([]*gf2n.GF2nElement, c.k+1, c.k+1)
	res := uint64(0)
	for i, item := range b {
		if item == '1' {
			res += 1 << (n - 1 - uint8(i)%n)
		}
		if i > 0 && uint8(i)%n == n-1 {
			coeff[uint8(i)/n] = c.field.NewGF2nElement(res)
			res = uint64(0)
		}
	}
	f := c.field.NewGF2nPoly(coeff)

	//construct {(x, f(x))} x\in [0, d*k+1]
	len := c.d*c.k + 1
	bSet := make([]*CFFBSetElement, len, len)
	for i := uint64(0); i < len; i++ {
		x := c.field.NewGF2nElement(i)
		bSet[i] = &CFFBSetElement{x, f.Eval(x)}
	}

	return bSet, nil
}

//convert a CFFBSetElement to a 128-bit string
func (e *CFFBSetElement) String() string {
	str := fmt.Sprintf("%064b", e.argument.Value())
	str += fmt.Sprintf("%064b", e.image.Value())
	return str
}
