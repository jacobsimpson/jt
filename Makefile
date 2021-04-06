
jt: *.go datetime/*.go parser/*.go ast/*.go antlrgen/*.go
	go build .

antlrgen/*.go: Program.g4 pparser/grammar.pigeon
	go generate ./...
	java \
		-Xmx500M \
		-cp tools/antlr-4.7.1-complete.jar \
		org.antlr.v4.Tool \
		-Dlanguage=Go \
		-o antlrgen \
		-package antlrgen \
		-no-listener \
		-visitor \
		Program.g4

test: jt
	go test ./...

tests: jt
	go test ./...
	tests/test tests/binary_eq_operator
	tests/test tests/binary_ge_operator
	tests/test tests/binary_gt_operator
	tests/test tests/binary_le_operator
	tests/test tests/binary_lt_operator
	tests/test tests/datetime_eq_operator
	tests/test tests/datetime_ge_operator
	tests/test tests/datetime_gt_operator
	tests/test tests/datetime_le_operator
	tests/test tests/datetime_lt_operator
	tests/test tests/hex_eq_operator
	tests/test tests/hex_ge_operator
	tests/test tests/hex_gt_operator
	tests/test tests/hex_le_operator
	tests/test tests/hex_lt_operator
	tests/test tests/double_eq_operator
	tests/test tests/double_ge_operator
	tests/test tests/double_gt_operator
	tests/test tests/double_le_operator
	tests/test tests/double_lt_operator
	tests/test tests/integer_eq_operator
	tests/test tests/integer_ge_operator
	tests/test tests/integer_gt_operator
	tests/test tests/integer_le_operator
	tests/test tests/integer_lt_operator
	tests/test tests/no_such_function
	tests/test tests/print_column
	tests/test tests/println_column
	tests/test tests/re_column_match_comparison_operator
	tests/test tests/re_line_match_comparison_operator
	tests/test tests/re_line_match_implicit
	tests/test tests/re_line_match_implicit_pipe_delimited
	tests/test tests/stdin
	tests/test tests/substring_column
	tests/test tests/substring_column_empty_end_range
	tests/test tests/substring_column_empty_start_range
	tests/test tests/substring_column_negative_range
	tests/test tests/substring_column_overlapping_range

install: jt
	cp ./jt /usr/bin/jt

clean::
	rm -Rf jt antlrgen
