package main

import (
	"fmt"

	"github.com/inovacc/utils/v2/crypto/hashing"
	"github.com/inovacc/utils/v2/crypto/password"
)

func main() {
	newHasher := hashing.NewHasher(hashing.SHA256)
	fmt.Println("SHA256 Hash:", newHasher.HashString("hello world"))

	passBcrypt, _ := password.HashPasswordBcrypt("secret")
	fmt.Println("Bcrypt Password Hash:", passBcrypt)

	passArgon2, _ := password.HashPasswordArgon2("secret", nil)
	fmt.Println("Argon2 Password Hash:", passArgon2)
}
