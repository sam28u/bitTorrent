package torrentfile

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"github.com/jackpal/bencode-go"
)

type bencodeTorrent struct {
	Announce string      `bencode:"announce"`
	Info     bencodeInfo `bencode:"info"`
}

type bencodeInfo struct {
	Pieces      string `bencode:"pieces"`
	PieceLength int    `bencode:"piece length"`
	Length      int    `bencode:"length"`
	Name        string `bencode:"name"`
}

func (i *bencodeInfo) hash() ([20]byte, error) {
	var buf bytes.Buffer
	err := bencode.Marshal(&buf, *i)
	if err != nil {
		return [20]byte{}, err
	}
	h := sha1.Sum(buf.Bytes())
	return h, nil
}

func (i *bencodeInfo) splitPieceHashes() ([][20]byte, error) {
	hashLen := 20 
	buf := []byte(i.Pieces)
	
	if len(buf)%hashLen != 0 {
		return nil, fmt.Errorf("received malformed pieces of length %d", len(buf))
	}
	
	numHashes := len(buf) / hashLen
	hashes := make([][20]byte, numHashes)

	for j := 0; j < numHashes; j++ {
		copy(hashes[j][:], buf[j*hashLen:(j+1)*hashLen])
	}
	return hashes, nil
}