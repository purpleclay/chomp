package chomp

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
