# Shamir's Secret Sharing

[![Go Reference](https://pkg.go.dev/badge/github.com/lafriks/go-shamir.svg)](https://pkg.go.dev/github.com/lafriks/go-shamir)
[![Coverage Status](https://coveralls.io/repos/github/lafriks/go-shamir/badge.svg?branch=main)](https://coveralls.io/github/lafriks/go-shamir?branch=main)

Based on [github.com/codahale/sss](https://github.com/codahale/sss)

A pure Go implementation of [Shamir's Secret Sharing algorithm](http://en.wikipedia.org/wiki/Shamir's_Secret_Sharing)

## Usage

```sh
go get -u github.com/lafriks/go-shamir
```

## Example

```go
package main

import (
    "fmt"

    "github.com/lafriks/go-shamir"
)

func main() {
    secret := []byte("example")

    // Split secret to 5 shares and require 3 shares to reconstruct secret
    shares, err := shamir.Split(secret, 5, 3)
    if err != nil {
        panic(err)
    }

    // Reconstruct secret from shares
    reconstructed, err := shamir.Combine(shares[0], shares[2], shares[4])
    if err != nil {
        panic(err)
    }

    // secret == reconstructed
}
```
