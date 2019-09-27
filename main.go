package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"strings"
)

var dependenciesStr string
var fromCommit string

func init() {
	flag.StringVar(&dependenciesStr, "d", "", "directories (comma-separated) to check for change")
	flag.StringVar(&fromCommit, "c", "", "commit-hash or branch or tag to start checking from")
}

func main() {
	flag.Parse()

	command := flag.Args()
	dependencies := strings.Split(dependenciesStr, ",")
	assertInput(dependencies, fromCommit, command)

	filesChanged, err := detectChangedFiles(fromCommit)
	if err != nil {
		log.Fatalf("could not run git-diff to detect changed files: %v\n", err)
	}

	if !requiresRun(filesChanged, dependencies) {
		log.Println("no change in any dependency, skipping run ...")
		return
	}

	log.Printf("running command: %v\n", strings.Join(command, " "))
	cmd := exec.Command(command[0], command[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
	cmd.Wait()
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
