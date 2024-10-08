package peers

import (
	"encoding/binary"
	"fmt"
	"net"
	"strconv"
)

const peerSize = 6 // 4 for ip and 2 for port 

type Peer struct {
	IP 		net.IP
	Port 	uint16
}

func FetchPeerList(peersBuffer []byte) ([]Peer , error) {
	
	if len(peersBuffer) % peerSize != 0 {
		err := fmt.Errorf("Received malformed peers")
		return nil, err
	}	

	numOfPeers := len(peersBuffer) / peerSize
	peers := make([]Peer , numOfPeers)

	for i := 0; i < numOfPeers; i++ {
		offset := i * peerSize
		peers[i].IP = net.IP(peersBuffer[offset : offset+4])
		peers[i].Port = binary.BigEndian.Uint16([]byte(peersBuffer[offset+4 : offset+6]))
	}
	return peers , nil
}


func (p Peer) String() string {
	return net.JoinHostPort(p.IP.String(), strconv.Itoa(int(p.Port)))
}