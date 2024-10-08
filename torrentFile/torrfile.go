package torrentfile

import (
	"bytes"
	"crypto/rand"
	"crypto/sha1"
	"fmt"
	"os"

	"github.com/jackpal/bencode-go"
	"github.com/pknrj/GolangTorrentClient/downloader"
)


const Port uint16 = 6881 

type bencodeInfo struct {
	Pieces      string `bencode:"pieces"`
	PieceLength int    `bencode:"piece length"`
	Length      int    `bencode:"length"`
	Name        string `bencode:"name"`
}

type bencodeTor struct {
	Announce string      `bencode:"announce"`
	Info     bencodeInfo `bencode:"info"`
}

type TorFile struct {
	Announce	string
	InfoHash	[20]byte
	PieceHashes	[][20]byte
	PieceLength int
	Length 		int
	Name 		string
}


func OpenTorFile(path string) (TorFile , error) {
	file , err := os.Open(path)
	if err != nil {
		return TorFile{} , err
	}	
	defer file.Close()

	bt := bencodeTor{}
	err = bencode.Unmarshal(file , &bt)

	if err != nil {
		return TorFile{} , err
	}

	return bt.bencodeToTor()

}

func (t *TorFile) DownloadTorFile(path string) error {
	var peerID [20]byte
	_, err := rand.Read(peerID[:])
	if err != nil {
		return err
	}

	peers, err := t.requestListOfPeers(peerID, Port)
	if err != nil {
		return err
	}

	torrent := downloader.Torrent{
		Peers:       peers,
		PeerID:      peerID,
		InfoHash:    t.InfoHash,
		PieceHashes: t.PieceHashes,
		PieceLength: t.PieceLength,
		Length:      t.Length,
		Name:        t.Name,
	}
	buf, err := torrent.Download()
	if err != nil {
		return err
	}

	outFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer outFile.Close()
	_, err = outFile.Write(buf)
	if err != nil {
		return err
	}
	return nil
}

func (bt *bencodeTor) bencodeToTor() (TorFile , error) {
	infohash , err := bt.Info.hash() 
	if err != nil {
		return TorFile{} , err 
	}

	pieceHashes , err := bt.Info.dividePieceStringHashes()	
	if err != nil {
		return TorFile{} , err
	}

	tf := TorFile {
		Announce: bt.Announce,
		InfoHash: infohash,
		PieceHashes : pieceHashes,
		PieceLength : bt.Info.PieceLength,
		Length: bt.Info.Length,
		Name : bt.Info.Name,
	}

	return tf , nil 
} 

func (bi *bencodeInfo) hash() ([20]byte, error) {
	var buf bytes.Buffer
	err := bencode.Marshal(&buf, *bi)
	if err != nil {
		return [20]byte{}, err
	}
	h := sha1.Sum(buf.Bytes())
	return h, nil
}


func (bi *bencodeInfo) dividePieceStringHashes() ([][20]byte, error) {
	// Length of SHA-1 hash
	hashLen := 20 
	buf := []byte(bi.Pieces)
	if len(buf)%hashLen != 0 {
		err := fmt.Errorf("Received malformed pieces of length %d", len(buf))
		return nil, err
	}
	numOfHashes := len(buf) / hashLen
	hashes := make([][20]byte, numOfHashes)

	for i := 0; i < numOfHashes; i++ {
		copy(hashes[i][:], buf[i*hashLen:(i+1)*hashLen])
	}
	return hashes, nil
}

