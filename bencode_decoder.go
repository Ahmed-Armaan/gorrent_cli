package main

import (
	"fmt"
	"os"
	"strconv"
)

type dictionary struct {
	key   interface{}
	value interface{}
}

func decode_list(str string) []interface{} {
	var decoded_list []interface{}
	var l int = 0
	for i := 0; i < len(str); {
		if str[i] == 'i' { // decode the i<num>e
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
		} else if str[i] == 'l' || str[i] == 'd' { // handling a list or dictionary inside a list
			l, e_left := i, 0
			for {
				i++
				if i >= len(str) {
					break
				}
				if str[i] == 'i' {
					e_left++
				}
				if str[i] == 'e' {
					if e_left > 0 {
						e_left--
					} else if e_left == 0 {
						if str[l] == 'l' {
							decoded_list = append(decoded_list, decode_list(str[l+1:i]))
						} else if str[i] == 'd' {
							decoded_list = append(decoded_list, decode_dictionary(str[l+1:i]))
						}
						break
					}
				}
			}
			i++
		} else {
		}
	}
	return decoded_list
}

func decode_dictionary(str string) []dictionary {
	items_list := decode_list(str)
	if len(items_list)%2 != 0 {
		fmt.Printf("\x1b[1;31mAn error occured\n")
		fmt.Printf("The dictionary does not have a key-value pair\n\x1b[1;0m")
		os.Exit(1)
	}
	decoded_dictionary := make([]dictionary, len(items_list)/2)
	for i := 0; i < len(items_list); i += 2 {
		decoded_dictionary[i/2].key = items_list[i]
		decoded_dictionary[i/2].value = items_list[i+1]
	}
	return decoded_dictionary
}
