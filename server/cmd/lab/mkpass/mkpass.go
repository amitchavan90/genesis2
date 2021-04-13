package main

import (
	"fmt"
	"genesis/crypto"
	"os"
	"path"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("generate password for db use")
		fmt.Printf("usage: %s abc\n", path.Base(os.Args[0]))
		return
	}
	in := os.Args[1]
	out := crypto.HashPassword(in)
	fmt.Println(out)
}
