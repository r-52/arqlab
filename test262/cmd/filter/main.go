package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"es6-interpreter/test262"
)

func main() {
	root := flag.String("root", "", "path to the cloned test262 repository")
	flag.Parse()

	if *root != "" {
		if _, err := test262.NewRunner(*root, ""); err != nil {
			fmt.Fprintf(os.Stderr, "warning: unable to validate test262 root: %v\n", err)
		}
	}

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 1024), 1024*1024)

	var cases []test262.TestCase
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		cases = append(cases, test262.TestCase{Path: line})
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to read input: %v\n", err)
		os.Exit(1)
	}

	filtered := test262.FilterAsync(cases)
	for _, tc := range filtered {
		fmt.Println(tc.Path)
	}
}
