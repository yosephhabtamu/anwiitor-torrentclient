package main

import (
	// "fmt"
	"encoding/binary"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"
)

func GenerateSeeder(fileName string) (filePath string, err error) {
	_, err = ReadToGenerateTorrentFile(fileName)
	if err != nil {
		log.Panic("error generating the torrent file", err)
	}
	filePath, err = filepath.Abs("./" + fileName + ".torrent")
	if err != nil {
		log.Panic("error generating file in this directory", err)
	}
	return

}

func CheckAvailability(conn net.Conn) (signal int, err error) {
	defer conn.Close()

	leecherSignal := 1 // 1 means request for piece availability
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(leecherSignal))
	conn.Write(buf)

	respBuf := make([]byte, 4)
	_, err = conn.Read(respBuf)
	if err != nil {
		log.Fatal("Error:", err)
	}

	// Convert the response to an integer
	resp := int(binary.BigEndian.Uint32(respBuf))
	signal = resp
	return
}

func awaitChoke() {
	time.Sleep(5 * time.Minute)
}

func ManageLeech(torrentFile string) (torrentStruct Torrent, err error) {
	var file *os.File
	torrentStruct, err = MyUnmarshall(torrentFile)
	if err != nil {
		log.Panic("error unmarshalling torrent file", err)
	}
	file, err = OpenExistingFile(torrentStruct.Name)
	if err != nil {
		log.Fatal("error handling file")
	}

	defer file.Close()

	peerIp := torrentStruct.Ip
	if len(peerIp) == 0 {
		log.Fatalf("No peers avaialable")
	}

	currPeer := peerIp[1]
	log.Print(currPeer.String())
	curr_conn, err := net.Dial("tcp", ":6882")
	if err != nil {
		log.Fatalf("Failed Connecting with peer")
	}

	var missingPieces []int
	missingPieces, err = FindMissingPieces(torrentStruct, file)
	log.Printf("%v", missingPieces)
	if err != nil {
		if err == io.EOF {
			log.Print("successful missing pieces")
			return
		}
		log.Panic("cant resume: ", err)
		return
	}
	// for i :=1; i<=  torrentStruct.PieceLength;i+=1 {
	// 	// seederSignal, err := CheckAvailability(curr_conn)
	// 	// if err != nil {
	// 	// 	log.Fatalf("Error recieving signals")
	// 	// }
	// 	// for seederSignal != 2 { // 2 means it's available
	// 	// 	log.Print("Getting choked by the seeder or sedder not available")
	// 	// 	awaitChoke()
	// 	// }

	// 	//send torrent request
	var data []byte
		data, err = ReceiveData(curr_conn)
		if err != nil {

			log.Fatal("error loading")
		}
	// 	log.Printf("%v", data)
	// 	bufSize := torrentStruct.BufSize[0]
	// 	if i == torrentStruct.PieceLength {
	// 		bufSize = torrentStruct.BufSize[1]
	// 	}
		WriteNthPiece(file, data, 0, binary.Size(data))

	// }


	return
}
