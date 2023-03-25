package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:10021")
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan struct{})
	go func() {
		io.Copy(os.Stdout, conn)
		done <- struct{}{}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		cmd, args := parseCommand(scanner.Text())
		fmt.Println(cmd, args)

		fmt.Fprintln(conn, scanner.Text())
	}
	conn.Close()
	<-done
}

func parseCommand(s string) (string, []string) {
	args := strings.Split(strings.Trim(s, " "), " ")
	cmd := strings.ToUpper(args[0])
	args = args[1:]
	return cmd, args
}
