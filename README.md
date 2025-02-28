# dump - A Simple File and Directory Content Dumper

`dump` is a command-line utility that displays the content of files and directories to standard output. perfect for quickly copy paste context into llm 

## Installation

1. Download the binary
2. Use it

### _For the lazy ones_

#### Mac
``` bash
curl -L "https://github.com/Jawkx/cmtbot/releases/download/<VERSION>/dump-darwin-amd64" -o dump && chmod +x dump && mkdir -p "$HOME/.local/bin" 2>/dev/null && mv dump "$HOME/.local/bin/dump" && if ! echo "$PATH" | grep -q "$HOME/.local/bin"; then echo "Please add $HOME/.local/bin to your PATH."; fi
```

#### Linux
``` bash
curl -L "https://github.com/Jawkx/cmtbot/releases/download/<VERSION>/dump-linux-amd64" -o dump && chmod +x dump && mkdir -p "$HOME/.local/bin" 2>/dev/null && mv dump "$HOME/.local/bin/dump" && if ! echo "$PATH" | grep -q "$HOME/.local/bin"; then echo "Please add $HOME/.local/bin to your PATH."; fi
```
_Not tested on anything other than Mac_

## Usage

```bash
dump [options] <file_path1> <file_path2> ...
```

## Arguments

*   `<file_path>`:  Path to a file or directory. Supports glob patterns.  You can specify multiple paths.

## Options

*   `-version` or `-v`: Display the version of the `dump` utility.
*   `-ignore="<pattern1>,<pattern2>,..."`:  Comma-separated list of patterns to ignore.  Patterns can include wildcards (`*`, `?`, `[]`).  Directory patterns should end with a forward slash (`/`). This flag and `-i` are mutually exclusive.
*   `-i="<pattern1>,<pattern2>,..."`:  Short form of `-ignore`. Comma-separated list of patterns to ignore. This flag and `-ignore` are mutually exclusive.
*   `-hidden`: Include hidden files and directories (those starting with a dot `.` ). By default, hidden files and directories are skipped.
*   `-help` or `-h`: Show the help information (this documentation).

## Ignore Patterns

Ignore patterns allow you to exclude specific files and directories from being dumped.  Here's how they work:

*   **Wildcards:**
    *   `*`: Matches any sequence of characters within a file or directory name.
    *   `?`: Matches any single character within a file or directory name.
    *   `[]`: Matches any character within the brackets.  For example, `[abc]` matches `a`, `b`, or `c`.  Ranges are also supported: `[a-z]` matches any lowercase letter.
*   **Directory Matching:** To match a directory, end the pattern with a forward slash (`/`). For instance, `temp/` will ignore the entire `temp` directory and all its contents.
*   **File Matching:** Patterns without a trailing slash will match files.  For example, `*.log` will ignore all files ending with `.log`.
*   **Path Matching:** The patterns are matched against both the base name of the file/directory and the full path.
*   **Precedence:** If multiple ignore patterns match a file or directory, the *first* matching pattern in the comma separated list takes precedence.
*   **Special Cases:**
    *  `.` or `./`: ignore the current directory itself.

## Examples

1.  **Dump a single file:**

    ```bash
    dump file.txt
    ```

2.  **Dump all files in a directory recursively:**

    ```bash
    dump src/
    ```

3.  **Dump all Go files in the current directory:**

    ```bash
    dump *.go
    ```

4.  **Ignore .log files and the temp directory:**

    ```bash
    dump -ignore="*.log,temp/" project/
    ```
    or equivalently:
    ```bash
    dump -i="*.log,temp/" project/
    ```

5.  **Include hidden files like .gitignore:**

    ```bash
    dump -hidden project/
    ```

6.  **Dump multiple files:**

    ```bash
    dump file1.txt file2.txt src/main.go
    ```

7.  **Dump all files in a directory, ignoring a specific subdirectory "vendor":**

    ```bash
    dump -ignore="vendor/" myproject/
    ```

8.  **Dump all .txt files, excluding those in a subdirectory named "drafts":**
    ```bash
    dump -i="drafts/" *.txt
    ```

9.  **Dump files but ignore the directory named data and all files ending in .bak**
    ```bash
    dump -i="data/,*.bak" /path/to/files
    ```

10. **Ignore files starting with 'tmp' (e.g., tmp1.txt, tmp_data.csv):**

    ```bash
    dump -i="tmp*" .
    ```

11. **Ignore any file or sub-directory inside a `build` directory, regardless of depth:**

    ```
    dump -i="build/*/" project
    ```

12. **Show the version:**

    ```bash
    dump -version
    ```

13. **Show Help:**

    ```bash
    dump -help
    ```

## Output Format

The output for each file is enclosed within `---FILE-START---` and `---FILE-END---` markers. The filename is displayed in a code block header.

## Error Handling

*   If no file or directory paths are specified, an error message is displayed, and the help information is shown.
*   If an invalid glob pattern is used, an error message is printed for that pattern.
*   If a file or directory cannot be accessed (e.g., due to permissions), an error message is printed.
*   If you use both `-ignore` and `-i`, the program will report an error and display usage.
*   If a glob pattern doesn't match any files, a warning is displayed.

## Exit Codes
* 0: Successful execution
* 1: Error (no path given, conflicting flags).
