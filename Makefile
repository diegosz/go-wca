GOOS := windows
GOARCH := amd64

.PHONY: clean all

RELEASE_DIR := build
REVISION    := $(shell git rev-parse --verify HEAD | cut -c-6)
VERSION     ?= latest

examples := $(wildcard example/*)
bins     := $(patsubst example/%,$(RELEASE_DIR)/%.exe,$(examples))

all: $(RELEASE_DIR) $(bins)

$(RELEASE_DIR)/%.exe: $(wildcard example/$*/*)
	go build -o "$@" -ldflags "-X main.revision=$(REVISION) -X main.version=$(VERSION)" github.com/diegosz/go-wca/example/$*

$(RELEASE_DIR):
	mkdir -p $(RELEASE_DIR)

clean:
	rm -rf $(RELEASE_DIR)/*.*
