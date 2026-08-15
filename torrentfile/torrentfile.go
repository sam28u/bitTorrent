package torrentfile

import (
	"github.com/jackpal/bencode-go"
	"os"
)

type TorrentFile struct {
	Announce    string
	InfoHash    [20]byte
	PieceHashes [][20]byte
	PieceLength int
	Length      int
	Name        string
}

func Open(path string) (*TorrentFile, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bto := bencodeTorrent{}
	err = bencode.Unmarshal(file, &bto)
	if err != nil {
		return nil, err
	}
	return bto.toTorrentFile()
}

func (bto *bencodeTorrent) toTorrentFile() (*TorrentFile, error) {
	infoHash, err := bto.Info.hash()
	if err != nil {
		return nil, err
	}

	pieceHashes, err := bto.Info.splitPieceHashes()
	if err != nil {
		return nil, err
	}

	t := &TorrentFile{
		Announce:    bto.Announce,
		InfoHash:    infoHash,
		PieceHashes: pieceHashes,
		PieceLength: bto.Info.PieceLength,
		Length:      bto.Info.Length,
		Name:        bto.Info.Name,
	}
	return t, nil
}
