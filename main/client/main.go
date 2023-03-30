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
	// reply from server
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
			args := strings.Split(cl, " ")
			c := args[0]
			args = args[1:]
			switch c {
			case "exit":
				cmd := fmt.Sprint(cmd.QUIT)
				fmt.Fprintln(conn, cmd)
			case "pwd":
				cmd := fmt.Sprint(cmd.PWD)
				fmt.Fprintln(conn, cmd)

				fmt.Println(strings.Split(<-reply, " ")[1])
			case "ls":
				listener, err := net.Listen("tcp", "localhost:10000")
				if err != nil {
					log.Println(err)
					continue
				}
				cmd := fmt.Sprint(cmd.LIST)
				fmt.Fprintf(conn, "%s %s\n", cmd, args[0])

				data_conn, err := listener.Accept()
				if err != nil {
					log.Println(err)
					continue
				}
				log.Println(<-reply)
				io.Copy(os.Stdin, data_conn)
				data_conn.Close()
				listener.Close()

				log.Println(<-reply)
			default:
				fmt.Println("command not found.")
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
