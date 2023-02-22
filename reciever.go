// Client

package main

import (
	"encoding/binary"
	"io"
	"log"
	"net"
)

const RetryLimit = 10

func ReceiveChunk(conn net.Conn, data []byte, start, end int) error {
	log.Printf("Recieve chunk")
	_, err := io.ReadFull(conn, data)
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

func ReceiveData(conn net.Conn, ChunkSize int) ([]byte, error) {
	// Receive the length of the data
	var dataLen int64
	err := binary.Read(conn, binary.LittleEndian, &dataLen)
	if err != nil {
		return nil, err
	}

	data := make([]byte, ChunkSize)
	ReceiveChunk(conn, data, 0, ChunkSize)
	// var wg sync.WaitGroup
	// var mu sync.Mutex
	// log.Printf("%d", ChunkSize)
	// log.Printf("%d", int(dataLen))
	// for i := 0; i < ChunkSize; i += ChunkSize {
	// 	j := (i + ChunkSize) - 1

	// 	wg.Add(1)
	// 	go func(i, j int) {
	// 		defer wg.Done()

	// 		// Keep track of the number of retries
	// 		retry := 0

	// 		for {
	// 			err := ReceiveChunk(conn, data, i, j)
	// 			if err == nil {
	// 				break
	// 			}

	// 			// If there was an error, check if we've reached the retry limit
	// 			retry++
	// 			if retry == RetryLimit {
	// 				mu.Lock()
	// 				defer mu.Unlock()
	// 				log.Printf("Error receiving chunk, retry limit reached: %v", err)
	// 				return
	// 			}
	// 		}
	// 	}(i, j)
	// }

	// wg.Wait()

	return data, nil
}

func StartSending(ChunkSize int) ([]byte, error) {
	conn, err := net.Dial("tcp", ":6882")
	if err != nil {
		log.Fatalf("Error dialing: %v", err)
	}
	defer conn.Close()

	data, err := ReceiveData(conn, ChunkSize)
	
	if err != nil {
		log.Fatalf("Error receiving data: %v", err)
	}
	return data, nil
}
