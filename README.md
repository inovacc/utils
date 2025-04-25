# utils [![Test](https://github.com/inovacc/utils/actions/workflows/test.yml/badge.svg)](https://github.com/inovacc/utils/actions/workflows/test.yml)

Collection of utility functions and packages for various cryptographic and encoding operations. The project is organized
into several packages, each serving a specific purpose.

## To install

```sh
go get -u github.com/inovacc/utils/v2
```

# Packages

### random

This package provides functions for generating random strings, integers, and bytes.

* RandomString(n int) string: Generates a random string of length n using a secure random number generator.
* RandomInt(min, max int) int,err: Generates a random integer between min and max using a secure random number
  generator.
* RandomBytes(n uint32) []byte,err: Generates a random byte slice of length n.

```go
// Generate a random string of length 10
randomString := random.RandomString(10)
fmt.Println("Random String:", randomString)

// Generate a random integer between 1 and 100
randomInt, _ := random.RandomInt(1, 100)
fmt.Println("Random Integer:", randomInt)

// Generate a random byte slice of length 16
randomBytes, _ := random.RandomBytes(16)
fmt.Println("Random Bytes:", randomBytes)
```

### encoder

This package provides functions for encoding and decoding data using various encoding schemes.

* Base64Encode(data []byte) string: Encodes the given data to a Base64 string.
* Base64Decode(data string) ([]byte, error): Decodes the given Base64 string to a byte slice.
* Base62Encode(data []byte) string: Encodes the given data to a Base62 string.
* Base62Decode(data string) ([]byte, error): Decodes the given Base62 string to a byte slice.
* Base58Encode(data []byte) string: Encodes the given data to a Base58 string.
* Base58Decode(data string) ([]byte, error): Decodes the given Base58 string to a byte slice.

```go
// Base62 Encoding and Decoding
data := []byte("Hello, World!")
encoded := encoder.Base64Encode(data)
fmt.Println("Base64 Encoded:", encoded)

// Base62 Decoding
decoded, err := encoder.Base64Decode(encoded)
if err != nil {
panic(err)
}
fmt.Println("Base64 Decoded:", string(decoded))
```

### cryptohelper/password

This package provides functions for hashing and verifying passwords using the Argon2ID algorithm.

* HashPassword(password string, p *Params) (string, error): Generates a secure Argon2ID hash for the given password and
  parameters.
* CheckPasswordHash(encoded, password string) (bool, error): Compares a plain-text password with a stored hash (JSON
  encoded).
*
* HashPassword(password string) (string, error): Hashes a password using bcrypt.
* CheckPasswordHash(password, hash string) bool: Checks if a plain-text password matches a bcrypt hash.

```go
password := "mySecurePassword"
params := &password.Params{
Memory:      64 * 1024,
Iterations:  3,
Parallelism: 2,
SaltLength:  16,
KeyLength:   32,
}

hash, err := password.HashPassword(password, params)
if err != nil {
panic(err)
}
fmt.Println("Hashed Password:", hash)

match, err := password.CheckPasswordHash(hash, password)
if err != nil {
panic(err)
}
fmt.Println("Password Match:", match)
```

```go
password := "mySecurePassword"
hash, err := bcrypt.HashPassword(password)
if err != nil {
panic(err)
}
fmt.Println("Hashed Password:", hash)

match := bcrypt.CheckPasswordHash(password, hash)
fmt.Println("Password Match:", match)
```

### crypto/hash

This package provides functions for hashing data using the SHA-256 algorithm.

* HashSHA256(data string) string: Hashes a string using SHA-256 and returns the hexadecimal representation.
* HashSHA256Bytes(data []byte) string: Hashes a byte slice using SHA-256 and returns the hexadecimal representation.

```go
data := "Hello, World!"
hashed := hash.Sha256(data)
fmt.Println("SHA-256 Hash:", hashed)
```

### file

This package provides functions for reading from and writing to files.

* WriteToFile(filename string, data []byte) error: Writes data to a file.
* ReadFromFile(filename string) ([]byte, error): Reads data from a file.

```go
filename := "myfile.txt"
data := []byte("Hello, World!")

if err = file.WriteToFile(filename, data); err != nil {
panic(err)
}

readData, err := file.ReadFromFile(filename)
if err != nil {
panic(err)
}
fmt.Println("Read Data:", string(readData))
```

### compress

This package provides functions for compressing and decompressing data using various algorithms.

* Compress(data []byte) ([]byte, error): Compresses data using the specified algorithm.
* Decompress(data []byte) ([]byte, error): Decompresses data using the specified algorithm.

```go
data := []byte("Hello, World!")
compressor := compress.NewCompress(compress.TypeZip)

compressedData, err := compressor.Compress(data)
if err != nil {
panic(err)
}
fmt.Println("Compressed Data:", compressedData)

decompressedData, err := compressor.Decompress(compressedData)
if err != nil {
panic(err)
}
fmt.Println("Decompressed Data:", string(decompressedData))
```

### data/country/country/br/cpf

This package provides functions for generating, formatting, and validating Brazilian CPF numbers.

* GenerateCPF() string: Generates a valid CPF number (unformatted).
* FormatCPF(cpf string) string: Formats a CPF as XXX.XXX.XXX-XX.
* UnformatCPF(cpf string) string: Removes all non-numeric characters from the CPF.
* ValidateCPF(value string) bool: Checks if a CPF number is valid.
* Origin(value string) string: Returns the issuing region based on the 9th digit.

```go
cpfNumber := cpf.GenerateCPF()
fmt.Println("Generated CPF:", cpfNumber)

formattedCPF := cpf.FormatCPF(cpfNumber)
fmt.Println("Formatted CPF:", formattedCPF)

isValid := cpf.ValidateCPF(cpfNumber)
fmt.Println("Is CPF Valid:", isValid)

origin := cpf.Origin(cpfNumber)
fmt.Println("CPF Origin:", origin)
```

### data/country/country/br/cnpj

This package provides functions for generating, formatting, and validating Brazilian CNPJ numbers.

* GenerateCNPJ() string: Generates a valid alphanumeric CNPJ.
* ValidateCNPJ(cnpj string) bool: Checks whether an alphanumeric CNPJ is valid.
* FormatCNPJ(cnpj string) string: Formats an alphanumeric CNPJ in the pattern "12.ABC.345/01DE-XX".
* UnformatCNPJ(cnpj string) string: Removes formatting from an alphanumeric CNPJ.

```go
cnpjNumber := cnpj.GenerateCNPJ()
fmt.Println("Generated CNPJ:", cnpjNumber)

formattedCNPJ := cnpj.FormatCNPJ(cnpjNumber)
fmt.Println("Formatted CNPJ:", formattedCNPJ)

isValid := cnpj.ValidateCNPJ(cnpjNumber)
fmt.Println("Is CNPJ Valid:", isValid)
```

### data/country/country/br/cep

This package provides functions for validating Brazilian postal codes (CEP).

* IsValidCEP(input string) bool: Validates a Brazilian postal code (CEP).

```go
cepCode := "12345678"
isValid := cep.IsValidCEP(cepCode)
fmt.Println("Is CEP Valid:", isValid)
```

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
panic(err)
}
fmt.Println("Generated Password:", generated)
```

### reflection

This package provides utilities to manipulate Go structs using reflection.

* MergeZeroFields(dst, src any): Merge fields from src into dst only if fields in dst are zero-valued.
* CopyExportedFields(dst, src any): Copies all exported fields from src to dst regardless of zero-value.
* ZeroStruct(v any): Sets all settable fields of the struct to their zero values.
* StructToMap(input any) map[string]any: Converts a struct into a map of its exported fields.

```go
type MyStruct struct {
Name string
Age  int
}

s1 := MyStruct{Name: "Alice"}
s2 := MyStruct{Age: 30}

reflection.MergeZeroFields(&s2, &s1)
fmt.Printf("%+v\n", s2) // Output: {Name:Alice Age:30}
```

### hash

This package provides hashing utilities.

* Sha256(string) string: Returns the SHA256 hash of a string.
* Sha256Bytes([]byte) string: Returns the SHA256 hash of byte slice.

### file

This package provides basic file I/O utilities.

* WriteToFile(filename string, data []byte) error: Saves data to a file.
* ReadFromFile(filename string) ([]byte, error): Reads data from a file.

```go
err := file.WriteToFile("example.txt", []byte("Hello"))
data, err := file.ReadFromFile("example.txt")
fmt.Println(string(data)) // Output: Hello
```

### encode

* EncodeGob(data any) ([]byte, error): Encodes any Go value into binary using gob.
* DecodeGob(data []byte, v any) error: Decodes gob data into a Go value.

## License

This project is licensed under the MIT License. See the LICENSE file for more details.
