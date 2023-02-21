package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
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

func manageLeech(torrentFile string)(torrentStruct Torrent,  err error){
	var file *os.File
	torrentStruct, err = myUnmarshall(torrentFile)
	if err!=nil{
		log.Panic("error unmarshalling torrent file", err)
	}
	file, err = OpenExistingFile(torrentStruct.Name)
	if err!= nil{
		log.Fatal("error handling file")
	}
	
	defer file.Close()
	

	var missingPieces []int
	missingPieces,err = findMissingPieces(torrentStruct,file ) 
		if err!=nil{
			log.Panic("cant resume", err)
			return 
	}
	for _, _ = range missingPieces {

       //send torrent request
	   receiveData()
    }
	return 
}
