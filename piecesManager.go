package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

func ReadNthPiece(n int, file *os.File, bufSize int) (buf []byte, err error) {
	buf = make([]byte, bufSize)
	_, err = file.ReadAt(buf, int64((n)))
	if err != nil {
		if err == io.EOF {
			return
		}
		fmt.Println("error reading the piece specified: ", n, err)
		return
	}
	return
}

func WriteNthPiece(file *os.File, bufData []byte, n int, bufSize int) (buf []byte) {
	_, err := file.WriteAt(bufData, int64((n)))
	if err != nil {
		log.Fatal("error writing to the piece specified: ", n, err)
		return
	}
	return
}

func StorePiecehash(loc int, hash string) (result map[string]string, err error) {
	result = make(map[string]string)
	result[strconv.Itoa(loc)] = hash
	return

}
