package connection

import (
	"bytes"
	"fmt"
	"net"
	"time"

	"github.com/pknrj/GolangTorrentClient/bitfields"
	"github.com/pknrj/GolangTorrentClient/peers"
)

// TCP connection with peers request

type Client struct {
	Conn     net.Conn
	Choked   bool
	Bitfield bitfields.Bitfield
	peer     peers.Peer
	infoHash [20]byte
	peerID   [20]byte
}


func NewClient(peer peers.Peer, peerID, infoHash [20]byte) (*Client, error) {
	
	conn, err := net.DialTimeout("tcp", peer.String(), 3*time.Second)
	if err != nil {
		return nil, err
	}

	_, err = completeHandshake(conn, infoHash, peerID)
	if err != nil {
		conn.Close()
		return nil, err
	}

	bf, err := receiveBitfield(conn)
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &Client{
		Conn:     conn,
		Choked:   true,
		Bitfield: bf,
		peer:     peer,
		infoHash: infoHash,
		peerID:   peerID,
	}, nil
}


func completeHandshake(conn net.Conn, infohash, peerID [20]byte) (*BitHandshakeInfo, error) {
	
	conn.SetDeadline(time.Now().Add(3 * time.Second))
	defer conn.SetDeadline(time.Time{})

	req := New(infohash, peerID)
	_, err := conn.Write(req.ToBuffer())
	if err != nil {
		return nil, err
	}

	// res, err := connection.ParseHandshakeResponse(conn)
	res, err := ParseHandshakeResponse(conn)
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(res.InfoHash[:], infohash[:]) {
		return nil, fmt.Errorf("Expected infohash %x but got %x", res.InfoHash, infohash)
	}
	return res, nil
}


func receiveBitfield(conn net.Conn) (bitfields.Bitfield, error) {
	conn.SetDeadline(time.Now().Add(5 * time.Second))
	defer conn.SetDeadline(time.Time{}) // Disable the deadline

	msg, err := ParseMessage(conn)
	if err != nil {
		return nil, err
	}
	if msg == nil {
		err := fmt.Errorf("Expected bitfield but got %s", msg)
		return nil, err
	}
	if msg.ID != Bitfield {
		err := fmt.Errorf("Expected bitfield but got ID %d", msg.ID)
		return nil, err
	}

	return msg.Payload, nil
}


// Read reads and consumes a message from the connection

func (c *Client) Read() (*Message, error) {
	msg, err := ParseMessage(c.Conn)
	return msg, err
}

// SendRequest sends a Request message to the peer

func (c *Client) SendRequest(index, begin, length int) error {
	req := FormatRequest(index, begin, length)
	_, err := c.Conn.Write(req.Serialize())
	return err
}

// SendInterested sends an Interested message to the peer

func (c *Client) SendInterested() error {
	msg := Message{ID: Interested}
	_, err := c.Conn.Write(msg.Serialize())
	return err
}

// SendNotInterested sends a NotInterested message to the peer

func (c *Client) SendNotInterested() error {
	msg := Message{ID: NotInterested}
	_, err := c.Conn.Write(msg.Serialize())
	return err
}

// SendUnchoke sends an Unchoke message to the peer

func (c *Client) SendUnchoke() error {
	msg := Message{ID: Unchoke}
	_, err := c.Conn.Write(msg.Serialize())
	return err
}

// SendHave sends a Have message to the peer

func (c *Client) SendHave(index int) error {
	msg := FormatHave(index)
	_, err := c.Conn.Write(msg.Serialize())
	return err
}