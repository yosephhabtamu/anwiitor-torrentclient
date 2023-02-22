package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	// "time"
)

func Write(sendingFile string, recieverFile string) (n int, err error) {
	file, err := os.Create(recieverFile)
	sender, err := os.Open(sendingFile)
	c, err := os.Stat(sendingFile)
	if err != nil {
		log.Fatal("error loading stat")
		return
	}
	if err != nil {
		fmt.Println("error creating file")
		return
	}

	defer file.Close()
	defer sender.Close()

	buf := make([]byte, c.Size())
	n, err = sender.Read(buf)
	if nil != err {
		return
	}
	n, err = file.Write(buf)
	return
}

func Read(filename string) (torrentStruct Torrent, err error) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal("error opening file")
	}
	stat, err := os.Stat(filename)
	if err != nil {
		fmt.Println("not found")
		return
	}
	torrentStruct.Size = stat.Size()
	torrentStruct.Name = stat.Name()
	torrentStruct.BufSize = append(torrentStruct.BufSize, 1000000)
	torrentStruct.BufSize = append(torrentStruct.BufSize, int(torrentStruct.Size)%1000000)

	defer file.Close()
	//test, err := os.Create("anwiishetua.mp4")
	torrentStruct.PieceLength = 0

	for i := -1; i <= int(stat.Size()); i += torrentStruct.BufSize[0] {

		var piece map[string]string
		torrentStruct.PieceLength += 1
		// log.Printf("%v", torrentStruct)
		readData, err := ReadNthPiece(i+1, file, torrentStruct.BufSize[0])
		if err != nil {
			if err == io.EOF {
				//reading the remaining data just like flush
				remainder := int(stat.Size()) % torrentStruct.BufSize[1]
				readData, err := ReadNthPiece((int(stat.Size()) - remainder), file, remainder)
				if err != nil {
					log.Fatal("error reading the final piece ")
				}
				hash := GeneratePieceHash(readData)
				piece, err = StorePiecehash(torrentStruct.PieceLength, hash)
				torrentStruct.Pieces = append(torrentStruct.Pieces, piece)
				if err != nil {
					log.Panic("error storing hash", err)
				}
				//WriteNthPiece(test, readData, (int(stat.Size()) - remainder), bufSize)

				break
			}
			fmt.Println("something went Wrong", err)
			break
		}
		hash := GeneratePieceHash(readData)
		piece, err = StorePiecehash(torrentStruct.PieceLength, hash)
		torrentStruct.Pieces = append(torrentStruct.Pieces, piece)
		if err != nil {
			fmt.Println("error hashing the piece")
			break
		}
		//infohash =

		//WriteNthPiece(test, readData, i+1, bufSize)
	}

	return
}

func ReadToGenerateTorrentFile(fileName string) (torrentFile os.File, err error) {
	name, err := os.Hostname()
	if err != nil {
		log.Print("error fetching hostname", err)
		return
	}

	ip, err := net.LookupIP(name)
	// var workingIps []net.IP
	// // timeout := (20 * time.Second)
	// for _, currIp := range ip {
	// 	conn, err := net.Dial("tcp", net.JoinHostPort(currIp.String(), "80"))
	// 	if err != nil {
	// 		fmt.Printf("%s: offline\n", currIp.String())
	// 		continue
	// 	}
	// 	workingIps = append(workingIps, currIp)
	// 	conn.Close()
	// 	fmt.Printf("%s: online\n", currIp.String())
	// }
	if err != nil {
		log.Panic("error fetching ip")
	}
	torrentStruct, err := Read(fileName)
	if err != nil {
		log.Panic("error while Reading", err)
	}
	torrentStruct.Ip = ip

	_, err = MyMarshall(fileName, torrentStruct)
	if err != nil {
		log.Panic("error marshalling", err)
	}

	return
}
