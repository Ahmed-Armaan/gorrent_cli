package utils

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"github.com/Ahmed-Armaan/gorrent_cli/bencode_decoder"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
)

func GetPeers(infoHash [20]byte, pieceLen int, announceURL string) net.Conn {
	myId, err := generatePeerID()
	if err != nil {
		fmt.Println("Error: peerId cannot be generated")
		return nil
	}

	params := url.Values{}
	params.Add("info_hash", string(infoHash[:]))
	params.Add("peer_id", myId)
	params.Add("port", "6881")
	params.Add("uploaded", "0")
	params.Add("downloaded", "0")
	params.Add("left", strconv.Itoa(pieceLen))
	params.Add("compact", "1")

	fullUrl := announceURL + "?" + params.Encode()
	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		fmt.Println("Error: cannot create a request")
		return nil
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Error: cannot call GET request")
		return nil
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response body")
		return nil
	}

	peers_decoded, _ := bencodedecoder.Decode(body, 0)

	peers_map, ok := peers_decoded.(map[string]any)
	if !ok {
		fmt.Println("Bencode err: Invalid bencode format")
		return nil
	}

	//	interval, ok := peers_map["interval"].(int)
	//	if !ok {
	//		fmt.Println("Error: Peers access invalid")
	//		return nil, nil, -1
	//	}

	peersEntry, ok := peers_map["peers"].([]byte)
	if !ok {
		fmt.Println("Error: Peers access invalid")
		return nil
	}

	peers, ports := extractPeersPort(peersEntry)
	if peers == nil || ports == nil {
		return nil
	}

	return handshake(peers, ports, infoHash[:], myId, pieceLen)
	//return peers, ports, interval
}

func generatePeerID() (string, error) {
	var alphabet = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]byte, 20)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		return "", err
	}

	for i := range b {
		b[i] = alphabet[int(b[i])%len(alphabet)]
	}
	return string(b), nil
}

func extractPeersPort(peersEntry []byte) ([]string, []string) {
	if len(peersEntry)%6 != 0 {
		fmt.Println("Error: peers and ports cannot be exracted")
		return nil, nil
	}

	var peers []string
	var ports []string

	for i := 0; i < len(peersEntry); i += 6 {
		peer := net.IPv4(peersEntry[i], peersEntry[i+1], peersEntry[i+2], peersEntry[i+3]).String()
		port := binary.BigEndian.Uint16(peersEntry[i+4 : i+6])
		peers = append(peers, peer)
		ports = append(ports, strconv.Itoa(int(port)))
	}

	return peers, ports
}

func handshake(ips []string, ports []string, infoHash []byte, myId string, pieceLen int) net.Conn {
	ip := ips[0]
	port := ports[0]
	address := net.JoinHostPort(ip, port)

	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("Error: handshake could not be made")
		return nil
	}

	protocol_str_len := byte(19)
	protocol_str := []byte("BitTorrent protocol")
	reserved := make([]byte, 8)
	handshake := append([]byte{protocol_str_len}, protocol_str...)
	handshake = append(handshake, reserved...)
	handshake = append(handshake, infoHash...)
	handshake = append(handshake, myId...)

	_, err = conn.Write(handshake)
	buffer := make([]byte, 68)
	_, err = conn.Read(buffer)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	//DownloadPiece(conn, pieceLen, 0)
	fmt.Printf("PEER ID = %x\n", (buffer[48:]))
	return conn
}
