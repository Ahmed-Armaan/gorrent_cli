package main

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
)

func check_error(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	file, err := os.Open(os.Args[1])
	check_error(err)
	defer file.Close()

	dc := byte_data{bufio.NewReader(file)}
	decoded, err := dc.decoder()
	check_error(err)

	m, ok := decoded.(map[string]interface{})
	if !ok {
		fmt.Println("expected data to be a dictionary, found %T", decoded)
		return
	}

	info, ok := m["info"].(map[string]interface{})
	if !ok {
		fmt.Println("expected info block to be a dictionary, found %T", m["info"])
		return
	}

	buffer := bytes.Buffer{}
	ec := encoded_data{&buffer}
	err = ec.encoder(info)
	check_error(err)

	h := sha1.New()
	io.Copy(h, &buffer)
	sum := h.Sum(nil)

	fmt.Printf("Tracker URL = %s\n", m["announce"])
	fmt.Println("length = ", info["length"])
	fmt.Println("Pieces length = ", info["piece length"])
	fmt.Printf("Hash = %x\n", sum)
	fmt.Printf("pieces hash = %x\n", info["pieces"])

	peers_data := get_peers(m, info, sum)
	for i := range peers_data {
		fmt.Printf("%s::%d\n", peers_data[i].ip.String(), peers_data[i].port)
	}

	conn := establish_connection(peers_data, sum)
	out_path := os.Args[2]

	data := download_piece(conn, info, 0)

	out_file, err := os.Create(out_path)
	check_error(err)
	defer out_file.Close()

	_, err = out_file.Write(data)
	check_error(err)

	fmt.Println("file downloaded")
}
