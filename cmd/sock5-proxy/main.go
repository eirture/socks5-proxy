package main

import (
	"flag"
	"io"
	"log"
	"net"

	"github.com/eirture/sock5-proxy/lib/socks5"
)

func Socks5Forward(client, target net.Conn) {
	log.Printf("forward: %s -> %s\n", client.RemoteAddr(), target.RemoteAddr())
	forward := func(src, dest net.Conn) {
		defer src.Close()
		defer dest.Close()
		io.Copy(src, dest)
	}
	go forward(client, target)
	go forward(target, client)
}

func process(client net.Conn) {
	if err := socks5.Auth(client); err != nil {
		log.Println("auth error:", err)
		client.Close()
		return
	}

	target, err := socks5.Connect(client)
	if err != nil {
		log.Println("connect error:", err)
		client.Close()
		return
	}

	Socks5Forward(client, target)
}

func main() {
	var addr = ""
	flag.StringVar(&addr, "address", ":1080", "")
	flag.Parse()

	server, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Listen failed: %v\n", err)
	}

	for {
		client, err := server.Accept()
		if err != nil {
			log.Printf("Accept failed: %v\n", err)
			continue
		}
		go process(client)
	}
}
