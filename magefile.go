// +build mage

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/target"
)

// Default target to run when none is specified
// If not set, running mage will list available targets
var Default = Build

func Build() error {
	mg.Deps(BuildGrammar)
	fmt.Println("Building...")
	return exec.Command("go", "build", "-o", "jt", ".").Run()
}

func BuildGrammar() error {
	newer, err := target.Path("antlrgen/program_lexer.go", "Program.g4")
	if !newer && err == nil {
		return nil
	}
	return exec.Command("java",
		"-Xmx500M",
		"-cp", "tools/antlr-4.7.1-complete.jar",
		"org.antlr.v4.Tool",
		"-Dlanguage=Go",
		"-o", "antlrgen",
		"-package", "antlrgen",
		"-no-listener",
		"-visitor",
		"Program.g4").Run()
}

func Test() error {
	mg.Deps(Build)
	fmt.Println("Testing...")
	fmt.Printf("%-60s", "go test ./...")
	if err := exec.Command("go", "test", "./...").Run(); err != nil {
		color.Red("FAIL")
		return err
	}
	color.Green("SUCCESS")

	fileList := []string{}
	err := filepath.Walk("tests", func(p string, f os.FileInfo, err error) error {
		if path.Base(p) == "test" && p != "tests/test" {
			fileList = append(fileList, p)
		}
		return nil
	})
	if err != nil {
		return err
	}

	for _, file := range fileList {
		file := path.Dir(file)
		fmt.Printf("%-60s", file)
		if err := exec.Command("tests/test", file).Run(); err != nil {
			color.Red("FAIL")
			return fmt.Errorf("    %s failed", file)
		}
		color.Green("SUCCESS")
	}
	return nil
}

// A custom install step if you need your bin someplace other than go/bin
func Install() error {
	mg.Deps(Build)
	fmt.Println("Installing...")
	return os.Rename("./jt", "/usr/bin/jt")
}

// Clean up after yourself
func Clean() {
	fmt.Println("Cleaning...")
	os.RemoveAll("jt")
	os.RemoveAll("antlrgen")
}
