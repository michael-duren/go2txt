package converter

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestSplitPatterns(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want []string
	}{
		{"empty", "", nil},
		{"single", "*.go", []string{"*.go"}},
		{"multiple", "*.go,*.md", []string{"*.go", "*.md"}},
		{"whitespace", " *.go , *.md ", []string{"*.go", "*.md"}},
		{"empty entries", "*.go,,*.md,", []string{"*.go", "*.md"}},
		{"only commas", ",,,", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := splitPatterns(tt.in)
			if len(tt.want) == 0 && len(got) == 0 {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("splitPatterns(%q) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}

func TestMatchAny(t *testing.T) {
	tests := []struct {
		name     string
		patterns []string
		file     string
		want     bool
		wantErr  bool
	}{
		{"empty patterns", nil, "foo.go", false, false},
		{"basename match", []string{"*.go"}, "src/foo.go", true, false},
		{"basename no match", []string{"*.md"}, "src/foo.go", false, false},
		{"full path match", []string{"src/*.go"}, "src/foo.go", true, false},
		{"multiple patterns first hits", []string{"*.go", "*.md"}, "foo.go", true, false},
		{"multiple patterns second hits", []string{"*.md", "*.go"}, "foo.go", true, false},
		{"none match", []string{"*.md", "*.txt"}, "foo.go", false, false},
		{"invalid pattern", []string{"[invalid"}, "foo.go", false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := matchAny(tt.patterns, tt.file)
			if (err != nil) != tt.wantErr {
				t.Fatalf("matchAny err = %v, wantErr = %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("matchAny = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRunner_Run_IncludeOnly(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "a.go", "package a")
	writeFile(t, dir, "b.md", "# md")
	writeFile(t, dir, "c.txt", "txt")

	r := NewRunner("", "*.go,*.md", false)
	var buf bytes.Buffer
	if err := r.Run(filesIn(dir, "a.go", "b.md", "c.txt"), &buf); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "a.go") || !strings.Contains(out, "b.md") {
		t.Errorf("expected a.go and b.md included; got:\n%s", out)
	}
	if strings.Contains(out, "c.txt") {
		t.Errorf("c.txt should be excluded by includeOnly; got:\n%s", out)
	}
}

func TestRunner_Run_Exclude(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "a.go", "package a")
	writeFile(t, dir, "b.lock", "lock")

	r := NewRunner("*.lock", "", false)
	var buf bytes.Buffer
	if err := r.Run(filesIn(dir, "a.go", "b.lock"), &buf); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "a.go") {
		t.Errorf("a.go should appear; got:\n%s", out)
	}
	if strings.Contains(out, "=== "+filepath.Join(dir, "b.lock")+" ===\n") {
		t.Errorf("b.lock should be excluded; got:\n%s", out)
	}
}

func TestRunner_Run_SkipsEmptyAndMissing(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "a.go", "package a")
	missing := filepath.Join(dir, "ghost.go")

	r := NewRunner("", "", false)
	var buf bytes.Buffer
	files := []string{"", filepath.Join(dir, "a.go"), missing}
	if err := r.Run(files, &buf); err != nil {
		t.Fatalf("Run err = %v", err)
	}
	if !strings.Contains(buf.String(), "a.go") {
		t.Errorf("a.go missing from output:\n%s", buf.String())
	}
	if strings.Contains(buf.String(), "ghost.go") {
		t.Errorf("missing file should be skipped silently")
	}
}

func TestRunner_Run_BinaryFileSkipped(t *testing.T) {
	dir := t.TempDir()
	bin := filepath.Join(dir, "bin.dat")
	if err := os.WriteFile(bin, []byte{0xff, 0xfe, 0xfd, 0x00}, 0o644); err != nil {
		t.Fatal(err)
	}

	r := NewRunner("", "", false)
	var buf bytes.Buffer
	if err := r.Run([]string{bin}, &buf); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "[BINARY FILE SKIPPED]") {
		t.Errorf("expected binary skip marker; got:\n%s", buf.String())
	}
}

func TestRunner_Run_LargeFileSkipped(t *testing.T) {
	dir := t.TempDir()
	big := filepath.Join(dir, "big.txt")
	f, err := os.Create(big)
	if err != nil {
		t.Fatal(err)
	}
	if err := f.Truncate(maxFileSize + 1); err != nil {
		t.Fatal(err)
	}
	f.Close()

	r := NewRunner("", "", false)
	var buf bytes.Buffer
	if err := r.Run([]string{big}, &buf); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "[LARGE FILE SKIPPED:") {
		t.Errorf("expected large-file marker; got:\n%s", buf.String())
	}
}

func TestRunner_Run_InvalidPatternErrors(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "a.go", "package a")

	r := NewRunner("[bad", "", false)
	var buf bytes.Buffer
	err := r.Run(filesIn(dir, "a.go"), &buf)
	if err == nil {
		t.Fatal("expected error from invalid pattern, got nil")
	}
}

func TestRunner_Run_WritesFileContent(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "a.go", "package a\n\nfunc Hi() {}\n")

	r := NewRunner("", "", false)
	var buf bytes.Buffer
	if err := r.Run(filesIn(dir, "a.go"), &buf); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "func Hi()") {
		t.Errorf("file body missing; got:\n%s", out)
	}
	if !strings.Contains(out, "=== ") {
		t.Errorf("file header missing; got:\n%s", out)
	}
}

func writeFile(t *testing.T, dir, name, body string) {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func filesIn(dir string, names ...string) []string {
	out := make([]string, len(names))
	for i, n := range names {
		out[i] = filepath.Join(dir, n)
	}
	return out
}
