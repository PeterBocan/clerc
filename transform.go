package main

import (
	"os"
	"path/filepath"
	"strings"
)

// CompileCommandEntry
type CompileCommandEntry struct {
	Directory string   `json:"directory"`
	Arguments []string `json:"arguments,omitempty"`
	File      string   `json:"file"`
	Command   string   `json:"command,omitempty"`
	Output    string   `json:"output,omitempty"`
}

// Transformer converts the Please structure into a flat list of entries
// of CompileCommands.
type Transformer struct {
	RootDirectory string
}

// Transform transforms the call graph with the config to
func (t Transformer) Transform(graph *PleaseGraph, config *PleaseConfig) []CompileCommandEntry {
	output := []CompileCommandEntry{}

	for pkgName, pkg := range graph.Packages {
		// skip the builtin _please package
		if pkgName == "_please" {
			continue
		}
		for targetName, target := range pkg.Targets {
			// filter out anything that doesn't have a #cc label
			if !strings.HasSuffix(targetName, "cc") {
				continue
			}

			// generate outputs for each file
			for _, input := range target.Inputs {
				output = append(output, CompileCommandEntry{
					File:      filepath.Join(t.RootDirectory, input),
					Directory: filepath.Join(t.RootDirectory, filepath.Dir(input)),
					Command:   t.expandCommand(target),
				})
			}
		}
	}

	return output
}

func (t Transformer) expandCommand(target *PleaseTarget) string {
	expandFunction := func(varName string) string {
		if strings.HasPrefix(varName, "TOOLS") {
			_, after, _ := strings.Cut(varName, "TOOLS_")
			toolName := strings.ToLower(after)
			return target.ToolForName(toolName)
		}

		if strings.HasPrefix(varName, "SRCS") {
		  return target.AllSources()
    }

		if varName == "OUTS" {
			return strings.Join(target.Outs, " ")
		}

		return ""
	}

  before, _, _ := strings.Cut(target.Command, "&&")
	command := os.Expand(before, expandFunction)
	return command 
}
