// Package shamir implements Shamir's Secret Sharing algorithm over GF(2^8).
//
// Shamir's Secret Sharing algorithm allows you to securely share a secret with
// N people, allowing the recovery of that secret if K of those people combine
// their shares.
//
// It begins by encoding a secret as a number (e.g., 42), and generating N
// random polynomial equations of degree K-1 which have an X-intercept equal to
// the secret. Given K=3, the following equations might be generated:
//
//     f1(x) =  78x^2 +  19x + 42
//     f2(x) = 128x^2 + 171x + 42
//     f3(x) = 121x^2 +   3x + 42
//     f4(x) =  91x^2 +  95x + 42
//     etc.
//
// These polynomials are then evaluated for values of X > 0:
//
//     f1(1) =  139
//     f2(2) =  896
//     f3(3) = 1140
//     f4(4) = 1783
//     etc.
//
// These (x, y) pairs are the shares given to the parties. In order to combine
// shares to recover the secret, these (x, y) pairs are used as the input points
// for Lagrange interpolation, which produces a polynomial which matches the
// given points. This polynomial can be evaluated for f(0), producing the secret
// value--the common x-intercept for all the generated polynomials.
//
// If fewer than K shares are combined, the interpolated polynomial will be
// wrong, and the result of f(0) will not be the secret.
//
// This package constructs polynomials over the field GF(2^8) for each byte of
// the secret, allowing for fast splitting and combining of anything which can
// be encoded as bytes.
//
// This package has not been audited by cryptography or security professionals.
package shamir

import (
	crand "crypto/rand"
	"encoding/binary"
	"errors"
	mrand "math/rand"
)

var (
	// ErrInvalidCount is returned when the count parameter is invalid.
	ErrInvalidCount = errors.New("shares must be more or equal to treshold but not more than 255")
	// ErrInvalidThreshold is returned when the threshold parameter is invalid.
	ErrInvalidThreshold = errors.New("treshold must be at least 2 but not more than 255")
	// ErrEmptySecret is returned when provided secret is empty.
	ErrEmptySecret = errors.New("secret can not be empty")
	// ErrInvalidShares is returned when not required minimum shares are provided or shares does not have same length.
	ErrInvalidShares = errors.New("at least 2 shares are required and must have same length")
)

type cryptoSource [8]byte

func (s *cryptoSource) Int63() int64 {
	_, err := crand.Read(s[:])
	if err != nil {
		panic(err)
	}
	return int64(binary.BigEndian.Uint64(s[:]) & (1<<63 - 1))
}

func (s *cryptoSource) Seed(seed int64) {
	panic("seed")
}

// Split the given secret into N shares of which K are required to recover the
// secret. Returns an array of shares.
func Split(secret []byte, n, k int) ([][]byte, error) {
	if k <= 1 || k > 255 {
		return nil, ErrInvalidThreshold
	}

	if n < k || n > 255 {
		return nil, ErrInvalidCount
	}

	if len(secret) == 0 {
		return nil, ErrEmptySecret
	}

	rnd := mrand.New(&cryptoSource{})
	cords := rnd.Perm(255)

	shares := make([][]byte, n)
	for i := range shares {
		shares[i] = make([]byte, len(secret)+1)
		shares[i][len(secret)] = byte(cords[i]) + 1
	}

	for i, b := range secret {
		p, err := generate(byte(k)-1, b)
		if err != nil {
			return nil, err
		}

		for j := 0; j < n; j++ {
			x := byte(cords[j]) + 1
			shares[j][i] = eval(p, x)
		}
	}

	return shares, nil
}

// Combine the given shares into the original secret.
func Combine(shares ...[]byte) ([]byte, error) {
	if len(shares) < 2 || len(shares[0]) < 2 {
		return nil, ErrInvalidShares
	}

	l := len(shares[0])
	c := make(map[byte]bool, len(shares))
	c[shares[0][l-1]] = true
	for i := 1; i < len(shares); i++ {
		if len(shares[i]) != l {
			return nil, ErrInvalidShares
		}
		if ok := c[shares[i][l-1]]; ok {
			return nil, ErrInvalidShares
		}
		c[shares[i][l-1]] = true
	}

	secret := make([]byte, l-1)

	points := make([]pair, len(shares))
	for i := range secret {
		for p, v := range shares {
			points[p] = pair{x: v[l-1], y: v[i]}
		}
		secret[i] = interpolate(points, 0)
	}

	return secret, nil
}
