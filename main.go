package main

import (
	"io"
	"log"
	"net"
	"strconv"
)

var address = "habr.com"
var ports = []int{80, 443}

func handleConnection(conn net.Conn, port int) {
	destination, err := net.Dial("tcp", address+":"+strconv.Itoa(port))
	if err != nil {
		log.Println("Error connecting to destination:", err)
		conn.Close()
		return
	}
	defer destination.Close()

	go func() {
		_, err := io.Copy(destination, conn)
		if err != nil {
			log.Println("Error copying from client to destination:", err)
		}
	}()

	_, err = io.Copy(conn, destination)
	if err != nil {
		log.Println("Error copying from destination to client:", err)
	}
}

func mdain() {
	for _, port := range ports {
		go func(port int) {
			listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
			if err != nil {
				log.Fatal("Error starting listener:", err)
			}
			defer listener.Close()

			log.Printf("Listening on port %d\n", port)

			for {
				conn, err := listener.Accept()
				if err != nil {
					log.Println("Error accepting connection:", err)
					continue
				}

				go handleConnection(conn, port)
			}
		}(port)
	}

	select {}
}
