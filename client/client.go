package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/armon/go-socks5"
	"github.com/moezakura/EscapeProxy/model"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net"
	"regexp"
	"strings"
)

var (
	proxyAddress  = flag.String("p", "8080", "proxy server address ex: proxy.mox:8080")
	serverAddress = flag.String("s", "8080", "gateway proxy server address ex: proxy.mox:8080")
	listenPort    = flag.String("l", "9999", "local socks port es: 8080")
	configPath    = flag.String("c", "", "./config.yaml")
)

func main() {
	flag.Parse()

	buf, err := ioutil.ReadFile(*configPath)
	if err != nil {
		panic(err)
	}
	var config model.ConfigYaml
	err = yaml.Unmarshal(buf, &config)
	if err != nil {
		panic(err)
	}


	reg := regexp.MustCompile(`HTTP/(1\.0|1\.1|2\.0) 200 Connection established`)
	proxy := *proxyAddress
	next_proxy := *serverAddress

	conf := &socks5.Config{
		Credentials: NewAuth(config.Users),
		Dial: func(ctx context.Context, network, addr string) (conn net.Conn, e error) {
			fmt.Printf("network: %s\n", network)
			fmt.Printf("addr: %s\n", addr)
			fmt.Printf("next: %s\n", next_proxy)

			n, e := net.Dial("tcp", proxy)
			if e != nil {
				return nil, e
			}
			num, err := n.Write([]byte("CONNECT " + next_proxy + " HTTP/1.1\r\n\r\n"))
			if err != nil {
				return nil, err
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

			jsonBytes, err := json.Marshal(model.ConnectPacket{
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

	if err := server.ListenAndServe("tcp", "127.0.0.1:"+*listenPort); err != nil {
		panic(err)
	}
}
