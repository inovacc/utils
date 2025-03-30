# utils [![Test](https://github.com/inovacc/utils/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/inovacc/utils/actions/workflows/test.yml)

Collection of utility functions and packages for various cryptographic and encoding operations. The project is organized into several packages, each serving a specific purpose.

# Packages

### crypto/rand

This package provides functions for generating random strings, integers, and bytes.

* RandomString(n int) string: Generates a random string of length n using a secure random number generator.
* RandomInt(min, max int) int: Generates a random integer between min and max using a secure random number generator.
* RandomBytes(n uint32) []byte: Generates a random byte slice of length n.

### encode

This package provides functions for encoding and decoding data using various encoding schemes.

* Base64Encode(data []byte) string: Encodes the given data to a Base64 string.
* Base64Decode(data string) ([]byte, error): Decodes the given Base64 string to a byte slice.
* Base62Encode(data []byte) string: Encodes the given data to a Base62 string.
* Base62Decode(data string) ([]byte, error): Decodes the given Base62 string to a byte slice.
* Base58Encode(data []byte) string: Encodes the given data to a Base58 string.
* Base58Decode(data string) ([]byte, error): Decodes the given Base58 string to a byte slice.

### crypto/password/argon2

This package provides functions for hashing and verifying passwords using the Argon2ID algorithm.

* HashPassword(password string, p *Params) (string, error): Generates a secure Argon2ID hash for the given password and parameters.
* CheckPasswordHash(encoded, password string) (bool, error): Compares a plain-text password with a stored hash (JSON encoded).

### crypto/password/bcrypt

This package provides functions for hashing and verifying passwords using the bcrypt algorithm.

* HashPassword(password string) (string, error): Hashes a password using bcrypt.
* CheckPasswordHash(password, hash string) bool: Checks if a plain-text password matches a bcrypt hash.

### crypto/hash

This package provides functions for hashing data using the SHA-256 algorithm.

* HashSHA256(data string) string: Hashes a string using SHA-256 and returns the hexadecimal representation.
* HashSHA256Bytes(data []byte) string: Hashes a byte slice using SHA-256 and returns the hexadecimal representation.

## To install

```sh
go get -u github.com/inovacc/utils/v2
```

## Usage
Import the necessary packages in your Go code and use the provided functions as needed. For example:

```go
package main

import (
	"fmt"
	"github.com/inovacc/utils/v2/crypto/rand"
)

func main() {
	randomString := rand.RandomString(10)
	fmt.Println("Random String:", randomString)
}
```

## License
This project is licensed under the MIT License. See the LICENSE file for more details.
