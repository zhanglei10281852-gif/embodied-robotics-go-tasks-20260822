GOFLAGS=-mod=mod
export GOTOOLCHAIN=local

build:
	go build ./...

test:
	go test ./... -count=1

race:
	go test -race ./... -count=1

vet:
	go vet ./...
