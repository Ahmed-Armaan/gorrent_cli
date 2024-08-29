package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"net"
)

type peer_message struct {
	message_length_prefix uint32
	message_id            byte
	index                 uint32
	begin                 uint32
	length                uint32
}

func receive_and_check_messsage_id(buf []byte, conn net.Conn, expected_id byte) bool {
	var Peer_message peer_message
	Peer_message.message_length_prefix = binary.BigEndian.Uint32(buf)
	payload_buf := make([]byte, Peer_message.message_length_prefix)

	_, err := conn.Read(payload_buf)
	check_error(err)
	Peer_message.message_id = payload_buf[0]

	fmt.Printf("received message : %v\n", Peer_message)
	if Peer_message.message_id != expected_id {
		fmt.Println("bitfied message unexpected")
		return false
	} else {
		return true
	}
}

func download_piece(conn net.Conn, info map[string]interface{}, index int) []byte {
	buf := make([]byte, 4)
	_, err := conn.Read(buf)
	check_error(err)

	if !receive_and_check_messsage_id(buf, conn, 5) {
		return nil
	}

	intrested_message := []byte{0, 0, 0, 1, 2}
	_, err = conn.Write(intrested_message)
	check_error(err)

	buf = make([]byte, 4)
	_, err = conn.Read(buf)
	check_error(err)

	if !receive_and_check_messsage_id(buf, conn, 1) {
		return nil
	}

	piece_size, ok := info["piece length"].(int)
	if !ok {
		fmt.Println("Pieces Length format unexpected")
		return nil
	}

	block_size := 16 * 1024
	block_cnt := int(math.Ceil(float64(piece_size) / float64(block_size)))

	var data []byte
	for i := 0; i < block_cnt; i++ {
		block_len := block_size
		if i == block_cnt-1 {
			block_len = piece_size - ((block_cnt - 1) * int(block_size))
		}

		Peer_message := peer_message{
			message_length_prefix: 13,
			message_id:            6,
			index:                 uint32(index),
			begin:                 uint32(i * int(block_size)),
			length:                uint32(block_len),
		}

		var buffer bytes.Buffer
		binary.Write(&buffer, binary.BigEndian, Peer_message)
		_, err = conn.Write(buffer.Bytes())
		check_error(err)

		buf = make([]byte, 4)
		_, err = conn.Read(buf)
		check_error(err)

		Peer_message = peer_message{}
		Peer_message.message_length_prefix = binary.BigEndian.Uint32(buf)
		payload := make([]byte, Peer_message.message_length_prefix)
		_, err = io.ReadFull(conn, payload)
		check_error(err)

		Peer_message.message_id = payload[0]
		fmt.Printf("received message : %v\n", Peer_message)

		data = append(data, payload[9:]...)
	}
	conn.Close()
	return data
}
