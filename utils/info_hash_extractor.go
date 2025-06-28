package utils

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"strconv"
)

func InfoHashExtractor(bencode []byte) [20]byte {
	infoDictionary := infoExtractor(bencode)
	if infoDictionary != nil {
		hash := sha1.Sum(infoDictionary)
		return hash
	}
	return [20]byte{}
}

func infoExtractor(bencode []byte) []byte {
	pattern := []byte("4:infod")
	startPos := bytes.Index(bencode, pattern)
	if startPos == -1 {
		fmt.Println("Error: invalid bencode, info dictionary not found")
		return nil
	}

	startPos += len(pattern) - 1
	//	if bencode[startPos] != 'd' {
	//		fmt.Println("Error: invald dictionary, info dictionary not available")
	//		return nil
	//	}

	depth := 1
	for i := startPos + 1; i < len(bencode); {
		switch bencode[i] {
		case 'd', 'l', 'i':
			depth++
			i++
		case 'e':
			depth--
			i++
			if depth == 0 {
				return bencode[startPos:i]
			}
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			lengthStart := i
			for bencode[i] != ':' {
				i++
			}
			length, _ := strconv.Atoi(string(bencode[lengthStart:i]))
			i++
			i += length
		default:
			i++
		}
	}

	fmt.Println("Error: invalid bencode, unmatched dictionary 'e'")
	return nil
}
