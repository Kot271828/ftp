package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"

	"ftp/cmd"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:10021")
	if err != nil {
		log.Fatal(err)
	}

	reply := make(chan string)
	done := make(chan struct{})
	// input from user
	go func() {
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			reply <- scanner.Text()
		}
		done <- struct{}{}
	}()

	// login
	fmt.Fprintf(conn, "%s annoymous\n", cmd.USER)
	//s := <-reply
	log.Println(<-reply)

	// user PI
	input := stdin()
	for {
		select {
		case cl := <-input:
			switch strings.Split(cl, " ")[0] {
			case "exit":
				cmd := fmt.Sprint(cmd.QUIT)
				fmt.Fprintln(conn, cmd)
			case "pwd":
				listener, err := net.Listen("tcp", "localhost:10000")
				if err != nil {
					log.Println(err)
					continue
				}
				cmd := fmt.Sprint(cmd.PWD)
				fmt.Fprintln(conn, cmd)

				data_conn, err := listener.Accept()
				if err != nil {
					log.Println(err)
					continue
				}
				io.Copy(os.Stdin, data_conn)
				data_conn.Close()
				listener.Close()

				log.Println(<-reply)
			}

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
