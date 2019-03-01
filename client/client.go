package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/armon/go-socks5"
	"github.com/moezakura/escape-proxy/model"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net"
	"regexp"
	"strings"
	"time"
)

var (
	config      model.ConfigYaml
	connectMode model.CONNECT_MODE = model.CONNECT_MODE_PROXY
)

func Client(configPath string) {
	buf, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(buf, &config)
	if err != nil {
		panic(err)
	}

	proxy := config.ProxyServer
	next_proxy := config.GatewayServer
	if len(proxy) == 0 {
		panic("config has no proxy setting.")
	}
	if len(next_proxy) == 0 {
		panic("config has no gateway setting.")
	}
	if len(next_proxy) == 0 {
		panic("config has no listen setting.")
	}

	conf := &socks5.Config{
		Credentials: NewAuth(config.Users),
		Dial: func(ctx context.Context, network, addr string) (conn net.Conn, e error) {
			fmt.Println("-------------")
			connectAddr := proxy
			fmt.Printf("CONNECT MODE: %s\n", connectMode.String())
			if connectMode == model.CONNECT_MODE_DIRECT {
				connectAddr = addr
			}

			isForceDirect := false
			for _, ipRange := range config.ExcludeIps {
				targetIp, _, err := net.SplitHostPort(addr)
				if err != nil {
					return nil, err
				}
				if isContainIp(ipRange, targetIp){
					fmt.Printf("EXCLUDE : %s in %s\n", ipRange, targetIp)
					connectAddr = addr
					isForceDirect = true
					break
				}
			}

			if !isForceDirect {
				printRoute(connectMode, addr)
			}else{
				fmt.Printf("FORCE ")
				printRoute(model.CONNECT_MODE_DIRECT, addr)
			}

			retryCount := 0
			_connectMode := connectMode
		retry:
			n, e := net.DialTimeout(network, connectAddr, time.Second*5)
			if e != nil {
				if config.AutoDirectConnect {
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
				}

				return nil, e
			}
			if retryCount > 0 {
				connectMode = _connectMode
			}

			if connectMode == model.CONNECT_MODE_PROXY && !isForceDirect {
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

	if err := server.ListenAndServe("tcp", config.Listen); err != nil {
		panic(err)
	}
}

func isContainIp(cidr, ip string) bool {
	_, cidrNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return false
	}

	targetIP := net.ParseIP(ip)
	return cidrNet.Contains(targetIP)
}

func printRoute(mode model.CONNECT_MODE, addr string) {
	proxy := config.ProxyServer
	next_proxy := config.GatewayServer

	if mode == model.CONNECT_MODE_DIRECT {
		fmt.Printf("ROUTE: localhost -> %s\n", addr)
	} else {
		fmt.Printf("ROUTE: localhost -> %s -> %s -> %s\n", proxy, next_proxy, addr)
	}
}

func proxyConnect(n net.Conn, addr string) (err error) {
	reg := regexp.MustCompile(`HTTP/(1\.0|1\.1|2\.0) 200 Connection established`)
	next_proxy := config.GatewayServer

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
