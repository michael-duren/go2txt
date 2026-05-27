package main

import (
	"flag"
	"fmt"
	"runtime/debug"
)

var version = ""

func init() {
	flag.Usage = printUsage
}

func getVersion() string {
	if version != "" {
		return version
	}
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "unknown"
	}
	if v := info.Main.Version; v != "" && v != "(devel)" {
		return v
	}
	var rev, mod string
	for _, s := range info.Settings {
		switch s.Key {
		case "vcs.revision":
			rev = s.Value
		case "vcs.modified":
			if s.Value == "true" {
				mod = "-dirty"
			}
		}
	}
	if rev != "" {
		if len(rev) > 7 {
			rev = rev[:7]
		}
		return rev + mod
	}
	return "unknown"
}

func printVersion() {
	fmt.Printf("go2txt %s\n", getVersion())
}

func printUsage() {
	fmt.Fprint(flag.CommandLine.Output(), `go2txt - dump a git repository to a single text file

Usage:
  go2txt [flags]

Flags:
  -o, -output       string   output file path (default: ../<repo>.txt)
  -e, -exclude      string   comma-separated glob patterns to exclude (e.g. *.jsx,*.ts)
  -D, -exclude-dirs string   comma-separated glob patterns of directories to exclude (e.g. node_modules,vendor)
  -i, -include      string   comma-separated glob patterns to include exclusively (e.g. *.go,*.md)
  -d, -dir          string   only include files under this directory (e.g. internal)
  -v, -verbose               verbose output
  -V, -version               print version and exit
  -h, -help                  show this help message

Examples:
  go2txt
  go2txt -o repo.txt
  go2txt -i "*.go,*.md"
  go2txt -d internal
  go2txt -e "*.lock,*.sum" -v

Notes:
  Must run inside a git repository. Only git-tracked files are included.
`)
}
