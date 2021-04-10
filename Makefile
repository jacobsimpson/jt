
jt: *.go datetime/*.go ast/*.go parser/grammar.go
	go build .

parser/grammar.go: parser/grammar.peg
	go generate ./...

test: jt
	go test ./...

tests: jt FORCE
	go test ./...
	tests/test

install: jt
	cp ./jt /usr/bin/jt

.PHONY: FORCE

clean::
	rm -Rf jt parser/grammar.go
