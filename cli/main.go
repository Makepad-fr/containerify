package main

import (
	"log"

	"github.com/containerd/containerd"
)

func main() {
	client, err := containerd.New("/run/containerd/containerd.sock")
	if err != nil {
		// handle error
		log.Fatalln("Conatinerd is not running")
	}
	log.Println("Containerd is running")
	defer client.Close()

}
