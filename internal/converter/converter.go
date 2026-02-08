package converter

import (
	"bufio"
	"fmt"
	"os"
	"unicode/utf8"
)

const (
	maxFileSize = 100 * 1024 * 1024 // 100MB
)

type RunConfig struct {
	// Files to ignore also supports glob syntax *.js
	ExcludedFiles []string

	// Output results for each file scanned
	Verbose bool
	// Git repository found
	Git bool
	// If a url is supplied make request to download first
	RemoteRepository bool
}

func NewRunConfig(excludedFiles []string, verbose, git, remoteRepository bool) *RunConfig {
	return &RunConfig{
		ExcludedFiles:    excludedFiles,
		Verbose:          verbose,
		Git:              git,
		RemoteRepository: remoteRepository,
	}
}

// Run process
func Run() {

}

// ProcessFile reads a file and writes its content to the writer if it meets criteria
func ProcessFile(filename string, writer *bufio.Writer) error {
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
