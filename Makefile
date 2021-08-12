default: build

bin:
	mkdir -p bin

build: bin
	go build -o bin/postier cmd/postier.go
