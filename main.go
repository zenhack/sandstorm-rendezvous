package main

import (
	"log"
	"os"
)

func usage() {
	log.Fatalf(
		"Usage: %v ( listen | connect ) wss://api-XXXX.sandstorm.example.com/...",
		os.Args[0],
	)
}

func main() {
	if len(os.Args) <= 1 {
		usage()
	}
	switch os.Args[1] {
	case "app":
		serverMain()
	case "listen":
		if len(os.Args) <= 2 {
			usage()
		}
		listenMain(os.Args[2])
	case "connect":
		if len(os.Args) <= 2 {
			usage()
		}
		connectMain(os.Args[2])
	default:
		usage()
	}
}
