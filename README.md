# go2txt

A lightweight command-line tool that converts your Git repository into a single text file, making it easy to share codebases with AI assistants, create documentation snapshots, or analyze your entire project at once.

## Features

- 📦 **Simple & Fast** - One command to export your entire repository
- 🎯 **Git-Aware** - Only processes tracked files (respects `.gitignore`)
- 🛡️ **Smart Filtering** - Automatically skips binary files and large files (>100MB)
- 🔍 **UTF-8 Validated** - Ensures only text files are included
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

This will create a `repo.txt` file in the parent directory containing:

- Repository name and generation timestamp
- Contents of all Git-tracked text files
- Clear section markers for each file

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
3. Filters out:
   - Binary files (non-UTF-8 content)
   - Large files (>100MB)
4. Concatenates all text files with clear section markers
5. Outputs to `repo.txt` in the parent directory

## Configuration

Currently, the tool uses sensible defaults:

- **Max file size**: 100MB
- **Output location**: `../repo.txt` (parent directory)
- **File validation**: UTF-8 encoding check

## Limitations

- Requires Git to be installed and repository initialized
- Output file is created in the parent directory (to avoid Git tracking)
- Binary files and files >100MB are automatically skipped
- No recursive submodule support (yet)

## Contributing

Contributions are welcome! Areas for improvement:

- [ ] Add test coverage
- [ ] Make output location configurable
- [ ] Add file pattern include/exclude flags
- [ ] Support for submodules
- [ ] Custom size limits via CLI flags

### Running from Source

```bash
git clone https://github.com/michael-duren/go2txt.git
cd go2txt
go run cmd/cli/main.go
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
