# Shamir's Secret Sharing

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
    shares, err := Split(5, 3, secret)
    if err != nil {
        panic(err)
    }

    // Reconstruct secret from shares
    reconstructed, err := Combine([][]byte{shares[0], shares[2], shares[4]})
    if err != nil {
        panic(err)
    }

    // secret == reconstructed
}
```
