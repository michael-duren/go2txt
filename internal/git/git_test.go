package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"testing"
)

func TestIsRepo(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git binary not available")
	}

	t.Run("true in repo", func(t *testing.T) {
		dir := initRepo(t)
		t.Chdir(dir)
		if !IsRepo() {
			t.Error("IsRepo() = false in initialized repo, want true")
		}
	})

	t.Run("false outside repo", func(t *testing.T) {
		dir := t.TempDir()
		t.Chdir(dir)
		if IsRepo() {
			t.Error("IsRepo() = true outside repo, want false")
		}
	})
}

func TestGetRepoName(t *testing.T) {
	dir := t.TempDir()
	want := filepath.Base(dir)

	t.Chdir(dir)
	if got := GetRepoName(); got != want {
		t.Errorf("GetRepoName() = %q, want %q", got, want)
	}
}

func TestGetFiles(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git binary not available")
	}

	t.Run("tracked files returned", func(t *testing.T) {
		dir := initRepo(t)
		writeAndAdd(t, dir, "a.go", "package a")
		writeAndAdd(t, dir, "sub/b.md", "# b")
		commit(t, dir, "init")

		t.Chdir(dir)
		got, err := GetFiles()
		if err != nil {
			t.Fatalf("GetFiles err = %v", err)
		}
		sort.Strings(got)
		want := []string{"a.go", "sub/b.md"}
		if !equalSlices(got, want) {
			t.Errorf("GetFiles() = %v, want %v", got, want)
		}
	})

	t.Run("empty repo returns nil", func(t *testing.T) {
		dir := initRepo(t)
		t.Chdir(dir)
		got, err := GetFiles()
		if err != nil {
			t.Fatalf("GetFiles err = %v", err)
		}
		if len(got) != 0 {
			t.Errorf("GetFiles() = %v, want empty", got)
		}
	})

	t.Run("untracked files excluded", func(t *testing.T) {
		dir := initRepo(t)
		writeAndAdd(t, dir, "tracked.go", "package x")
		commit(t, dir, "init")
		writeFile(t, dir, "untracked.go", "package x")

		t.Chdir(dir)
		got, err := GetFiles()
		if err != nil {
			t.Fatalf("GetFiles err = %v", err)
		}
		for _, f := range got {
			if f == "untracked.go" {
				t.Errorf("untracked.go should not appear in GetFiles output")
			}
		}
	})
}

func initRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	runGit(t, dir, "init", "-b", "main")
	runGit(t, dir, "config", "user.email", "test@example.com")
	runGit(t, dir, "config", "user.name", "Test")
	runGit(t, dir, "config", "commit.gpgsign", "false")
	return dir
}

func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v failed: %v\n%s", args, err, out)
	}
}

func writeFile(t *testing.T, dir, rel, body string) {
	t.Helper()
	p := filepath.Join(dir, rel)
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func writeAndAdd(t *testing.T, dir, rel, body string) {
	t.Helper()
	writeFile(t, dir, rel, body)
	runGit(t, dir, "add", rel)
}

func commit(t *testing.T, dir, msg string) {
	t.Helper()
	runGit(t, dir, "commit", "-m", msg)
}

func equalSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
