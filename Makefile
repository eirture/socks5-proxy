

.PHONY: build
build:
	@mkdir -p bin
	go build -o bin/socks5-proxy ./cmd/...

.PHONY: clean
clean:
	rm -r ./bin/*
