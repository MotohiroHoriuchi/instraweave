// Package diff provides unified diff generation between two text strings.
package diff

import (
	"fmt"
	"strings"
)

type editOp int

const (
	opEqual  editOp = iota
	opInsert editOp = iota
	opDelete editOp = iota
)

type edit struct {
	op   editOp
	text string
}

// compute returns the edit script transforming a into b using LCS.
func compute(a, b []string) []edit {
	m, n := len(a), len(b)

	// dp[i][j] = LCS length of a[:i] and b[:j]
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if a[i-1] == b[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
			} else if dp[i-1][j] >= dp[i][j-1] {
				dp[i][j] = dp[i-1][j]
			} else {
				dp[i][j] = dp[i][j-1]
			}
		}
	}

	// Trace back
	edits := make([]edit, 0, m+n)
	i, j := m, n
	for i > 0 || j > 0 {
		if i > 0 && j > 0 && a[i-1] == b[j-1] {
			edits = append(edits, edit{opEqual, a[i-1]})
			i--
			j--
		} else if j > 0 && (i == 0 || dp[i][j-1] >= dp[i-1][j]) {
			edits = append(edits, edit{opInsert, b[j-1]})
			j--
		} else {
			edits = append(edits, edit{opDelete, a[i-1]})
			i--
		}
	}

	// Reverse
	for l, r := 0, len(edits)-1; l < r; l, r = l+1, r-1 {
		edits[l], edits[r] = edits[r], edits[l]
	}
	return edits
}

type hunkLine struct {
	op   editOp
	text string
}

type hunk struct {
	oldStart, oldCount int
	newStart, newCount int
	lines              []hunkLine
}

func (h *hunk) String() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "@@ -%d,%d +%d,%d @@\n", h.oldStart, h.oldCount, h.newStart, h.newCount)
	for _, l := range h.lines {
		switch l.op {
		case opEqual:
			sb.WriteString(" " + l.text + "\n")
		case opInsert:
			sb.WriteString("+" + l.text + "\n")
		case opDelete:
			sb.WriteString("-" + l.text + "\n")
		}
	}
	return sb.String()
}

const contextLines = 3

// buildHunks groups edits into hunks with contextLines lines of context.
func buildHunks(edits []edit) []hunk {
	if len(edits) == 0 {
		return nil
	}

	// First pass: tag each edit with old/new line numbers
	type taggedEdit struct {
		op       editOp
		text     string
		oldLine  int
		newLine  int
	}

	tagged := make([]taggedEdit, len(edits))
	oldLine, newLine := 1, 1
	for i, e := range edits {
		tagged[i] = taggedEdit{op: e.op, text: e.text, oldLine: oldLine, newLine: newLine}
		if e.op == opEqual || e.op == opDelete {
			oldLine++
		}
		if e.op == opEqual || e.op == opInsert {
			newLine++
		}
	}

	// Find changed edit indices
	changed := make([]bool, len(tagged))
	anyChanged := false
	for i, t := range tagged {
		if t.op != opEqual {
			changed[i] = true
			anyChanged = true
		}
	}
	if !anyChanged {
		return nil
	}

	// Expand context around changed lines
	inHunk := make([]bool, len(tagged))
	for i, c := range changed {
		if !c {
			continue
		}
		start := i - contextLines
		if start < 0 {
			start = 0
		}
		end := i + contextLines
		if end >= len(tagged) {
			end = len(tagged) - 1
		}
		for k := start; k <= end; k++ {
			inHunk[k] = true
		}
	}

	// Build hunks from contiguous inHunk ranges
	var hunks []hunk
	i := 0
	for i < len(tagged) {
		if !inHunk[i] {
			i++
			continue
		}
		// Start of a new hunk
		start := i
		for i < len(tagged) && inHunk[i] {
			i++
		}
		end := i

		h := hunk{}
		h.oldStart = tagged[start].oldLine
		h.newStart = tagged[start].newLine

		for _, t := range tagged[start:end] {
			h.lines = append(h.lines, hunkLine{op: t.op, text: t.text})
			if t.op == opEqual || t.op == opDelete {
				h.oldCount++
			}
			if t.op == opEqual || t.op == opInsert {
				h.newCount++
			}
		}
		// For a pure-insert hunk (new file), old start is conventionally 0
		if h.oldCount == 0 {
			h.oldStart = 0
		}
		hunks = append(hunks, h)
	}
	return hunks
}

// splitLines splits content into lines, preserving content without trailing newlines.
func splitLines(s string) []string {
	if s == "" {
		return nil
	}
	lines := strings.Split(s, "\n")
	// If the string ends with \n, the last element is empty - remove it
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	return lines
}

// Unified generates a unified diff string between oldContent and newContent.
// Returns empty string if there are no differences.
func Unified(oldName, newName, oldContent, newContent string) string {
	oldLines := splitLines(oldContent)
	newLines := splitLines(newContent)

	edits := compute(oldLines, newLines)
	hunks := buildHunks(edits)

	if len(hunks) == 0 {
		return ""
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "--- %s\n", oldName)
	fmt.Fprintf(&sb, "+++ %s\n", newName)
	for _, h := range hunks {
		sb.WriteString(h.String())
	}
	return sb.String()
}
