package main

import (
	"fmt"
	"log"
	"os"
	"github.com/pknrj/GolangTorrentClient/torrentFile"
)


func main(){
	fmt.Println("Bittorrent Client in Go !!!!")

	inPath := os.Args[1]
	outPath := os.Args[2]

	tf, err := torrentfile.OpenTorFile(inPath)
	if err != nil {
		log.Fatal(err)
	}

	err = tf.DownloadTorFile(outPath)
	if err != nil {
		log.Fatal(err)
	}
}



// parsing .torrent file (trackers information) --> done
// retrieving peers from tracker 

// TCP connection 
// Bit-torrent two way handshake 
// Exchange messages to download pieces 

