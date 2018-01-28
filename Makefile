
ANTLR=java -Xmx500M -cp "antlr-4.7-complete.jar" org.antlr.v4.Tool

build: jt

jt: parser/program_listener.go *.go */*.go
	go build
	go test ./...

parser/program_listener.go: parser/program.g4
	$(ANTLR) -Dlanguage=Go parser/Program.g4

test: tests/re_line_match_implicit \
	  tests/re_line_match_implicit_pipe_delimited \
	  tests/re_line_match_comparison_operator \
      tests/re_column_match_comparison_operator

tests/re_line_match_implicit::
	tests/test tests/re_line_match_implicit

tests/re_line_match_implicit_pipe_delimited::
	tests/test tests/re_line_match_implicit_pipe_delimited

tests/re_line_match_comparison_operator::
	tests/test tests/re_line_match_comparison_operator

tests/re_column_match_comparison_operator::
	tests/test tests/re_column_match_comparison_operator

clean::
	rm -f \
		parser/*.tokens \
		parser/program_*.go \
		jt
