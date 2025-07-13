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

func GetPeers(infoHash [20]byte, fileLen int, announceURL string) ([]string, []string, int) {
	peerId, err := generatePeerID()
	if err != nil {
		fmt.Println("Error: peerId cannot be generated")
		return nil, nil, -1
	}

	params := url.Values{}
	params.Add("info_hash", string(infoHash[:]))
	params.Add("peer_id", peerId)
	params.Add("port", "6881")
	params.Add("uploaded", "0")
	params.Add("downloaded", "0")
	params.Add("left", strconv.Itoa(fileLen))
	params.Add("compact", "1")

	fullUrl := announceURL + "?" + params.Encode()
	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		fmt.Println("Error: cannot create a request")
		return nil, nil, -1
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Error: cannot call GET request")
		return nil, nil, -1
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response body")
		return nil, nil, -1
	}

	peers_decoded, _ := bencodedecoder.Decode(body, 0)

	peers_map, ok := peers_decoded.(map[string]any)
	if !ok {
		fmt.Println("Bencode err: Invalid bencode format")
		return nil, nil, -1
	}

	interval, ok := peers_map["interval"].(int)
	if !ok {
		fmt.Println("Error: Peers access invalid")
		return nil, nil, -1
	}

	peersEntry, ok := peers_map["peers"].([]byte)
	if !ok {
		fmt.Println("Error: Peers access invalid")
		return nil, nil, -1
	}

	peers, ports := extractPeersPort(peersEntry)
	if peers == nil || ports == nil {
		return nil, nil, -1
	}
	return peers, ports, interval
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
