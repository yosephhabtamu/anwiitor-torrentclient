// Client

package main

import (
	"encoding/binary"
	"io"
	"log"
	"net"
	"sync"
)

const ChunkSize = 1024
const RetryLimit = 3

func receiveChunk(conn net.Conn, data []byte, start, end int) error {
	_, err := io.ReadFull(conn, data[start:end])
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

func receiveData(conn net.Conn) ([]byte, error) {
	// Receive the length of the data
	var dataLen int64
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
				err := receiveChunk(conn, data, i, j)
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

func startSending() {
	conn, err := net.Dial("tcp", "192.168.1.4:6882")
	if err != nil {
		log.Fatalf("Error dialing: %v", err)
	}
	defer conn.Close()

	data, err := receiveData(conn)
	if err != nil {
		log.Fatalf("Error receiving data: %v", err)
	}

	log.Printf("Received data: %s", data)
}
