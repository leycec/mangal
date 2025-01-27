package util

import (
	"fmt"
	"github.com/metafates/mangal/constant"
	"github.com/samber/lo"
	"golang.org/x/exp/constraints"
	"golang.org/x/term"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

// PadZero pads a number with leading zeros.
func PadZero(s string, l int) string {
	return strings.Repeat("0", Max(l-len(s), 0)) + s
}

// replacers is a list of regexp.Regexp pairs that will be used to sanitize filenames.
var replacers = []lo.Tuple2[*regexp.Regexp, string]{
	{regexp.MustCompile(`[\\/<>:;"'|?!*{}#%&^+,~\s]`), "_"},
	{regexp.MustCompile(`__+`), "_"},
	{regexp.MustCompile(`^[_\-.]+|[_\-.]+$`), ""},
}

// SanitizeFilename will remove all invalid characters from a path.
func SanitizeFilename(filename string) string {
	for _, re := range replacers {
		filename = re.A.ReplaceAllString(filename, re.B)
	}

	return filename
}

// Quantity returns formatted quantity.
// Example:
//
//	Quantity(1, "manga") -> "1 manga"
//	Quantity(2, "manga") -> "2 mangas"
func Quantity(count int, thing string) string {
	thing = strings.TrimSuffix(thing, "s")
	if count == 1 {
		return fmt.Sprintf("%d %s", count, thing)
	}

	return fmt.Sprintf("%d %ss", count, thing)
}

// TerminalSize returns the dimensions of the given terminal.
func TerminalSize() (width, height int, err error) {
	return term.GetSize(int(os.Stdout.Fd()))
}

// FileStem returns the file name without the extension.
func FileStem(path string) string {
	return strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
}

// ClearScreen clears the terminal screen.
func ClearScreen() {
	run := func(name string, args ...string) error {
		command := exec.Command(name, args...)
		command.Stdout = os.Stdout
		command.Stdin = os.Stdin
		command.Stderr = os.Stderr
		return command.Run()
	}

	switch runtime.GOOS {
	case constant.Linux, constant.Darwin:
		err := run("tput", "clear")
		if err != nil {
			_ = run("clear")
		}
	case constant.Windows:
		_ = run("cls")
	}
}

// ReGroups parses the string with the given regular expression and returns the
// group values defined in the expression.
func ReGroups(pattern *regexp.Regexp, str string) (groups map[string]string) {
	groups = make(map[string]string)
	match := pattern.FindStringSubmatch(str)

	for i, name := range pattern.SubexpNames() {
		if i > 0 && i <= len(match) {
			groups[name] = match[i]
		}
	}

	return
}

// Ignore calls function and explicitely ignores error
func Ignore(f func() error) {
	_ = f()
}

// Max returns the maximum value of the given items.
func Max[T constraints.Ordered](items ...T) (max T) {
	for _, item := range items {
		if item > max {
			max = item
		}
	}

	return
}

// Min returns the minimum value of the given items.
func Min[T constraints.Ordered](items ...T) (min T) {
	min = items[0]
	for _, item := range items {
		if item < min {
			min = item
		}
	}

	return
}

// PrintErasable prints a string that can be erased by calling a returned function.
func PrintErasable(msg string) (eraser func()) {
	_, _ = fmt.Fprintf(os.Stdout, "\r%s", msg)

	return func() {
		_, _ = fmt.Fprintf(os.Stdout, "\r%s\r", strings.Repeat(" ", len(msg)))
	}
}

// Capitalize returns a string with the first letter capitalized.
func Capitalize(s string) string {
	return strings.ToUpper(s[:1]) + s[1:]
}

func CompareVersions(a, b string) (int, error) {
	type version struct {
		major, minor, patch int
	}

	parse := func(s string) (version, error) {
		var v version
		_, err := fmt.Sscanf(strings.TrimPrefix(s, "v"), "%d.%d.%d", &v.major, &v.minor, &v.patch)
		return v, err
	}

	av, err := parse(a)
	if err != nil {
		return 0, err
	}

	bv, err := parse(b)
	if err != nil {
		return 0, err
	}

	for _, pair := range []lo.Tuple2[int, int]{
		{av.major, bv.major},
		{av.minor, bv.minor},
		{av.patch, bv.patch},
	} {
		if pair.A > pair.B {
			return 1, nil
		}

		if pair.A < pair.B {
			return -1, nil
		}
	}

	return 0, nil
}
