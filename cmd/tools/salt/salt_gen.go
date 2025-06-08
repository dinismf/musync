package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
)

func main() {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		_, err := fmt.Fprintf(os.Stderr, "Failed to generate random salt: %v\n", err)
		if err != nil {
			return
		}
		os.Exit(1)
	}
	fmt.Println(hex.EncodeToString(salt))
}
