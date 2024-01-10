/*
Copyright (c) 2023 - 2024 Purple Clay

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package chomp

// Pair will scan the input text and match each [Combinator] in turn. Both combinators
// must match. The result of each will be returned in the slice in execution order.
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

// SepPair will scan the input text and match each [Combinator] in turn. All combinators
// must match. The result of the separator combinator is discarded and not included
// within the returned slice.
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
			return rem, nil, ParserError{Err: err, Type: "pair_sep"}
		}

		rem, _, err = sep(rem)
		if err != nil {
			return rem, nil, ParserError{Err: err, Type: "pair_sep"}
		}

		rem, out2, err := c2(rem)
		if err != nil {
			return rem, nil, ParserError{Err: err, Type: "pair_sep"}
		}

		var ext []string
		ext = combine(ext, out1, out2)

		return rem, ext, nil
	}
}

// Repeat will scan the input text and repeat the [Combinator] the defined number
// of times. Each combinator must match, with the output of each contained in
// the returned slice.
//
//	chomp.Repeat(chomp.Parentheses(), 2)("(Hello)(World)(!)")
//	// ("(!)", []string{"(Hello)", "(World)"}, nil)
func Repeat[T Result](c Combinator[T], n int) Combinator[[]string] {
	return func(s string) (string, []string, error) {
		var ext []string
		var err error

		rem := s
		for i := 0; i < n; i++ {
			var out T
			if rem, out, err = c(rem); err != nil {
				return rem, nil, ParserError{Err: err, Type: "repeat"}
			}
			ext = combine(ext, out)
		}

		return rem, ext, nil
	}
}

// Delimited will match a series of combinators against the input text. The left
// and right combinators are used to match a delimited sequence and are discarded.
// Only the text between the delimiters is extracted.
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

// QuoteDouble will match any text delimited (or surrounded) by a pair
// of "double quotes". The delimiters are discarded.
//
//	chomp.DoubleQuote()(`"Hello, World!"`)
//	// ("", "Hello, World!", nil)
func QuoteDouble() Combinator[string] {
	return func(s string) (string, string, error) {
		return Delimited(Tag("\""), Until("\""), Tag("\""))(s)
	}
}

// QuoteSingle will match any text delimited (or surrounded) by a pair
// of 'single quotes'. The delimiters are discarded.
//
//	chomp.QuoteSingle()("'Hello, World!'")
//	// ("", "Hello, World!", nil)
func QuoteSingle() Combinator[string] {
	return func(s string) (string, string, error) {
		return Delimited(Tag("'"), Until("'"), Tag("'"))(s)
	}
}

// BracketSquare will match any text delimited (or surrounded) by a pair
// of [square brackets]. The delimiters are discarded.
//
//	chomp.BracketSquare()("[Hello, World!]")
//	// ("", "Hello, World!", nil)
func BracketSquare() Combinator[string] {
	return func(s string) (string, string, error) {
		return Delimited(Tag("["), Until("]"), Tag("]"))(s)
	}
}

// Parentheses will match any text delimited (or surrounded) by a pair
// of (parentheses). The delimiters are discarded.
//
//	chomp.Parentheses()("(Hello, World!)")
//	// ("", "Hello, World!", nil)
func Parentheses() Combinator[string] {
	return func(s string) (string, string, error) {
		return Delimited(Tag("("), Until(")"), Tag(")"))(s)
	}
}

// BracketAngled will match any text delimited (or surrounded) by a pair
// of <angled brackets>. The delimiters are discarded.
//
//	chomp.BracketAngled()("<Hello, World!>")
//	// ("", "Hello, World!", nil)
func BracketAngled() Combinator[string] {
	return func(s string) (string, string, error) {
		return Delimited(Tag("<"), Until(">"), Tag(">"))(s)
	}
}

// First will match the input text against a series of combinators. Matching
// stops as soon as the first combinator succeeds. One combinator must match.
// For better performance, try and order the combinators from most to least
// likely to match.
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

// All will match the input text against a series of combinators. All
// combinators must match in the order provided.
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
