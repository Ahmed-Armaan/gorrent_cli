package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func decode_number(str string) string {
	number_lenght := strings.Index(str, "e")
	return str[1:number_lenght]
}

func decode_string(str string) string {
	str_arr := strings.Split(str, ":")
	len_string, _ := strconv.ParseInt(str_arr[0], 10, 64)
	final_string := str_arr[1][:len_string]
	return final_string
}

func decode_list(str string) []string {
	var decoded_list []string
	var l int = 0
	for i := 0; i < len(str); {
		if str[i] >= 'i' { // decode the i<num>e
			l = i
			for {
				i++
				if i > len(str) {
					break
				}
				if str[i] == 'e' {
					decoded_list = append(decoded_list, str[l+1:i])
					break
				}
			}
			i++
		} else if str[i] >= '0' && str[i] <= '9' { // decode the <num>:<string>
			l = i
			for {
				i++
				if i > len(str) {
					break
				}
				if str[i] == ':' {
					str_len, err := strconv.ParseInt(string(str[l:i]), 10, 64)
					if err != nil {
						fmt.Printf("An error encountered")
						os.Exit(1)
					}
					if i+1+int(str_len) > len(str) {
						fmt.Printf("\x1b[1;31mAn error occured\n")
						fmt.Printf("String length: index [%d] out of bound\n\x1b[1;0m", i+1+int(str_len))
					}
					decoded_list = append(decoded_list, str[i+1:i+1+int(str_len)])
					i = i + int(str_len)
					break
				}
			}
			i++
		} else {
			i++
		}
	}
	return decoded_list
}

func decode_dictionary(str string) map[string]string {
	decoded_dictionary := make(map[string]string)
	item_list := decode_list(str)
	for i := 0; i < len(item_list); i += 2 {
		decoded_dictionary[item_list[i]] = item_list[i+1]
	}
	return decoded_dictionary
}
