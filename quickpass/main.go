package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
)

var CharacterSet = map[string]string{
	"LowerCase": "abcdefghijklmnopqrstuvwxyz",
	"UpperCase": "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
	"Numbers":   "0123456789",
	"Symbols":   "-_+@\\|<>?",
}

func main() {
	var input int
	fmt.Printf("Enter the password length: ")
	fmt.Scanln(&input)

	if input < 4 {
		log.Fatal("Password must be greater than 4 characters")
	}
	var l, u, s, n string
	fmt.Printf("Include lowercase? (y/n): ")
	fmt.Scanln(&l)

	fmt.Printf("Include uppercase? (y/n): ")
	fmt.Scanln(&u)

	fmt.Printf("Include numbers? (y/n): ")
	fmt.Scanln(&n)

	fmt.Printf("Include symbols? (y/n): ")
	fmt.Scanln(&s)

	var characterPool string
	if l == "y" {
		characterPool += CharacterSet["LowerCase"]
	}
	if u == "y" {
		characterPool += CharacterSet["UpperCase"]
	}
	if n == "y" {
		characterPool += CharacterSet["Numbers"]
	}
	if s == "y" {
		characterPool += CharacterSet["Symbols"]
	}

	if len(characterPool) == 0 {
		log.Fatalf("Number of characters can't be empty")
	}

	var password string
	for i := 0; i < input; i++ {
		p, err := rand.Int(rand.Reader, big.NewInt(int64(len(characterPool))))
		if err != nil {
			log.Print("Error while generating the password")
		}
		password += string(characterPool[p.Int64()])
	}
	fmt.Println(password)
}
