package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/michael-duren/go2txt/internal/converter"
	"github.com/michael-duren/go2txt/internal/git"
)

func main() {
	var (
		outputFile    string
		excludedFiles string
		includeOnly   string
		verbose       bool
	)

	flag.StringVar(&outputFile, "output", "", "output file path (default: ../<repo>.txt)")
	flag.StringVar(&outputFile, "o", "", "short for -output")
	flag.StringVar(&excludedFiles, "exclude", "", "comma-separated glob patterns to exclude (e.g. *.jsx,*.ts)")
	flag.StringVar(&excludedFiles, "e", "", "short for -exclude")
	flag.StringVar(&includeOnly, "include", "", "comma-separated glob patterns to include exclusively (e.g. *.go,*.md)")
	flag.StringVar(&includeOnly, "i", "", "short for -include")
	flag.BoolVar(&verbose, "verbose", false, "verbose output")
	flag.BoolVar(&verbose, "v", false, "short for -verbose")
	flag.Parse()

	if !git.IsRepo() {
		fmt.Fprintln(os.Stderr, "Error: Not in a git repository")
		os.Exit(1)
	}

	repoName := git.GetRepoName()
	if outputFile == "" {
		outputFile = filepath.Join("..", repoName+".txt")
	}

	out, err := os.Create(outputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output file: %v\n", err)
		os.Exit(1)
	}
	defer out.Close()

	writer := bufio.NewWriter(out)
	defer writer.Flush()

	fmt.Fprintf(writer, "Repository: %s\n", repoName)
	fmt.Fprintf(writer, "Generated: %s\n", time.Now().Format(time.RFC1123))
	fmt.Fprintln(writer, "================================")

	files, err := git.GetFiles()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting git files: %v\n", err)
		os.Exit(1)
	}

	c := converter.NewRunner(excludedFiles, includeOnly, verbose)
	if err := c.Run(files, writer); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if err := writer.Flush(); err != nil {
		fmt.Fprintf(os.Stderr, "Error flushing output: %v\n", err)
		os.Exit(1)
	}

	if info, err := os.Stat(outputFile); err == nil {
		fmt.Printf("Done! Output: %s (%.2f MB)\n", outputFile, float64(info.Size())/(1024*1024))
	}
}
