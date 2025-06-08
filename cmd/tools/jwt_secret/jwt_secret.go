package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
)

func main() {
	secret := make([]byte, 32)
	_, err := rand.Read(secret)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to generate random secret: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(base64.RawURLEncoding.EncodeToString(secret))
}
