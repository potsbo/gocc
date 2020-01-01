BIN := bin/gocc
${BIN}: main.go
	go build -o $@ $<

.PHONY: test
test: ${BIN}
	script/test
