test: vet lint
	go test -v ./...

vet:
	go vet ./...

lint:
	golangci-lint run

install:
	go build -o $(GOPATH)/bin/stickers ./cmd/cli/...
