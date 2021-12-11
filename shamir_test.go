package shamir

import (
	"bytes"
	"errors"
	"testing"
)

func TestSplit_valid(t *testing.T) {
	secret := []byte("example")

	shares, err := Split(secret, 5, 3)
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

	if _, err := Split(secret, 5, 0); !errors.Is(err, ErrInvalidThreshold) {
		t.Fatal("Expected ErrInvalidCount")
	}

	if _, err := Split(secret, 0, 3); !errors.Is(err, ErrInvalidCount) {
		t.Fatal("Expected ErrInvalidCount")
	}

	if _, err := Split(secret, 400, 300); !errors.Is(err, ErrInvalidThreshold) {
		t.Fatal("Expected ErrInvalidCount")
	}

	if _, err := Split(secret, 300, 3); !errors.Is(err, ErrInvalidCount) {
		t.Fatal("Expected ErrInvalidCount")
	}

	if _, err := Split(nil, 5, 3); !errors.Is(err, ErrEmptySecret) {
		t.Fatal("Expected ErrEmptySecret")
	}
}

func TestCombine_valid(t *testing.T) {
	secret := []byte("example")

	shares, err := Split(secret, 5, 3)
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
				reconstructed, err := Combine(shares[i], shares[j], shares[k])
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}

				if !bytes.Equal(secret, reconstructed) {
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

	if _, err := Combine([]byte("exam"), []byte("ple")); err == nil {
		t.Fatal("Expected ErrInvalidShares")
	}

	if _, err := Combine([]byte("a"), []byte("b")); err == nil {
		t.Fatal("Expected ErrInvalidShares")
	}

	if _, err := Combine([]byte("aa"), []byte("aa")); err == nil {
		t.Fatal("Expected ErrInvalidShares")
	}
}
