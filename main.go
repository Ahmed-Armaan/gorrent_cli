package main

import (
	"fmt"
	"os"

	"github.com/Ahmed-Armaan/gorrent_cli/bencode_decoder"
)

func main() {
	input := os.Args[1]
	decoded, _ := bencodedecoder.Decode(input, 0)
	fmt.Println(decoded)
}
