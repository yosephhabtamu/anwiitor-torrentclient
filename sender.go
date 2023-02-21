// Server

package main

import (
	"encoding/binary"
	"io"
	"log"
	"net"
)



func sendChunk(conn net.Conn, data []byte, start, end int) error {
	_, err := conn.Write(data[start:end])
	if err != nil {
		return err
	}

	// Wait for a confirmation from the client that the chunk was received
	var confirmation int32
	err = binary.Read(conn, binary.LittleEndian, &confirmation)
	if err != nil {
		return err
	}

	if int(confirmation) != end-start {
		return io.ErrShortWrite
	}

	return nil
}

func sendData(conn net.Conn, data []byte) error {
	// Send the length of the data
	err := binary.Write(conn, binary.LittleEndian, int64(len(data)))
	if err != nil {
		return err
	}

	for i := 0; i < len(data); i += ChunkSize {
		j := i + ChunkSize
		if j > len(data) {
			j = len(data)
		}

		// Keep track of the number of retries
		retry := 0

		for {
			err := sendChunk(conn, data, i, j)
			if err == nil {
				break
			}

			// If there was an error, check if we've reached the retry limit
			retry++
			if retry == RetryLimit {
				log.Printf("Error sending chunk, retry limit reached: %v", err)
				return err
			}
		}
	}

	return nil
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	data := []byte("124")

	err := sendData(conn, data)
	if err != nil {
		log.Printf("Error sending data: %v", err)
	}
}

func startListen() {
	ln, err := net.Listen("tcp", ":6882")
	if err != nil {
		log.Fatalf("Error listening: %v", err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}

		go handleConnection(conn)
	}
}
