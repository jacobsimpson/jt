
jt: *.go datetime/*.go ast/*.go parser/grammar.go
	go build .

parser/grammar.go: parser/grammar.peg
	go generate ./...

test: jt
	go test ./...

tests: jt tests/test tests/*/*
	go test ./...
	tests/test

install: jt
	cp ./jt /usr/bin/jt

clean::
	rm -Rf jt parser/grammar.go
