package chomp_test

import (
	"testing"
	"unicode"

	"github.com/purpleclay/chomp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// assertFailureNonConsuming asserts rule 1 of the combinator contract:
// on failure, a combinator must return the original input unchanged, and
// the zero value for its extracted type.
func assertFailureNonConsuming[S any](t *testing.T, name, input string, run func(string) (string, S, error)) {
	t.Helper()

	rem, ext, err := run(input)
	require.Errorf(t, err, "%s: expected failure for input %q", name, input)

	var zero S
	assert.Equalf(t, input, rem, "%s: on failure, rem must equal the original input", name)
	assert.Equalf(t, zero, ext, "%s: on failure, ext must be the zero value", name)
}

// assertSuccessIsPrefix asserts rule 2 of the combinator contract: on
// success, the extracted text is exactly the consumed prefix of the input,
// i.e. input == ext + rem. Only applicable to non-transforming Combinator[string].
func assertSuccessIsPrefix(t *testing.T, name, input string, run func(string) (string, string, error)) {
	t.Helper()

	rem, ext, err := run(input)
	require.NoErrorf(t, err, "%s: expected success for input %q", name, input)
	assert.Equalf(t, input, ext+rem, "%s: input must equal ext+rem", name)
}

func TestContract_FailureNonConsuming(t *testing.T) {
	t.Parallel()

	// basic.go
	t.Run("Char", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "Char", "", chomp.Char('a'))
		assertFailureNonConsuming(t, "Char", "xyz", chomp.Char('a'))
	})
	t.Run("AnyChar", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "AnyChar", "", chomp.AnyChar())
	})
	t.Run("Satisfy", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "Satisfy", "hello", chomp.Satisfy(unicode.IsUpper))
	})
	t.Run("Tag", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "Tag", "World", chomp.Tag("Hello"))
	})
	t.Run("TagNoCase", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "TagNoCase", "World", chomp.TagNoCase("hello"))
	})
	t.Run("Any", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "Any", "xyz", chomp.Any("eH"))
	})
	t.Run("Not", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "Not", "oxyz", chomp.Not("ol"))
	})
	t.Run("OneOf", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "OneOf", "xyz", chomp.OneOf("eH"))
	})
	t.Run("NoneOf", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "NoneOf", "epic", chomp.NoneOf("loWrd!e"))
	})
	t.Run("Until", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "Until", "Hello there", chomp.Until("World"))
	})
	t.Run("Take", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "Take", "Hi", chomp.Take(5))
	})
	t.Run("TakeUntil1", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "TakeUntil1", ",World!", chomp.TakeUntil1(","))
	})
	t.Run("Escaped", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "Escaped", "123",
			chomp.Escaped(chomp.While(chomp.IsLetter), '\\', chomp.OneOf(`"n\`)))
	})
	t.Run("EscapedTransform", func(t *testing.T) {
		t.Parallel()
		transform := func(s string) (string, string, error) {
			return s, "", chomp.CombinatorParseError{Text: s, Type: "test_transform"}
		}
		assertFailureNonConsuming(t, "EscapedTransform", "123",
			chomp.EscapedTransform(chomp.While(chomp.IsLetter), '\\', transform))
	})

	// predicate.go
	t.Run("While", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "While", "123", chomp.While(chomp.IsLetter))
	})
	t.Run("WhileN", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "WhileN", "1a", chomp.WhileN(chomp.IsDigit, 2))
	})
	t.Run("WhileNM", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "WhileNM", "123", chomp.WhileNM(chomp.IsLetter, 1, 5))
	})
	t.Run("WhileNot", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "WhileNot", "123abc", chomp.WhileNot(chomp.IsDigit))
	})
	t.Run("WhileNotN", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "WhileNotN", "123abc", chomp.WhileNotN(chomp.IsDigit, 1))
	})
	t.Run("WhileNotNM", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "WhileNotNM", "abcdefgh", chomp.WhileNotNM(chomp.IsLetter, 1, 8))
	})
	t.Run("Alpha", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "Alpha", "123", chomp.Alpha())
	})
	t.Run("Digit", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "Digit", "abc", chomp.Digit())
	})
	t.Run("Alphanumeric", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "Alphanumeric", "!!!", chomp.Alphanumeric())
	})
	t.Run("Space", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "Space", "abc", chomp.Space())
	})
	t.Run("Multispace", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "Multispace", "abc", chomp.Multispace())
	})
	t.Run("HexDigit", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "HexDigit", "xyz", chomp.HexDigit())
	})
	t.Run("OctalDigit", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "OctalDigit", "89", chomp.OctalDigit())
	})
	t.Run("BinaryDigit", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "BinaryDigit", "234", chomp.BinaryDigit())
	})
	t.Run("Newline", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "Newline", "Hello", chomp.Newline())
	})
	t.Run("Tab", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "Tab", "Hello", chomp.Tab())
	})
	t.Run("NotLineEnding", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "NotLineEnding", "\nHello", chomp.NotLineEnding())
	})
	t.Run("AnyDigit", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "AnyDigit", "abc", chomp.AnyDigit())
	})
	t.Run("AnyLetter", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "AnyLetter", "123", chomp.AnyLetter())
	})
	t.Run("AnyAlphanumeric", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "AnyAlphanumeric", "!!!", chomp.AnyAlphanumeric())
	})
	t.Run("AnyHexDigit", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "AnyHexDigit", "xyz", chomp.AnyHexDigit())
	})
	t.Run("AnyOctalDigit", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "AnyOctalDigit", "89", chomp.AnyOctalDigit())
	})
	t.Run("AnyBinaryDigit", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "AnyBinaryDigit", "234", chomp.AnyBinaryDigit())
	})

	// parser.go
	t.Run("Crlf", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "Crlf", "Hello", chomp.Crlf())
	})
	t.Run("LineEnding", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "LineEnding", "Hello", chomp.LineEnding())
	})

	// sequence.go
	t.Run("Pair_FirstFails", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "Pair", "xyz", chomp.Pair(chomp.Tag("Hello,"), chomp.Tag(" World")))
	})
	t.Run("Pair_SecondFails", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "Pair", "Hello,World",
			chomp.Pair(chomp.Tag("Hello,"), chomp.Tag(" World")))
	})
	t.Run("SepPair_FirstFails", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "SepPair", "xyz",
			chomp.SepPair(chomp.Tag("Hello"), chomp.Tag(", "), chomp.Tag("World")))
	})
	t.Run("SepPair_SepFails", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "SepPair", "HelloWorld",
			chomp.SepPair(chomp.Tag("Hello"), chomp.Tag(", "), chomp.Tag("World")))
	})
	t.Run("SepPair_SecondFails", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "SepPair", "Hello, xyz",
			chomp.SepPair(chomp.Tag("Hello"), chomp.Tag(", "), chomp.Tag("World")))
	})
	t.Run("Repeat", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "Repeat", "ababxy", chomp.Repeat(chomp.Tag("ab"), 3))
	})
	t.Run("RepeatRange", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "RepeatRange", "abxy", chomp.RepeatRange(chomp.Tag("ab"), 3, 5))
	})
	t.Run("Delimited_LeftFails", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "Delimited", `no quotes`,
			chomp.Delimited(chomp.Tag(`"`), chomp.Until(`"`), chomp.Tag(`"`)))
	})
	t.Run("Delimited_MiddleFails", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "Delimited", `"unterminated`,
			chomp.Delimited(chomp.Tag(`"`), chomp.Tag("unterminated"), chomp.Tag(`"`)))
	})
	t.Run("Delimited_RightFails", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "Delimited", `"Hello, World!`,
			chomp.Delimited(chomp.Tag(`"`), chomp.Until(`"`), chomp.Tag(`"`)))
	})
	t.Run("QuoteDouble", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "QuoteDouble", "no quotes", chomp.QuoteDouble())
	})
	t.Run("QuoteSingle", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "QuoteSingle", "no quotes", chomp.QuoteSingle())
	})
	t.Run("BracketSquare", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "BracketSquare", "no brackets", chomp.BracketSquare())
	})
	t.Run("Parentheses", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "Parentheses", "no parens", chomp.Parentheses())
	})
	t.Run("BracketAngled", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "BracketAngled", "no angles", chomp.BracketAngled())
	})
	t.Run("First", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "First", "xyz", chomp.First(chomp.Tag("Hello"), chomp.Tag("World")))
	})
	t.Run("All", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "All", "abx", chomp.All(chomp.Tag("a"), chomp.Tag("b"), chomp.Tag("c")))
	})
	t.Run("Many", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "Many", "xyz", chomp.Many(chomp.Tag("a")))
	})
	t.Run("ManyN", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "ManyN", "aax", chomp.ManyN(chomp.Tag("a"), 3))
	})
	t.Run("Prefixed_PreFails", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "Prefixed", `Hello, World!"`,
			chomp.Prefixed(chomp.Tag("Hello"), chomp.Tag(`"`)))
	})
	t.Run("Prefixed_CFails", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "Prefixed", `"Goodbye, World!"`,
			chomp.Prefixed(chomp.Tag("Hello"), chomp.Tag(`"`)))
	})
	t.Run("Suffixed_CFails", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "Suffixed", "Goodbye, World!",
			chomp.Suffixed(chomp.Tag("Hello"), chomp.Tag(", ")))
	})
	t.Run("Suffixed_SufFails", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "Suffixed", "Hello World!",
			chomp.Suffixed(chomp.Tag("Hello"), chomp.Tag(", ")))
	})
	t.Run("SeparatedList", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "SeparatedList", "123", chomp.SeparatedList(chomp.Alpha(), chomp.Tag(",")))
	})
	t.Run("ManyTill_NoElements", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "ManyTill", "END", chomp.ManyTill(chomp.AnyChar(), chomp.Tag("END")))
	})
	t.Run("ManyTill_ElementFails", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "ManyTill", "ab123",
			chomp.ManyTill(chomp.AnyLetter(), chomp.Tag("END")))
	})
	t.Run("ManyTill0_ElementFails", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "ManyTill0", "ab123",
			chomp.ManyTill0(chomp.AnyLetter(), chomp.Tag("END")))
	})
	t.Run("FoldMany", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "FoldMany", "abc",
			chomp.FoldMany(chomp.AnyDigit(), 0, func(acc int, val string) int { return acc + int(val[0]-'0') }))
	})
	t.Run("ManyCount", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "ManyCount", "123", chomp.ManyCount(chomp.AnyLetter()))
	})
	t.Run("LengthCount_LengthFails", func(t *testing.T) {
		t.Parallel()
		lengthParser := chomp.Map(chomp.AnyDigit(), func(s string) uint { return uint(s[0] - '0') })
		assertFailureNonConsuming(t, "LengthCount", "abcdef", chomp.LengthCount(lengthParser, chomp.AnyLetter()))
	})
	t.Run("LengthCount_ElementFails", func(t *testing.T) {
		t.Parallel()
		lengthParser := chomp.Map(chomp.AnyDigit(), func(s string) uint { return uint(s[0] - '0') })
		assertFailureNonConsuming(t, "LengthCount", "3ab123", chomp.LengthCount(lengthParser, chomp.AnyLetter()))
	})
	t.Run("Fill", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "Fill", "abc", chomp.Fill(chomp.AnyLetter(), 5))
	})
	t.Run("Verify", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "Verify", "Hi", chomp.Verify(chomp.Alpha(), func(s string) bool { return len(s) >= 3 }))
	})
	t.Run("Recognize", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "Recognize", "xyz", chomp.Recognize(chomp.Tag("Hello")))
	})
	t.Run("Consumed", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "Consumed", "xyz", chomp.Consumed(chomp.Tag("Hello")))
	})
	t.Run("Eof", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "Eof", "remaining", chomp.Eof())
	})
	t.Run("AllConsuming", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "AllConsuming", "Hello, World!", chomp.AllConsuming(chomp.Tag("Hello")))
	})
	t.Run("Cut", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "Cut", "World", chomp.Cut(chomp.Tag("Hello")))
	})
	t.Run("PeekNot", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "PeekNot", "Hello, World!", chomp.PeekNot(chomp.Tag("Hello")))
	})

	// modifier.go
	t.Run("Map", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "Map", "xyz",
			chomp.Map(chomp.While(chomp.IsDigit), func(in string) int { return len(in) }))
	})
	t.Run("S", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "S", "xyz", chomp.S(chomp.Tag("Hello")))
	})
	t.Run("I", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "I", "xyz",
			chomp.I(chomp.SepPair(chomp.Tag("Hello"), chomp.Tag(", "), chomp.Tag("World")), 1))
	})
	t.Run("I_IndexOutOfBounds", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "I", "Hello, World!",
			chomp.I(chomp.SepPair(chomp.Tag("Hello"), chomp.Tag(", "), chomp.Tag("World")), 5))
	})
	t.Run("Peek", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "Peek", "World", chomp.Peek(chomp.Tag("Hello")))
	})
	t.Run("Flatten", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "Flatten", "xyz", chomp.Flatten(chomp.Many(chomp.Parentheses())))
	})
}

func TestContract_SuccessIsPrefix(t *testing.T) {
	t.Parallel()

	t.Run("Char", func(t *testing.T) {
		t.Parallel()
		assertSuccessIsPrefix(t, "Char", "xyz", chomp.Char('x'))
	})
	t.Run("AnyChar", func(t *testing.T) {
		t.Parallel()
		assertSuccessIsPrefix(t, "AnyChar", "abc", chomp.AnyChar())
	})
	t.Run("Satisfy", func(t *testing.T) {
		t.Parallel()
		assertSuccessIsPrefix(t, "Satisfy", "Hello", chomp.Satisfy(unicode.IsUpper))
	})
	t.Run("Tag", func(t *testing.T) {
		t.Parallel()
		assertSuccessIsPrefix(t, "Tag", "Hello, World!", chomp.Tag("Hello"))
	})
	t.Run("TagNoCase", func(t *testing.T) {
		t.Parallel()
		assertSuccessIsPrefix(t, "TagNoCase", "HELLO, World!", chomp.TagNoCase("hello"))
	})
	t.Run("Any", func(t *testing.T) {
		t.Parallel()
		assertSuccessIsPrefix(t, "Any", "Hello, World!", chomp.Any("eH"))
	})
	t.Run("Not", func(t *testing.T) {
		t.Parallel()
		assertSuccessIsPrefix(t, "Not", "Hello, World!", chomp.Not("ol"))
	})
	t.Run("OneOf", func(t *testing.T) {
		t.Parallel()
		assertSuccessIsPrefix(t, "OneOf", "Hello, World!", chomp.OneOf("!,eH"))
	})
	t.Run("NoneOf", func(t *testing.T) {
		t.Parallel()
		assertSuccessIsPrefix(t, "NoneOf", "Hello, World!", chomp.NoneOf("loWrd!e"))
	})
	t.Run("Until", func(t *testing.T) {
		t.Parallel()
		assertSuccessIsPrefix(t, "Until", "Hello, World!", chomp.Until("World"))
	})
	t.Run("Take", func(t *testing.T) {
		t.Parallel()
		assertSuccessIsPrefix(t, "Take", "Hello, World!", chomp.Take(5))
	})
	t.Run("TakeUntil1", func(t *testing.T) {
		t.Parallel()
		assertSuccessIsPrefix(t, "TakeUntil1", "Hello, World!", chomp.TakeUntil1(","))
	})
	t.Run("Escaped", func(t *testing.T) {
		t.Parallel()
		assertSuccessIsPrefix(t, "Escaped", `Hello\"World`,
			chomp.Escaped(chomp.While(chomp.IsLetter), '\\', chomp.OneOf(`"n\`)))
	})
	t.Run("While", func(t *testing.T) {
		t.Parallel()
		assertSuccessIsPrefix(t, "While", "Hello, World!", chomp.While(chomp.IsLetter))
	})
	t.Run("WhileN", func(t *testing.T) {
		t.Parallel()
		assertSuccessIsPrefix(t, "WhileN", "123abc", chomp.WhileN(chomp.IsDigit, 2))
	})
	t.Run("WhileNM", func(t *testing.T) {
		t.Parallel()
		assertSuccessIsPrefix(t, "WhileNM", "abc123", chomp.WhileNM(chomp.IsLetter, 1, 5))
	})
	t.Run("WhileNot", func(t *testing.T) {
		t.Parallel()
		assertSuccessIsPrefix(t, "WhileNot", "abc123", chomp.WhileNot(chomp.IsDigit))
	})
	t.Run("WhileNotN", func(t *testing.T) {
		t.Parallel()
		assertSuccessIsPrefix(t, "WhileNotN", "abc123", chomp.WhileNotN(chomp.IsDigit, 1))
	})
	t.Run("WhileNotNM", func(t *testing.T) {
		t.Parallel()
		assertSuccessIsPrefix(t, "WhileNotNM", "2024abcd", chomp.WhileNotNM(chomp.IsLetter, 1, 8))
	})
	t.Run("Alpha", func(t *testing.T) {
		t.Parallel()
		assertSuccessIsPrefix(t, "Alpha", "Hello123", chomp.Alpha())
	})
	t.Run("Digit", func(t *testing.T) {
		t.Parallel()
		assertSuccessIsPrefix(t, "Digit", "123abc", chomp.Digit())
	})
	t.Run("Newline", func(t *testing.T) {
		t.Parallel()
		assertSuccessIsPrefix(t, "Newline", "\nHello", chomp.Newline())
	})
	t.Run("Tab", func(t *testing.T) {
		t.Parallel()
		assertSuccessIsPrefix(t, "Tab", "\tHello", chomp.Tab())
	})
	t.Run("NotLineEnding", func(t *testing.T) {
		t.Parallel()
		assertSuccessIsPrefix(t, "NotLineEnding", "Hello, World!\nNext line", chomp.NotLineEnding())
	})
	t.Run("AnyDigit", func(t *testing.T) {
		t.Parallel()
		assertSuccessIsPrefix(t, "AnyDigit", "123", chomp.AnyDigit())
	})
	t.Run("AnyLetter", func(t *testing.T) {
		t.Parallel()
		assertSuccessIsPrefix(t, "AnyLetter", "Hello", chomp.AnyLetter())
	})
	t.Run("Crlf", func(t *testing.T) {
		t.Parallel()
		assertSuccessIsPrefix(t, "Crlf", "\r\nHello", chomp.Crlf())
	})
	t.Run("LineEnding", func(t *testing.T) {
		t.Parallel()
		assertSuccessIsPrefix(t, "LineEnding", "\r\nHello", chomp.LineEnding())
	})
	t.Run("Recognize", func(t *testing.T) {
		t.Parallel()
		assertSuccessIsPrefix(t, "Recognize", "Hello, World!",
			chomp.Recognize(chomp.SepPair(chomp.Alpha(), chomp.Tag(", "), chomp.Alpha())))
	})
	t.Run("Rest", func(t *testing.T) {
		t.Parallel()
		assertSuccessIsPrefix(t, "Rest", "Hello, World!", chomp.Rest())
	})
	t.Run("PeekNot", func(t *testing.T) {
		t.Parallel()
		assertSuccessIsPrefix(t, "PeekNot", "World!", chomp.PeekNot(chomp.Tag("Hello")))
	})
}
