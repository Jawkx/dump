# dump - File/Directory Content Dumper

`dump` displays file and directory content to standard output, ideal for LLM input.

## Installation

1.  Download the binary.
2.  Use it.

### Quick Install

#### Mac

```bash
curl -L "https://github.com/Jawkx/dump/releases/download/<VERSION>/dump-darwin-amd64" -o dump && chmod +x dump && mkdir -p "$HOME/.local/bin" 2>/dev/null && mv dump "$HOME/.local/bin/dump" && if ! echo "$PATH" | grep -q "$HOME/.local/bin"; then echo "Please add $HOME/.local/bin to your PATH."; fi
```

#### Linux

```bash
curl -L "https://github.com/Jawkx/dump/releases/download/<VERSION>/dump-linux-amd64" -o dump && chmod +x dump && mkdir -p "$HOME/.local/bin" 2>/dev/null && mv dump "$HOME/.local/bin/dump" && if ! echo "$PATH" | grep -q "$HOME/.local/bin"; then echo "Please add $HOME/.local/bin to your PATH."; fi
```

## Usage

```bash
dump [options] <file_path1> <file_path2> ...
```

## Arguments

*   `<file_path>`: Path to file/directory. Supports glob patterns. Multiple paths allowed.

## Options

*   `-version`, `-v`: Show version.
*   `-ignore="<pattern1>,<pattern2>,..."`, `-i="<pattern1>,<pattern2>,..."`: Ignore patterns (mutually exclusive). Use commas to separate. Directory patterns end with `/`. Supports wildcards (`*`, `?`, `[]`).
*   `-hidden`: Include hidden files/directories (starting with `.`).
*   `-help`, `-h`: Show help.

## Ignore Patterns

*   `*`: Matches any sequence of characters.
*   `?`: Matches any single character.
*   `[]`: Matches any character within brackets (e.g., `[abc]`, `[a-z]`).
*   End directory patterns with `/` (e.g., `temp/`).
*   Patterns match both base name and full path.
*   First matching pattern takes precedence.

## Examples

```bash
dump file.txt                      # Dump single file
dump src/                         # Dump directory recursively
dump *.go                          # Dump all Go files
dump -i="*.log,temp/" project/     # Ignore .log files and temp directory
dump -hidden project/             # Include hidden files
dump file1.txt file2.txt src/main.go # Dump multiple files
dump -i="vendor/" myproject/       # Ignore vendor directory
dump -i="drafts/" *.txt             # Exclude drafts subdirectory from .txt files
dump -i="data/,*.bak" /path/to/files  # Ignore data dir and .bak files
dump -i="tmp*" .                   # Ignore tmp files (tmp1.txt, tmp_data.csv)
dump -i="build/*/" project       # Ignore files inside a build directory
dump -version                      # Show version
dump -help                         # Show help
```

## Output Format

`````
```
---FILE-START---
```
< File content>
```
---FILE-END---
```
`````
