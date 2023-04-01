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
	"ftp/reply"
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
		switch c {
		case cmd.USER:
			reply.Send(conn, "230")
		case cmd.PWD:
			reply.Send257(conn, "257", cwd)
		case cmd.LIST:
			data_conn, err := net.Dial("tcp", "localhost:10000")
			if err != nil {
				reply.Send(conn, "421")
				break
			}
			reply.Send(conn, "125")
			ls(data_conn, cwd, args[0])
			data_conn.Close()
			reply.Send(conn, "250")
		case cmd.RETR:
			data_conn, err := net.Dial("tcp", "localhost:10000")
			if err != nil {
				reply.Send(conn, "421")
				break
			}
			reply.Send(conn, "125")
			cp(data_conn, cwd, args[0])
			data_conn.Close()
			reply.Send(conn, "250")
		case cmd.QUIT:
			conn.Close()
			break
		case cmd.UNKNOWN:
			reply.Send(conn, "502")
		case cmd.NOOP:
			reply.Send(conn, "200")
		}

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
