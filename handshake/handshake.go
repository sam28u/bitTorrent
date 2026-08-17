package handshake

import "bytes"

type Handshake struct {
	Pstr     string
	InfoHash [20]byte
	PeerID   [20]byte
}

func new(infoHash, peerID [20]byte) *Handshake {
	return &Handshake{
		Pstr:     "BitTorrent Protocol",
		InfoHash: infoHash,
		PeerID:   peerID,
	}
}

func (h *Handshake) Serialize() []byte {
	var buf bytes.Buffer
	buf.WriteByte(byte(len(h.Pstr)))
	buf.WriteString(h.Pstr)
	buf.Write(make([]byte, 8))
	buf.Write(h.InfoHash[:])
	buf.Write(h.PeerID[:])
	return buf.Bytes()
}