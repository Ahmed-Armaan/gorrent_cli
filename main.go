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
	lenght, ok := metaDataInfo["piece length"].(int)
	if !ok {
		fmt.Println("Bencode err: Invalid bencode format")
		return
	}
	url, ok := metaData["announce"].([]byte)
	if !ok {
		fmt.Println("2Bencode err: Invalid bencode format")
		return
	}
	//piecesHash := utils.PiecesHashExtractor(metaDataInfo)

	conn := utils.GetPeers(hash, lenght, string(url))
	downloadedPiece := utils.DownloadPiece(conn, lenght, 0)

	err = os.WriteFile("piece_0.bin", downloadedPiece, 0644)
	if err != nil {
		fmt.Println("Failed to save file:", err)
		return
	}
	fmt.Println("Piece saved to piece_0.bin")
	//fmt.Println(interval)
	//for i := range peers {
	//	fmt.Printf("%s, %s\n", ports[i], peers[i])
	//}
}
