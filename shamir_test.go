package shamir

import (
	"bytes"
	"errors"
	"testing"
)

func TestSplit_valid(t *testing.T) {
	secret := []byte("example")

	shares, err := Split(5, 3, secret)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(shares) != 5 {
		t.Fatalf("Expected 5 shares but got %d", len(shares))
	}

	for _, share := range shares {
		if len(share) != len(secret)+1 {
			t.Fatalf("Invalid split: %v", shares)
		}
	}
}

func TestSplit_invalid(t *testing.T) {
	secret := []byte("example")

	if _, err := Split(5, 0, secret); !errors.Is(err, ErrInvalidThreshold) {
		t.Fatal("Expected ErrInvalidCount")
	}

	if _, err := Split(0, 3, secret); !errors.Is(err, ErrInvalidCount) {
		t.Fatal("Expected ErrInvalidCount")
	}

	if _, err := Split(400, 300, secret); !errors.Is(err, ErrInvalidThreshold) {
		t.Fatal("Expected ErrInvalidCount")
	}

	if _, err := Split(300, 3, secret); !errors.Is(err, ErrInvalidCount) {
		t.Fatal("Expected ErrInvalidCount")
	}

	if _, err := Split(5, 3, nil); !errors.Is(err, ErrEmptySecret) {
		t.Fatal("Expected ErrEmptySecret")
	}
}

func TestCombine_valid(t *testing.T) {
	secret := []byte("example")

	shares, err := Split(5, 3, secret)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			if j == i {
				continue
			}
			for k := 0; k < 5; k++ {
				if k == i || k == j {
					continue
				}
				reconstructed, err := Combine([][]byte{shares[i], shares[j], shares[k]})
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}

				if bytes.Compare(secret, reconstructed) != 0 {
					t.Fatalf("Expected '%s' but got '%s' (i:%d j:%d k:%d shares: %v)", secret, reconstructed, i, j, k, shares)
				}
			}
		}
	}
}

func TestCombine_invalid(t *testing.T) {
	if _, err := Combine(nil); !errors.Is(err, ErrInvalidShares) {
		t.Fatal("Expected ErrInvalidShares")
	}

	if _, err := Combine([][]byte{[]byte("exam"),[]byte("ple")}); err == nil {
		t.Fatal("Expected ErrInvalidShares")
	}

	if _, err := Combine([][]byte{[]byte("a"),[]byte("b")}); err == nil {
		t.Fatal("Expected ErrInvalidShares")
	}

	if _, err := Combine([][]byte{[]byte("aa"),[]byte("aa")}); err == nil {
		t.Fatal("Expected ErrInvalidShares")
	}
}