// Server

package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"strconv"
)

func sendChunk(conn net.Conn, data []byte, start, end int, ChunkSize int) error {
	_, err := conn.Write(data)
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

func sendData(conn net.Conn, data []byte, ChunkSize int) error {
	// Send the length of the data
	err := binary.Write(conn, binary.LittleEndian, int64(len(data)))
	if err != nil {
		return err
	}

	for i := 0; i < 1; i += ChunkSize {
		j := i + ChunkSize
		if j > len(data) {
			j = len(data)
		}

		// Keep track of the number of retries
		retry := 0

		for {
			err := sendChunk(conn, data, i, j, ChunkSize)
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

func handleSignal(conn net.Conn) (err error) {
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		log.Printf("Error reading request")
		return
	}

	signal, err := strconv.Atoi(string(buf[:n]))
	log.Printf("%d", signal)
	if err != nil {
		log.Printf("Received request was not a signal")
	} else {
		if signal == 1 {
			response := []byte(strconv.Itoa(2))
			if _, err = conn.Write(response); err != nil {
				log.Printf("Error sending signal")
				return
			}
		}
	}
	return
}

func handleConnection(conn net.Conn, torrentStruct Torrent, filename string) {
	defer conn.Close()

	// if err := handleSignal(conn); err != nil {
	// 	data := []byte("124")
	// 	err := sendData(conn, data)
	// 	if err != nil {
	// 		log.Printf("Error sending data: %v", err)
	// 	}
	// }

	data, _ := fileToByteArray(filename)
	err := sendData(conn, data, binary.Size(data))
	if err != nil {
		log.Printf("Error sending data: %v", err)
	}
}

func StartListen(torrentStruct Torrent, filename string) {
	ln, err := net.Listen("tcp", ":6882")
	if err != nil {
		log.Fatalf("Error listening: %v", err)
	}
	defer ln.Close()
	fmt.Print("Listening for leechers on port 6882")

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}

		go handleConnection(conn, torrentStruct, filename)
	}
}

func fileToByteArray(filename string) ([]byte, error) {
    // Read the entire file into memory
    fileContents, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, err
    }

    // Convert the file contents to a byte array
    byteArray := []byte(fileContents)

    return byteArray, nil
}