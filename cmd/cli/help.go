package main

import (
	"flag"
	"fmt"
)

func init() {
	flag.Usage = printUsage
}

func printUsage() {
	fmt.Fprint(flag.CommandLine.Output(), `go2txt - dump a git repository to a single text file

Usage:
  go2txt [flags]

Flags:
  -o, -output    string   output file path (default: ../<repo>.txt)
  -e, -exclude   string   comma-separated glob patterns to exclude (e.g. *.jsx,*.ts)
  -i, -include   string   comma-separated glob patterns to include exclusively (e.g. *.go,*.md)
  -v, -verbose            verbose output
  -h, -help               show this help message

Examples:
  go2txt
  go2txt -o repo.txt
  go2txt -i "*.go,*.md"
  go2txt -e "*.lock,*.sum" -v

Notes:
  Must run inside a git repository. Only git-tracked files are included.
`)
}
