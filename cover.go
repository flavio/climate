// Copyright (C) 2015 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
)

const defaultThreshold = 80.0

type cover struct {
	set       bool
	open      bool
	threshold float64
}

func (c *cover) setOptions() {
	flag.BoolVar(&c.set, "cover", false, "Use the cover tool.")
	flag.BoolVar(&c.open, "open", true,
		"Open the results of the cover tool on the browser.")
	flag.Float64Var(&c.threshold,
		"threshold", defaultThreshold, "The accepted code coverage threshold.")
}

func (c *cover) installed() bool {
	return packageExists("golang.org/x/tools/cover")
}

func (c *cover) isSet() bool {
	// Returns true when it has been either explicitly set, or an option has
	// actually been modified from its default value.
	return c.set || c.open == false || c.threshold != defaultThreshold
}

func (c *cover) run(pack string) bool {
	printBackendStatus("cover")

	cmd := exec.Command("go", "test", "-race", "-coverprofile=c.out", "-covermode=count", pack)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	defer os.Remove("c.out")

	// Usually when this happens is because the package doesn't build or
	// something like that. In this case, it's better to print a custom message
	// than what err.Error() says (which is not very useful here...).
	if err != nil {
		printResult("cover", "`go test` failed!\n", errored)
		return false
	}

	re := regexp.MustCompile("coverage:\\s?(.+)%")
	matches := re.FindStringSubmatch(out.String())
	if len(matches) == 2 {
		thresh, err := strconv.ParseFloat(matches[1], 64)
		if err != nil {
			printResult("cover", err.Error(), errored)
			return false
		}
		if c.open {
			cmd := exec.Command("go", "tool", "cover", "-html=c.out")
			cmd.Run()
		}
		if thresh < c.threshold {
			str := fmt.Sprintf("Coverage required: %.2f%%, got: %.2f%%\n",
				c.threshold, thresh)
			printResult("cover", str, failed)
			return false
		}
	}

	printResult("cover", "", ok)
	return true
}
