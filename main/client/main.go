package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"

	"ftp/cmd"
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
			cmd, args := cmd.Parse(cl)
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
