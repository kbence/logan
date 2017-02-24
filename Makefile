.PHONY: all test dependencies targets

SOURCES = $(shell find . -type d -name '*.go')
TARGETS = \
	linux-386 \
	linux-amd64 \
	darwin-amd64

PROGRAMS = $(foreach target,$(TARGETS),target/logan-$(target))
TARGET_DIR = ./target

all: test targets

dependencies:
	go get -v

test: dependencies
	find . -type d -not -path '*/.git*' | xargs -n 1 go test

logan: $(SOURCES)
	go build -v .

$(TARGET_DIR):
	mkdir -p $(TARGET_DIR)

targets: $(PROGRAMS)

define TARGET_template =
	target/logan-$(1): $(TARGET_DIR) $(SOURCES)
endef

$(PROGRAMS): $(SOURCES)
	GOOS=$$(echo "$@" | xargs basename | cut -d- -f2) \
	GOARCH=$$(echo "$@" | xargs basename | cut -d- -f3) \
	go build -v -o "$@" .
