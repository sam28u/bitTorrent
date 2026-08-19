package main

import (
	"bittorrent/p2p"
	"bittorrent/torrentfile"
	"bittorrent/tracker"
	"fmt"
	"log"
)

func main() {
	// 1. Parse the torrent file (Phase 1)
	tf, err := torrentfile.Open("debian-13.6.0-amd64-netinst.iso.torrent") // Make sure this matches your file name!
	if err != nil {
		log.Fatalf("Failed to open torrent file: %v", err)
	}
	fmt.Println("--- Torrent Metafile Parsed ---")

	// 2. Contact Tracker (Phase 2)
	peerID, err := tracker.GeneratePeerID()
	if err != nil {
		log.Fatalf("Failed to generate Peer ID: %v", err)
	}

	fmt.Println("\n--- Contacting Tracker ---")
	peers, err := tracker.RequestPeers(tf, peerID, 6881)
	if err != nil {
		log.Fatalf("Tracker request failed: %v", err)
	}
	fmt.Printf("Success! Found %d peers.\n", len(peers))

	// 3. Connect to Peers (Phase 3)
	fmt.Println("\n--- Initiating TCP Handshakes ---")

	// Loop through the peers we got from the tracker
	for _, peer := range peers {
		fmt.Printf("Dialing %s:%d...\n", peer.IP.String(), peer.Port)

		// Try to connect and handshake
		client, err := p2p.Connect(peer, tf.InfoHash, peerID)
		if err != nil {
			fmt.Printf("  [-] Failed: %v\n", err)
			continue // Skip to the next peer in the list
		}

		// If we get here, the 68-byte handshake was a complete success!
		fmt.Printf("  [+] SUCCESS! Handshake completed with %s:%d\n", peer.IP.String(), peer.Port)

		// For now, let's just close the connection and stop after our first success.
		// (In the final phase, we will keep it open to start downloading pieces!)
		client.Conn.Close()
		break
	}
}
