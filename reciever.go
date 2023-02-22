// Client

package main

import (
	"encoding/binary"
	"io"
	"log"
	"net"
	"sync"
)

const ChunkSize = 1000000
const RetryLimit = 10

func ReceiveChunk(conn net.Conn, data []byte, start, end int) error {
	log.Printf("Recieve chunk")
	_, err := io.ReadFull(conn, data[start:end])
	log.Printf("%v", data)
	if err != nil {
		log.Fatal(err)
	}

	// Send a confirmation to the server that the chunk was received
	confirmation := int32(end - start)
	err = binary.Write(conn, binary.LittleEndian, confirmation)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func ReceiveData(conn net.Conn, torrentStruct Torrent) ([]byte, error) {
	// Receive the length of the data
	dataLen := torrentStruct.Size
	err := binary.Read(conn, binary.LittleEndian, &dataLen)
	if err != nil {
		return nil, err
	}

	data := make([]byte, dataLen)
	var wg sync.WaitGroup
	var mu sync.Mutex

	for i := 0; i < int(dataLen); i += ChunkSize {
		j := i + ChunkSize
		if j > int(dataLen) {
			j = int(dataLen)
		}

		wg.Add(1)
		go func(i, j int) {
			defer wg.Done()

			// Keep track of the number of retries
			retry := 0

			for {
				err := ReceiveChunk(conn, data, i, j)
				if err == nil {
					break
				}

				// If there was an error, check if we've reached the retry limit
				retry++
				if retry == RetryLimit {
					mu.Lock()
					defer mu.Unlock()
					log.Printf("Error receiving chunk, retry limit reached: %v", err)
					return
				}
			}
		}(i, j)
	}

	wg.Wait()

	return data, nil
}

func StartSending(torrentStruct Torrent) (data []byte, err error) {
	conn, err := net.Dial("tcp", "192.168.1.4:6882")
	if err != nil {
		log.Fatalf("Error dialing: %v", err)
	}
	defer conn.Close()

	data, err = ReceiveData(conn, torrentStruct)
	log.Print(data)
	if err != nil {
		log.Fatalf("Error receiving data: %v", err)
	}

	log.Printf("Received data: %s", data)
	return
}
