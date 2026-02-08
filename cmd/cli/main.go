package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/michael-duren/go2txt/internal/converter"
	"github.com/michael-duren/go2txt/internal/git"
)

const (
	outputFile = "../repo.txt"
)

func main() {
	var output string
	flag.StringVar(&output, "output", "", "output file")

	if !git.IsRepo() {
		fmt.Println("Error: Not in a git repository")
		os.Exit(1)
	}

	out, err := os.Create(outputFile)
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		os.Exit(1)
	}
	defer out.Close()

	writer := bufio.NewWriter(out)
	defer writer.Flush()

	repoName := git.GetRepoName()
	fmt.Fprintf(writer, "Repository: %s\n", repoName)
	fmt.Fprintf(writer, "Generated: %s\n", time.Now().Format(time.RFC1123))
	fmt.Fprintln(writer, "================================")

	// Get list of git-tracked files
	files, err := git.GetFiles()
	if err != nil {
		fmt.Printf("Error getting git files: %v\n", err)
		os.Exit(1)
	}

	for _, file := range files {
		if err := converter.ProcessFile(file, writer); err != nil {
			fmt.Printf("Warning: Error processing %s: %v\n", file, err)
		}
	}

	if info, err := os.Stat(outputFile); err == nil {
		fmt.Printf("Done! Output size: %.2f MB\n", float64(info.Size())/(1024*1024))
	}
}

func parseFlags() {

}
