BIN := bin/gocc
${BIN}: main.go token/token.go
	go build -o $@ $<

.PHONY: test
test: ${BIN}
	script/test
