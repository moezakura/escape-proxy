GO111MODULE=on

build:
	go build -o dist/escape-proxy

build-all-platform:
	make build-mac
	make build-windows

build-mac:
	env GOOS=darwin go build -o dist/escape-proxy-mac

build-windows:
	env GOOS=windows go build -o dist/escape-proxy-windows.exe

build-linux:
	env GOOS=linux go build -o dist/escape-proxy-linux
