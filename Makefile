BIN_NAME=docker-event-hooks
VERSION?=development

all: clean build

build:
	go build -ldflags="-s -w -X main.version=$(VERSION) -X main.binName=$(BIN_NAME)" -o $(BIN_NAME) -v

clean:
	rm -rf $(BIN_NAME)

.PHONY: all build clean