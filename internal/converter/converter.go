package converter

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

const (
	maxFileSize = 100 * 1024 * 1024 // 100MB
)

type Runner struct {
	excludedFiles []string
	excludedDirs  []string
	includeOnly   []string
	verbose       bool
}

func NewRunner(excludedFiles, excludedDirs, includeOnly string, verbose bool) *Runner {
	return &Runner{
		excludedFiles: splitPatterns(excludedFiles),
		excludedDirs:  splitPatterns(excludedDirs),
		includeOnly:   splitPatterns(includeOnly),
		verbose:       verbose,
	}
}

func splitPatterns(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := parts[:0]
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}

func (r *Runner) Run(files []string, writer io.Writer) error {
	for _, file := range files {
		if file == "" {
			continue
		}

		included, err := r.isIncludedFile(file)
		if err != nil {
			return err
		}
		if !included {
			continue
		}

		excludedDir, err := r.inExcludedDirs(file)
		if err != nil {
			return err
		}
		if excludedDir {
			if r.verbose {
				fmt.Println("Excluded dir:", file)
			}
			continue
		}

		excluded, err := r.inExcludedFiles(file)
		if err != nil {
			return err
		}
		if excluded {
			if r.verbose {
				fmt.Println("Excluded:", file)
			}
			continue
		}

		if err := r.processFile(file, writer); err != nil {
			return err
		}
	}
	return nil
}

func (r *Runner) isIncludedFile(file string) (bool, error) {
	if len(r.includeOnly) == 0 {
		return true, nil
	}
	return matchAny(r.includeOnly, file)
}

func (r *Runner) inExcludedFiles(file string) (bool, error) {
	if len(r.excludedFiles) == 0 {
		return false, nil
	}
	return matchAny(r.excludedFiles, file)
}

// inExcludedDirs returns true if any directory component of file matches an
// excluded-dir pattern, or if the file's parent path matches a pattern.
func (r *Runner) inExcludedDirs(file string) (bool, error) {
	if len(r.excludedDirs) == 0 {
		return false, nil
	}
	dir := filepath.Dir(file)
	if dir == "." || dir == string(filepath.Separator) {
		return false, nil
	}
	for _, p := range r.excludedDirs {
		ok, err := filepath.Match(p, dir)
		if err != nil {
			return false, fmt.Errorf("invalid pattern %q: %w", p, err)
		}
		if ok {
			return true, nil
		}
		for comp := range strings.SplitSeq(dir, string(filepath.Separator)) {
			if comp == "" {
				continue
			}
			ok, err := filepath.Match(p, comp)
			if err != nil {
				return false, fmt.Errorf("invalid pattern %q: %w", p, err)
			}
			if ok {
				return true, nil
			}
		}
	}
	return false, nil
}

// matchAny returns true if file matches any glob pattern.
// Patterns match against either the basename or the full path.
func matchAny(patterns []string, file string) (bool, error) {
	base := filepath.Base(file)
	for _, p := range patterns {
		ok, err := filepath.Match(p, base)
		if err != nil {
			return false, fmt.Errorf("invalid pattern %q: %w", p, err)
		}
		if ok {
			return true, nil
		}
		ok, err = filepath.Match(p, file)
		if err != nil {
			return false, fmt.Errorf("invalid pattern %q: %w", p, err)
		}
		if ok {
			return true, nil
		}
	}
	return false, nil
}

func (r *Runner) processFile(filename string, writer io.Writer) error {
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
		if r.verbose {
			fmt.Println("Skipped large file:", filename)
		}
		return nil
	}

	content, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	if !utf8.Valid(content) {
		fmt.Fprintf(writer, "\n\n=== %s === [BINARY FILE SKIPPED]\n", filename)
		if r.verbose {
			fmt.Println("Skipped binary file:", filename)
		}
		return nil
	}

	fmt.Fprintf(writer, "\n\n=== %s ===\n", filename)
	if _, err := writer.Write(content); err != nil {
		return err
	}
	if r.verbose {
		fmt.Println("Processed:", filename)
	}
	return nil
}
