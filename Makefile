MAIN := ./main.go
OUTPUT := suggester-gateway

PKG := $(shell find . -type f -name '*.go')
DEPS := $(shell find . -type f -name '*.mod')

$(OUTPUT): $(MAIN) $(PKG) $(DEPS)
	go build -o $(OUTPUT) $(MAIN) -ldflags "-s -w"

build: $(OUTPUT)

build-alpine: $(MAIN) $(PKG) $(DEPS)
	CGO_ENABLED=0 GOOS=linux go build -o $(OUTPUT) -installsuffix cgo -ldflags "-s -w" $(MAIN)

deps: $(DEPS)
	go get

clean:
	go clean
	rm -f $(OUTPUT)

.PHONY: dev
dev:
	go run $(MAIN) -- -c config.toml
