BINARY_NAME=goph-keeper

hello:
	echo "Hello"

build:
	GOARCH=amd64 GOOS=darwin go build -o ${BINARY_NAME}-darwin-server cmd/server/server.go
	GOARCH=amd64 GOOS=darwin go build -o ${BINARY_NAME}-darwin-client cmd/client/client.go
	GOARCH=amd64 GOOS=linux go build -o ${BINARY_NAME}-linux-server cmd/server/server.go
	GOARCH=amd64 GOOS=linux go build -o ${BINARY_NAME}-linux-client cmd/client/client.go
	GOARCH=amd64 GOOS=windows go build -o ${BINARY_NAME}-windows-server cmd/server/server.go
	GOARCH=amd64 GOOS=windows go build -o ${BINARY_NAME}-windows-client cmd/client/client.go

clean:
	go clean
	rm ${BINARY_NAME}-darwin-server
	rm ${BINARY_NAME}-darwin-client
	rm ${BINARY_NAME}-linux-server
	rm ${BINARY_NAME}-linux-client
	rm ${BINARY_NAME}-windows-server
	rm ${BINARY_NAME}-windows-client

test:
	go test ./...

test_coverage:
	go test ./... -coverprofile=coverage.out

dep:
	go mod download

vet:
	go vet

lint:
	golangci-lint run --enable-all