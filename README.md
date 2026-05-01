# go2txt

A lightweight command-line tool that converts your Git repository into a single text file, making it easy to share codebases with AI assistants, create documentation snapshots, or analyze your entire project at once.

## Features

- ❌**NO external depdencies** - Just vanilla go
- 📦 **Simple & Fast** - One command to export your entire repository
- 🎯 **Git-Aware** - Only processes tracked files (respects `.gitignore`)
- 🛡️ **Smart Filtering** - Automatically skips binary files and large files (>100MB)
- 🔍 **UTF-8 Validated** - Ensures only text files are included
- 🎛️ **Glob Filters** - Include or exclude files by pattern
- 📊 **Progress Feedback** - Shows file processing and final output size

## Installation

### Install with Go

```bash
go install github.com/michael-duren/go2txt/cmd/cli@latest
```

Make sure `$GOPATH/bin` (usually `~/go/bin`) is in your `PATH`.

### Download Pre-built Binary

Download the latest release for your platform from the [releases page](https://github.com/michael-duren/go2txt/releases).

## Usage

Navigate to any Git repository and run:

```bash
cd /path/to/your/repo
go2txt
```

This creates `../<repo>.txt` (one directory above the repo, named after it) containing:

- Repository name and generation timestamp
- Contents of all Git-tracked text files
- Clear section markers for each file

### Options

| Flag                  | Default          | Description                                                       |
|-----------------------|------------------|-------------------------------------------------------------------|
| `-output`, `-o`       | `../<repo>.txt`  | Output file path                                                  |
| `-include`, `-i`      | _(none)_         | Comma-separated glob patterns to include exclusively (e.g. `*.go,*.md`) |
| `-exclude`, `-e`      | _(none)_         | Comma-separated glob patterns to exclude (e.g. `*.jsx,*.ts`)      |
| `-verbose`, `-v`      | `false`          | Print each file as it is processed, excluded, or skipped          |

Patterns use Go's [`filepath.Match`](https://pkg.go.dev/path/filepath#Match) syntax and are tested against both the file basename and full path. When `-include` is set, only matching files are processed; `-exclude` is then applied on top.

### Examples

```bash
# Write to a custom location
go2txt -o /tmp/myrepo.txt

# Only include Go and Markdown files
go2txt -i '*.go,*.md'

# Exclude generated and lock files
go2txt -e '*.pb.go,*.lock,package-lock.json'

# Verbose mode
go2txt -v
```

### Example Output

```
Repository: my-awesome-project
Generated: Tue, 27 Jan 2026 09:38:01 CST
================================


=== main.go ===
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}


=== README.md ===
# My Awesome Project
...
```

## Use Cases

- **AI Context** - Share your entire codebase with ChatGPT, Claude, or other AI assistants
- **Code Reviews** - Create snapshots for comprehensive reviews
- **Documentation** - Generate project overviews and documentation
- **Analysis** - Feed your codebase to analysis tools that expect single-file input
- **Archiving** - Create text-based backups of your code

## How It Works

1. Verifies you're in a Git repository
2. Lists all Git-tracked files using `git ls-files`
3. Applies `-include` / `-exclude` glob filters
4. Filters out:
   - Binary files (non-UTF-8 content)
   - Large files (>100MB)
5. Concatenates remaining text files with clear section markers
6. Writes to the path given by `-output` (default `../<repo>.txt`, one level above the repo)

## Limitations

- Requires Git to be installed and repository initialized
- Binary files and files >100MB are automatically skipped
- No recursive submodule support (yet)

## Contributing

Contributions are welcome! Areas for improvement:

- [ ] Add test coverage
- [ ] Support for submodules
- [ ] Custom size limits via CLI flags
- [ ] Stronger binary detection (null-byte sniff)

### Running from Source

```bash
git clone https://github.com/michael-duren/go2txt.git
cd go2txt
go run ./cmd/cli
```

## Credits

Inspired by [git2txt](https://github.com/addyosmani/git2txt) by Addy Osmani, reimplemented for local repository support with enhanced filtering.

## License

MIT License - See LICENSE file for details

## Author

**Michael Duren** - [GitHub](https://github.com/michael-duren)

---

<div align="center">

**⭐ If you find this tool useful, please consider giving it a star!**

</div>
