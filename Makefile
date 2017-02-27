.PHONY: all test dependencies targets

SOURCES := $(shell find . -type f -name '*.go')
TARGETS := \
	linux-386 \
	linux-amd64 \
	darwin-amd64

TARGET_DIR := target
PROGRAMS := $(foreach target,$(TARGETS),$(TARGET_DIR)/logan-$(target))

all: test targets

dependencies:
	go get -v

test: dependencies
	find . -type d -not -path '*/.git*' -a -not -path '*/target*' | xargs -n 1 go test

logan: $(SOURCES)
	go build -v .

$(TARGET_DIR):
	mkdir -p $(TARGET_DIR)

targets: $(TARGET_DIR) $(PROGRAMS)

$(PROGRAMS): $(SOURCES)
	GOOS=$$(echo "$@" | xargs basename | cut -d- -f2) \
	GOARCH=$$(echo "$@" | xargs basename | cut -d- -f3) \
	go build -v -o "$@" .
