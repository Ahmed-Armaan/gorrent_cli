package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
)

func establish_connection(peers_data []peer_ip_port, info_hash []byte) net.Conn {
	ip := peers_data[0].ip
	port := peers_data[0].port
	address := net.JoinHostPort(ip.String(), strconv.Itoa(int(port)))

	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	protocol_str_len := byte(19)
	protocol_str := []byte("BitTorrent protocol")
	reserved := make([]byte, 8)
	peer_id := []byte("00112233445566778899")
	handshake := append([]byte{protocol_str_len}, protocol_str...)
	handshake = append(handshake, reserved...)
	handshake = append(handshake, info_hash...)
	handshake = append(handshake, peer_id...)

	_, err = conn.Write(handshake)
	buffer := make([]byte, 68)
	_, err = conn.Read(buffer)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("PEER ID = %x\n", (buffer[48:]))
	return conn
}

//length of the protocol string (BitTorrent protocol) which is 19 (1 byte)
//the string BitTorrent protocol (19 bytes)
//eight reserved bytes, which are all set to zero (8 bytes)
//sha1 infohash (20 bytes) (NOT the hexadecimal representation, which is 40 bytes long)
//peer id (20 bytes) (you can use 00112233445566778899 for this challenge)
