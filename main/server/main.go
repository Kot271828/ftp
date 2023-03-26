package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"

	"ftp/cmd"
)

func main() {
	run(context.Background())
}

func run(ctx context.Context) {
	listener, err := net.Listen("tcp", "localhost:10021")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		child_ctx, cancel := context.WithCancel(ctx)
		defer cancel()
		handleConn(child_ctx, conn)
	}
}

func handleConn(ctx context.Context, conn net.Conn) {
	userName := conn.RemoteAddr().String()
	log.Printf("%s's connection is opened.\n", userName)

	scanner := bufio.NewScanner(conn)
	for {
		if !scanner.Scan() {
			break
		}

		// parse
		c, args := cmd.Parse(scanner.Text())
		log.Println("Recieve:", c, args)

		// handle command
		var replyCode string
		switch c {
		case cmd.USER:
			replyCode = "200"
		case cmd.QUIT:
			conn.Close()
			break
		}
		fmt.Fprintf(conn, "%s\n", replyCode)
	}
	log.Printf("%s's connection is closed.\n", userName)
}
