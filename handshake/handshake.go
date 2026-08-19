package handshake

import (
	"bytes"
	"fmt"
	"io"
)

type Handshake struct {
	Pstr     string
	InfoHash [20]byte
	PeerID   [20]byte
}

func New(infoHash, peerID [20]byte) *Handshake {
	return &Handshake{
		Pstr:     "BitTorrent protocol",
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

func Read(r io.Reader) (*Handshake , error){
	buf := make([]byte , 1)
	_, err := io.ReadFull(r , buf)
	if err != nil {
		return nil , err
	}
	pstrlen := int(buf[0])
	if pstrlen == 0{
		return nil , fmt.Errorf("pstrlen length cannot be 0")
	}
	HandshakeBuf := make([]byte , pstrlen+48)
	_, err = io.ReadFull(r , HandshakeBuf)
	if err != nil {
		return nil , err
	}
	var infoHash , peerID [20]byte
	copy(infoHash[:] , HandshakeBuf[pstrlen + 8 : pstrlen + 28])
	copy(peerID[:] , HandshakeBuf[pstrlen + 28 : pstrlen + 48])
	return &Handshake{
		Pstr: string(HandshakeBuf[0 : pstrlen]),
		InfoHash: infoHash,
		PeerID: peerID,
	} , nil
}