package torrentfile

import (
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/jackpal/bencode-go"
	"github.com/pknrj/GolangTorrentClient/peers"
)



type trackerResponseBencode struct {
	Interval 	int 	`bencode:"interval"`
	Peers 		string 	`bencode:"peers"`
}


func (tr *TorFile) requestListOfPeers(peerID [20]byte , port uint16) ([]peers.Peer, error){
	url, err := tr.getTrackerUrl(peerID, port)
	if err != nil {
		return nil, err
	}

	reqClient := &http.Client{Timeout: 15 * time.Second}
	resp, err := reqClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	trackerResp := trackerResponseBencode{}
	err = bencode.Unmarshal(resp.Body, &trackerResp)
	if err != nil {
		return nil, err
	}
	return peers.FetchPeerList([]byte(trackerResp.Peers))
}


func (tr *TorFile) getTrackerUrl(peerID [20]byte , port uint16)	(string,error){
	baseUrl, err := url.Parse(tr.Announce)

	if err != nil {
		return "",err
	}

	urlParams := url.Values {
		"info_hash":  []string{string(tr.InfoHash[:])},
        "peer_id":    []string{string(peerID[:])},
        "port":       []string{strconv.Itoa(int(port))},
        "uploaded":   []string{"0"},
        "downloaded": []string{"0"},
        "compact":    []string{"1"},
        "left":       []string{strconv.Itoa(tr.Length)},
	}

	baseUrl.RawQuery = urlParams.Encode()
	return baseUrl.String(),nil
}


