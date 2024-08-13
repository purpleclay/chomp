package chomp

import "fmt"

// MappedCombinator is a function capable of converting the output from a [Combinator]
// into any given type. Upon success, it will return the unparsed text, along with the
// mapped value. All combinators are strict and must parse its input. Any failure to do
// so should raise a [CombinatorParseError]. It is designed for exclusive use by the
// [Map] function
type MappedCombinator[S any, T Result] func(string) (string, S, error)

// Map the result of a [Combinator] to any other type
//
//	chomp.Map(
//		chomp.While(chomp.IsDigit),
//		func (in string) int { return len(in) })("123456")
//	// ("", 6, nil)
func Map[S any, T Result](c Combinator[T], mapper func(in T) S) MappedCombinator[S, T] {
	return func(s string) (string, S, error) {
		var mapped S

		rem, out, err := c(s)
		if err != nil {
			return rem, mapped, err
		}

		mapped = mapper(out)
		return rem, mapped, nil
	}
}

// Opt allows a combinator to be optional. Any error returned by the underlying
// combinator will be swallowed. The parsed text will not be modified if the
// underlying combinator did not run.
//
//	chomp.Opt(chomp.Tag("Hey"))("Hello, World!")
//	// ("Hello, World!", "", nil)
func Opt[T Result](c Combinator[T]) Combinator[T] {
	return func(s string) (string, T, error) {
		rem, out, _ := c(s)
		return rem, out, nil
	}
}

// S wraps the result of the inner combinator within a string slice.
// Combinators of differing return types can be successfully chained
// together while using this conversion combinator.
//
//	chomp.S(chomp.Until(","))("Hello, World!")
//	// (", World!", []string{"Hello"}, nil)
func S(c Combinator[string]) Combinator[[]string] {
	return func(s string) (string, []string, error) {
		rem, ext, err := c(s)
		if err != nil {
			return rem, nil, err
		}

		return rem, []string{ext}, err
	}
}

// I extracts and returns a single string from the result of the inner combinator.
// Combinators of differing return types can be successfully chained together while
// using this conversion combinator.
//
//	chomp.I(chomp.SepPair(
//		chomp.Tag("Hello"),
//		chomp.Tag(", "),
//		chomp.Tag("World")), 1)("Hello, World!")
//	// ("!", "World", nil)
func I(c Combinator[[]string], i int) Combinator[string] {
	return func(s string) (string, string, error) {
		rem, ext, err := c(s)
		if err != nil {
			return rem, "", err
		}

		if i < 0 || i >= len(ext) {
			return rem, "", ParserError{
				Err:  fmt.Errorf("index %d is out of bounds within string slice of %d elements", i, len(ext)),
				Type: "i",
			}
		}

		return rem, ext[i], nil
	}
}

// Peek will scan the text and apply the parser without consuming any of the input.
// Useful if you need to lookahead.
//
//	chomp.Peek(chomp.Tag("Hello"))("Hello, World!")
//	// ("Hello, World!", "Hello", nil)
//
//	chomp.Peek(
//		chomp.Many(chomp.Suffixed(chomp.Tag(" "), chomp.Until(" "))),
//	)("Hello and Good Morning!")
//	// ("Hello and Good Morning!", []string{"Hello", "and", "Good"}, nil)
func Peek[T Result](c Combinator[T]) Combinator[T] {
	return func(s string) (string, T, error) {
		_, ext, err := c(s)
		return s, ext, err
	}
}
