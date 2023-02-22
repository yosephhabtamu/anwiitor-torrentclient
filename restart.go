package main

import (
	"io"
	"log"
	"os"
	"strconv"
)

//func checkIntegrity

//func request for the missing pieces

func CheckPieceIntegrity(piece []byte, hash string) (valid bool) {
	valid = GeneratePieceHash(piece) == hash
	return
}

func FindMissingPieces(torrentStruct Torrent, file *os.File) (missingPieces []int, err error) {

	var byteArray = make([]byte, torrentStruct.Size)
	pieceCounter := 0
	for {
		pieceCounter += 1
		bufSize := torrentStruct.BufSize[0]
		// if torrentStruct.PieceLength == pieceCounter {
		// 	bufSize = torrentStruct.BufSize[1]
		// }
		byteArray, err = ReadNthPiece(pieceCounter, file, bufSize)
		if err != nil {
			if err == io.EOF {
				byteArray, err = ReadNthPiece(pieceCounter, file, (int(torrentStruct.Size) - (torrentStruct.BufSize[0] * pieceCounter)))
				pieceState := CheckPieceIntegrity(byteArray, torrentStruct.Pieces[pieceCounter][strconv.Itoa(pieceCounter)])
				if !pieceState {
					missingPieces = append(missingPieces, pieceCounter)
				}
				return
			}
			log.Fatalf("can not read the pieces")
		}
		pieceState := CheckPieceIntegrity(byteArray, torrentStruct.Pieces[pieceCounter][strconv.Itoa(pieceCounter)])
		if !pieceState {
			missingPieces = append(missingPieces, pieceCounter)
		}
		// if pieceCounter == torrentStruct.PieceLength {
		// 	byteArray, err = ReadNthPiece(pieceCounter, file, (int(torrentStruct.Size) - (torrentStruct.BufSize[1] * pieceCounter)))
		// 	pieceState := CheckPieceIntegrity(byteArray, torrentStruct.Pieces[pieceCounter][strconv.Itoa(pieceCounter)])
		// 	if !pieceState {
		// 		missingPieces = append(missingPieces, pieceCounter)
		// 	}

		// 	return
		// }
	}

}

func OpenExistingFile(fileName string) (file *os.File, err error) {
	entries, err := os.ReadDir("./")
	if err != nil {
		log.Fatalf("error reading current directory check the previlage")
	}
	for _, e := range entries {
		if e.Name() == fileName {
			file, err = os.Open(e.Name())
			return
		}
	}
	file, err = os.Create(fileName)
	return
}
