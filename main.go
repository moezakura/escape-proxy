package main

import (
	"errors"
	"flag"
	"github.com/moezakura/escape-proxy/client"
	"github.com/moezakura/escape-proxy/server"
	"os"
)

const (
	SERVER_MODE_STRING = "server"
	CLIENT_MODE_STRING = "client"
)

func main() {
	flag.Parse()
	if len(os.Args) < 2 {
		panic("Client or server sub command not specified.")
	}

	serverArg := flag.NewFlagSet("server", flag.ExitOnError)
	serverPort := serverArg.String("s", "443", "server listen port ex) -s 443")
	if os.Args[1] == "server" {
		if parseArgs(serverArg) != nil{
			return
		}
	}

	clientArg := flag.NewFlagSet("client", flag.ExitOnError)
	configPath := clientArg.String("c", "./config.yaml", "client config file path")
	if os.Args[1] == "client" {
		if parseArgs(clientArg) != nil{
			return
		}
	}

	subCommand := flag.Args()[0]
	if subCommand == SERVER_MODE_STRING {
		server.Server(*serverPort)
	} else if subCommand == CLIENT_MODE_STRING {
		client.Client(*configPath)
	}
}

func parseArgs(set *flag.FlagSet) error{
	if err := set.Parse(os.Args[2:]); err != nil {
		return errors.New("")
	}
	return nil
}