These handwritten notes form a complete, rock-solid architectural blueprint for a strict BEP 3 client. The logic mapped out here handles exactly what is needed to get a robust Go implementation off the ground without overcomplicating the spec.

Here is the structured alignment of your notes, organized into the exact development phases you will need to code.
Phase 1: Metafile Parsing & Bencoding

Before hitting the network, the client must decode the .torrent file.

    Bencode Data Types: You will need to parse strings (e.g., 4:spam), integers (e.g., i3e, i-3e), lists (e.g., l4:spam4:eggse), and dictionaries (e.g., d3:cow3:moo4:spam4:eggse).

    File Structure: Metafiles are UTF-8 encoded and contain two main keys: the announce URL and the info dictionary.

    The Info Dictionary: This dictionary requires extracting the name, the piece length (which is mostly always a power of 2), and the pieces string, which is subdivided into lengths of 20 bytes representing the SHA-1 hash at the corresponding index.

    File Modes: The dictionary will contain a length key for a single file, or a files dictionary for a set of files.

Phase 2: Tracker Communication

The client must announce itself to the tracker to get a list of peers.

    Request Parameters: The HTTP request to the tracker requires the info_hash, a 20-length peer_id, the ip, the port, the total amount uploaded, the total amount downloaded, the number of bytes left to download, and an optional event (Started, completed, or stopped).

    Tracker Response: Trackers return a compact representation of the peer list.

    Byte Translation: When parsing the compact peer list, encoding/binary must be used to translate the 2-byte chunks into port numbers.

Phase 3: TCP Connection & The Handshake

This phase establishes the foundational peer-to-peer wire stream.

    Dialing Peers: Use net.DialTimeout to dial out to peer IPs to ensure you don't hang on dead nodes.

    Handshake Structure: The handshake is a 68-byte sequence: a 1 byte prefix, 19 bytes containing the "BitTorrent protocol" name, 8 empty reserved slots for extra features, 20 bytes for the info-hash to identify the torrent, and 20 bytes for the peer-id.

    Main Go Flow: Read the 68-byte handshake, validate the protocol, info-hash, and peer-id, send your own handshake back, and then enter the main message loop.

    Checksums: The crypto/sha1 package is used to calculate the info-hash and feed byte arrays into SHA-1 to generate and compare checksums.

Phase 4: The Message Loop & Pipelining

Once connected, the client reads a never-ending stream of length-prefixed messages over the net.Conn stream.

    Deterministic Reading: Because network packets fragment, use io.ReadFull rather than a standard read to ensure the full message is pulled from the TCP stream.

    Message Lengths: First, read the 4-byte length prefix using encoding/binary. If the length is 0, it is a keep-alive message.

    Message Parsing: If the length is greater than 0, read the message and parse the message ID and payload.

    Message IDs:

        0 (choke), 1 (unchoke), 2 (interested), and 3 (not interested) have no payload.

        4 is have, 5 is bitfield (which is only sent as the first message), 6 is req, 7 is piece, 8 is cancel, and 9 is port.

Phase 5: State Management & Buffer Logic

Efficient downloading requires managing the connection state properly.

    Pipelining Requests: TCP takes time to send and acknowledge data. Instead of sending a request and waiting for a response, keep several requests in flight continuously to keep the TCP connection busy.

    Handling Chokes: A choke message indicates the peer is "not gonna send data anymore".

    Request Cancellation: Pipelining makes it easy to cancel pending requests; if a sender sends a choke message, you can discard all the requests at once instead of waiting for their status.

    Buffer Logic: If the buffer space is not full, store the request or data; if it is full, queue it in the memory associated with that specific connection.