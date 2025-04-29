BIN_NAME=mess

COMPILE_NAME=$(BIN_NAME)
BUILD_DIR=build

GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean

all: clean build

build:
	$(GOBUILD) -o $(BUILD_DIR)/$(COMPILE_NAME) ./cmd/mess

clean:
	$(GOCLEAN)
	-rm -r $(BUILD_DIR)

install: build
	cp $(BUILD_DIR)/$(COMPILE_NAME) ${HOME}/.local/bin/$(BIN_NAME)

uninstall:
	rm ${HOME}/.local/bin/$(BIN_NAME)

.PHONY: build clean install
