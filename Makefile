BINARY_NAME=server

build:
	GOARCH=amd64 GOOS=darwin go build -o ${BINARY_NAME}-mac cmd/server/main.go
	GOARCH=amd64 GOOS=linux go build -o ${BINARY_NAME}-linux cmd/server/main.go
	GOARCH=amd64 GOOS=windows go build -o ${BINARY_NAME}-windows cmd/server/main.go

run: build
	./${BINARY_NAME}

clean:
	go clean
	rm ${BINARY_NAME}-mac
	rm ${BINARY_NAME}-linux
	rm ${BINARY_NAME}-windows

test:
	go test ./...

test_coverage:
	go test ./... -coverprofile=coverage.out

dep:
	go mod download

vet:
	go vet
