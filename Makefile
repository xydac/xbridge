.PHONY: build test clean run

build:
	go build -o xbridge ./cmd/xbridge

test:
	go test ./...

clean:
	rm -f xbridge
	rm -rf DerivedData

run:
	go run ./cmd/xbridge serve

# Cross-compilation
build-darwin-arm64:
	GOOS=darwin GOARCH=arm64 go build -o xbridge-darwin-arm64 ./cmd/xbridge

build-darwin-amd64:
	GOOS=darwin GOARCH=amd64 go build -o xbridge-darwin-amd64 ./cmd/xbridge

build-linux-amd64:
	GOOS=linux GOARCH=amd64 go build -o xbridge-linux-amd64 ./cmd/xbridge

build-all: build-darwin-arm64 build-darwin-amd64 build-linux-amd64
