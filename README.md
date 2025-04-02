# utils [![Test](https://github.com/inovacc/utils/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/inovacc/utils/actions/workflows/test.yml)

Collection of utility functions and packages for various cryptographic and encoding operations. The project is organized
into several packages, each serving a specific purpose.

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

* HashPassword(password string, p *Params) (string, error): Generates a secure Argon2ID hash for the given password and
  parameters.
* CheckPasswordHash(encoded, password string) (bool, error): Compares a plain-text password with a stored hash (JSON
  encoded).

### crypto/password/bcrypt

This package provides functions for hashing and verifying passwords using the bcrypt algorithm.

* HashPassword(password string) (string, error): Hashes a password using bcrypt.
* CheckPasswordHash(password, hash string) bool: Checks if a plain-text password matches a bcrypt hash.

### crypto/hash

This package provides functions for hashing data using the SHA-256 algorithm.

* HashSHA256(data string) string: Hashes a string using SHA-256 and returns the hexadecimal representation.
* HashSHA256Bytes(data []byte) string: Hashes a byte slice using SHA-256 and returns the hexadecimal representation.

### file

This package provides functions for reading from and writing to files.

* WriteToFile(filename string, data []byte) error: Writes data to a file.
* ReadFromFile(filename string) ([]byte, error): Reads data from a file.

```go
package main

import (
  "fmt"
  "github.com/inovacc/utils/v2/file"
)

func main() {
  filename:= "myfile.txt"
  data, err := file.ReadFromFile(filename)
  if err != nil {
    panic(err)
  }

  fmt.Println(data)
}
```

### compress

This package provides functions for compressing and decompressing data using various algorithms.

* Compress(data []byte) ([]byte, error): Compresses data using the specified algorithm.
* Decompress(data []byte) ([]byte, error): Decompresses data using the specified algorithm.

```go
package main

import (
  "fmt"
  "github.com/inovacc/utils/v2/compress"
)

func main() {
  v := compress.NewCompress(compress.TypeZip, []byte("test"))
  data, err := v.Compress()
  if err != nil {
    panic(err)
  }

  fmt.Println(data)
}
```

### data/country/country/br/cpf

This package provides functions for generating, formatting, and validating Brazilian CPF numbers.

* GenerateCPF() string: Generates a valid CPF number (unformatted).
* FormatCPF(cpf string) string: Formats a CPF as XXX.XXX.XXX-XX.
* UnformatCPF(cpf string) string: Removes all non-numeric characters from the CPF.
* ValidateCPF(value string) bool: Checks if a CPF number is valid.
* Origin(value string) string: Returns the issuing region based on the 9th digit.

### data/country/country/br/cnpj

This package provides functions for generating, formatting, and validating Brazilian CNPJ numbers.

* GenerateCNPJ() string: Generates a valid alphanumeric CNPJ.
* ValidateCNPJ(cnpj string) bool: Checks whether an alphanumeric CNPJ is valid.
* FormatCNPJ(cnpj string) string: Formats an alphanumeric CNPJ in the pattern "12.ABC.345/01DE-XX".
* UnformatCNPJ(cnpj string) string: Removes formatting from an alphanumeric CNPJ.

### data/country/country/br/cep

This package provides functions for validating Brazilian postal codes (CEP).

* IsValidCEP(input string) bool: Validates a Brazilian postal code (CEP).

### rand/password

This package provides functions for generating random passwords with various options.

* NewPassword(opts ...Option) *Password: Creates a new Password instance with the given options.
* WithLength(length int) Option: Sets the length of the password.
* WithNumbers(enabled bool) Option: Enables or disables numbers in the password.
* WithSpecial(enabled bool) Option: Enables or disables special characters in the password.
* WithLower(enabled bool) Option: Enables or disables lowercase letters in the password.
* WithUpper(enabled bool) Option: Enables or disables uppercase letters in the password.
* (p *Password) Generate() (string, error): Generates a password based on the specified options.

```go
newPassword := password.NewPassword(
  password.WithLength(16),
  password.WithNumbers(true),
  password.WithSpecial(true),
  password.WithLower(true),
  password.WithUpper(true),
)

generated, err := newPassword.Generate()
if err != nil {
  ...
}
```

## To install

```sh
go get -u github.com/inovacc/utils/v2
```

## Usage

Import the necessary packages in your Go code and use the provided functions as needed. For example:

```go
randomString := rand.RandomString(10)
```

## License

This project is licensed under the MIT License. See the LICENSE file for more details.
