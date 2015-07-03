// Copyright (C) 2015 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package main

import (
	"bytes"
	"flag"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type errcheck struct {
	set bool
}

func (ec *errcheck) setOptions() {
	flag.BoolVar(&ec.set, "errcheck", false, "Use the errcheck tool.")
}

func (ec *errcheck) installed() bool {
	return packageExists("github.com/kisielk/errcheck")
}

func (ec *errcheck) isSet() bool {
	return ec.set
}

func (ec *errcheck) run() bool {
	printBackendStatus("errcheck")

	path := getPath()
	if path == "" {
		printResult("Could not locate package", failed)
		return false
	}

	// The errcheck has two possible errors:
	//   - exit(1): there are some checks that are not passing.
	//   - exit(2): the package could not be parsed.
	cmd := exec.Command("errcheck", path)
	var out, errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	err := cmd.Run()

	if err != nil {
		if str := out.String(); str != "" {
			printResult(str, failed)
		} else {
			printResult(errOut.String(), failed)
		}
		return false
	}
	printResult("", ok)
	return true
}

func getPath() string {
	if path := getPathFrom("GOPATH"); path != "" {
		return path
	}
	return getPathFrom("GOROOT")
}

func getPathFrom(env string) string {
	base := os.Getenv(env)
	if base == "" {
		return ""
	}

	// We add the separator so it gets matched on the `strings.Replace` call
	// and the final path doesn't start with the separator.
	path := filepath.Join(base, "src") + string(filepath.Separator)
	abs, err := filepath.Abs(".")
	if err != nil || !strings.HasPrefix(abs, path) {
		return ""
	}
	return strings.Replace(abs, path, "", 1)
}
