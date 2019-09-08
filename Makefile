all: build

build:
	mkdir -p bin
	go build -o bin/rkyv ./cmd/rkyv

.PHONY = build
