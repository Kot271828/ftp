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
	"strconv"
	"strings"

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
	reply.Send(conn, "200")

	userName := conn.RemoteAddr().String()
	log.Printf("%s's connection is opened.\n", userName)
	cwd, _ := filepath.Abs("./test_dir/server_dir/")
	data_conn_address, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:10000")

	scanner := bufio.NewScanner(conn)
	for {
		if !scanner.Scan() {
			break
		}

		// parse
		c, args := cmd.Parse(scanner.Text())
		log.Println("Recieve:", scanner.Text())

		// handle command
		switch c {
		case cmd.USER:
			reply.Send(conn, "230")
		case cmd.PORT:
			addrs := strings.Split(args[0], ",")
			i, _ := strconv.Atoi(addrs[4])
			j, _ := strconv.Atoi(addrs[5])
			addr := fmt.Sprintf("%s.%s.%s.%s:%s", addrs[0], addrs[1], addrs[2], addrs[3], strconv.Itoa(i*256+j))
			var err error
			data_conn_address, err = net.ResolveTCPAddr("tcp", addr)
			if err != nil {
				log.Fatalln(err)
			}
			reply.Send(conn, "200")
		case cmd.PWD:
			reply.Send257(conn, "257", cwd)
		case cmd.LIST:
			if len(args) != 1 {
				reply.Send(conn, "500")
				break
			}
			data_conn, err := net.Dial("tcp", data_conn_address.String())
			if err != nil {
				reply.Send(conn, "421")
				break
			}
			reply.Send(conn, "125")
			ls(data_conn, cwd, args[0])
			data_conn.Close()
			reply.Send(conn, "250")
		case cmd.RETR:
			data_conn, err := net.Dial("tcp", data_conn_address.String())
			if err != nil {
				reply.Send(conn, "421")
				break
			}
			reply.Send(conn, "125")
			cp(data_conn, cwd, args[0])
			data_conn.Close()
			reply.Send(conn, "250")
		case cmd.TYPE:
			t := args[0]
			if t != "A" {
				reply.Send(conn, "504")
				break
			}
			if (len(args) >= 2) && (args[1] != "N") {
				reply.Send(conn, "504")
				break
			}
			reply.Send(conn, "200")
		case cmd.MODE:
			m := args[0]
			if m == "S" {
				reply.Send(conn, "200")
			} else {
				reply.Send(conn, "504")
			}
		case cmd.STRU:
			s := args[0]
			if s == "F" {
				reply.Send(conn, "200")
			} else {
				reply.Send(conn, "504")
			}
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
