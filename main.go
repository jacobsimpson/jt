package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/jacobsimpson/jt/parser"
	flag "github.com/ogier/pflag"
)

const VERSION = "0.0.1"

func execName() string {
	return filepath.Base(os.Args[0])
}

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTION]... {script-only-if-no-other-script} [input-file]...\n", execName())
		fmt.Fprintln(os.Stderr)
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, `
If no -e, --expression, -f, or --file option is given, then the first
non-option argument is taken as the %s script to interpret.  All remaining
arguments are names of input files; if no input files are specified, then the
standard input is read.
`, execName())
	}
}

func main() {
	var scriptFile string
	var expressions string
	var inputFiles []string
	var version bool

	flag.StringVarP(&expressions, "expression", "e", expressions,
		"add the script to the commands to be execute")
	flag.StringVarP(&scriptFile, "file", "f", scriptFile,
		"add the contents of script-file to the commands to be execute")
	flag.BoolVar(&version, "version", version,
		"output version information and exit")
	flag.Parse()

	if version {
		fmt.Printf("%s %s\n", execName(), VERSION)
		os.Exit(0)
	}

	args := flag.Args()
	if scriptFile != "" {
		data, err := ioutil.ReadFile(scriptFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: couldn't open file %s: No such file or directory", execName(), scriptFile)
		}
		expressions = string(data)
		inputFiles = args
	} else if expressions != "" {
		inputFiles = args
	} else if len(args) > 0 {
		expressions = args[0]
		inputFiles = args[1:]
	} else {
		flag.Usage()
		os.Exit(1)
	}

	if err := execute(expressions, inputFiles); err != nil {
		os.Exit(1)
	}
}

func execute(expressions string, inputFiles []string) error {
	var result error

	input := antlr.NewInputStream(expressions)
	lexer := parser.NewProgramLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	parser := parser.NewProgramParser(stream)
	parser.BuildParseTrees = true
	tree := parser.Program()
	interpreter := NewInterpreterListener()
	antlr.ParseTreeWalkerDefault.Walk(interpreter, tree)

	for _, f := range inputFiles {
		if err := processFile(interpreter, f); err != nil {
			result = err
		}
	}

	return result
}

func processFile(interpreter *InterpreterListener, fileName string) error {
	f, err := os.Open(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: can't read %s: No such file or directory", execName(), fileName)
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	lineNumber := 0
	for scanner.Scan() {
		if anyExpressionMatches(interpreter, scanner.Text(), lineNumber) {
			fmt.Println(scanner.Text())
		}
		lineNumber++
	}
	return nil
}

func anyExpressionMatches(interpreter *InterpreterListener, line string, lineNumber int) bool {

	for _, expression := range interpreter.Expressions {
		if expression.Matcher.Match(line, lineNumber) {
			return true
		}
	}
	return false
}
