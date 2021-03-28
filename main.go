package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/jacobsimpson/jt/ast"
	"github.com/jacobsimpson/jt/debug"

	// For some reason if this import is done without the alias, golang assumes
	// this is imported as `listener`. No idea why.

	flag "github.com/spf13/pflag"
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
	var rules string
	var inputFiles []string
	var version bool
	var verbose int

	flag.CountVarP(&verbose, "verbose", "v",
		"increase output for debugging purposes")
	flag.StringVarP(&rules, "expression", "e", rules,
		"add the script to the commands to be execute")
	flag.StringVarP(&scriptFile, "file", "f", scriptFile,
		"add the contents of script-file to the commands to be execute")
	flag.BoolVar(&version, "version", version,
		"output version information and exit")
	flag.Parse()

	debug.SetLevel(verbose)
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
		rules = string(data)
		inputFiles = args
	} else if rules != "" {
		inputFiles = args
	} else if len(args) > 0 {
		rules = args[0]
		inputFiles = args[1:]
	} else {
		flag.Usage()
		os.Exit(1)
	}

	if err := execute(rules, inputFiles); err != nil {
		os.Exit(1)
	}
}

func execute(rules string, inputFiles []string) error {
	var result error

	ast, err := parse(rules)
	if err != nil {
		return err
	}

	debug.Debug("ast = %s\n", ast)

	if len(inputFiles) == 0 {
		return processReader(ast, os.Stdin)
	} else {
		for _, f := range inputFiles {
			if err := processFile(ast, f); err != nil {
				result = err
			}
		}
	}

	return result
}

func processFile(interpreter *ast.Program, fileName string) error {
	f, err := os.Open(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: can't read %s: No such file or directory", execName(), fileName)
		return err
	}
	defer f.Close()

	return processReader(interpreter, f)
}

func processReader(interpreter *ast.Program, reader io.Reader) error {
	scanner := bufio.NewScanner(reader)

	lineNumber := 0
	for scanner.Scan() {
		if applyRules(interpreter, scanner.Text(), lineNumber) {
		}
		lineNumber++
	}
	return nil
}

func applyRules(interpreter *ast.Program, line string, lineNumber int) bool {
	environment := make(map[string]string)

	environment["%0"] = line
	environment["%#"] = fmt.Sprintf("%d", lineNumber)
	for i, c := range strings.Split(line, " ") {
		environment[fmt.Sprintf("%%%d", i+1)] = c
	}
	debug.Debug("Line %d splits as %s", lineNumber, environment)

	debug.Info("There are %d rules", len(interpreter.Rules))
	for _, rule := range interpreter.Rules {
		debug.Info("    Evaluating: %s\n", rule)
		result, err := rule.Evaluate(environment)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not evaluate %q: %v", rule, err)
		} else if b, ok := result.(bool); ok && b {
			debug.Info("        Executing block\n")
			rule.Execute(environment)
		}
	}
	return true
}
