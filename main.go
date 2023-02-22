package main

import (
	"log"
	"os"
)

func main() {

	args := os.Args[1:]
	if len(args) < 2 || len(args) > 2 {
		log.Panic("invalid argument please read the doc")
	}

	if args[0] == "seed" {
		 path,torrentStruct, err := GenerateSeeder(args[1])
		if err != nil {
			log.Panic("seeding Failed", err)
			return
		}
		log.Print("the generated torrent file is located at: \n", path)

		StartListen(torrentStruct)
	}
	if args[0] == "leech" {
		//leech(args[1])
		ManageLeech(args[1])
		log.Print("not now")
	}
	//Write("steven.mp4", "anwii.mp4")
}
