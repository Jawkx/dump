package main

import (
	"path/filepath"
	"strings"
)

func parseIgnorePatterns(ignoreStr string) []string {
	if ignoreStr == "" {
		return nil
	}

	patterns := strings.Split(ignoreStr, ",")
	for i := range patterns {
		patterns[i] = strings.TrimSpace(patterns[i])
	}
	return patterns
}

func shouldIgnore(path string, isDir bool, ignorePatterns []string, includeHidden bool) bool {
	if !includeHidden {
		baseName := filepath.Base(path)
		if len(baseName) > 0 && baseName[0] == '.' {
			return true
		}
	}

	if len(ignorePatterns) == 0 {
		return false
	}

	path = filepath.ToSlash(path)

	testPath := path
	if isDir && !strings.HasSuffix(testPath, "/") {
		testPath += "/"
	}

	for _, pattern := range ignorePatterns {
		pattern = filepath.ToSlash(pattern)

		if pattern == "." || pattern == "./" {
			if path == "." || path == "./" {
				return true
			}
		}

		isDirectoryPattern := strings.HasSuffix(pattern, "/")
		if isDirectoryPattern && isDir {
			dirPattern := strings.TrimSuffix(pattern, "/")

			if matched, _ := filepath.Match(dirPattern, strings.TrimSuffix(testPath, "/")); matched {
				return true
			}

			if strings.HasPrefix(path, pattern) {
				return true
			}
		}

		matched, err := filepath.Match(pattern, filepath.Base(path))
		if err == nil && matched && !isDir {
			return true
		}

		matched, err = filepath.Match(pattern, path)
		if err == nil && matched {
			return true
		}

		if strings.Contains(pattern, "*") || strings.Contains(pattern, "?") ||
			strings.Contains(pattern, "[") {
			if ok, _ := filepath.Match(pattern, path); ok {
				return true
			}

			if strings.Contains(pattern, "*/") {
				parts := strings.Split(pattern, "*/")
				if len(parts) >= 2 && strings.HasPrefix(path, parts[0]) {
					restPattern := strings.Join(parts[1:], "*/")
					restPath := strings.TrimPrefix(path, parts[0])

					dirs := strings.Split(restPath, "/")
					for i := range dirs {
						subPath := strings.Join(dirs[i:], "/")
						if ok, _ := filepath.Match(restPattern, subPath); ok {
							return true
						}
					}
				}
			}
		}
	}

	return false
}

func containsGlobPattern(path string) bool {
	return strings.ContainsAny(path, "*?[")
}
