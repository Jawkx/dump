# dump - A Simple File and Directory Content Dumper

`dump` is a command-line utility that displays the contents of files and directories. It supports glob patterns and selective ignoring of files/directories

## Usage

```bash
dump [options] <file_path1> <file_path2> ...
```

**Arguments:**

*   `<file_path1> <file_path2> ...`: One or more paths to files or directories. Supports glob patterns (wildcards).  If a directory is provided, `dump` recursively displays the contents of all files within it (subject to any ignore patterns). You can use relative or absolute paths.

**Options:**

*   `-version`: Displays the `dump` utility version.
*   `-ignore="<pattern1>,<pattern2>,..."`:  A comma-separated list of patterns to ignore.  *No spaces are allowed within the pattern string itself,* but leading/trailing spaces around each individual pattern will be trimmed.  Patterns can be:
    *   **File patterns:** `*.log`, `config.ini`, `temp*.txt`
    *   **Directory patterns:** `temp/`, `.git/`, `__pycache__/` (**Note:** Trailing `/` is crucial for directory matches).
    *   **Relative paths:** `src/ignored_file.txt`, `data/private/`
    *    **Current directory**: `.`, `./`
    *   Glob patterns are fully supported within ignore patterns (e.g., `-ignore="logs/*/*.log"`).

*   `--help`: Displays help information (this option was added for completeness).

**Key Points about `-ignore`:**

*   **Comma Separation:** Patterns *must* be separated by commas (`,`).
*   **Directory Trailing Slash:** Use a trailing slash (`/`) to ignore directories (e.g., `temp/`). This distinguishes it from a file named "temp".
*   **Globbing:**  Glob patterns (`*`, `?`, `[]`) are supported in both file and directory ignore patterns.
*   **Relative Path Matching:** Ignore patterns are matched against paths *relative to the directory you specify as an argument* to `dump`.  (e.g., `dump -ignore="temp/" src/` ignores `src/temp/`).
*   **Current Directory**: The pattern `.` or `./` will ignore current directory


## Examples

```bash
# Dump a single file:
dump myfile.txt

# Dump multiple files:
dump file1.txt file2.go file3.html

# Dump a directory (recursively):
dump my_project/

# Dump all .go files in a directory using a glob pattern:
dump src/*.go

# Dump all .js files recursively (requires shell support for ** globbing):
dump "**/*.js"

# Ignore .log files:
dump -ignore="*.log" my_project/

# Ignore a directory (temp/) and all .tmp files:
dump -ignore="temp/,*.tmp" my_project/

# Ignore multiple directories and files:
dump -ignore=".git/,logs/,*.bak,config.ini" my_project/

# Ignore a directory and its contents using a glob pattern:
dump -ignore="test_data/*" .

# Ignore a directory and a specific file, relatively:
dump -ignore="src/temp/,src/main_test.go" src/

# Ignore the current working directory:
dump -ignore="." my_project/

# Show the version:
dump -version

# Ignore files in sub-sub directories:
dump -ignore="logs/*/*.log" .  # Ignores logs/date1/x.log, logs/date2/y.log, etc.

# show help
dump --help
```

## Output Format

The output for each file is wrapped in `--- FILE-START ---` and `--- FILE-END ---` markers, followed by the file content:

**Example Files:**

`fileA.md`:

```
fileA-Content
```

`fileB.md`:

```
fileB-Content
```

**Command:**

```bash
dump ./fileA.md ./fileB.md
```

**Output:**

````
---FILE-START---
``` fileA.md
fileA-Content
```
---FILE-END---

--- FILE-START ---
``` fileB.md
fileB-Content
```
---FILE-END---
````

