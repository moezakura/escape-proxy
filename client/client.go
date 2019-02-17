package main

import (
	".."
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/armon/go-socks5"
	"net"
	"regexp"
	"strings"
)

var (
	proxyAddress  = flag.String("p", "", "proxy server address ex: proxy.mox:8080")
	serverAddress = flag.String("s", "", "gateway proxy server address ex: proxy.mox:8080")
)

func main() {
	flag.Parse()

	reg := regexp.MustCompile(`HTTP/(1\.0|1\.1|2\.0) 200 Connection established`)
	proxy := *proxyAddress
	next_proxy := *serverAddress

	conf := &socks5.Config{
		Dial: func(ctx context.Context, network, addr string) (conn net.Conn, e error) {
			fmt.Printf("network: %s\n", network)
			fmt.Printf("addr: %s\n", addr)
			fmt.Printf("next: %s\n", next_proxy)

			n, e := net.Dial("tcp", proxy)
			num, err := n.Write([]byte("CONNECT " + next_proxy + " HTTP/1.1\r\n\r\n"))
			if err != nil {
				return nil, e
			}

			buff := make([]byte, 1000)
			num, err = n.Read(buff)
			if err != nil {
				return nil, e
			}

			res := strings.Replace(string(buff[:num]), "\r\n", "", -1)
			fmt.Println(res)
			if !reg.MatchString(res) {
				fmt.Println("NG!")
				return nil, errors.New("access error")
			}

			jsonBytes, err := json.Marshal(escapeProxy.ConnectPacket{
				Addr: addr,
			})
			lengthPacket := fmt.Sprintf("%0500d", len(jsonBytes))
			_, err = n.Write([]byte(lengthPacket))
			num, err = n.Write(jsonBytes)

			fmt.Println("OK!")
			return n, nil
		},
	}
	server, err := socks5.New(conf)
	if err != nil {
		panic(err)
	}

	// Create SOCKS5 proxy on localhost port 8000
	if err := server.ListenAndServe("tcp", "127.0.0.1:9999"); err != nil {
		panic(err)
	}
}
