package main

import (
	"fmt"
	"log"

	"bittorrent/torrentfile"
)

func main() {
	tf, err := torrentfile.Open("debian-11.6.0-amd64-netinst.iso.torrent")
	if err != nil {
		log.Fatalf("Failed to open torrent file: %v", err)
	}

	fmt.Println("--- Torrent Metafile Parsed Successfully ---")
	fmt.Printf("Name:        %s\n", tf.Name)
	fmt.Printf("Announce:    %s\n", tf.Announce)
	fmt.Printf("File Length: %d bytes\n", tf.Length)
	fmt.Printf("Piece Len:   %d bytes\n", tf.PieceLength)
	fmt.Printf("InfoHash:    %x\n", tf.InfoHash)
	fmt.Printf("Piece Count: %d\n", len(tf.PieceHashes))
}