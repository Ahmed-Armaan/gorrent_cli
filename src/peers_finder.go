package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
)

type peer_ip_port struct {
	ip   net.IP
	port byte
}

func get_peers(torrent map[string]interface{}, info map[string]interface{}, info_hash []byte) []peer_ip_port {
	query := url.Values{}
	query.Add("peer_id", "00112233445566778899")
	query.Add("port", "6881")
	query.Add("uploaded", "0")
	query.Add("downloaded", "0")
	var peer_data []peer_ip_port

	length, ok := info["length"].(int)
	if !ok {
		fmt.Println("Dictionary value inconsistent for 'length'")
		return peer_data
	}
	query.Add("left", strconv.Itoa(length))
	query.Add("compact", "1")

	tracker_url, ok := torrent["announce"].(string)
	if !ok {
		fmt.Println("Dictionary value inconsistent for 'announce'")
		return peer_data
	}

	Info_hash := url.QueryEscape(string(info_hash))

	res, err := http.Get(tracker_url + "?" + query.Encode() + "&info_hash=" + Info_hash)
	if err != nil {
		fmt.Println("Error fetching peers:", err)
		return peer_data
	}
	defer res.Body.Close()

	dc := byte_data{bufio.NewReader(res.Body)}
	decoded, err := dc.decoder()
	check_error(err)

	fmt.Println(decoded)
	peer_map, ok := decoded.(map[string]interface{})
	if !ok {
		fmt.Println("Receivedd data invalid")
		return peer_data
	}

	peers, _ := peer_map["peers"].(string)
	peer_bytes := []byte(peers)

	for i := 0; i+5 < len(peer_bytes); i += 6 {
		ip := net.IPv4(peer_bytes[i], peer_bytes[i+1], peer_bytes[i+2], peer_bytes[i+3])
		port := (peer_bytes[i+4])<<8 + (peer_bytes[i+5])

		peer_data = append(peer_data, peer_ip_port{
			ip:   ip,
			port: port,
		})
	}

	for i := range peer_data {
		fmt.Printf("%s::%d\n", peer_data[i].ip.String(), peer_data[i].port)
	}
	return peer_data
}
