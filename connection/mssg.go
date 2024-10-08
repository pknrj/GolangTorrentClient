package connection

import (
	"encoding/binary"
	"fmt"
	"io"
)

// Choke chokes the receiver
// Unchoke unchokes the receiver
// Interested expresses interest in receiving data
// NotInterested expresses disinterest in receiving data
// Have alerts the receiver that the sender has downloaded a piece
// Bitfield encodes which pieces that the sender has downloaded
// Request requests a block of data from the receiver
// Piece delivers a block of data to fulfill a request
// Cancel cancels a request



const (
	Choke uint8 = 0
	Unchoke uint8 = 1
	Interested uint8 = 2
	NotInterested uint8 = 3
	Have uint8 = 4
	Bitfield uint8 = 5
	Request uint8 = 6
	Piece uint8 = 7
	Cancel uint8 = 8
)



type Message struct {
	ID 			uint8
	Payload		[]byte 
}


func (m *Message) Serialize() []byte {
	if m == nil {				// keep alive message 
		return make([]byte , 4)
	}

	length := uint32(len(m.Payload) + 1) // payload + id 
	buf := make([]byte, 4+length) // length of payload + length 

	binary.BigEndian.PutUint32(buf[0:4], length)
	
	buf[4] = byte(m.ID)
	copy(buf[5:], m.Payload)
	
	return buf
}

func ParseMessage(r io.Reader) (*Message, error) {
	
	lengthBuf := make([]byte, 4)
	_, err := io.ReadFull(r, lengthBuf)
	if err != nil {
		return nil, err
	}
	length := binary.BigEndian.Uint32(lengthBuf)

	// keep-alive message
	if length == 0 {
		return nil, nil
	}

	messageBuf := make([]byte, length)
	_, err = io.ReadFull(r, messageBuf)
	if err != nil {
		return nil, err
	}

	m := Message{
		ID:      uint8(messageBuf[0]),
		Payload: messageBuf[1:],
	}

	return &m, nil
}


func FormatRequest(index, begin, length int) *Message {
	payload := make([]byte, 12)
	binary.BigEndian.PutUint32(payload[0:4], uint32(index))
	binary.BigEndian.PutUint32(payload[4:8], uint32(begin))
	binary.BigEndian.PutUint32(payload[8:12], uint32(length))
	return &Message{ID: Request, Payload: payload}
}

func FormatHave(index int) *Message {
	payload := make([]byte, 4)
	binary.BigEndian.PutUint32(payload, uint32(index))
	return &Message{ID: Have, Payload: payload}
}

func ParsePiece(index int, buf []byte, msg *Message) (int, error) {
	if msg.ID != Piece {
		return 0, fmt.Errorf("Expected PIECE (ID %d), got ID %d", Piece, msg.ID)
	}
	if len(msg.Payload) < 8 {
		return 0, fmt.Errorf("Payload too short. %d < 8", len(msg.Payload))
	}
	parsedIndex := int(binary.BigEndian.Uint32(msg.Payload[0:4]))
	if parsedIndex != index {
		return 0, fmt.Errorf("Expected index %d, got %d", index, parsedIndex)
	}
	begin := int(binary.BigEndian.Uint32(msg.Payload[4:8]))
	if begin >= len(buf) {
		return 0, fmt.Errorf("Begin offset too high. %d >= %d", begin, len(buf))
	}
	data := msg.Payload[8:]
	if begin+len(data) > len(buf) {
		return 0, fmt.Errorf("Data too long [%d] for offset %d with length %d", len(data), begin, len(buf))
	}
	copy(buf[begin:], data)
	return len(data), nil
}

func ParseHave(msg *Message) (int, error) {
	if msg.ID != Have {
		return 0, fmt.Errorf("Expected HAVE (ID %d), got ID %d", Have, msg.ID)
	}
	if len(msg.Payload) != 4 {
		return 0, fmt.Errorf("Expected payload length 4, got length %d", len(msg.Payload))
	}
	index := int(binary.BigEndian.Uint32(msg.Payload))
	return index, nil
}

func (m *Message) name() string {
	if m == nil {
		return "KeepAlive"
	}
	switch m.ID {
	case Choke:
		return "Choke"
	case Unchoke:
		return "Unchoke"
	case Interested:
		return "Interested"
	case NotInterested:
		return "NotInterested"
	case Have:
		return "Have"
	case Bitfield:
		return "Bitfield"
	case Request:
		return "Request"
	case Piece:
		return "Piece"
	case Cancel:
		return "Cancel"
	default:
		return fmt.Sprintf("Unknown#%d", m.ID)
	}
}

func (m *Message) String() string {
	if m == nil {
		return m.name()
	}
	return fmt.Sprintf("%s [%d]", m.name(), len(m.Payload))
}