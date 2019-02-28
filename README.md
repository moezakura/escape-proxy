# Escape Proxy

## what is EscapeProxy?
It is an application for TCP communication on squid.  
In this application it is possible to build socks 5 proxy via squid.

### Note
In this application, access may be made by bypassing the http proxy unless prohibited by setting.

## How to install and build
```
go get github.com/moezakura/escape-proxy 
cd $GOPATH/github.com/moezakura/escape-proxy
make build
```

## Usage
### client mode
```
./escape-proxy client -c config.yaml
```

### server mode
```
./escape-proxy server -s [listen port]
```

## Config
Use yaml for config.  

```
# authMode[true|false]
auth: true
# Whether to allow access automatically bypassing the Proxy [true|false]
auto_direct_connect: true
# http proxy host[hostName:port]
proxy: http.proxy:3128 
# escape proxy server[hostName:port]
# recommend 443 (Many squids prohibit the connect method except for the 443 port.)
gateway: escape-proxy.server:443
# listen socks5 proxy[bindAddress:port]
listen: localhost:9999
# It is specified when auth mode is true. (authentication of listening socks 5)
users:
# auth user id
  - id: mox  
# auth user password
    password: kafu-chino_ha_sekaide_itiban_kawaii
```

## License 
- MIT License
- Copyright (c) 2019 Moezakura