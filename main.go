package main

import (
	"fmt"
	"log"
	"bittorrent/torrentfile"
	"bittorrent/tracker"
)

func main() {
	tf, err := torrentfile.Open("debian-13.6.0-amd64-netinst.iso.torrent")
	if err != nil {
		log.Fatalf("Failed to open torrent file: %v", err)
	}
	fmt.Println("--- Torrent Metafile Parsed Successfully ---")
	fmt.Printf("Announce URL: %s\n", tf.Announce)
	peerID, err := tracker.GeneratePeerID()
	if err != nil {
		log.Fatalf("Failed to generate Peer ID: %v", err)
	}
	fmt.Println("\n--- Contacting Tracker ---")
	const port = 6881 
	peers, err := tracker.RequestPeers(tf, peerID, port)
	if err != nil {
		log.Fatalf("Tracker request failed: %v", err)
	}
	fmt.Printf("Success! Found %d peers.\n", len(peers))
	for i, peer := range peers {
		fmt.Printf("Peer %d: %s:%d\n", i+1, peer.IP.String(), peer.Port)
		if i >= 4 {
			fmt.Printf("... and %d more\n", len(peers)-5)
			break
		}
	}
}