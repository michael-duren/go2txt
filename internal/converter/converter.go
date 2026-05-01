package converter

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"unicode/utf8"
)

const (
	maxFileSize = 100 * 1024 * 1024 // 100MB
)

type Runner struct {
	// Files to ignore also supports glob syntax *.js
	excludedFiles []string

	// A file type to include only like *.md
	includeOnly string

	// Output results for each file scanned
	verbose bool
	// git repository found
	git bool
	// If a url is supplied make request to download first
	remoteRepository string
}

func NewRunner(excludedFiles string, verbose, git bool, remoteRepository string) *Runner {
	var excluded []string
	if excludedFiles != "" {
		excluded = strings.Split(excludedFiles, ",")
	}

	return &Runner{
		excludedFiles:    excluded,
		verbose:          verbose,
		git:              git,
		remoteRepository: remoteRepository,
	}
}

// Run process
func (r Runner) Run(files []string, writer *bufio.Writer) error {
	for _, file := range files {
		inExcluded, err := r.inExcludedFiles(file)
		if err != nil {
			return err
		}
		if inExcluded {
			continue
		}

		included, err := r.isIncludedFile(file)
		if err != nil {
			return err
		}
		if !included {
			continue
		}

		if err := r.processFile(file, writer); err != nil {
			return err
		}
	}
	return nil
}

func (r Runner) isIncludedFile(file string) (bool, error) {
	if r.includeOnly == "" {
		return true, nil
	}

	rg, err := regexp.Compile(r.includeOnly)
	if err != nil {
		return false, fmt.Errorf("can't use include with: %s", r.includeOnly)
	}
	return rg.Match([]byte(file)), nil
}

func (r Runner) inExcludedFiles(file string) (bool, error) {
	if r.includeOnly != "" {
		return false, nil
	}

	for _, ex := range r.excludedFiles {
		rg, err := regexp.Compile(ex)
		if err != nil {
			return false, err
		}
		if rg.Match([]byte(file)) {
			return true, nil
		}
	}
	return false, nil
}

// ProcessFile reads a file and writes its content to the writer if it meets criteria
func (r Runner) processFile(filename string, writer *bufio.Writer) error {
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

func isUTF8(filePath string) (bool, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return false, fmt.Errorf("error reading file: %w", err)
	}

	return utf8.Valid(content), nil
}
