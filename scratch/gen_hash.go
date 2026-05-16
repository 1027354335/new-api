package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	password := "daklwdjawlkdjikwdaodjawio1*dkwa1"
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	fmt.Printf("Password: %s\nHash: %s\n", password, string(hash))
}
