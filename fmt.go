// Copyright (C) 2015 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package main

import (
	"bytes"
	"flag"
	"os/exec"
)

type gofmt struct {
	set bool
}

func (f *gofmt) setOptions() {
	flag.BoolVar(&f.set, "fmt", false, "Use the fmt tool.")
}

func (f *gofmt) installed() bool {
	return true
}

func (f *gofmt) isSet() bool {
	return f.set
}

func (f *gofmt) run() bool {
	printBackendStatus("fmt")

	cmd := exec.Command("gofmt", "-d", "-e", ".")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()

	if err != nil {
		printResult(err.Error(), errored)
		return false
	}

	// And print out the results.
	results := out.String()
	if results == "" {
		printResult("", ok)
		return true
	}
	printResult(out.String(), failed)
	return false
}
