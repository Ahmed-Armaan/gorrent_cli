package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	argc := len(os.Args)
	argv := os.Args

	if argc > 1 && argv[1] == "decode" {
		//Todo: open torrent file and extract the metadata
		torrent_data := argv[2]
		if torrent_data[:1] == "i" {
			decoded_number := decode_list(argv[2])
			fmt.Printf("%s\n", decoded_number)
		} else if _, err := strconv.Atoi(strings.Split(torrent_data, ":")[0]); err == nil {
			decoded_string := decode_list(argv[2])
			fmt.Printf("%s\n", decoded_string)
		} else if torrent_data[0] == 'l' && torrent_data[len(torrent_data)-1] == 'e' {
			decoded_list := decode_list(torrent_data[1 : len(torrent_data)-1])
			fmt.Printf("%s\n", decoded_list)
		} else if torrent_data[0] == 'd' && torrent_data[len(torrent_data)-1] == 'e' {
			decoded_dictionary := decode_dictionary(torrent_data[1 : (int)(len(argv[2]))-1])
			for i := 0; i < len(decoded_dictionary); i++ {
				fmt.Print(decoded_dictionary[i].key)
				fmt.Print(":")
				fmt.Println(decoded_dictionary[i].value)
			}
			fmt.Printf("\n")
		} else {
			fmt.Printf("\x1b[1;31mAn error occured\n")
			fmt.Printf("The provided bencode is not of proper format\n\x1b[1;0m")
			os.Exit(1)
		}
	}
}
