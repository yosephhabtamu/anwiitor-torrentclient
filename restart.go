package main

import (
	"log"
	"os"
	"strconv"
)

//func checkIntegrity

//func request for the missing pieces

func checkPieceIntegrity(piece []byte, hash string ) (valid bool){
	valid = generatePieceHash(piece) == hash
	return 
}


func findMissingPieces(torrentStruct Torrent, file *os.File)(missingPieces[]int, err error){

	var byteArray = make([]byte,torrentStruct.Size )
	pieceCounter := 0
	for {
		pieceCounter+=1
		byteArray, err= ReadNthPiece(pieceCounter, file,torrentStruct.BufSize,)
		if err!=nil{
			log.Fatalf("can not read the pieces")
		}
		pieceState := checkPieceIntegrity(byteArray, torrentStruct.Pieces[pieceCounter][strconv.Itoa(pieceCounter)])
		if !pieceState{
			missingPieces = append(missingPieces, pieceCounter)
		}
		if pieceCounter == torrentStruct.PieceLength{
			byteArray, err= ReadNthPiece(pieceCounter, file,(int(torrentStruct.Size)-(torrentStruct.BufSize* pieceCounter)))
			pieceState := checkPieceIntegrity(byteArray, torrentStruct.Pieces[pieceCounter][strconv.Itoa(pieceCounter)])
		if !pieceState{
			missingPieces = append(missingPieces, pieceCounter)
		}
			
			return
		}
	}
	

}

func OpenExistingFile(fileName string)(file *os.File, err error){
	entries , err := os.ReadDir("./")
	if err!=nil{
		log.Fatalf("error reading current directory check the previlage")
	} 
	for _,e := range entries{
		if e.Name() == fileName{
			file, err = os.Open(e.Name())
			return 
		}
	}
	file, err = os.Create(fileName)
	return
}
