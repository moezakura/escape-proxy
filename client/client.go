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
	"time"
)

var (
	proxyAddress                     = flag.String("p", "8080", "proxy server address ex: proxy.mox:8080")
	serverAddress                    = flag.String("s", "8080", "gateway proxy server address ex: proxy.mox:8080")
	listenPort                       = flag.String("l", "9999", "local socks port es: 8080")
	configPath                       = flag.String("c", "", "./config.yaml")
	connectMode   model.CONNECT_MODE = model.CONNECT_MODE_PROXY
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

	proxy := *proxyAddress
	next_proxy := *serverAddress

	conf := &socks5.Config{
		Credentials: NewAuth(config.Users),
		Dial: func(ctx context.Context, network, addr string) (conn net.Conn, e error) {
			fmt.Println("-------------")
			connectAddr := proxy
			fmt.Printf("CONNECT MODE: %s\n", connectMode.String())
			if connectMode == model.CONNECT_MODE_DIRECT {
				connectAddr = addr
			}
			printRoute(connectMode, addr)

			retryCount := 0
			_connectMode := connectMode
		retry:
			n, e := net.DialTimeout(network, connectAddr, time.Second * 5)
			if e != nil {
				retryCount++
				if retryCount < 3 {
					if connectMode == model.CONNECT_MODE_PROXY {
						connectAddr = addr
						_connectMode = model.CONNECT_MODE_DIRECT
						fmt.Printf("!CHANGE [%d]", retryCount)
						printRoute(_connectMode, addr)
						goto retry
					} else if connectMode == model.CONNECT_MODE_DIRECT {
						connectAddr = next_proxy
						_connectMode = model.CONNECT_MODE_PROXY
						fmt.Printf("!CHANGE [%d]", retryCount)
						printRoute(_connectMode, addr)
						goto retry
					}
				}

				return nil, e
			}
			if retryCount > 0 {
				connectMode = _connectMode
			}

			if connectMode == model.CONNECT_MODE_PROXY {
				err := proxyConnect(n, addr)
				if err != nil {
					return nil, err
				}
			}

			fmt.Println("CONNECT: OK!")
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

func printRoute(mode model.CONNECT_MODE, addr string) {
	proxy := *proxyAddress
	next_proxy := *serverAddress

	if mode == model.CONNECT_MODE_DIRECT {
		fmt.Printf("ROUTE: localhost -> %s\n", addr)
	} else {
		fmt.Printf("ROUTE: localhost -> %s -> %s -> %s\n", proxy, next_proxy, addr)
	}
}

func proxyConnect(n net.Conn, addr string) (err error) {
	reg := regexp.MustCompile(`HTTP/(1\.0|1\.1|2\.0) 200 Connection established`)
	next_proxy := *serverAddress

	num, err := n.Write([]byte("CONNECT " + next_proxy + " HTTP/1.1\r\n\r\n"))
	if err != nil {
		return err
	}

	buff := make([]byte, 1000)
	num, err = n.Read(buff)
	if err != nil {
		return err
	}

	res := strings.Replace(string(buff[:num]), "\r\n", "", -1)
	if !reg.MatchString(res) {
		fmt.Println("CONNECT: NG!")
		return errors.New("access error")
	}

	jsonBytes, err := json.Marshal(model.ConnectPacket{
		Addr: addr,
	})
	lengthPacket := fmt.Sprintf("%0500d", len(jsonBytes))
	_, err = n.Write([]byte(lengthPacket))
	num, err = n.Write(jsonBytes)
	return nil
}
