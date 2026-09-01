package chomp

import "strings"

// State carries the original input alongside a cursor into it. It is the
// threading value passed between every [Combinator], so any point in a
// parse can recover its absolute position with no external bookkeeping.
type State struct {
	input string
	pos   int
}

// NewState returns a [State] positioned at the start of input.
func NewState(input string) State {
	return State{input: input}
}

// Rest returns the unconsumed suffix of input from the current position.
func (s State) Rest() string {
	return s.input[s.pos:]
}

// Advance returns a copy of s moved forward by n bytes.
func (s State) Advance(n int) State {
	s.pos += n
	return s
}

// Pos returns the current byte offset into the original input.
func (s State) Pos() int {
	return s.pos
}

// since returns the text consumed between start and s.
func (s State) since(start State) string {
	return s.input[start.pos:s.pos]
}

// lineTerminators is the set of bytes that end a line, matching the
// terminators [Eol] itself recognises: LF, CRLF, and a bare CR.
const lineTerminators = "\r\n"

// lineInfo returns the 1-based line, 1-based rune column, and the text of
// the line containing the cursor. Line boundaries match [Eol]'s policy
// (LF, CRLF, or a bare CR each end a line), so a bare CR consumed upstream
// is still counted as its own line rather than folded into the next one.
func (s State) lineInfo() (line, col int, text string) {
	start := strings.LastIndexAny(s.input[:s.pos], lineTerminators) + 1
	end := len(s.input)
	if idx := strings.IndexAny(s.input[s.pos:], lineTerminators); idx != -1 {
		end = s.pos + idx
	}
	text = s.input[start:end]

	line = 1
	for i := 0; i < start; i++ {
		switch s.input[i] {
		case '\n':
			line++
		case '\r':
			line++
			if i+1 < start && s.input[i+1] == '\n' {
				i++ // \r\n is a single line break, don't double count
			}
		}
	}

	col = 1
	for range s.input[start:s.pos] {
		col++
	}

	return line, col, text
}

// Position returns the 1-based line and rune-aware column of the cursor
// within the original input.
func (s State) Position() (line, col int) {
	line, col, _ = s.lineInfo()
	return
}
