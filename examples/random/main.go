package main

import (
	"fmt"

	"github.com/inovacc/utils/v2/random/password"
	"github.com/inovacc/utils/v2/random/random"
)

func main() {
	// generate strong password
	randPass := password.NewPassword(
		password.WithLength(16),
		password.WithNumbers(),
		password.WithSpecial(),
		password.WithLower(),
		password.WithUpper(),
	)
	fmt.Println("Strong Password:", randPass)

	randNum, _ := random.RandomInt(1, 100)
	fmt.Println("Random Number:", randNum)
}
