package main

import (
	"bufio"
	"fmt"
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
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
		done <- struct{}{}
	}()

	input := stdin()
	for {
		select {
		case cl := <-input:
			cmd, args := parseCommand(cl)
			fmt.Println(cmd, args)

			fmt.Fprintln(conn, cl)
		case <-done:
			conn.Close()
			return
		}
	}
}

func stdin() <-chan string {
	stdin := make(chan string)
	scanner := bufio.NewScanner(os.Stdin)
	go func() {
		for scanner.Scan() {
			stdin <- scanner.Text()
		}
	}()
	return stdin
}

func parseCommand(s string) (string, []string) {
	args := strings.Split(strings.Trim(s, " "), " ")
	cmd := strings.ToUpper(args[0])
	args = args[1:]
	return cmd, args
}
