package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/purpleclay/chomp"
)

const (
	// git diff header delimiter > @@ ... @@
	hdrDelim = "@@"
	// prefix for lines added
	addPrefix = "+"
	// prefix for lines removed
	remPrefix = "-"
)

var (
	red   = lipgloss.NewStyle().Foreground(lipgloss.Color("#b34139"))
	green = lipgloss.NewStyle().Foreground(lipgloss.Color("#29b337"))
)

type FileDiff struct {
	Path   string
	Chunks []DiffChunk
}

func (d FileDiff) String() string {
	var buf strings.Builder
	fmt.Fprintf(&buf, "path: %s\n", d.Path)
	for _, chunk := range d.Chunks {
		buf.WriteString(chunk.String())
	}

	return buf.String()
}

type DiffChunk struct {
	Added   DiffChange
	Removed DiffChange
}

func (d DiffChunk) String() string {
	var buf strings.Builder

	buf.WriteString("(")
	buf.WriteString(red.Render("-"))
	fmt.Fprintf(&buf, "%d,%d ", d.Removed.LineNo, d.Removed.Count)
	buf.WriteString(green.Render("+"))
	fmt.Fprintf(&buf, "%d,%d", d.Added.LineNo, d.Added.Count)
	buf.WriteString(")\n")

	if d.Removed.Change != "" {
		buf.WriteString(red.Render(d.Removed.Change))
		buf.WriteString("\n")
	}

	if d.Added.Change != "" {
		buf.WriteString(green.Render(d.Added.Change))
		buf.WriteString("\n")
	}

	return buf.String()
}

type DiffChange struct {
	LineNo int
	Count  int
	Change string
}

func Parse(in string) (FileDiff, error) {
	rem, path, err := diffPath()(in)
	if err != nil {
		return FileDiff{}, err
	}

	rem, _, err = chomp.Until(hdrDelim)(rem)
	if err != nil {
		return FileDiff{}, err
	}

	chunks, err := diffChunks(rem)
	if err != nil {
		return FileDiff{}, err
	}

	return FileDiff{
		Path:   path,
		Chunks: chunks,
	}, nil
}

// diffPath parses the diff header line and extracts the file path.
// Uses Recognize to capture the raw path text from "a/path" format.
func diffPath() chomp.Combinator[string] {
	return func(s string) (string, string, error) {
		// diff --git a/scan/scanner.go b/scan/scanner.go
		rem, _, err := chomp.Tag("diff --git ")(s)
		if err != nil {
			return s, "", err
		}

		// Use Recognize to capture "a/path" as raw text, then extract path after "/"
		var rawPath string
		rem, rawPath, err = chomp.Recognize(
			chomp.Pair(chomp.Tag("a/"), chomp.Until(" ")),
		)(rem)
		if err != nil {
			return rem, "", err
		}

		// Strip the "a/" prefix
		path := rawPath[2:]

		rem, _, err = chomp.Eol()(rem)
		return rem, path, err
	}
}

func diffChunks(in string) ([]DiffChunk, error) {
	_, chunks, err := chomp.Map(chomp.Many(diffChunk()),
		func(in []string) []DiffChunk {
			var diffChunks []DiffChunk

			for i := range len(in) / 6 {
				// 0: removed line
				// 1: removed count
				// 2: added line
				// 3: added count
				// 4: removed lines
				// 5: added lines
				idx := i * 6
				chunk := DiffChunk{
					Removed: DiffChange{
						LineNo: mustInt(in[idx]),
						Count:  mustInt(in[idx+1]),
						Change: in[idx+4],
					},
					Added: DiffChange{
						LineNo: mustInt(in[idx+2]),
						Count:  mustInt(in[idx+3]),
						Change: in[idx+5],
					},
				}

				diffChunks = append(diffChunks, chunk)
			}

			return diffChunks
		},
	)(in)

	return chunks, err
}

func mustInt(in string) int {
	out, _ := strconv.Atoi(in)
	return out
}

func diffChunk() chomp.Combinator[[]string] {
	return func(s string) (string, []string, error) {
		/*
			@@ -25 +3,3 @@ package scan
			-import "bytes"
			+import (
			+       "bytes"
			+)
		*/
		rem, changes, err := chomp.Delimited(
			chomp.Tag(hdrDelim+" "),
			chomp.SepPair(diffChunkHeaderChange(remPrefix), chomp.Tag(" "), diffChunkHeaderChange(addPrefix)),
			chomp.Eol(),
		)(s)
		if err != nil {
			return rem, nil, err
		}

		var removed string
		rem, removed, err = chomp.Map(
			chomp.ManyN(chomp.Prefixed(chomp.Eol(), chomp.Tag(remPrefix)), 0),
			func(in []string) string { return strings.Join(in, "\n") },
		)(rem)
		if err != nil {
			return rem, nil, err
		}

		var added string
		rem, added, err = chomp.Map(
			chomp.ManyN(chomp.Prefixed(chomp.Eol(), chomp.Tag(addPrefix)), 0),
			func(in []string) string { return strings.Join(in, "\n") },
		)(rem)
		if err != nil {
			return rem, nil, err
		}

		return rem, append(changes, removed, added), nil
	}
}

// diffChunkHeaderChange parses line number and optional count from chunk header.
// Uses Verify to ensure line numbers are valid positive integers.
func diffChunkHeaderChange(prefix string) chomp.Combinator[[]string] {
	return func(s string) (string, []string, error) {
		// Matches patterns like "-25" or "+3,3"
		rem, _, err := chomp.Tag(prefix)(s)
		if err != nil {
			return rem, nil, err
		}

		return chomp.All(
			// Line number - verify it's a valid positive integer
			chomp.Verify(chomp.Digit(), func(s string) bool {
				n, err := strconv.Atoi(s)
				return err == nil && n >= 0
			}),
			// Optional count (,N) - defaults to empty string if not present
			chomp.Opt(chomp.Prefixed(chomp.Digit(), chomp.Tag(","))),
		)(rem)
	}
}

func main() {
	// generated by: git diff -U0 -- scan/scanner.go
	diff := `diff --git a/scan/scanner.go b/scan/scanner.go
index fdf8e52..02d20e5 100644
--- a/scan/scanner.go
+++ b/scan/scanner.go
@@ -25 +3,3 @@ package scan
-import "bytes"
+import (
+       "bytes"
+)
@@ -57,0 +38,25 @@ func eat(prefix byte, data []byte) []byte {
+
+// DiffLines is a split function for a [bufio.Scanner] that splits a git diff output
+// into multiple blocks of text, each prefixed by the diff --git marker. Each block
+// of text will be stripped of any leading and trailing whitespace. If the git diff
+// marker isn't detected, the entire block of text is returned, with any leading and
+// trailing whitespace stripped
+func DiffLines() func(data []byte, atEOF bool) (advance int, token []byte, err error) {
+       prefix := []byte("\ndiff --git")
+
+       return func(data []byte, atEOF bool) (advance int, token []byte, err error) {
+               if atEOF && len(data) == 0 {
+                       return 0, nil, nil
+               }
+
+               if i := bytes.Index(data, prefix); i >= 0 {
+                       return i + 1, bytes.TrimSpace(data[:i]), nil
+               }
+
+               if atEOF {
+                       return len(data), bytes.TrimSpace(data), nil
+               }
+
+               return 0, nil, nil
+       }
+}`

	fileDiff, err := Parse(diff)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(fileDiff)
}
