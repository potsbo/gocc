BIN := bin/gocc
SRC=$(wildcard *.go) $(wildcard */*.go)

${BIN}: ${SRC}
	go build -o $@ .

.PHONY: test
test: ${BIN}
	script/test
