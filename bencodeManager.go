package main

import (
	"log"
	"net"
	"os"

	"github.com/jackpal/bencode-go"
)

type Torrent struct {
	Name      string
	Ip        []net.IP
	InfoHash  string
	BufSize   []int
	PieceLength int
	Size      int64
	Pieces    []map[string]string
}

// writing to the torrent file
func MyMarshall(fileName string, torrentStruct Torrent) (torrentFile os.File, err error) {

	file, err := os.Create(fileName + ".torrent")
	if err != nil {
		log.Panic("error creating file")
	}
	defer torrentFile.Close()
	bencode.Marshal(file, torrentStruct)
	return

}

func MyUnmarshall(fileName string) (torrentStruct Torrent, err error) {
	file, err := os.Open(fileName)

	if err != nil {
		log.Panic("error opening file for unmarshall", err)
	}

	err = bencode.Unmarshal(file, &torrentStruct)
	return
}
