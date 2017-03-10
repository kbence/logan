.PHONY: all test dependencies targets

SOURCES := $(shell find . -type f -name '*.go')
PEGS := $(shell find . -type f -name '*.peg')

TARGETS := \
	linux-386 \
	linux-amd64 \
	darwin-amd64

TARGET_DIR := target
PROGRAMS := $(foreach target,$(TARGETS),$(TARGET_DIR)/logan-$(target))
GENERATED_PEGS := $(PEGS:.peg=.peg.go)

all: parsers test targets

dependencies:
	go get -v

parsers: $(GENERATED_PEGS)

test: dependencies
	find . -type d -not -path '*/.git*' -a -not -path '*/target*' -a -not -path '*/docs*' | \
		xargs -n 1 go test

logan: $(SOURCES)
	go build -v .

$(TARGET_DIR):
	mkdir -p $(TARGET_DIR)

targets: $(TARGET_DIR) $(PROGRAMS)

$(PROGRAMS): $(SOURCES)
	GOOS=$$(echo "$@" | xargs basename | cut -d- -f2) \
	GOARCH=$$(echo "$@" | xargs basename | cut -d- -f3) \
	go build -v -o "$@" .

%.peg.go: %.peg
	peg $<
