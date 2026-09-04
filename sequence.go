package chomp

import (
	"errors"
	"fmt"
	"log/slog"
)

// Pair will scan the input text and match each [Combinator] in turn.
// Both combinators must match.
//
//	chomp.Pair(chomp.Tag("Hello,"), chomp.Tag(" World")).Run("Hello, World!")
//	// ("!", Tuple2[string, string]{First: "Hello,", Second: " World"}, nil)
func Pair[A, B any](c1 Combinator[A], c2 Combinator[B]) Combinator[Tuple2[A, B]] {
	return func(s State) (State, Tuple2[A, B], error) {
		var def Tuple2[A, B]

		rem, out1, err := c1(s)
		if err != nil {
			return s, def, ParserError{Err: err, Type: "pair"}
		}

		rem, out2, err := c2(rem)
		if err != nil {
			return s, def, ParserError{Err: err, Type: "pair"}
		}

		return rem, Tuple2[A, B]{First: out1, Second: out2}, nil
	}
}

// SepPair will scan the input text and match each [Combinator], discarding
// the separator's output. All combinators must match.
//
//	chomp.SepPair(
//		chomp.Tag("Hello"),
//		chomp.Tag(", "),
//		chomp.Tag("World")).Run("Hello, World!")
//	// ("!", Tuple2[string, string]{First: "Hello", Second: "World"}, nil)
func SepPair[A, U, B any](c1 Combinator[A], sep Combinator[U], c2 Combinator[B]) Combinator[Tuple2[A, B]] {
	return func(s State) (State, Tuple2[A, B], error) {
		var def Tuple2[A, B]

		rem, out1, err := c1(s)
		if err != nil {
			return s, def, ParserError{Err: err, Type: "sep_pair"}
		}

		rem, _, err = sep(rem)
		if err != nil {
			return s, def, ParserError{Err: err, Type: "sep_pair"}
		}

		rem, out2, err := c2(rem)
		if err != nil {
			return s, def, ParserError{Err: err, Type: "sep_pair"}
		}

		return rem, Tuple2[A, B]{First: out1, Second: out2}, nil
	}
}

// Repeat will scan the input text and match the combinator the defined
// number of times. Every execution must match.
//
//	chomp.Repeat(chomp.Parentheses(), 2).Run("(Hello)(World)(!)")
//	// ("(!)", []string{"Hello", "World"}, nil)
func Repeat[T any](c Combinator[T], n int) Combinator[[]T] {
	return func(s State) (State, []T, error) {
		if n < 0 {
			return s, nil, ParserError{Err: fmt.Errorf("chomp: count must be non-negative, got %d", n), Type: "repeat"}
		}
		var ext []T

		rem := s
		for i := range n {
			tmpRem, out, err := c(rem)
			if err != nil {
				return s, nil, RangedParserError{
					Err:  err,
					Exec: RangeExecution(i, n),
					Type: "repeat",
				}
			}
			rem = tmpRem
			ext = append(ext, out)
		}

		return rem, ext, nil
	}
}

// RepeatRange will scan the input text and match the [Combinator] between
// a minimum and maximum number of times. It must match the expected minimum
// number of times.
//
//	chomp.RepeatRange(chomp.OneOf("Hleo"), 1, 8).Run("Hello, World!")
//	// (", World!", []string{"H", "e", "l", "l", "o"}, nil)
func RepeatRange[T any](c Combinator[T], n, m int) Combinator[[]T] {
	return func(s State) (State, []T, error) {
		if n < 0 || m < 0 {
			return s, nil, ParserError{Err: fmt.Errorf("chomp: count must be non-negative, got n=%d, m=%d", n, m), Type: "repeat_range"}
		}
		var ext []T

		if n > m {
			n, m = m, n
		}

		rem := s
		for i := range m {
			tmpRem, out, err := c(rem)
			if err != nil {
				if i+1 > n {
					break
				}
				return s, nil, RangedParserError{
					Err:  err,
					Exec: RangeExecution(i, n, m),
					Type: "repeat_range",
				}
			}
			rem = tmpRem
			ext = append(ext, out)
		}

		return rem, ext, nil
	}
}

// Delimited will match a series of combinators against the input text. All
// must match, with the delimiters being discarded.
//
//	chomp.Delimited(
//		chomp.Tag("'"),
//		chomp.Tag("Hello, World!"),
//		chomp.Tag("'")).Run("'Hello, World!'")
//	// ("", "Hello, World!", nil)
func Delimited[T, U, V any](left Combinator[T], str Combinator[U], right Combinator[V]) Combinator[U] {
	return func(s State) (State, U, error) {
		var def U

		rem, _, err := left(s)
		if err != nil {
			return s, def, ParserError{Err: err, Type: "delimited"}
		}

		rem, ext, err := str(rem)
		if err != nil {
			return s, def, ParserError{Err: err, Type: "delimited"}
		}

		rem, _, err = right(rem)
		if err != nil {
			return s, def, ParserError{Err: err, Type: "delimited"}
		}

		return rem, ext, nil
	}
}

// QuoteDouble will match any text delimited (or surrounded) by a
// pair of "double quotes".
//
//	chomp.QuoteDouble().Run(`"Hello, World!"`)
//	// ("", "Hello, World!", nil)
func QuoteDouble() Combinator[string] {
	return func(s State) (State, string, error) {
		return Delimited(Tag("\""), Until("\""), Tag("\""))(s)
	}
}

// QuoteSingle will match any text delimited (or surrounded) by a
// pair of 'single quotes'.
//
//	chomp.QuoteSingle().Run("'Hello, World!'")
//	// ("", "Hello, World!", nil)
func QuoteSingle() Combinator[string] {
	return func(s State) (State, string, error) {
		return Delimited(Tag("'"), Until("'"), Tag("'"))(s)
	}
}

// BracketSquare will match any text delimited (or surrounded) by
// a pair of [square brackets].
//
//	chomp.BracketSquare().Run("[Hello, World!]")
//	// ("", "Hello, World!", nil)
func BracketSquare() Combinator[string] {
	return func(s State) (State, string, error) {
		return Delimited(Tag("["), Until("]"), Tag("]"))(s)
	}
}

// Parentheses will match any text delimited (or surrounded) by
// a pair of (parentheses).
//
//	chomp.Parentheses().Run("(Hello, World!)")
//	// ("", "Hello, World!", nil)
func Parentheses() Combinator[string] {
	return func(s State) (State, string, error) {
		return Delimited(Tag("("), Until(")"), Tag(")"))(s)
	}
}

// BracketAngled will match any text delimited (or surrounded) by
// a pair of <angled brackets>.
//
//	chomp.BracketAngled().Run("<Hello, World!>")
//	// ("", "Hello, World!", nil)
func BracketAngled() Combinator[string] {
	return func(s State) (State, string, error) {
		return Delimited(Tag("<"), Until(">"), Tag(">"))(s)
	}
}

// First will match the input text against a series of [Combinator]s.
// Matching stops as soon as the first combinator succeeds. One [Combinator]
// must match. For better performance, try and order the combinators from
// most to least likely to match.
//
// If a [CutError] is encountered during parsing, backtracking stops immediately
// and the error is propagated. This allows [Cut] to commit to a parsing path.
//
//	chomp.First(
//		chomp.Tag("Good Morning"),
//		chomp.Tag("Hello")).Run("Good Morning, World!")
//	// (", World!", "Good Morning", nil)
func First[T any](c ...Combinator[T]) Combinator[T] {
	return func(s State) (State, T, error) {
		for _, comb := range c {
			rem, ext, err := comb(s)
			if err == nil {
				return rem, ext, nil
			}

			// Check for CutError - stop backtracking immediately
			if isCut(err) {
				var out T
				return s, out, err
			}
		}

		var out T
		return s, out, CombinatorParseError{State: s, Type: "first"}
	}
}

// isCut reports whether err is, or wraps, a [CutError]. Written as a
// manual Unwrap loop rather than errors.As: errors.As's target any
// parameter defeats escape analysis, forcing a heap allocation on every
// call regardless of whether a CutError is found - costly on First's hot
// backtracking path, where this runs once per discarded alternative.
func isCut(err error) bool {
	for err != nil {
		if _, ok := err.(CutError); ok {
			return true
		}
		err = errors.Unwrap(err)
	}
	return false
}

// All will match the input text against a series of [Combinator]s.
// All combinators must match in the order provided.
//
//	chomp.All(
//		chomp.Tag("Hello"),
//		chomp.Until("W"),
//		chomp.Tag("World!")).Run("Hello, World!")
//	// ("", []string{"Hello", ", ", "World!"}, nil)
func All[T any](c ...Combinator[T]) Combinator[[]T] {
	return func(s State) (State, []T, error) {
		var ext []T
		var err error

		rem := s
		for _, comb := range c {
			var out T
			if rem, out, err = comb(rem); err != nil {
				return s, nil, ParserError{Err: err, Type: "all"}
			}
			ext = append(ext, out)
		}

		return rem, ext, nil
	}
}

// Many will scan the input text, and it must match the [Combinator] at least
// once. This [Combinator] is greedy and will continuously execute until the first
// failed match. It is the equivalent of calling [ManyN] with an argument of 1.
// See [ManyN] for its zero-width matching behaviour.
//
//	chomp.Many(chomp.OneOf("Ho")).Run("Hello, World!")
//	// ("ello, World!", []string{"H"}, nil)
func Many[T any](c Combinator[T]) Combinator[[]T] {
	return ManyN(c, 1)
}

// ManyN will scan the input text and match the [Combinator] a minimum number
// of times. This [Combinator] is greedy and will continuously execute until
// the first failed match. An iteration that succeeds without consuming input
// stops the loop instead of being counted, so it can never repeat forever.
//
//	chomp.ManyN(chomp.OneOf("W"), 0).Run("Hello, World!")
//	// ("Hello, World!", nil, nil)
func ManyN[T any](c Combinator[T], n int) Combinator[[]T] {
	return func(s State) (State, []T, error) {
		if n < 0 {
			return s, nil, ParserError{Err: fmt.Errorf("chomp: count must be non-negative, got %d", n), Type: "many_n"}
		}
		var ext []T
		var err error
		var count int

		rem := s
		for {
			var out T
			var tmpRem State

			if tmpRem, out, err = c(rem); err != nil {
				break
			}
			if tmpRem.Pos() == rem.Pos() {
				// zero-width success: cannot make progress, stop to avoid looping forever
				err = CombinatorParseError{State: rem, Type: "many_n"}
				break
			}
			rem = tmpRem
			ext = append(ext, out)
			count++
		}

		if count < n {
			return s, nil, RangedParserError{
				Err:  err,
				Exec: RangeExecution(count, n),
				Type: "many_n",
			}
		}

		return rem, ext, nil
	}
}

// Preceded will scan the input text for a defined prefix and discard it
// before matching the remaining text against the [Combinator]. Both
// combinators must match.
//
//	chomp.Preceded(
//		chomp.Tag(`"`),
//		chomp.Tag("Hello")).Run(`"Hello, World!"`)
//	// (`, World!"`, "Hello", nil)
func Preceded(pre, c Combinator[string]) Combinator[string] {
	return func(s State) (State, string, error) {
		rem, _, err := pre(s)
		if err != nil {
			return s, "", err
		}

		rem, ext, err := c(rem)
		if err != nil {
			return s, "", err
		}

		return rem, ext, nil
	}
}

// Terminated will scan the input text against the [Combinator] before
// matching a suffix and discarding it. Both combinators must match.
//
//	chomp.Terminated(
//		chomp.Tag("Hello"),
//		chomp.Tag(", ")).Run("Hello, World!")
//	// ("World!", "Hello", nil)
func Terminated(c, suf Combinator[string]) Combinator[string] {
	return func(s State) (State, string, error) {
		rem, ext, err := c(s)
		if err != nil {
			return s, "", err
		}

		rem, _, err = suf(rem)
		if err != nil {
			return s, "", err
		}

		return rem, ext, nil
	}
}

// SeparatedList will scan the input text and match the [Combinator] separated
// by the provided separator. At least one element must match. The separator
// output is discarded. If a separator and element together succeed without
// consuming input, iteration stops instead of repeating forever.
//
//	chomp.SeparatedList(chomp.Alpha(), chomp.Tag(",")).Run("a,b,c,")
//	// (",", []string{"a", "b", "c"}, nil)
func SeparatedList[T, U any](c Combinator[T], sep Combinator[U]) Combinator[[]T] {
	return func(s State) (State, []T, error) {
		var ext []T
		var err error

		rem := s

		// First element (required)
		var out T
		if rem, out, err = c(rem); err != nil {
			return s, nil, RangedParserError{
				Err:  err,
				Exec: RangeExecution(0, 1),
				Type: "separated_list",
			}
		}
		ext = append(ext, out)

		// Subsequent elements (sep + element pairs)
		for {
			// Try separator - if fails, we're done
			var sepRem State
			if sepRem, _, err = sep(rem); err != nil {
				break
			}

			// Parse element after separator
			var tmpRem State
			if tmpRem, out, err = c(sepRem); err != nil {
				break
			}

			if tmpRem.Pos() == rem.Pos() {
				// zero-width success: separator and element together made
				// no progress, stop to avoid looping forever
				break
			}

			rem = tmpRem
			ext = append(ext, out)
		}

		return rem, ext, nil
	}
}

// SeparatedList0 will scan the input text and match the [Combinator] separated
// by the provided separator. Zero or more elements may match. The separator
// output is discarded. If a separator and element together succeed without
// consuming input, iteration stops instead of repeating forever.
//
//	chomp.SeparatedList0(chomp.Alpha(), chomp.Tag(",")).Run("123")
//	// ("123", []string{}, nil)
func SeparatedList0[T, U any](c Combinator[T], sep Combinator[U]) Combinator[[]T] {
	return func(s State) (State, []T, error) {
		// Try to match first element
		rem, out, err := c(s)
		if err != nil {
			// If first element fails, return empty (0 is allowed)
			return s, []T{}, nil
		}

		var ext []T
		ext = append(ext, out)

		// Subsequent elements (sep + element pairs)
		for {
			// Try separator - if fails, we're done
			var sepRem State
			if sepRem, _, err = sep(rem); err != nil {
				break
			}

			// Parse element after separator
			var tmpRem State
			if tmpRem, out, err = c(sepRem); err != nil {
				break
			}

			if tmpRem.Pos() == rem.Pos() {
				// zero-width success: separator and element together made
				// no progress, stop to avoid looping forever
				break
			}

			rem = tmpRem
			ext = append(ext, out)
		}

		return rem, ext, nil
	}
}

// ManyTill will scan the input text, matching the [Combinator] repeatedly until
// the terminator matches. The terminator is consumed but not included in the
// result. At least one element must match before the terminator. If an
// element succeeds without consuming input, the terminator can never be
// reached, so parsing fails instead of looping forever.
//
//	chomp.ManyTill(chomp.AnyChar(), chomp.Tag("END")).Run("abcEND")
//	// ("", []string{"a", "b", "c"}, nil)
func ManyTill[T, U any](c Combinator[T], term Combinator[U]) Combinator[[]T] {
	return func(s State) (State, []T, error) {
		var ext []T
		var err error
		var count uint

		rem := s
		for {
			// Check for terminator first
			if tmpRem, _, termErr := term(rem); termErr == nil {
				if count == 0 {
					return s, nil, RangedParserError{
						Err:  CombinatorParseError{State: s, Type: "many_till"},
						Exec: RangeExecution(0, 1),
						Type: "many_till",
					}
				}
				return tmpRem, ext, nil
			}

			// Parse element
			var out T
			var tmpRem State
			if tmpRem, out, err = c(rem); err != nil {
				return s, nil, ParserError{Err: err, Type: "many_till"}
			}

			if tmpRem.Pos() == rem.Pos() {
				// zero-width success: the terminator can never be reached
				return s, nil, ParserError{
					Err:  CombinatorParseError{State: rem, Type: "many_till"},
					Type: "many_till",
				}
			}

			rem = tmpRem
			ext = append(ext, out)
			count++
		}
	}
}

// ManyTill0 will scan the input text, matching the [Combinator] repeatedly until
// the terminator matches. The terminator is consumed but not included in the
// result. Zero or more elements may match before the terminator. If an
// element succeeds without consuming input, the terminator can never be
// reached, so parsing fails instead of looping forever.
//
//	chomp.ManyTill0(chomp.AnyChar(), chomp.Tag("END")).Run("END")
//	// ("", nil, nil)
func ManyTill0[T, U any](c Combinator[T], term Combinator[U]) Combinator[[]T] {
	return func(s State) (State, []T, error) {
		var ext []T
		var err error

		rem := s
		for {
			if tmpRem, _, termErr := term(rem); termErr == nil {
				return tmpRem, ext, nil
			}

			var out T
			var tmpRem State
			if tmpRem, out, err = c(rem); err != nil {
				return s, nil, ParserError{Err: err, Type: "many_till_0"}
			}

			if tmpRem.Pos() == rem.Pos() {
				// zero-width success: the terminator can never be reached
				return s, nil, ParserError{
					Err:  CombinatorParseError{State: rem, Type: "many_till_0"},
					Type: "many_till_0",
				}
			}

			rem = tmpRem
			ext = append(ext, out)
		}
	}
}

// FoldMany will scan the input text, matching the [Combinator] repeatedly and
// accumulating results using the provided reducer function. At least one element
// must match. An iteration that succeeds without consuming input stops the
// loop instead of being counted, so it can never repeat forever.
//
//	chomp.FoldMany(chomp.AnyDigit(), 0, func(acc int, val string) int {
//	    n, _ := strconv.Atoi(val)
//	    return acc + n
//	}).Run("123abc")
//	// ("abc", 6, nil)
func FoldMany[S, T any](c Combinator[T], init S, reducer func(S, T) S) Combinator[S] {
	return func(s State) (State, S, error) {
		acc := init
		var err error
		var count uint

		rem := s
		for {
			var out T
			var tmpRem State
			if tmpRem, out, err = c(rem); err != nil {
				break
			}
			if tmpRem.Pos() == rem.Pos() {
				// zero-width success: cannot make progress, stop to avoid looping forever
				err = CombinatorParseError{State: rem, Type: "fold_many"}
				break
			}
			rem = tmpRem
			acc = reducer(acc, out)
			count++
		}

		if count == 0 {
			return s, init, RangedParserError{
				Err:  err,
				Exec: RangeExecution(0, 1),
				Type: "fold_many",
			}
		}

		return rem, acc, nil
	}
}

// FoldMany0 will scan the input text, matching the [Combinator] repeatedly and
// accumulating results using the provided reducer function. Zero or more elements
// may match. An iteration that succeeds without consuming input stops the
// loop instead of being counted, so it can never repeat forever.
//
//	chomp.FoldMany0(chomp.AnyDigit(), 0, func(acc int, val string) int {
//	    n, _ := strconv.Atoi(val)
//	    return acc + n
//	}).Run("abc")
//	// ("abc", 0, nil)
func FoldMany0[S, T any](c Combinator[T], init S, reducer func(S, T) S) Combinator[S] {
	return func(s State) (State, S, error) {
		acc := init

		rem := s
		for {
			tmpRem, out, err := c(rem)
			if err != nil {
				break
			}
			if tmpRem.Pos() == rem.Pos() {
				// zero-width success: cannot make progress, stop to avoid looping forever
				break
			}
			rem = tmpRem
			acc = reducer(acc, out)
		}

		return rem, acc, nil
	}
}

// ManyCount will scan the input text and count the number of times the
// [Combinator] matches. At least one match is required. Results are not
// stored, making this memory efficient for counting. An iteration that
// succeeds without consuming input stops counting instead of repeating
// forever.
//
//	chomp.ManyCount(chomp.AnyLetter()).Run("abc123")
//	// ("123", 3, nil)
func ManyCount[T any](c Combinator[T]) Combinator[int] {
	return func(s State) (State, int, error) {
		var count int
		var err error

		rem := s
		for {
			var tmpRem State
			if tmpRem, _, err = c(rem); err != nil {
				break
			}
			if tmpRem.Pos() == rem.Pos() {
				// zero-width success: cannot make progress, stop to avoid looping forever
				err = CombinatorParseError{State: rem, Type: "many_count"}
				break
			}
			rem = tmpRem
			count++
		}

		if count == 0 {
			return s, 0, RangedParserError{
				Err:  err,
				Exec: RangeExecution(0, 1),
				Type: "many_count",
			}
		}

		return rem, count, nil
	}
}

// ManyCount0 will scan the input text and count the number of times the
// [Combinator] matches. Zero or more matches are allowed. Results are not
// stored, making this memory efficient for counting. An iteration that
// succeeds without consuming input stops counting instead of repeating
// forever.
//
//	chomp.ManyCount0(chomp.AnyLetter()).Run("123")
//	// ("123", 0, nil)
func ManyCount0[T any](c Combinator[T]) Combinator[int] {
	return func(s State) (State, int, error) {
		var count int

		rem := s
		for {
			tmpRem, _, err := c(rem)
			if err != nil {
				break
			}
			if tmpRem.Pos() == rem.Pos() {
				// zero-width success: cannot make progress, stop to avoid looping forever
				break
			}
			rem = tmpRem
			count++
		}

		return rem, count, nil
	}
}

// LengthCount will first parse a length value using the length combinator,
// then apply the element combinator that exact number of times.
//
//	chomp.LengthCount(
//	    chomp.Map(chomp.AnyDigit(), func(s string) int {
//	        n, _ := strconv.Atoi(s)
//	        return n
//	    }),
//	    chomp.AnyLetter(),
//	).Run("3abc")
//	// ("", []string{"a", "b", "c"}, nil)
func LengthCount[T any](length Combinator[int], c Combinator[T]) Combinator[[]T] {
	return func(s State) (State, []T, error) {
		rem, count, err := length(s)
		if err != nil {
			return s, nil, ParserError{Err: err, Type: "length_count"}
		}

		rem, ext, err := Repeat(c, count)(rem)
		if err != nil {
			return s, nil, err
		}

		return rem, ext, nil
	}
}

// Verify validates the parsed result against a predicate function without
// modifying the output. If the predicate returns false, the combinator fails.
// Useful for semantic validation of parsed data.
//
//	chomp.Verify(chomp.Alpha(), func(s string) bool {
//	    return len(s) >= 3
//	}).Run("Hello, World!")
//	// (", World!", "Hello", nil)
//
//	chomp.Verify(chomp.Alpha(), func(s string) bool {
//	    return len(s) >= 10
//	}).Run("Hello, World!")
//	// ("Hello, World!", "", error)
func Verify[T any](c Combinator[T], predicate func(T) bool) Combinator[T] {
	return func(s State) (State, T, error) {
		var def T

		rem, out, err := c(s)
		if err != nil {
			return s, def, err
		}

		if !predicate(out) {
			return s, def, CombinatorParseError{State: s, Type: "verify"}
		}

		return rem, out, nil
	}
}

// Recognize returns the consumed input as the output, regardless of the
// inner parser's result. Useful for capturing complex patterns as text.
//
//	chomp.Recognize(chomp.SepPair(
//	    chomp.Alpha(),
//	    chomp.Tag(", "),
//	    chomp.Alpha())).Run("Hello, World!")
//	// ("!", "Hello, World", nil)
func Recognize[T any](c Combinator[T]) Combinator[string] {
	return func(s State) (State, string, error) {
		rem, _, err := c(s)
		if err != nil {
			return s, "", err
		}

		consumed := rem.since(s)
		return rem, consumed, nil
	}
}

// Consumed provides both the raw consumed text and the parsed output as a tuple.
// Enables access to both representations simultaneously.
//
//	chomp.Consumed(chomp.SepPair(
//	    chomp.Alpha(),
//	    chomp.Tag(", "),
//	    chomp.Alpha())).Run("Hello, World!")
//	// ("!", Tuple2[string, Tuple2[string, string]]{First: "Hello, World", Second: Tuple2[string, string]{First: "Hello", Second: "World"}}, nil)
func Consumed[T any](c Combinator[T]) Combinator[Tuple2[string, T]] {
	return func(s State) (State, Tuple2[string, T], error) {
		var def Tuple2[string, T]

		rem, out, err := c(s)
		if err != nil {
			return s, def, err
		}

		consumed := rem.since(s)

		return rem, Tuple2[string, T]{First: consumed, Second: out}, nil
	}
}

// Eof matches only when at the end of input, returning an empty string
// on success. Prevents partial parsing by ensuring no input remains.
//
//	chomp.Eof().Run("")
//	// ("", "", nil)
//
//	chomp.Eof().Run("remaining")
//	// ("remaining", "", error)
func Eof() Combinator[string] {
	return func(s State) (State, string, error) {
		if s.Rest() == "" {
			return s, "", nil
		}
		return s, "", CombinatorParseError{State: s, Type: "eof"}
	}
}

// AllConsuming ensures the entire input is consumed by the inner parser,
// failing if any text remains unparsed.
//
//	chomp.AllConsuming(chomp.Tag("Hello")).Run("Hello")
//	// ("", "Hello", nil)
//
//	chomp.AllConsuming(chomp.Tag("Hello")).Run("Hello, World!")
//	// ("Hello, World!", "", error)
func AllConsuming[T any](c Combinator[T]) Combinator[T] {
	return func(s State) (State, T, error) {
		var def T

		rem, out, err := c(s)
		if err != nil {
			return s, def, err
		}

		if rem.Rest() != "" {
			return s, def, CombinatorParseError{
				State: rem,
				Type:  "all_consuming",
			}
		}

		return rem, out, nil
	}
}

// Rest returns all remaining unconsumed input as a string value.
// Always succeeds, even with empty input.
//
//	chomp.Rest().Run("Hello, World!")
//	// ("", "Hello, World!", nil)
//
//	chomp.Rest().Run("")
//	// ("", "", nil)
func Rest() Combinator[string] {
	return func(s State) (State, string, error) {
		rest := s.Rest()
		return s.Advance(len(rest)), rest, nil
	}
}

// Value returns a fixed value upon parser success, discarding the actual
// parse result. Useful for assigning semantic meaning to parsed tokens.
//
//	chomp.Value(chomp.Tag("true"), true).Run("true")
//	// ("", true, nil)
//
//	chomp.Value(chomp.Tag("false"), false).Run("false")
//	// ("", false, nil)
func Value[S, T any](c Combinator[T], val S) Combinator[S] {
	return func(s State) (State, S, error) {
		var def S

		rem, _, err := c(s)
		if err != nil {
			return s, def, err
		}

		return rem, val, nil
	}
}

// Cond conditionally applies a parser based on a boolean flag. If the
// condition is true, the parser is applied. Otherwise, it returns an
// empty result without consuming input. Enables optional parsing logic.
//
//	chomp.Cond(true, chomp.Tag("Hello")).Run("Hello, World!")
//	// (", World!", "Hello", nil)
//
//	chomp.Cond(false, chomp.Tag("Hello")).Run("Hello, World!")
//	// ("Hello, World!", "", nil)
func Cond[T any](cond bool, c Combinator[T]) Combinator[T] {
	return func(s State) (State, T, error) {
		var def T
		if !cond {
			return s, def, nil
		}
		return c(s)
	}
}

// CutError is a fatal parsing error that prevents backtracking past the
// decision point. Used with [Cut] to improve error messaging.
type CutError struct {
	// Err contains the underlying error that caused the cut.
	Err error
}

// Error delegates to the inner error's Error(), for the same reason as
// [ParserError.Error]. Fatal-vs-recoverable is a distinction for
// [errors.As]/[errors.Is] on the CutError type itself, not the rendered
// string.
func (e CutError) Error() string {
	if e.Err == nil {
		return "chomp: cut error with no underlying cause"
	}
	return e.Err.Error()
}

// Unwrap returns the inner error.
func (e CutError) Unwrap() error {
	return e.Err
}

// LogValue delegates to the inner error's LogValue if it implements
// [slog.LogValuer], otherwise falls back to its Error() string.
func (e CutError) LogValue() slog.Value {
	if lv, ok := e.Err.(slog.LogValuer); ok {
		return lv.LogValue()
	}
	return slog.StringValue(e.Error())
}

// Cut converts recoverable parsing errors into fatal failures, preventing
// backtracking past decision points. Improves error messaging by committing
// to a parsing path once the cut point is reached.
//
//	// Without Cut, First would try the second alternative
//	// With Cut, once "if" matches, failure is fatal
//	chomp.First(
//	    chomp.All(
//	        chomp.Tag("if"),
//	        chomp.Cut(chomp.Tag("("))),
//	    chomp.S(chomp.Tag("identifier"))).Run("if x")
//	// ("if x", nil, CutError{...})
func Cut[T any](c Combinator[T]) Combinator[T] {
	return func(s State) (State, T, error) {
		rem, out, err := c(s)
		if err != nil {
			var def T
			return s, def, CutError{Err: err}
		}
		return rem, out, nil
	}
}

// PeekNot succeeds when the inner parser fails without consuming input.
// Implements negative lookahead for validation. On success, returns an
// empty string without consuming any input. Pairs with [Peek] for
// positive lookahead.
//
//	chomp.PeekNot(chomp.Tag("Hello")).Run("World!")
//	// ("World!", "", nil)
//
//	chomp.PeekNot(chomp.Tag("Hello")).Run("Hello, World!")
//	// ("Hello, World!", "", error)
func PeekNot[T any](c Combinator[T]) Combinator[string] {
	return func(s State) (State, string, error) {
		_, _, err := c(s)
		if err == nil {
			return s, "", CombinatorParseError{State: s, Type: "peek_not"}
		}
		return s, "", nil
	}
}
