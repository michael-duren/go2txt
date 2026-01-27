package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	maxFileSize = 100 * 1024 * 1024 // 100MB
	// create outside of directory to avoid any
	// cylclic iterations
	outputFile = "../repo.txt"
)

func main() {
	if !isGitRepo() {
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

	repoName := filepath.Base(getCurrentDir())
	fmt.Fprintf(writer, "Repository: %s\n", repoName)
	fmt.Fprintf(writer, "Generated: %s\n", time.Now().Format(time.RFC1123))
	fmt.Fprintln(writer, "================================")

	// Get list of git-tracked files
	files, err := getGitFiles()
	if err != nil {
		fmt.Printf("Error getting git files: %v\n", err)
		os.Exit(1)
	}

	for _, file := range files {
		if err := processFile(file, writer); err != nil {
			fmt.Printf("Warning: Error processing %s: %v\n", file, err)
		}
	}

	if info, err := os.Stat(outputFile); err == nil {
		fmt.Printf("Done! Output size: %.2f MB\n", float64(info.Size())/(1024*1024))
	}
}

func isGitRepo() bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	return cmd.Run() == nil
}

func getCurrentDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return "unknown"
	}
	return dir
}

func getGitFiles() ([]string, error) {
	cmd := exec.Command("git", "ls-files")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	return lines, nil
}

func isUTF8(filePath string) (bool, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return false, fmt.Errorf("error reading file: %w", err)
	}

	return utf8.Valid(content), nil
}

func processFile(filename string, writer *bufio.Writer) error {
	info, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	if info.Size() > maxFileSize {
		fmt.Fprintf(writer, "\n\n=== %s === [LARGE FILE SKIPPED: %.2f MB]\n",
			filename, float64(info.Size())/(1024*1024))
		fmt.Println("Skipped large file:", filename)
		return nil
	}

	fileIsUTF8, err := isUTF8(filename)
	if !fileIsUTF8 || err != nil {
		fmt.Fprintf(writer, "\n\n=== %s === [BINARY FILE SKIPPED]\n", filename)
		fmt.Println("Skipped binary file:", filename)
		return nil
	}

	fmt.Fprintf(writer, "\n\n=== %s ===\n", filename)

	content, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	_, err = writer.Write(content)
	return err
}
