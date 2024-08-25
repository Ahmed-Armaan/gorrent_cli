package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"unicode"
)

type byte_data struct {
	*bufio.Reader
}

var Err_improper_file = errors.New("Improper bencode type")
var Err_non_string_key = errors.New("The key must be a string")

func (b *byte_data) decoder() (interface{}, error) {
	next_byte, err := b.ReadByte()
	if err != nil {
		return nil, err
	}

	if next_byte == 'i' {
		str, err := b.ReadString('e')
		if err != nil {
			return nil, err
		}

		return strconv.Atoi(str[:len(str)-1])
	} else if unicode.IsDigit(rune(next_byte)) {
		str, err := b.ReadString(':')
		if err != nil {
			return nil, err
		}

		str = string(next_byte) + str[:len(str)-1]
		len, _ := strconv.Atoi(str)

		s := make([]byte, len)
		_, err = b.Read(s)
		if err != nil {
			return nil, err
		}

		return string(s), nil
	} else if next_byte == 'l' {
		var list []interface{}

		for {
			item, err := b.decoder()
			if err != nil {
				return nil, err
			}

			if item == nil {
				break
			}
			list = append(list, item)
		}
		return list, nil
	} else if next_byte == 'd' {
		m := make(map[string]interface{})

		for {
			key, err := b.decoder()
			if err != nil {
				return nil, err
			}
			if key == nil {
				break
			}

			strkey, ok := key.(string)
			if !ok {
				return nil, Err_non_string_key
			}

			val, err := b.decoder()
			if err != nil {
				return nil, err
			}
			if val == nil {
				break
			}

			m[strkey] = val
		}
		return m, nil
	} else if next_byte == 'e' {
		return nil, nil
	} else {
		return nil, Err_improper_file
	}
}

type encoded_data struct {
	*bytes.Buffer
}

func (b *encoded_data) encoder(val interface{}) error {
	switch v := val.(type) {
	case string:
		b.WriteString(fmt.Sprintf("%d:%s", len(v), v))
		return nil
	case int:
		b.WriteString(fmt.Sprintf("i%de", v))
		return nil
	case []interface{}:
		b.WriteByte('l')
		for _, item := range v {
			b.encoder(item)
		}
		b.WriteByte('e')

		return nil
	case map[string]interface{}:
		keys := make([]string, 0, len(v))
		for k := range v {
			keys = append(keys, k)
		}

		sort.Strings(keys)

		b.WriteByte('d')

		for _, k := range keys {
			b.encoder(k)
			b.encoder(v[k])
		}

		b.WriteByte('e')

		return nil
	default:
		return Err_improper_file
	}
}
