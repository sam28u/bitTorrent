package tracker

import (
	"bittorrent/torrentfile"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/jackpal/bencode-go"
)

type Peer struct {
	IP   net.IP
	Port uint16
}

type bencodeTrackerResp struct {
	Interval int    `bencode:"interval"`
	Peers    string `bencode:"peers"`
}

func GeneratePeerID() ([20]byte, error) {
	var peerID [20]byte
	_, err := rand.Read(peerID[:])
	if err != nil {
		return [20]byte{}, err
	}
	return peerID, err
}

func buildTrackerURL(t *torrentfile.TorrentFile , peerID [20]byte , port uint16) (string , error){
	base , err := url.Parse(t.Announce)
	if err != nil {
		return "" , err
	}
	params := url.Values{
		"info_hash" : []string{string(t.InfoHash[:])},
		"peer_id" : []string{string(peerID[:])},
		"port" : []string{strconv.Itoa(int(port))},
		"uploaded" : []string{"0"},
		"downloaded" : []string{"0"},
		"compact" : []string{"1"},
		"left" : []string{strconv.Itoa(t.Length)},
	}
	base.RawQuery = params.Encode()
	return base.String() , nil
}

func RequestPeers(t *torrentfile.TorrentFile , peerID [20]byte, port uint16)([]Peer ,error){
	url , err := buildTrackerURL(t , peerID , port)
	if err != nil{
		return nil , err
	} 
	client := &http.Client{Timeout : 15 * time.Second}

	resp , err := client.Get(url)
	if err != nil{
		return nil , err
	}

	defer resp.Body.Close()

	trackerResp := bencodeTrackerResp{}
	err = bencode.Unmarshal(resp.Body , &trackerResp)
	if err != nil{
		return nil , err
	}
	return parsePeers([]byte(trackerResp.Peers))
}

func parsePeers(peersBin []byte) ([]Peer , error){
	const peerSize = 6

	if len(peersBin)%peerSize != 0{
		return nil , fmt.Errorf("Received malformed peers list")
	}

	numPeers := len(peersBin) / peerSize
	peers := make([]Peer , numPeers)

	for i := 0; i < numPeers; i++ {
		offset := i * peerSize
		peers[i].IP = net.IP(peersBin[offset : offset + 4])
		peers[i].Port = binary.BigEndian.Uint16(peersBin[offset + 4 : offset + 6])
	}
	return peers , nil
}