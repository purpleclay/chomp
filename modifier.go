package chomp

import (
	"fmt"
	"strings"
)

// MappedCombinator is a function capable of converting the output from a [Combinator]
// into any given type. Upon success, it will return the unparsed [State], along with the
// mapped value. All combinators are strict and must parse its input. Any failure to
// so should raise a [CombinatorParseError]. It is designed for exclusive use by the
// [Map] function
type MappedCombinator[S any, T Result] func(State) (State, S, error)

// Run applies c against input, the sole string-in/string-out entry point.
// Unlike raw invocation c(state), any returned error is passed through
// [Finalize] first.
func (c MappedCombinator[S, T]) Run(input string) (string, S, error) { //nolint:ireturn // S is caller-supplied and not an open interface.
	return run[S](c, input)
}

// Map the result of a [Combinator] to any other type
//
//	chomp.Map(
//		chomp.While(chomp.IsDigit),
//		func (in string) int { return len(in) }).Run("123456")
//	// ("", 6, nil)
func Map[S any, T Result](c Combinator[T], mapper func(in T) S) MappedCombinator[S, T] {
	return func(s State) (State, S, error) {
		var mapped S

		rem, out, err := c(s)
		if err != nil {
			return s, mapped, err
		}

		mapped = mapper(out)
		return rem, mapped, nil
	}
}

// Opt allows a [Combinator] to be optional by discarding its returned
// error and not modifying the input text upon failure. The inner
// combinator's remainder is never trusted on failure, even if it
// partially consumed the input before erroring.
//
//	chomp.Opt(chomp.Tag("Hey")).Run("Hello, World!")
//	// ("Hello, World!", "", nil)
func Opt[T Result](c Combinator[T]) Combinator[T] {
	return func(s State) (State, T, error) {
		rem, out, err := c(s)
		if err != nil {
			var def T
			return s, def, nil
		}

		return rem, out, nil
	}
}

// S wraps the result of the inner [Combinator] within a string slice.
// Combinators of differing return types can be successfully chained
// together while using this conversion combinator.
//
//	chomp.S(chomp.Until(",")).Run("Hello, World!")
//	// (", World!", []string{"Hello"}, nil)
func S(c Combinator[string]) Combinator[[]string] {
	return func(s State) (State, []string, error) {
		rem, ext, err := c(s)
		if err != nil {
			return s, nil, err
		}

		return rem, []string{ext}, nil
	}
}

// I extracts and returns a single string from the result of the inner
// [Combinator]. Combinators of differing return types can be successfully
// chained together while using this conversion combinator.
//
//	chomp.I(chomp.SepPair(
//		chomp.Tag("Hello"),
//		chomp.Tag(", "),
//		chomp.Tag("World")), 1).Run("Hello, World!")
//	// ("!", "World", nil)
func I(c Combinator[[]string], i int) Combinator[string] {
	return func(s State) (State, string, error) {
		rem, ext, err := c(s)
		if err != nil {
			return s, "", err
		}

		if i < 0 || i >= len(ext) {
			return s, "", ParserError{
				Err:  fmt.Errorf("index %d is out of bounds within string slice of %d elements", i, len(ext)),
				Type: "i",
			}
		}

		return rem, ext[i], nil
	}
}

// Peek will scan the text and apply the [Combinator] without consuming
// any input. Useful if you need to look ahead.
//
//	chomp.Peek(chomp.Tag("Hello")).Run("Hello, World!")
//	// ("Hello, World!", "Hello", nil)
//
//	chomp.Peek(
//		chomp.Many(chomp.Suffixed(chomp.Until(" "), chomp.Tag(" "))),
//	).Run("Hello and Good Morning!")
//	// ("Hello and Good Morning!", []string{"Hello", "and", "Good"}, nil)
func Peek[T Result](c Combinator[T]) Combinator[T] {
	return func(s State) (State, T, error) {
		_, ext, err := c(s)
		return s, ext, err
	}
}

// Flatten the output from a [Combinator] by joining all extracted values
// into a string.
//
//	chomp.Flatten(
//		chomp.Many(chomp.Parentheses()),
//	).Run("(H)(el)(lo), World!")
//	// (", World!", "Hello", nil)
func Flatten(c Combinator[[]string]) Combinator[string] {
	return func(s State) (State, string, error) {
		rem, ext, err := c(s)
		if err != nil {
			return s, "", ParserError{Err: err, Type: "flatten"}
		}
		return rem, strings.Join(ext, ""), nil
	}
}
