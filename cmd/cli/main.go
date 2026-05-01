package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/michael-duren/go2txt/internal/converter"
	"github.com/michael-duren/go2txt/internal/git"
)


func main() {
	outputFile := flag.String("output", "repo.txt", "output file")
	excludedFiles := flag.String("excluded", "", "excluded files, ex: *.jsx,*.ts")
	verbose := flag.Bool("verbose", false, "verbose output")
	gitRepo := flag.Bool("git", true, "is git repository")
	remoteRepo := flag.String("remote-repository", "", "is a remote repo")
	flag.Parse()

	if !git.IsRepo() {
		fmt.Println("Error: Not in a git repository")
		os.Exit(1)
	}

	outputPath := path.Join("..", *outputFile)
	out, err := os.Create(outputPath)
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

	c := converter.NewRunner(*excludedFiles, *verbose, *gitRepo, *remoteRepo)
	err = c.Run(files, writer)

	if err != nil {
		fmt.Printf("an error ocurred: %v\n", err)
	}

	if info, err := os.Stat(outputPath); err == nil {
		fmt.Printf("Done! Output size: %.2f MB\n", float64(info.Size())/(1024*1024))
	}
}
