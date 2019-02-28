package server

import (
	"encoding/json"
	"fmt"
	"github.com/moezakura/escape-proxy/model"
	"golang.org/x/sync/errgroup"
	"net"
	"strconv"
)

const (
	BUFFER_SIZE = 0xFFFF
)

func Server(port string) {
	listen, _ := net.Listen("tcp", "0.0.0.0:"+port)
	fmt.Printf("Listen 0.0.0.0:%s\n", port)

	for {
		conn, _ := listen.Accept()
		go func() {
			defer func() {
				_ = conn.Close()
			}()

			buff := make([]byte, 500)
			num, err := conn.Read(buff)
			if err != nil {
				return
			}
			length, err := strconv.Atoi(string(buff[:num]))
			if err != nil {
				return
			}

			buff = make([]byte, length)
			num, err = conn.Read(buff)
			if err != nil {
				return
			}
			var connectPacket model.ConnectPacket
			err = json.Unmarshal(buff[:num], &connectPacket)
			if err != nil {
				return
			}

			connectAddr := connectPacket.Addr
			remoteConn, err := net.Dial("tcp", connectAddr)
			if err != nil {
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
