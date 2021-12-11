package shamir

import (
	"testing"
)

func TestPolynomialEval_Invalid(t *testing.T) {
	r := eval([]byte{1}, 0)
	if r != 1 {
		t.Errorf("Invalid result: %v", r)
	}
}
