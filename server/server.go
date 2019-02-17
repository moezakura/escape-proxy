package main

import (
	".."
	"encoding/json"
	"flag"
	"golang.org/x/sync/errgroup"
	"net"
	"strconv"
)

const (
	BUFFER_SIZE = 0xFFFF
)

var (
	port = flag.String("p", "", "port ex:443")
)

func main() {
	flag.Parse()

	listen, _ := net.Listen("tcp", "0.0.0.0:"+*port)

	for {
		conn, _ := listen.Accept()
		go func() {
			buff := make([]byte, 500)
			num, err := conn.Read(buff)
			if err != nil {
				_ = conn.Close()
				return
			}
			length, err := strconv.Atoi(string(buff[:num]))
			if err != nil {
				_ = conn.Close()
				return
			}

			buff = make([]byte, length)
			num, err = conn.Read(buff)
			if err != nil {
				_ = conn.Close()
				return
			}
			var connectPacket escapeProxy.ConnectPacket
			err = json.Unmarshal(buff[:num], &connectPacket)
			if err != nil {
				_ = conn.Close()
				return
			}

			connectAddr := connectPacket.Addr
			remoteConn, err := net.Dial("tcp", connectAddr)
			if err != nil {
				_ = conn.Close()
				return
			}
			defer func() {
				_ = remoteConn.Close()
			}()

			var eg errgroup.Group
			eg.Go(func() error { return relay(&eg, conn, remoteConn) })
			eg.Go(func() error { return relay(&eg, remoteConn, conn) })

			if eg.Wait() != nil {
				_ = conn.Close()
			}
		}()
	}
}

func relay(eg *errgroup.Group, fromConn, toConn net.Conn) error {
	buff := make([]byte, BUFFER_SIZE)
	for {
		n, err := fromConn.Read(buff)
		if err != nil {
			return err
		}
		b := buff[:n]

		n, err = toConn.Write(b)
		if err != nil {
			return err
		}
	}
}
