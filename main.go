package main

import (
	"fmt"
	"os"

	"github.com/Ahmed-Armaan/gorrent_cli/bencode_decoder"
	"github.com/Ahmed-Armaan/gorrent_cli/utils"
)

func main() {
	input := os.Args[1]
	data, err := os.ReadFile(input)
	if err != nil {
		fmt.Println("File err: File does not exist or cannot be opened")
		return
	}

	decoded, _ := bencodedecoder.Decode(data, 0)
	metaData, ok := decoded.(map[string]any)
	if !ok {
		fmt.Println("Bencode err: Invalid bencode format")
		return
	}
	metaDataInfo, ok := metaData["info"].(map[string]any)
	if !ok {
		fmt.Println("Bencode err: Invalid bencode format")
		return
	}

	hash := utils.InfoHashExtractor(data)
	fmt.Printf("%x\n", hash)
	piecesHash := utils.PiecesHashExtractor(metaDataInfo)
	for i := range piecesHash {
		fmt.Printf("%x\n", piecesHash[i])
	}
}
