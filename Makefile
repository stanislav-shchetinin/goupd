BINARY := goupd
PKG := ./...

.PHONY: all build test cover fmt fmt-check vet lint tidy run clean

all: fmt vet test build

build:
	go build -o $(BINARY) ./cmd/goupd

test:
	go test -race -cover $(PKG)

cover:
	go test -coverprofile=coverage.out $(PKG)
	go tool cover -func=coverage.out

fmt:
	gofmt -w .

fmt-check:
	@out="$$(gofmt -l .)"; if [ -n "$$out" ]; then echo "Not gofmt'ed:"; echo "$$out"; exit 1; fi

vet:
	go vet $(PKG)

lint:
	golangci-lint run

tidy:
	go mod tidy

run:
	go run ./cmd/goupd $(ARGS)

clean:
	rm -f $(BINARY) coverage.out
