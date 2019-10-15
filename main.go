package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

var dependenciesStr string
var fromCommit string
var verbose bool
var version bool

func init() {
	flag.StringVar(&dependenciesStr, "d", "", "directories (comma-separated) to check for change")
	flag.StringVar(&fromCommit, "c", "", "commit-hash or branch or tag to start checking from")
	flag.BoolVar(&verbose, "v", false, "increase verbosity for debugging")
	flag.BoolVar(&version, "version", false, "show version")
}

func main() {
	flag.Parse()

	if version {
		fmt.Println("v0.10.1")
		return
	}

	command := flag.Args()
	dependencies := strings.Split(dependenciesStr, ",")
	assertInput(dependencies, fromCommit, command)

	hash, ok := dereferenceCommit(fromCommit)
	if verbose {
		log.Printf("dereferenced %q to %q\n", fromCommit, hash)
	}
	if ok {
		filesChanged, err := detectChangedFiles(hash)
		if err != nil {
			log.Fatalf("could not run git-diff to detect changed files: %v\n", err)
		}

		if !requiresRun(filesChanged, dependencies) {
			if verbose {
				log.Println("no change in any dependency, skipping run ...")
			}
			return
		}
	} else if verbose {
		log.Printf("could not dereference %q, skipping check ...\n", fromCommit)
	}

	if verbose {
		log.Printf("running command: %v\n", strings.Join(command, " "))
	}
	cmd := exec.Command(command[0], command[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				os.Exit(status.ExitStatus())
				return
			}
		} else {
			log.Fatalf("run error: %v", err)
		}
	}
}

func dereferenceCommit(fromCommit string) (string, bool) {
	out, err := exec.Command("git", "rev-parse", "--verify", fromCommit).Output()
	if err != nil {
		return "", false
	}
	lines := strings.Split(string(out), "\n")
	return lines[0], true
}

func assertInput(deps []string, prevCommit string, command []string) {
	if deps == nil || len(deps) == 0 {
		log.Fatal("E: at least 1 dependency (-d) required!")
	}
	if len(prevCommit) == 0 {
		log.Fatal("E: -c required!")
	}
	if command == nil || len(command) == 0 {
		log.Fatal("E: command required!")
	}
}

func detectChangedFiles(fromCommit string) ([]string, error) {
	out, err := exec.Command("git", "diff", "--name-only", fromCommit).Output()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(out), "\n")
	return lines, nil
}

func requiresRun(changedFiles []string, dependencies []string) bool {
	for _, changedFile := range changedFiles {
		for _, dependency := range dependencies {
			if strings.HasPrefix(changedFile, dependency) {
				return true
			}
		}
	}
	return false
}
