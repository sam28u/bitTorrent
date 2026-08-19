package p2p

import (
	"bytes"
	"fmt"
	"net"
	"strconv"
	"time"

	"bittorrent/handshake"
	"bittorrent/tracker"
)

type Client struct {
	Conn     net.Conn
	Peer     tracker.Peer
	Infohash [20]byte
	PeerID   [20]byte
}

func Connect(peer tracker.Peer, infohash, peerID [20]byte) (*Client, error) {
	address := net.JoinHostPort(peer.IP.String(), strconv.Itoa(int(peer.Port)))
	
	conn, err := net.DialTimeout("tcp", address, 3*time.Second)
	if err != nil {
		return nil, err
	}
	
	req := handshake.New(infohash, peerID)
	_, err = conn.Write(req.Serialize())
	if err != nil {
		conn.Close()
		return nil, err
	}
	
	res, err := handshake.Read(conn)
	if err != nil {
		conn.Close()
		return nil, err
	}
	
	if !bytes.Equal(res.InfoHash[:], infohash[:]) {
		conn.Close()
		return nil, fmt.Errorf("infohash mismatch")
	}
	
	return &Client{
		Conn:     conn,
		Peer:     peer,
		Infohash: infohash,
		PeerID:   peerID,
	}, nil
}