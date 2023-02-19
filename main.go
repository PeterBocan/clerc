package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func detectRepoRootDirectory(wd string) (string, error) {
	for ; wd != "/"; wd = filepath.Dir(wd) {
		stat, _ := os.Stat(filepath.Join(wd, ".plzconfig"))
		if stat.Mode().IsRegular() {
			return wd, nil
		}
	}
	return "", errors.New("the root directory not found")
}

func main() {

  workingDirectory, err := os.Getwd()
	if err != nil {
		fmt.Printf("could not detect the working directory: %s\n", workingDirectory)
		os.Exit(1)
	}

	if realWorkignDirectoryPath, err := os.Readlink(workingDirectory); err == nil {
		workingDirectory = realWorkignDirectoryPath
	}

	rootDirectory, err := detectRepoRootDirectory(workingDirectory)
	if err != nil {
		fmt.Printf("could not detect .plzconfig: %s\n", err)
		os.Exit(1)
	}

	err = os.Chdir(rootDirectory)
	if err != nil {
		fmt.Printf("could not change directory: %s\n", err)
		os.Exit(1)
	}

	plzExec := PleaseRunner{}
	config, err := plzExec.Config()
	if err != nil {
		fmt.Printf("could not read the Please config: %s\n", err)
		os.Exit(1)
	}

	// target := os.Args[1]
	graph, err := plzExec.Graph("//testdata/simple_binary:simple_binary")
	if err != nil {
		fmt.Printf("could not read the graph for: %s", err)
		os.Exit(1)
	}

	transformer := Transformer{
		RootDirectory: rootDirectory,
	}
	output := transformer.Transform(graph, config)
	out, _ := json.Marshal(output)
	fmt.Printf("%s", out)
}
