> This tool will stop being maintained, a successor was created using the experience I gained from using this tool, head over to [ctxcat](https://github.com/Jawkx/ctxcat) to check that out

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
Since the output is std out, you can pipe it into other program to customize what you want to do with it 

Using with [mods](https://github.com/charmbracelet/mods)

``` bash
dump ./* | mods "what does this directory does"
```

Or just pipe it to your clipboard 

``` bash
# For linux
dump ./* | xclip -selection clipboard 

# For macos
dump ./* | pbcopy
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

## Configuration

`dump` can be configured using a TOML file. It searches for configuration files in the following locations:

*   `~/.config/dump.toml`
*   `~/.config/dump/config.toml`

If a configuration file is found, `dump` will load the settings from it. If no configuration file is found, `dump` will use default values.

### Configuration Options

The following options can be configured in the TOML file:

*   `file_start`:  A string that is printed before the content of each file.
*   `file_end`: A string that is printed after the content of each file.
*   `code_start`: A string that is printed before the code content of each file. Useful for specifying things like the language for syntax highlighting in markdown.
*   `code_end`:   A string that is printed after the code content of each file.

All of the options support Go templates with the following data

*   `.FilePath`: The path to the file.
*   `.FileExt`:  The file extension.
*   `.Content`: The content of the file.

### Default Configuration

And best configuration honestly

```toml
file_start = "--- FILE-START: {{.FilePath}} ---"
file_end = "--- FILE-END ---"
code_start = "``` {{.FileExt}}"
code_end = "```"
```
