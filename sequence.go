package chomp

// Pair will scan the input text and match each [Combinator] in turn.
// Both combinators must match.
//
//	chomp.Pair(chomp.Tag("Hello,"), chomp.Tag(" World"))("Hello, World!")
//	// ("!", []string{"Hello,", " World"}, nil)
func Pair[T, U Result](c1 Combinator[T], c2 Combinator[U]) Combinator[[]string] {
	return func(s string) (string, []string, error) {
		rem, out1, err := c1(s)
		if err != nil {
			return rem, nil, ParserError{Err: err, Type: "pair"}
		}

		rem, out2, err := c2(rem)
		if err != nil {
			return rem, nil, ParserError{Err: err, Type: "pair"}
		}

		var ext []string
		ext = combine(ext, out1, out2)

		return rem, ext, nil
	}
}

func combine(in []string, elems ...any) []string {
	for _, e := range elems {
		switch t := e.(type) {
		case string:
			in = append(in, t)
		case []string:
			in = append(in, t...)
		}
	}

	return in
}

// SepPair will scan the input text and match each [Combinator], discarding
// the separator's output. All combinators must match.
//
//	chomp.SepPair(
//		chomp.Tag("Hello"),
//		chomp.Tag(", "),
//		chomp.Tag("World"))("Hello, World!")
//	// ("!", []string{"Hello", "World"}, nil)
func SepPair[T, U, V Result](c1 Combinator[T], sep Combinator[U], c2 Combinator[V]) Combinator[[]string] {
	return func(s string) (string, []string, error) {
		rem, out1, err := c1(s)
		if err != nil {
			return rem, nil, ParserError{Err: err, Type: "sep_pair"}
		}

		rem, _, err = sep(rem)
		if err != nil {
			return rem, nil, ParserError{Err: err, Type: "sep_pair"}
		}

		rem, out2, err := c2(rem)
		if err != nil {
			return rem, nil, ParserError{Err: err, Type: "sep_pair"}
		}

		var ext []string
		ext = combine(ext, out1, out2)

		return rem, ext, nil
	}
}

// Repeat will scan the input text and match the combinator the defined
// number of times. Every execution must match.
//
//	chomp.Repeat(chomp.Parentheses(), 2)("(Hello)(World)(!)")
//	// ("(!)", []string{"(Hello)", "(World)"}, nil)
func Repeat[T Result](c Combinator[T], n uint) Combinator[[]string] {
	return func(s string) (string, []string, error) {
		var ext []string
		var err error

		rem := s
		for i := range n {
			var out T
			if rem, out, err = c(rem); err != nil {
				return rem, nil, RangedParserError{
					Err:  err,
					Exec: RangeExecution(i, n),
					Type: "repeat",
				}
			}
			ext = combine(ext, out)
		}

		return rem, ext, nil
	}
}

// RepeatRange will scan the input text and match the [Combinator] between
// a minimum and maximum number of times. It must match the expected minimum
// number of times.
//
//	chomp.RepeatRange(chomp.OneOf("Hleo"), 1, 8)("Hello, World!")
//	// (", World!", []string{"H", "e", "l", "l", "o"}, nil)
func RepeatRange[T Result](c Combinator[T], n, m uint) Combinator[[]string] {
	return func(s string) (string, []string, error) {
		var ext []string
		var err error

		if n > m {
			n, m = m, n
		}

		rem := s
		for i := range m {
			var out T
			if rem, out, err = c(rem); err != nil {
				if i+1 > n {
					break
				}
				return rem, nil, RangedParserError{
					Err:  err,
					Exec: RangeExecution(i, n, m),
					Type: "repeat_range",
				}
			}
			ext = combine(ext, out)
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
//		chomp.Tag("'"))("'Hello, World!'")
//	// ("", "Hello, World!", nil)
func Delimited[T, U, V Result](left Combinator[T], str Combinator[U], right Combinator[V]) Combinator[U] {
	return func(s string) (string, U, error) {
		var def U

		rem, _, err := left(s)
		if err != nil {
			return s, def, ParserError{Err: err, Type: "delimited"}
		}

		rem, ext, err := str(rem)
		if err != nil {
			return rem, def, ParserError{Err: err, Type: "delimited"}
		}

		rem, _, err = right(rem)
		if err != nil {
			return rem, def, ParserError{Err: err, Type: "delimited"}
		}

		return rem, ext, nil
	}
}

// QuoteDouble will match any text delimited (or surrounded) by a
// pair of "double quotes".
//
//	chomp.DoubleQuote()(`"Hello, World!"`)
//	// ("", "Hello, World!", nil)
func QuoteDouble() Combinator[string] {
	return func(s string) (string, string, error) {
		return Delimited(Tag("\""), Until("\""), Tag("\""))(s)
	}
}

// QuoteSingle will match any text delimited (or surrounded) by a
// pair of 'single quotes'.
//
//	chomp.QuoteSingle()("'Hello, World!'")
//	// ("", "Hello, World!", nil)
func QuoteSingle() Combinator[string] {
	return func(s string) (string, string, error) {
		return Delimited(Tag("'"), Until("'"), Tag("'"))(s)
	}
}

// BracketSquare will match any text delimited (or surrounded) by
// a pair of [square brackets].
//
//	chomp.BracketSquare()("[Hello, World!]")
//	// ("", "Hello, World!", nil)
func BracketSquare() Combinator[string] {
	return func(s string) (string, string, error) {
		return Delimited(Tag("["), Until("]"), Tag("]"))(s)
	}
}

// Parentheses will match any text delimited (or surrounded) by
// a pair of (parentheses).
//
//	chomp.Parentheses()("(Hello, World!)")
//	// ("", "Hello, World!", nil)
func Parentheses() Combinator[string] {
	return func(s string) (string, string, error) {
		return Delimited(Tag("("), Until(")"), Tag(")"))(s)
	}
}

// BracketAngled will match any text delimited (or surrounded) by
// a pair of <angled brackets>.
//
//	chomp.BracketAngled()("<Hello, World!>")
//	// ("", "Hello, World!", nil)
func BracketAngled() Combinator[string] {
	return func(s string) (string, string, error) {
		return Delimited(Tag("<"), Until(">"), Tag(">"))(s)
	}
}

// First will match the input text against a series of [Combinator]s.
// Matching stops as soon as the first combinator succeeds. One [Combinator]
// must match. For better performance, try and order the combinators from
// most to least likely to match.
//
//	chomp.First(
//		chomp.Tag("Good Morning"),
//		chomp.Tag("Hello"))("Good Morning, World!")
//	// (" ,World!", "Good Morning", nil)
func First[T Result](c ...Combinator[T]) Combinator[T] {
	return func(s string) (string, T, error) {
		for _, comb := range c {
			if rem, ext, err := comb(s); err == nil {
				return rem, ext, nil
			}
		}

		var out T
		return s, out, CombinatorParseError{Text: s, Type: "first"}
	}
}

// All will match the input text against a series of [Combinator]s.
// All combinators must match in the order provided.
//
//	chomp.All(
//		chomp.Tag("Hello"),
//		chomp.Until("W"),
//		chomp.Tag("World!"))("Hello, World!")
//	// ("", []string{"Hello", ", ", "World!"}, nil)
func All[T Result](c ...Combinator[T]) Combinator[[]string] {
	return func(s string) (string, []string, error) {
		var ext []string
		var err error

		rem := s
		for _, comb := range c {
			var out T
			if rem, out, err = comb(rem); err != nil {
				return rem, nil, ParserError{Err: err, Type: "all"}
			}
			ext = combine(ext, out)
		}

		return rem, ext, nil
	}
}

// Many will scan the input text, and it must match the [Combinator] at least
// once. This [Combinator] is greedy and will continuously execute until the first
// failed match. It is the equivalent of calling [ManyN] with an argument of 1.
//
//	chomp.Many(one.Of("Ho"))("Hello, World!")
//	// ("ello, World!", []string{"H"}, nil)
func Many[T Result](c Combinator[T]) Combinator[[]string] {
	return ManyN(c, 1)
}

// ManyN will scan the input text and match the [Combinator] a minimum number
// of times. This [Combinator] is greedy and will continuously execute until
// the first failed match.
//
//	chomp.ManyN(chomp.OneOf("W"), 0)("Hello, World!")
//	// ("Hello, World!", nil, nil)
func ManyN[T Result](c Combinator[T], n uint) Combinator[[]string] {
	return func(s string) (string, []string, error) {
		var ext []string
		var err error
		var count uint

		rem := s
		for {
			var out T
			var tmpRem string

			if tmpRem, out, err = c(rem); err != nil {
				break
			}
			rem = tmpRem
			ext = combine(ext, out)
			count++
		}

		if count < n {
			return rem, nil, RangedParserError{
				Err:  err,
				Exec: RangeExecution(count, n),
				Type: "many_n",
			}
		}

		return rem, ext, nil
	}
}

// Prefixed will scan the input text for a defined prefix and discard it
// before matching the remaining text against the [Combinator]. Both
// combinators must match.
//
//	chomp.Prefixed(
//		chomp.Tag("Hello"),
//		chomp.Tag(`"`))(`"Hello, World!"`)
//	// (`, World!"`, "Hello", nil)
func Prefixed(c, pre Combinator[string]) Combinator[string] {
	return func(s string) (string, string, error) {
		rem, _, err := pre(s)
		if err != nil {
			return rem, "", err
		}

		return c(rem)
	}
}

// Suffixed will scan the input text against the [Combinator] before matching a
// suffix and discarding it. Both combinators must match.
//
//	chomp.Suffixed(
//		chomp.Tag("Hello"),
//		chomp.Tag(", "))("Hello, World!")
//	// ("World!", "Hello", nil)
func Suffixed(c, suf Combinator[string]) Combinator[string] {
	return func(s string) (string, string, error) {
		rem, ext, err := c(s)
		if err != nil {
			return rem, "", err
		}

		rem, _, err = suf(rem)
		if err != nil {
			return rem, "", err
		}

		return rem, ext, nil
	}
}

// SeparatedList will scan the input text and match the [Combinator] separated
// by the provided separator. At least one element must match. The separator
// output is discarded.
//
//	chomp.SeparatedList(chomp.Alpha(), chomp.Tag(","))("a,b,c,")
//	// (",", []string{"a", "b", "c"}, nil)
func SeparatedList[T, U Result](c Combinator[T], sep Combinator[U]) Combinator[[]string] {
	return func(s string) (string, []string, error) {
		var ext []string
		var err error

		rem := s

		// First element (required)
		var out T
		if rem, out, err = c(rem); err != nil {
			return rem, nil, RangedParserError{
				Err:  err,
				Exec: RangeExecution(0, 1),
				Type: "separated_list",
			}
		}
		ext = combine(ext, out)

		// Subsequent elements (sep + element pairs)
		for {
			// Try separator - if fails, we're done
			var tmpRem string
			if tmpRem, _, err = sep(rem); err != nil {
				break
			}

			// Parse element after separator
			if tmpRem, out, err = c(tmpRem); err != nil {
				break
			}
			rem = tmpRem
			ext = combine(ext, out)
		}

		return rem, ext, nil
	}
}

// SeparatedList0 will scan the input text and match the [Combinator] separated
// by the provided separator. Zero or more elements may match. The separator
// output is discarded.
//
//	chomp.SeparatedList0(chomp.Alpha(), chomp.Tag(","))("123")
//	// ("123", []string{}, nil)
func SeparatedList0[T, U Result](c Combinator[T], sep Combinator[U]) Combinator[[]string] {
	return func(s string) (string, []string, error) {
		// Try to match first element
		rem, out, err := c(s)
		if err != nil {
			// If first element fails, return empty (0 is allowed)
			return s, []string{}, nil
		}

		var ext []string
		ext = combine(ext, out)

		// Subsequent elements (sep + element pairs)
		for {
			// Try separator - if fails, we're done
			var tmpRem string
			if tmpRem, _, err = sep(rem); err != nil {
				break
			}

			// Parse element after separator
			if tmpRem, out, err = c(tmpRem); err != nil {
				break
			}
			rem = tmpRem
			ext = combine(ext, out)
		}

		return rem, ext, nil
	}
}

// ManyTill will scan the input text, matching the [Combinator] repeatedly until
// the terminator matches. The terminator is consumed but not included in the
// result. At least one element must match before the terminator.
//
//	chomp.ManyTill(chomp.AnyChar(), chomp.Tag("END"))("abcEND")
//	// ("", []string{"a", "b", "c"}, nil)
func ManyTill[T, U Result](c Combinator[T], term Combinator[U]) Combinator[[]string] {
	return func(s string) (string, []string, error) {
		var ext []string
		var err error
		var count uint

		rem := s
		for {
			// Check for terminator first
			if tmpRem, _, termErr := term(rem); termErr == nil {
				if count == 0 {
					return rem, nil, RangedParserError{
						Err:  CombinatorParseError{Text: s, Type: "many_till"},
						Exec: RangeExecution(0, 1),
						Type: "many_till",
					}
				}
				return tmpRem, ext, nil
			}

			// Parse element
			var out T
			var tmpRem string
			if tmpRem, out, err = c(rem); err != nil {
				return rem, nil, ParserError{Err: err, Type: "many_till"}
			}
			rem = tmpRem
			ext = combine(ext, out)
			count++
		}
	}
}

// ManyTill0 will scan the input text, matching the [Combinator] repeatedly until
// the terminator matches. The terminator is consumed but not included in the
// result. Zero or more elements may match before the terminator.
//
//	chomp.ManyTill0(chomp.AnyChar(), chomp.Tag("END"))("END")
//	// ("", []string{}, nil)
func ManyTill0[T, U Result](c Combinator[T], term Combinator[U]) Combinator[[]string] {
	return func(s string) (string, []string, error) {
		var ext []string
		var err error

		rem := s
		for {
			if tmpRem, _, termErr := term(rem); termErr == nil {
				return tmpRem, ext, nil
			}

			var out T
			var tmpRem string
			if tmpRem, out, err = c(rem); err != nil {
				return rem, nil, ParserError{Err: err, Type: "many_till_0"}
			}
			rem = tmpRem
			ext = combine(ext, out)
		}
	}
}

// FoldMany will scan the input text, matching the [Combinator] repeatedly and
// accumulating results using the provided reducer function. At least one element
// must match.
//
//	chomp.FoldMany(chomp.Digit(), 0, func(acc int, val string) int {
//	    n, _ := strconv.Atoi(val)
//	    return acc + n
//	})("123abc")
//	// ("abc", 6, nil)
func FoldMany[S any, T Result](c Combinator[T], init S, reducer func(S, T) S) MappedCombinator[S, T] {
	return func(s string) (string, S, error) {
		acc := init
		var err error
		var count uint

		rem := s
		for {
			var out T
			var tmpRem string
			if tmpRem, out, err = c(rem); err != nil {
				break
			}
			rem = tmpRem
			acc = reducer(acc, out)
			count++
		}

		if count == 0 {
			return rem, init, RangedParserError{
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
// may match.
//
//	chomp.FoldMany0(chomp.Digit(), 0, func(acc int, val string) int {
//	    n, _ := strconv.Atoi(val)
//	    return acc + n
//	})("abc")
//	// ("abc", 0, nil)
func FoldMany0[S any, T Result](c Combinator[T], init S, reducer func(S, T) S) MappedCombinator[S, T] {
	return func(s string) (string, S, error) {
		acc := init

		rem := s
		for {
			tmpRem, out, err := c(rem)
			if err != nil {
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
// stored, making this memory efficient for counting.
//
//	chomp.ManyCount(chomp.Alpha())("abc123")
//	// ("123", 3, nil)
func ManyCount[T Result](c Combinator[T]) MappedCombinator[uint, T] {
	return func(s string) (string, uint, error) {
		var count uint
		var err error

		rem := s
		for {
			var tmpRem string
			if tmpRem, _, err = c(rem); err != nil {
				break
			}
			rem = tmpRem
			count++
		}

		if count == 0 {
			return rem, 0, RangedParserError{
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
// stored, making this memory efficient for counting.
//
//	chomp.ManyCount0(chomp.Alpha())("123")
//	// ("123", 0, nil)
func ManyCount0[T Result](c Combinator[T]) MappedCombinator[uint, T] {
	return func(s string) (string, uint, error) {
		var count uint

		rem := s
		for {
			tmpRem, _, err := c(rem)
			if err != nil {
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
//	    chomp.Map(chomp.Digit(), func(s string) uint {
//	        n, _ := strconv.ParseUint(s, 10, 64)
//	        return uint(n)
//	    }),
//	    chomp.Alpha(),
//	)("3abc")
//	// ("", []string{"a", "b", "c"}, nil)
func LengthCount[T Result](length MappedCombinator[uint, string], c Combinator[T]) Combinator[[]string] {
	return func(s string) (string, []string, error) {
		rem, count, err := length(s)
		if err != nil {
			return rem, nil, ParserError{Err: err, Type: "length_count"}
		}

		return Repeat(c, count)(rem)
	}
}

// Fill will scan the input text and match the [Combinator] exactly n times,
// populating the result slice. All n matches must succeed.
//
//	chomp.Fill(chomp.Alpha(), 3)("abcdef")
//	// ("def", []string{"a", "b", "c"}, nil)
func Fill[T Result](c Combinator[T], n uint) Combinator[[]string] {
	return Repeat(c, n)
}
