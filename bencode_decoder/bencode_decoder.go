package bencodedecoder

import (
	"fmt"
	"strconv"
	"unicode"
)

func Decode(bencode []byte, pos int) (any, int) {
	switch {
	case unicode.IsDigit(rune(bencode[pos])):
		var colonPos int
		for i, c := range bencode[pos:] {
			if c == ':' {
				colonPos = i
				break
			}
		}

		stringSize, err := strconv.Atoi(string(bencode[pos : pos+colonPos]))
		if err != nil {
			fmt.Println("Decoding error occured")
			return err, -1
		}

		return bencode[pos+colonPos+1 : pos+colonPos+1+stringSize], pos + colonPos + 1 + stringSize

	case bencode[pos] == 'i':
		num, nextPos, sign := 0, pos, 1
		if bencode[pos+1] == '-' {
			sign = -1
			pos += 2
		} else {
			pos++
		}

		for i, c := range bencode[pos:] {
			nextPos = i
			if c == 'e' {
				break
			}
			num *= 10
			num += int(c - '0')
		}

		return num * sign, pos + nextPos + 1

	case bencode[pos] == 'l':
		var list []any
		pos++
		for pos < len(bencode) && bencode[pos] != 'e' {
			nextElement, nextPos := Decode(bencode, pos)
			if nextPos == -1 {
				fmt.Println("List Decoding error")
				return nil, -1
			}

			list = append(list, nextElement)
			pos = nextPos
			if pos >= len(bencode) {
				fmt.Println("Invalid list: 'e' not found")
				return nil, -1
			}
		}

		return list, pos + 1

	case bencode[pos] == 'd':
		dictionary := make(map[string]any)
		kvPair := make([]any, 0, 2)
		pos++

		for pos < len(bencode) && bencode[pos] != 'e' {
			nextElement, nextPos := Decode(bencode, pos)
			if nextPos == -1 {
				fmt.Println("Dictionary decoding error")
				return nil, -1
			}

			kvPair = append(kvPair, nextElement)
			//if len(kvPair) == 2 {
			//	keyStr, ok := kvPair[0].(string)
			//	if !ok {
			//		fmt.Println("Dictionary key is not a string")
			//		return nil, -1
			//	}
			//	dictionary[keyStr] = kvPair[1]
			//	kvPair = kvPair[:0]
			//}

			if len(kvPair) == 2 {
				key, ok := kvPair[0].([]byte)
				if !ok {
					fmt.Println("Dictionary key is not a string")
					return nil, -1
				}
				dictionary[string(key)] = kvPair[1]
				kvPair = kvPair[:0]
			}

			pos = nextPos
			if pos > len(bencode) {
				fmt.Println("Invalid dictionary: 'e' not found")
				return nil, -1
			}
		}

		if pos >= len(bencode) || bencode[pos] != 'e' {
			fmt.Println("Invalid dictionary: missing 'e'")
			return nil, -1
		}

		return dictionary, pos + 1

	default:
		fmt.Println("An error occured: Bencode structure improper")
		return nil, -1
	}
}
