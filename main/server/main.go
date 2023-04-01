package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"

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
	cwd, _ := filepath.Abs("./test_dir/server_dir/")

	scanner := bufio.NewScanner(conn)
	for {
		if !scanner.Scan() {
			break
		}

		// parse
		c, args := cmd.Parse(scanner.Text())
		log.Println("Recieve:", c, args)

		// handle command
		var reply string
		switch c {
		case cmd.USER:
			reply = "230"
		case cmd.PWD:
			reply = fmt.Sprintf("257 %s created.", cwd)
		case cmd.LIST:
			data_conn, err := net.Dial("tcp", "localhost:10000")
			if err != nil {
				reply = "421"
				break
			}
			fmt.Fprintf(conn, "%s\n", "125")
			ls(data_conn, cwd, args[0])
			data_conn.Close()
			reply = "250"
		case cmd.RETR:
			data_conn, err := net.Dial("tcp", "localhost:10000")
			if err != nil {
				reply = "421"
				break
			}
			fmt.Fprintf(conn, "%s\n", "125")
			cp(data_conn, cwd, args[0])
			data_conn.Close()
			reply = "250"
		case cmd.QUIT:
			conn.Close()
			break
		case cmd.UNKNOWN:
			reply = "502"
		case cmd.NOOP:
			reply = "200"
		}
		fmt.Fprintf(conn, "%s\n", reply)
	}
	log.Printf("%s's connection is closed.\n", userName)
}

func ls(w io.Writer, cwd, arg string) {
	var p string
	if filepath.IsAbs(arg) {
		p = arg
	} else {
		p = filepath.Join(cwd, arg)
		p, _ = filepath.Abs(p)
	}
	matches, err := filepath.Glob(fmt.Sprintf("%s/*", p))
	if err != nil {
		log.Println(err)
		return
	}
	for _, match := range matches {
		fmt.Fprintln(w, match)
	}

}

func cp(w io.Writer, cwd, arg string) {
	var p string
	if filepath.IsAbs(arg) {
		p = arg
	} else {
		p = filepath.Join(cwd, arg)
		p, _ = filepath.Abs(p)
	}
	f, err := os.Open(p)
	if err != nil {
		log.Println(err)
		return
	}
	io.Copy(w, f)
	f.Close()
}
