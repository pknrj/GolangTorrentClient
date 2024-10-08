package connection

import (
	"fmt"
	"io"
)


const protocolIdentifier = "BitTorrent protocol"

type BitHandshakeInfo struct {
	Pstr     string 
	InfoHash [20]byte
	PeerID   [20]byte
}

// handshake request info 

//\x13BitTorrent protocol\x00\x00\x00\x00\x00\x00\x00\x00\x86\xd4\xc8\x00\x24\xa4\x69\xbe\x4c\x50\xbc\x5a\x10\x2c\xf7\x17\x80\x31\x00\x74-TR2940-k8hj0wgej6ch

// len(protocol_identifier)BitTorrent protocol\8 reserved bytes \ infohash \ peerid 


func New(peerId [20]byte , infohash [20]byte) *BitHandshakeInfo{
	return &BitHandshakeInfo{
		Pstr: protocolIdentifier,
		InfoHash: infohash,
		PeerID: peerId,
	}
}


func (bh *BitHandshakeInfo) ToBuffer () []byte {
	buf := make([]byte, 68)
	
	buf[0] = byte(len(bh.Pstr))
	
	bytesOcc := 1 
	
	bytesOcc += copy(buf[bytesOcc:], bh.Pstr)
	bytesOcc += copy(buf[bytesOcc:], make([]byte, 8))
	bytesOcc += copy(buf[bytesOcc:], bh.InfoHash[:])
	bytesOcc += copy(buf[bytesOcc:], bh.PeerID[:])
	
	return buf
}

func ParseHandshakeResponse(r io.Reader) (*BitHandshakeInfo , error ) {
	
	lengthBuf := make([]byte, 1)
	_, err := io.ReadFull(r, lengthBuf)
	
	if err != nil {
		return nil, err
	}
	
	pstrlen := int(lengthBuf[0])

	if pstrlen == 0 {
		err := fmt.Errorf("protocolIdentifier cant be nil")
		return nil, err
	}

	handshakeBuf := make([]byte, 48+pstrlen)

	_, err = io.ReadFull(r, handshakeBuf)

	if err != nil {
		return nil, err
	}

	var infoHash, peerID [20]byte

	copy(infoHash[:], handshakeBuf[pstrlen+8:pstrlen+8+20])
	copy(peerID[:], handshakeBuf[pstrlen+8+20:])

	bh := BitHandshakeInfo{
		Pstr:     string(handshakeBuf[0:pstrlen]),
		InfoHash: infoHash,
		PeerID:   peerID,
	}

	return &bh, nil
}