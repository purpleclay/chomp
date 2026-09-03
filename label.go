package chomp

import "strings"

// Label attaches a grammar-level name to any failure beneath c, so an
// error speaks in terms the grammar's author chose rather than internal
// combinator names.
//
// Nested Labels chain outermost-first: Label("manifest", Label("version",
// c)) records Labels: []string{"manifest", "version"} on the underlying
// [CombinatorParseError], rendered by [CombinatorParseError.Error] and
// [CombinatorParseError.Snippet] as "... while parsing manifest > version",
// and by [CombinatorParseError.LogValue] as the "context" field. Any
// newline in name is collapsed to a space, preserving Error's single-line
// guarantee regardless of what the caller passes.
//
//	chomp.Label("version", chomp.Digit()).Run("abc")
//	// (..., "", chomp.CombinatorParseError{..., Labels: []string{"version"}})
func Label[T any](name string, c Combinator[T]) Combinator[T] {
	name = labelReplacer.Replace(name)

	return func(s State) (State, T, error) {
		rem, out, err := c(s)
		if err != nil {
			var def T
			return s, def, addLabel(err, name)
		}

		return rem, out, nil
	}
}

// labelReplacer collapses newlines in a label name to a space, so a label
// can never introduce a newline into Error's single-line guarantee.
var labelReplacer = strings.NewReplacer("\r\n", " ", "\n", " ", "\r", " ")

// addLabel prepends name to the Labels of the [CombinatorParseError]
// reached by walking down through this package's own wrapper error types,
// rebuilding the same chain around the updated leaf so Unwrap and
// errors.As/errors.Is behaviour is unaffected. Errors outside this closed
// set (e.g. [I]'s index-out-of-bounds error) are returned unchanged, since
// they carry no position to attach a label to.
func addLabel(err error, name string) error {
	switch e := err.(type) {
	case CombinatorParseError:
		e.Labels = append([]string{name}, e.Labels...)
		return e
	case ParserError:
		e.Err = addLabel(e.Err, name)
		return e
	case RangedParserError:
		e.Err = addLabel(e.Err, name)
		return e
	case CutError:
		e.Err = addLabel(e.Err, name)
		return e
	default:
		return err
	}
}
