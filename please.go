package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/please-build/gcfg"
)

type PleaseGraph struct {
	Packages map[string]*PleasePackage `json:"packages,omitempty"`
}

func (g *PleaseGraph) String() string {
	var sb strings.Builder

	for key, value := range g.Packages {
		sb.WriteString("Package{name=")
		sb.WriteString(key)
		sb.WriteString(", targets=[")
		sb.WriteString(value.String())
		sb.WriteString("]}")
	}

	return sb.String()
}

type PleasePackage struct {
	Targets map[string]*PleaseTarget `json:"targets,omitempty"`
}

func (p *PleasePackage) String() string {
	var sb strings.Builder
	for key, target := range p.Targets {
		sb.WriteString("Target{name=")
		sb.WriteString(key)
		sb.WriteString(", {")
		sb.WriteString(fmt.Sprintf("%+v", target))
		sb.WriteString("} }")
	}
	return sb.String()
}

// PleaseTarget maps one-to-one to the target definition provided by Please.
type PleaseTarget struct {
	Inputs   []string    `json:"inputs,omitempty"`
	Outs     []string    `json:"outs,omitempty"`
	Deps     []string    `json:"deps,omitempty"`
	Data     []string    `json:"data,omitempty"`
	Labels   []string    `json:"labels,omitempty"`
	Requires []string    `json:"requires,omitempty"`
	Command  string      `json:"command,omitempty"`
	Sources  interface{} `json:"srcs,omitempty"`
	Tools    interface{} `json:"tools,omitempty"`
	Hash     string      `json:"hash,omitempty"`
	Binary   bool        `json:"binary,omitempty"`
}

// ToolForName returns the concatenated list of tools in a single string. 
func (t *PleaseTarget) ToolForName(name string) string {
	tools, ok := t.Tools.(map[string]interface{})
	if !ok {
		return ""
	}

	listOfTools, ok := tools[name].([]interface{})
	if !ok {
		return ""
	}

	strTools := []string{}
	for _, tool := range listOfTools {
		val, ok := tool.(string)
		if ok {
			strTools = append(strTools, val)
		}
	}

	return strings.Join(strTools, " ")
}

// AllSources returns a list of sources concatenated by a whitespace.
func (t *PleaseTarget) AllSources() string {
	tools, ok := t.Sources.(map[string]interface{})
	if !ok {
		return ""
	}

	listOfTools, ok := tools["srcs"].([]interface{})
	if !ok {
		return ""
	}

	strTools := []string{}
	for _, tool := range listOfTools {
		val, ok := tool.(string)
		if ok {
			strTools = append(strTools, val)
		}
	}

	return strings.Join(strTools, " ")
}

// PleaseConfig contains the configuration for C/C++ rules.
type PleaseConfig struct {
	Cpp struct {
		CcTool             string
		CppTool            string
		LdTool             string
		ArTool             string
		LinkWithLdtool     string
		DefaultOptCflags   string
		DefaultDbgCflags   string
		DefaultOptCppFlags string
		DefaultDbgCppFlags string
		DefaltLdFlags      string
		PkgconfigPath      string
		Coverage           bool
		Testmain           string
		Clangmodules       bool
		Dsymtool           string
	}
}

// PleaseRunner runs the Please commands
type PleaseRunner struct{}

// Graph returns a dependency graph for the target
func (PleaseRunner) Graph(target string) (*PleaseGraph, error) {
	stdout := bytes.NewBuffer(nil)
	cmd := exec.Command("plz", "query", "graph", target)
	cmd.Stdout = stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	graph := &PleaseGraph{}
	if err := json.Unmarshal(stdout.Bytes(), graph); err != nil {
		return nil, err
	}
	return graph, nil
}

// Config returns the C/C++ configuration from .plzconfig
func (PleaseRunner) Config() (*PleaseConfig, error) {
	stdout := bytes.NewBuffer(nil)
	cmd := exec.Command("plz", "query", "config")
	cmd.Stdout = stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	config := &PleaseConfig{}
	err := gcfg.ReadInto(config, stdout)
	return config, gcfg.FatalOnly(err)
}
