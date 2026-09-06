package chomp_test

import (
	"errors"
	"strconv"
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
	assert.NotContainsf(t, err.Error(), "\n", "%s: Error() must not contain a newline", name)
}

// assertNegativeCountFailsGracefully asserts a count-taking combinator
// given a negative count fails the same way [assertFailureNonConsuming]
// checks, plus one more: the error must not be a CombinatorParseError,
// since a negative count is a constructor/programming error, not a parse
// failure at some position in the input - the same class of error as I's
// out-of-bounds index.
func assertNegativeCountFailsGracefully[S any](t *testing.T, name, input string, run func(string) (string, S, error)) {
	t.Helper()

	rem, ext, err := run(input)
	require.Errorf(t, err, "%s: expected failure for negative count", name)

	var zero S
	assert.Equalf(t, input, rem, "%s: on failure, rem must equal the original input", name)
	assert.Equalf(t, zero, ext, "%s: on failure, ext must be the zero value", name)
	assert.NotContainsf(t, err.Error(), "\n", "%s: Error() must not contain a newline", name)

	var pe chomp.CombinatorParseError
	assert.Falsef(t, errors.As(err, &pe), "%s: negative count error must not be a CombinatorParseError", name)
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
		assertFailureNonConsuming(t, "Char", "", chomp.Char('a').Run)
		assertFailureNonConsuming(t, "Char", "xyz", chomp.Char('a').Run)
	})
	t.Run("AnyChar", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "AnyChar", "", chomp.AnyChar().Run)
	})
	t.Run("Satisfy", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "Satisfy", "hello", chomp.Satisfy(unicode.IsUpper).Run)
	})
	t.Run("Tag", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "Tag", "World", chomp.Tag("Hello").Run)
	})
	t.Run("TagNoCase", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "TagNoCase", "World", chomp.TagNoCase("hello").Run)
	})
	t.Run("IsA", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "IsA", "xyz", chomp.IsA("eH").Run)
	})
	t.Run("IsNot", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "IsNot", "oxyz", chomp.IsNot("ol").Run)
	})
	t.Run("OneOf", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "OneOf", "xyz", chomp.OneOf("eH").Run)
	})
	t.Run("NoneOf", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "NoneOf", "epic", chomp.NoneOf("loWrd!e").Run)
	})
	t.Run("Until", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "Until", "Hello there", chomp.Until("World").Run)
	})
	t.Run("Take", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "Take", "Hi", chomp.Take(5).Run)
	})
	t.Run("TakeUntil1", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "TakeUntil1", ",World!", chomp.TakeUntil1(",").Run)
	})
	t.Run("Escaped", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "Escaped", "123",
			chomp.Escaped(chomp.While(chomp.IsLetter), '\\', chomp.OneOf(`"n\`)).Run)
	})
	t.Run("EscapedTransform", func(t *testing.T) {
		t.Parallel()
		transform := func(s chomp.State) (chomp.State, string, error) {
			return s, "", chomp.CombinatorParseError{State: s}
		}
		assertFailureNonConsuming(t, "EscapedTransform", "123",
			chomp.EscapedTransform(chomp.While(chomp.IsLetter), '\\', transform).Run)
	})

	// predicate.go
	t.Run("While", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "While", "123", chomp.While(chomp.IsLetter).Run)
	})
	t.Run("WhileN", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "WhileN", "1a", chomp.WhileN(chomp.IsDigit, 2).Run)
	})
	t.Run("WhileNM", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "WhileNM", "123", chomp.WhileNM(chomp.IsLetter, 1, 5).Run)
	})
	t.Run("WhileNot", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "WhileNot", "123abc", chomp.WhileNot(chomp.IsDigit).Run)
	})
	t.Run("WhileNotN", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "WhileNotN", "123abc", chomp.WhileNotN(chomp.IsDigit, 1).Run)
	})
	t.Run("WhileNotNM", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "WhileNotNM", "abcdefgh", chomp.WhileNotNM(chomp.IsLetter, 1, 8).Run)
	})
	t.Run("Alpha", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "Alpha", "123", chomp.Alpha().Run)
	})
	t.Run("Digit", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "Digit", "abc", chomp.Digit().Run)
	})
	t.Run("Alphanumeric", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "Alphanumeric", "!!!", chomp.Alphanumeric().Run)
	})
	t.Run("Space", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "Space", "abc", chomp.Space().Run)
	})
	t.Run("Multispace", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "Multispace", "abc", chomp.Multispace().Run)
	})
	t.Run("HexDigit", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "HexDigit", "xyz", chomp.HexDigit().Run)
	})
	t.Run("OctalDigit", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "OctalDigit", "89", chomp.OctalDigit().Run)
	})
	t.Run("BinaryDigit", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "BinaryDigit", "234", chomp.BinaryDigit().Run)
	})
	t.Run("Newline", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "Newline", "Hello", chomp.Newline().Run)
	})
	t.Run("Tab", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "Tab", "Hello", chomp.Tab().Run)
	})
	t.Run("NotLineEnding", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "NotLineEnding", "\nHello", chomp.NotLineEnding().Run)
	})
	t.Run("AnyDigit", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "AnyDigit", "abc", chomp.AnyDigit().Run)
	})
	t.Run("AnyLetter", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "AnyLetter", "123", chomp.AnyLetter().Run)
	})
	t.Run("AnyAlphanumeric", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "AnyAlphanumeric", "!!!", chomp.AnyAlphanumeric().Run)
	})
	t.Run("AnyHexDigit", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "AnyHexDigit", "xyz", chomp.AnyHexDigit().Run)
	})
	t.Run("AnyOctalDigit", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "AnyOctalDigit", "89", chomp.AnyOctalDigit().Run)
	})
	t.Run("AnyBinaryDigit", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "AnyBinaryDigit", "234", chomp.AnyBinaryDigit().Run)
	})

	// parser.go
	t.Run("Crlf", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "Crlf", "Hello", chomp.Crlf().Run)
	})
	t.Run("LineEnding", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "LineEnding", "Hello", chomp.LineEnding().Run)
	})

	// sequence.go
	t.Run("Pair_FirstFails", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "Pair", "xyz", chomp.Pair(chomp.Tag("Hello,"), chomp.Tag(" World")).Run)
	})
	t.Run("Pair_SecondFails", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "Pair", "Hello,World",
			chomp.Pair(chomp.Tag("Hello,"), chomp.Tag(" World")).Run)
	})
	t.Run("SepPair_FirstFails", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "SepPair", "xyz",
			chomp.SepPair(chomp.Tag("Hello"), chomp.Tag(", "), chomp.Tag("World")).Run)
	})
	t.Run("SepPair_SepFails", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "SepPair", "HelloWorld",
			chomp.SepPair(chomp.Tag("Hello"), chomp.Tag(", "), chomp.Tag("World")).Run)
	})
	t.Run("SepPair_SecondFails", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "SepPair", "Hello, xyz",
			chomp.SepPair(chomp.Tag("Hello"), chomp.Tag(", "), chomp.Tag("World")).Run)
	})
	t.Run("Repeat", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "Repeat", "ababxy", chomp.Repeat(chomp.Tag("ab"), 3).Run)
	})
	t.Run("RepeatRange", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "RepeatRange", "abxy", chomp.RepeatRange(chomp.Tag("ab"), 3, 5).Run)
	})
	t.Run("Delimited_LeftFails", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "Delimited", `no quotes`,
			chomp.Delimited(chomp.Tag(`"`), chomp.Until(`"`), chomp.Tag(`"`)).Run)
	})
	t.Run("Delimited_MiddleFails", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "Delimited", `"unterminated`,
			chomp.Delimited(chomp.Tag(`"`), chomp.Tag("unterminated"), chomp.Tag(`"`)).Run)
	})
	t.Run("Delimited_RightFails", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "Delimited", `"Hello, World!`,
			chomp.Delimited(chomp.Tag(`"`), chomp.Until(`"`), chomp.Tag(`"`)).Run)
	})
	t.Run("QuoteDouble", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "QuoteDouble", "no quotes", chomp.QuoteDouble().Run)
	})
	t.Run("QuoteSingle", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "QuoteSingle", "no quotes", chomp.QuoteSingle().Run)
	})
	t.Run("BracketSquare", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "BracketSquare", "no brackets", chomp.BracketSquare().Run)
	})
	t.Run("Parentheses", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "Parentheses", "no parens", chomp.Parentheses().Run)
	})
	t.Run("BracketAngled", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "BracketAngled", "no angles", chomp.BracketAngled().Run)
	})
	t.Run("First", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "First", "xyz", chomp.First(chomp.Tag("Hello"), chomp.Tag("World")).Run)
	})
	t.Run("All", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "All", "abx", chomp.All(chomp.Tag("a"), chomp.Tag("b"), chomp.Tag("c")).Run)
	})
	t.Run("Many", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "Many", "xyz", chomp.Many(chomp.Tag("a")).Run)
	})
	t.Run("ManyN", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "ManyN", "aax", chomp.ManyN(chomp.Tag("a"), 3).Run)
	})
	t.Run("Preceded_PreFails", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "Preceded", `Hello, World!"`,
			chomp.Preceded(chomp.Tag(`"`), chomp.Tag("Hello")).Run)
	})
	t.Run("Preceded_CFails", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "Preceded", `"Goodbye, World!"`,
			chomp.Preceded(chomp.Tag(`"`), chomp.Tag("Hello")).Run)
	})
	t.Run("Terminated_CFails", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "Terminated", "Goodbye, World!",
			chomp.Terminated(chomp.Tag("Hello"), chomp.Tag(", ")).Run)
	})
	t.Run("Terminated_SufFails", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "Terminated", "Hello World!",
			chomp.Terminated(chomp.Tag("Hello"), chomp.Tag(", ")).Run)
	})
	t.Run("SeparatedList", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "SeparatedList", "123", chomp.SeparatedList(chomp.Alpha(), chomp.Tag(",")).Run)
	})
	t.Run("ManyTill_NoElements", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "ManyTill", "END", chomp.ManyTill(chomp.AnyChar(), chomp.Tag("END")).Run)
	})
	t.Run("ManyTill_ElementFails", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "ManyTill", "ab123",
			chomp.ManyTill(chomp.AnyLetter(), chomp.Tag("END")).Run)
	})
	t.Run("ManyTill0_ElementFails", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "ManyTill0", "ab123",
			chomp.ManyTill0(chomp.AnyLetter(), chomp.Tag("END")).Run)
	})
	t.Run("FoldMany", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "FoldMany", "abc",
			chomp.FoldMany(chomp.AnyDigit(), 0, func(acc int, val string) int { return acc + int(val[0]-'0') }).Run)
	})
	t.Run("ManyCount", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "ManyCount", "123", chomp.ManyCount(chomp.AnyLetter()).Run)
	})
	t.Run("LengthCount_LengthFails", func(t *testing.T) {
		t.Parallel()
		lengthParser := chomp.Map(chomp.AnyDigit(), func(s string) int { return int(s[0] - '0') })
		assertFailureNonConsuming(t, "LengthCount", "abcdef", chomp.LengthCount(lengthParser, chomp.AnyLetter()).Run)
	})
	t.Run("LengthCount_ElementFails", func(t *testing.T) {
		t.Parallel()
		lengthParser := chomp.Map(chomp.AnyDigit(), func(s string) int { return int(s[0] - '0') })
		assertFailureNonConsuming(t, "LengthCount", "3ab123", chomp.LengthCount(lengthParser, chomp.AnyLetter()).Run)
	})
	t.Run("Verify", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "Verify", "Hi", chomp.Verify(chomp.Alpha(), func(s string) bool { return len(s) >= 3 }).Run)
	})
	t.Run("Recognize", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "Recognize", "xyz", chomp.Recognize(chomp.Tag("Hello")).Run)
	})
	t.Run("Consumed", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "Consumed", "xyz", chomp.Consumed(chomp.Tag("Hello")).Run)
	})
	t.Run("Eof", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "Eof", "remaining", chomp.Eof().Run)
	})
	t.Run("AllConsuming", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "AllConsuming", "Hello, World!", chomp.AllConsuming(chomp.Tag("Hello")).Run)
	})
	t.Run("Cut", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "Cut", "World", chomp.Cut(chomp.Tag("Hello")).Run)
	})
	t.Run("PeekNot", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "PeekNot", "Hello, World!", chomp.PeekNot(chomp.Tag("Hello")).Run)
	})

	// modifier.go
	t.Run("Map", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "Map", "xyz",
			chomp.Map(chomp.While(chomp.IsDigit), func(in string) int { return len(in) }).Run)
	})
	t.Run("MapRes", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "MapRes", "xyz", chomp.MapRes(chomp.While(chomp.IsDigit), strconv.Atoi).Run)
	})
	t.Run("S", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "S", "xyz", chomp.S(chomp.Tag("Hello")).Run)
	})
	t.Run("I", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "I", "xyz",
			chomp.I(chomp.All(chomp.Tag("Hello"), chomp.Tag(", World")), 1).Run)
	})
	t.Run("I_IndexOutOfBounds", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "I", "Hello, World!",
			chomp.I(chomp.All(chomp.Tag("Hello"), chomp.Tag(", World!")), 5).Run)
	})
	t.Run("Peek", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "Peek", "World", chomp.Peek(chomp.Tag("Hello")).Run)
	})
	t.Run("Flatten", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "Flatten", "xyz", chomp.Flatten(chomp.Many(chomp.Parentheses())).Run)
	})

	// label.go
	t.Run("Label", func(t *testing.T) {
		t.Parallel()
		assertFailureNonConsuming(t, "Label", "xyz", chomp.Label("greeting", chomp.Tag("Hello")).Run)
	})
}

func TestContract_SuccessIsPrefix(t *testing.T) {
	t.Parallel()

	t.Run("Char", func(t *testing.T) {
		t.Parallel()
		assertSuccessIsPrefix(t, "Char", "xyz", chomp.Char('x').Run)
	})
	t.Run("AnyChar", func(t *testing.T) {
		t.Parallel()
		assertSuccessIsPrefix(t, "AnyChar", "abc", chomp.AnyChar().Run)
	})
	t.Run("Satisfy", func(t *testing.T) {
		t.Parallel()
		assertSuccessIsPrefix(t, "Satisfy", "Hello", chomp.Satisfy(unicode.IsUpper).Run)
	})
	t.Run("Tag", func(t *testing.T) {
		t.Parallel()
		assertSuccessIsPrefix(t, "Tag", "Hello, World!", chomp.Tag("Hello").Run)
	})
	t.Run("TagNoCase", func(t *testing.T) {
		t.Parallel()
		assertSuccessIsPrefix(t, "TagNoCase", "HELLO, World!", chomp.TagNoCase("hello").Run)
	})
	t.Run("IsA", func(t *testing.T) {
		t.Parallel()
		assertSuccessIsPrefix(t, "IsA", "Hello, World!", chomp.IsA("eH").Run)
	})
	t.Run("IsNot", func(t *testing.T) {
		t.Parallel()
		assertSuccessIsPrefix(t, "IsNot", "Hello, World!", chomp.IsNot("ol").Run)
	})
	t.Run("OneOf", func(t *testing.T) {
		t.Parallel()
		assertSuccessIsPrefix(t, "OneOf", "Hello, World!", chomp.OneOf("!,eH").Run)
	})
	t.Run("NoneOf", func(t *testing.T) {
		t.Parallel()
		assertSuccessIsPrefix(t, "NoneOf", "Hello, World!", chomp.NoneOf("loWrd!e").Run)
	})
	t.Run("Until", func(t *testing.T) {
		t.Parallel()
		assertSuccessIsPrefix(t, "Until", "Hello, World!", chomp.Until("World").Run)
	})
	t.Run("Take", func(t *testing.T) {
		t.Parallel()
		assertSuccessIsPrefix(t, "Take", "Hello, World!", chomp.Take(5).Run)
	})
	t.Run("TakeUntil1", func(t *testing.T) {
		t.Parallel()
		assertSuccessIsPrefix(t, "TakeUntil1", "Hello, World!", chomp.TakeUntil1(",").Run)
	})
	t.Run("Escaped", func(t *testing.T) {
		t.Parallel()
		assertSuccessIsPrefix(t, "Escaped", `Hello\"World`,
			chomp.Escaped(chomp.While(chomp.IsLetter), '\\', chomp.OneOf(`"n\`)).Run)
	})
	t.Run("While", func(t *testing.T) {
		t.Parallel()
		assertSuccessIsPrefix(t, "While", "Hello, World!", chomp.While(chomp.IsLetter).Run)
	})
	t.Run("WhileN", func(t *testing.T) {
		t.Parallel()
		assertSuccessIsPrefix(t, "WhileN", "123abc", chomp.WhileN(chomp.IsDigit, 2).Run)
	})
	t.Run("WhileNM", func(t *testing.T) {
		t.Parallel()
		assertSuccessIsPrefix(t, "WhileNM", "abc123", chomp.WhileNM(chomp.IsLetter, 1, 5).Run)
	})
	t.Run("WhileNot", func(t *testing.T) {
		t.Parallel()
		assertSuccessIsPrefix(t, "WhileNot", "abc123", chomp.WhileNot(chomp.IsDigit).Run)
	})
	t.Run("WhileNotN", func(t *testing.T) {
		t.Parallel()
		assertSuccessIsPrefix(t, "WhileNotN", "abc123", chomp.WhileNotN(chomp.IsDigit, 1).Run)
	})
	t.Run("WhileNotNM", func(t *testing.T) {
		t.Parallel()
		assertSuccessIsPrefix(t, "WhileNotNM", "2024abcd", chomp.WhileNotNM(chomp.IsLetter, 1, 8).Run)
	})
	t.Run("Alpha", func(t *testing.T) {
		t.Parallel()
		assertSuccessIsPrefix(t, "Alpha", "Hello123", chomp.Alpha().Run)
	})
	t.Run("Digit", func(t *testing.T) {
		t.Parallel()
		assertSuccessIsPrefix(t, "Digit", "123abc", chomp.Digit().Run)
	})
	t.Run("Newline", func(t *testing.T) {
		t.Parallel()
		assertSuccessIsPrefix(t, "Newline", "\nHello", chomp.Newline().Run)
	})
	t.Run("Tab", func(t *testing.T) {
		t.Parallel()
		assertSuccessIsPrefix(t, "Tab", "\tHello", chomp.Tab().Run)
	})
	t.Run("NotLineEnding", func(t *testing.T) {
		t.Parallel()
		assertSuccessIsPrefix(t, "NotLineEnding", "Hello, World!\nNext line", chomp.NotLineEnding().Run)
	})
	t.Run("AnyDigit", func(t *testing.T) {
		t.Parallel()
		assertSuccessIsPrefix(t, "AnyDigit", "123", chomp.AnyDigit().Run)
	})
	t.Run("AnyLetter", func(t *testing.T) {
		t.Parallel()
		assertSuccessIsPrefix(t, "AnyLetter", "Hello", chomp.AnyLetter().Run)
	})
	t.Run("Crlf", func(t *testing.T) {
		t.Parallel()
		assertSuccessIsPrefix(t, "Crlf", "\r\nHello", chomp.Crlf().Run)
	})
	t.Run("LineEnding", func(t *testing.T) {
		t.Parallel()
		assertSuccessIsPrefix(t, "LineEnding", "\r\nHello", chomp.LineEnding().Run)
	})
	t.Run("Recognize", func(t *testing.T) {
		t.Parallel()
		assertSuccessIsPrefix(t, "Recognize", "Hello, World!",
			chomp.Recognize(chomp.SepPair(chomp.Alpha(), chomp.Tag(", "), chomp.Alpha())).Run)
	})
	t.Run("Rest", func(t *testing.T) {
		t.Parallel()
		assertSuccessIsPrefix(t, "Rest", "Hello, World!", chomp.Rest().Run)
	})
	t.Run("PeekNot", func(t *testing.T) {
		t.Parallel()
		assertSuccessIsPrefix(t, "PeekNot", "World!", chomp.PeekNot(chomp.Tag("Hello")).Run)
	})
	t.Run("Label", func(t *testing.T) {
		t.Parallel()
		assertSuccessIsPrefix(t, "Label", "Hello, World!", chomp.Label("greeting", chomp.Tag("Hello")).Run)
	})
}

func TestNegativeCountFailsGracefully(t *testing.T) {
	t.Parallel()

	t.Run("Take", func(t *testing.T) {
		t.Parallel()
		assertNegativeCountFailsGracefully(t, "Take", "abc", chomp.Take(-1).Run)
	})
	t.Run("WhileN", func(t *testing.T) {
		t.Parallel()
		assertNegativeCountFailsGracefully(t, "WhileN", "abc", chomp.WhileN(chomp.IsDigit, -1).Run)
	})
	t.Run("WhileNM", func(t *testing.T) {
		t.Parallel()
		assertNegativeCountFailsGracefully(t, "WhileNM", "abc", chomp.WhileNM(chomp.IsDigit, -1, 5).Run)
	})
	t.Run("WhileNotN", func(t *testing.T) {
		t.Parallel()
		assertNegativeCountFailsGracefully(t, "WhileNotN", "abc", chomp.WhileNotN(chomp.IsDigit, -1).Run)
	})
	t.Run("WhileNotNM", func(t *testing.T) {
		t.Parallel()
		assertNegativeCountFailsGracefully(t, "WhileNotNM", "abc", chomp.WhileNotNM(chomp.IsDigit, -1, 5).Run)
	})
	t.Run("Repeat", func(t *testing.T) {
		t.Parallel()
		assertNegativeCountFailsGracefully(t, "Repeat", "abc", chomp.Repeat(chomp.AnyLetter(), -1).Run)
	})
	t.Run("RepeatRange", func(t *testing.T) {
		t.Parallel()
		assertNegativeCountFailsGracefully(t, "RepeatRange", "abc", chomp.RepeatRange(chomp.AnyLetter(), -1, 5).Run)
	})
	t.Run("ManyN", func(t *testing.T) {
		t.Parallel()
		assertNegativeCountFailsGracefully(t, "ManyN", "abc", chomp.ManyN(chomp.AnyLetter(), -1).Run)
	})
}

func TestRepeatRangeInvertedBoundsFailsGracefully(t *testing.T) {
	t.Parallel()

	assertNegativeCountFailsGracefully(t, "RepeatRange", "abc", chomp.RepeatRange(chomp.AnyLetter(), 8, 1).Run)
}
