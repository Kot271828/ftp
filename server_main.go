package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

func main() {
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
		handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {

		// parse
		cmd := parseCommand(scanner.Text())

		// handle command
		for i, arg := range cmd {
			fmt.Fprintln(conn, i, arg)
		}
	}
}

func parseCommand(s string) []string {
	return strings.Split(strings.Trim(s, " "), " ")
}
