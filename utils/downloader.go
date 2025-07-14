package utils

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"net"
)

type RequestMessage struct {
	Index  uint32
	Begin  uint32
	Length uint32
}

func expectMessage(conn net.Conn, expectedID byte) bool {
	buf := make([]byte, 4)
	_, err := io.ReadFull(conn, buf)
	if err != nil {
		fmt.Println("Error: failed to read payload")
		return false
	}

	length := binary.BigEndian.Uint32(buf)
	if length == 0 {
		fmt.Println("Error: received message empty")
		return false
	}

	payload := make([]byte, length)
	_, err = io.ReadFull(conn, payload)
	if err != nil {
		fmt.Println("Error: failed to read payload")
		return false
	}

	msgID := payload[0]
	fmt.Printf("Received message ID: %d (expected %d)\n", msgID, expectedID)
	return msgID == expectedID
}

func DownloadPiece(conn net.Conn, pieceLength int, pieceIndex int) []byte {
	msgStatus := expectMessage(conn, 5)
	if msgStatus == false {
		return nil
	}

	// Send "interested" (ID = 2)
	interested := []byte{0, 0, 0, 1, 2}
	if _, err := conn.Write(interested); err != nil {
		return nil
	}

	// Wait for "unchoke" (ID = 1)
	msgStatus = expectMessage(conn, 1)
	if msgStatus == false {
		fmt.Println("Peer did not unchoke")
		return nil
	}

	blockSize := 16 * 1024
	blockCount := int(math.Ceil(float64(pieceLength) / float64(blockSize)))
	var pieceData []byte

	for i := range blockCount {
		currBlockSize := blockSize
		if i == blockCount-1 {
			currBlockSize = pieceLength - (i * blockSize)
		}

		req := RequestMessage{
			Index:  uint32(pieceIndex),
			Begin:  uint32(i * blockSize),
			Length: uint32(currBlockSize),
		}

		var buffer bytes.Buffer
		binary.Write(&buffer, binary.BigEndian, uint32(13))
		buffer.WriteByte(6) // message ID
		binary.Write(&buffer, binary.BigEndian, req)

		if _, err := conn.Write(buffer.Bytes()); err != nil {
			fmt.Println("Error: failed to send request")
			return nil
		}

		lengthBuf := make([]byte, 4)
		if _, err := io.ReadFull(conn, lengthBuf); err != nil {
			fmt.Println("Error: failed to read response length")
			return nil
		}

		length := binary.BigEndian.Uint32(lengthBuf)
		payload := make([]byte, length)
		if _, err := io.ReadFull(conn, payload); err != nil {
			fmt.Println("Error: failed to read response payload")
			return nil
		}

		if payload[0] != 7 {
			fmt.Println("Error: unexpected message ID")
			return nil
		}

		pieceData = append(pieceData, payload[9:]...) // skip ID, index, begin
	}

	conn.Close()
	return pieceData
}
