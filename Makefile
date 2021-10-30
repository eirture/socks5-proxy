
version ?= unknown

.PHONY: build
build:
	@mkdir -p bin
	go build -o bin/socks5-proxy ./cmd/...

.PHONY: pkg
pkg: clean
	make GOOS=darwin
	cd bin && zip socks5-proxy-${version}-darwin-amd64.zip socks5-proxy
	cd ../
	make GOOS=linux
	cd bin && tar zcf socks5-proxy-${version}-linux-amd64.tar.gz socks5-proxy

.PHONY: clean
clean:
	rm -r ./bin/*
