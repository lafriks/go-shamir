package shamir

import (
	"crypto/rand"
	"io"
)

// evaluate the polynomial at the given point
func eval(p []byte, x byte) (result byte) {
	if x == 0 {
		return p[0]
	}

	// Horner's scheme
	for i := 1; i <= len(p); i++ {
		result = mul(result, x) ^ p[len(p)-i]
	}
	return
}

// builds a degree-length polynomial with the given x-intercept, using mid as its random middle
// coefficients and top as its (non-zero) leading coefficient
func buildPolynomial(degree, x byte, mid []byte, top byte) []byte {
	result := make([]byte, degree+1)
	result[0] = x
	copy(result[1:degree], mid)
	result[degree] = top
	return result
}

// generates the random coefficients needed to build num degree-length polynomials: num*(degree-1)
// middle coefficients and num non-zero leading coefficients. The randomness is bulk-read up front
// so this costs a handful of crypto/rand reads total, instead of two per polynomial.
func generateCoefficients(num int, degree byte) (mid, top []byte, err error) {
	mid = make([]byte, num*int(degree-1))
	if _, err = io.ReadFull(rand.Reader, mid); err != nil {
		return nil, nil, err
	}

	top = make([]byte, num)
	if _, err = io.ReadFull(rand.Reader, top); err != nil {
		return nil, nil, err
	}

	// the leading coefficient can't be zero, or else it's a lower degree polynomial
	var buf [1]byte
	for i, b := range top {
		for b == 0 {
			if _, err = io.ReadFull(rand.Reader, buf[:]); err != nil {
				return nil, nil, err
			}
			b = buf[0]
		}
		top[i] = b
	}

	return mid, top, nil
}

// an input/output pair
type pair struct {
	x, y byte
}

// Lagrange interpolation
func interpolate(points []pair, x byte) (value byte) {
	for i, a := range points {
		weight := byte(1)
		for j, b := range points {
			if i != j {
				top := x ^ b.x
				bottom := a.x ^ b.x
				factor := div(top, bottom)
				weight = mul(weight, factor)
			}
		}
		value = value ^ mul(weight, a.y)
	}
	return
}
