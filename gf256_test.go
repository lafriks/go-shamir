package shamir

import (
	"testing"
)

func TestDivisionZero(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The division by zero did not panic")
		}
	}()

	_ = div(0, 0)
}
