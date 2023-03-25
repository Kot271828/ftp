package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"strings"
)

func main() {
	run(context.Background())
}

func run(ctx context.Context) {
	listener, err := net.Listen("tcp", "localhost:8000")
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
		fmt.Fprintf(conn, "%s >> ", userName)
		if !scanner.Scan() {
			break
		}

		// parse
		cmd := parseCommand(scanner.Text())

		// handle command
		if cmd[0] == "QUIT" {
			conn.Close()
			break
		}
		fmt.Fprintf(conn, "\tRecieve %s command.\n", cmd[0])
	}
	log.Printf("%s's connection is closed.\n", userName)
}

func parseCommand(s string) []string {
	return strings.Split(strings.Trim(s, " "), " ")
}
