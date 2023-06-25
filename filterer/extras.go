package filterer

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

////

var ansiEscape = regexp.MustCompile(`\x1B(?:[@-Z\\-_]|\[[0-?]*[ -/]*[@-~])`)

// Strip strips a string of any of it's ansi sequences.
func Strip(text string) string {
	return ansiEscape.ReplaceAllString(text, "")
}

////

func Read() (string, error) {
	if IsEmpty() {
		return "", fmt.Errorf("stdin is empty")
	}

	reader := bufio.NewReader(os.Stdin)
	var b strings.Builder

	for {
		r, _, err := reader.ReadRune()
		if err != nil && err == io.EOF {
			break
		}
		_, err = b.WriteRune(r)
		if err != nil {
			return "", fmt.Errorf("failed to write rune: %w", err)
		}
	}

	return b.String(), nil
}

// IsEmpty returns whether stdin is empty.
func IsEmpty() bool {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return true
	}

	if stat.Mode()&os.ModeNamedPipe == 0 && stat.Size() == 0 {
		return true
	}

	return false
}

////

// List returns a list of all files in the current directory.
// It ignores the .git directory.
func List() []string {
	var files []string
	err := filepath.Walk(".",
		func(path string, info os.FileInfo, err error) error {
			if shouldIgnore(path) || info.IsDir() || err != nil {
				return nil //nolint:nilerr
			}
			files = append(files, path)
			return nil
		})

	if err != nil {
		return []string{}
	}
	return files
}

var defaultIgnorePatterns = []string{"node_modules", ".git", "."}

func shouldIgnore(path string) bool {
	for _, prefix := range defaultIgnorePatterns {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}
	return false
}

////

const StatusAborted = 130

// ErrAborted is the error to return when a gum command is aborted by Ctrl + C.
var ErrAborted = fmt.Errorf("aborted")
