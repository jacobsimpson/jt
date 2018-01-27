
ANTLR=java -Xmx500M -cp "antlr-4.7-complete.jar" org.antlr.v4.Tool

build: jt

jt: parser/Program.g4 *.go
	$(ANTLR) -Dlanguage=Go parser/Program.g4
	go build
	go test

test: tests/re_line_match_implicit \
	  tests/re_line_match_implicit_pipe_delimited \
	  tests/re_line_match_comparison_operator

tests/re_line_match_implicit::
	tests/test tests/re_line_match_implicit

tests/re_line_match_implicit_pipe_delimited::
	tests/test tests/re_line_match_implicit_pipe_delimited

tests/re_line_match_comparison_operator::
	tests/test tests/re_line_match_comparison_operator

clean::
	rm -f \
		parser/*.tokens \
		parser/program_*.go \
		jt
