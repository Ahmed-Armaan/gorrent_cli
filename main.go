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

		if argv[2][:1] == "i" {
			decoded_number, _ := strconv.ParseInt(decode_number(argv[2]), 10, 64)
			fmt.Printf("%d\n", decoded_number)
		} else if _, err := strconv.Atoi(strings.Split(argv[2], ":")[0]); err == nil {
			decoded_string := decode_string(argv[2])
			fmt.Printf("%s\n", decoded_string)
		} else if argv[2][0] == 'l' && argv[2][len(argv[2])-1] == 'e' {
			decoded_list := decode_list(argv[2][1 : len(argv[2])-1])
			fmt.Printf("%s\n", decoded_list)
		} else if argv[2][0] == 'd' && argv[2][len(argv[2])-1] == 'e' {
			decoded_dictionary := decode_dictionary(argv[2][1 : len(argv[2])-1])
			for key, value := range decoded_dictionary {
				fmt.Printf("%s:%v\t", key, value)
			}
			fmt.Printf("\n")
		} else {
			fmt.Printf("\x1b[1;31mAn error occured\n")
			fmt.Printf("The provided bencode is not of proper format\n\x1b[1;0m")
			os.Exit(1)
		}
	}
}
