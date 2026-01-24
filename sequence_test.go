package chomp_test

import (
	"testing"

	"github.com/purpleclay/chomp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPair(t *testing.T) {
	t.Parallel()

	rem, ext, err := chomp.Pair(chomp.Tag("Hello,"), chomp.Tag(" World"))("Hello, World!")

	require.NoError(t, err)
	assert.Equal(t, "!", rem)
	require.Len(t, ext, 2)
	assert.Equal(t, "Hello,", ext[0])
	assert.Equal(t, " World", ext[1])
}

func TestSepPair(t *testing.T) {
	t.Parallel()

	rem, ext, err := chomp.SepPair(chomp.Tag("Hello"), chomp.Tag(", "), chomp.Tag("World"))("Hello, World!")

	require.NoError(t, err)
	assert.Equal(t, "!", rem)
	require.Len(t, ext, 2)
	assert.Equal(t, "Hello", ext[0])
	assert.Equal(t, "World", ext[1])
}

func TestRepeat(t *testing.T) {
	t.Parallel()

	rem, ext, err := chomp.Repeat(chomp.QuoteDouble(), 3)(`"Batman""ジョーカー""Two Face""ベイン"`)

	require.NoError(t, err)
	assert.Equal(t, `"ベイン"`, rem)
	require.Len(t, ext, 3)
	assert.Equal(t, "Batman", ext[0])
	assert.Equal(t, "ジョーカー", ext[1])
	assert.Equal(t, "Two Face", ext[2])
}

func TestRepeatRange(t *testing.T) {
	t.Parallel()

	rem, ext, err := chomp.RepeatRange(
		chomp.Pair(chomp.Until(","), chomp.Opt(chomp.Tag(","))), 1, 3)("Batman,Joker,")

	require.NoError(t, err)
	assert.Empty(t, rem)
	require.Len(t, ext, 4)
	assert.Equal(t, "Batman", ext[0])
	assert.Equal(t, "Joker", ext[2])
}

func TestDelimited(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		delim []string
		rem   string
		ext   string
	}{
		{
			name:  "Ascii",
			input: "#Hello and Good Morning@",
			delim: []string{"#", "Hello and Good Morning", "@"},
			rem:   "",
			ext:   "Hello and Good Morning",
		},
		{
			name:  "Unicode",
			input: "┃こんにちは、おはよう║",
			delim: []string{"┃", "こんにちは、おはよう", "║"},
			rem:   "",
			ext:   "こんにちは、おはよう",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			rem, ext, err := chomp.Delimited(
				chomp.Tag(tt.delim[0]),
				chomp.Tag(tt.delim[1]),
				chomp.Tag(tt.delim[2]))(tt.input)

			require.NoError(t, err)
			assert.Equal(t, tt.rem, rem)
			assert.Equal(t, tt.ext, ext)
		})
	}
}

func TestFirst(t *testing.T) {
	t.Parallel()

	rem, ext, err := chomp.First(chomp.Tag("Light"), chomp.Tag("Dark"))("Dark Knight")

	require.NoError(t, err)
	assert.Equal(t, " Knight", rem)
	assert.Equal(t, "Dark", ext)
}

func TestAll(t *testing.T) {
	t.Parallel()

	rem, ext, err := chomp.All(
		chomp.QuoteDouble(),
		chomp.Until("("),
		chomp.Parentheses())(`"Hello and Good Morning" (こんにちは、おはよう)`)

	require.NoError(t, err)
	assert.Empty(t, rem)
	require.Len(t, ext, 3)
	assert.Equal(t, "Hello and Good Morning", ext[0])
	assert.Equal(t, " ", ext[1])
	assert.Equal(t, "こんにちは、おはよう", ext[2])
}

func TestMany(t *testing.T) {
	t.Parallel()

	rem, ext, err := chomp.Many(chomp.OneOf("はんにこち"))("こんにちは、おはよう")

	require.NoError(t, err)
	assert.Equal(t, "、おはよう", rem)
	require.Len(t, ext, 5)
	assert.Equal(t, "こ", ext[0])
	assert.Equal(t, "ん", ext[1])
	assert.Equal(t, "に", ext[2])
	assert.Equal(t, "ち", ext[3])
	assert.Equal(t, "は", ext[4])
}

func TestManyNoMatches(t *testing.T) {
	t.Parallel()

	_, _, err := chomp.Many(chomp.OneOf("eHl"))("Good Morning")

	require.EqualError(t, err, "(many_n) parser failed [count: 0 min: 1]. (one_of) combinator failed to parse text 'Good Morning' with input 'eHl'")
}

func TestManyN(t *testing.T) {
	t.Parallel()

	rem, ext, err := chomp.ManyN(chomp.OneOf("eHl"), 2)("Hello and Good Morning")

	require.NoError(t, err)
	assert.Equal(t, "o and Good Morning", rem)
	require.Len(t, ext, 4)
	assert.Equal(t, "H", ext[0])
	assert.Equal(t, "e", ext[1])
	assert.Equal(t, "l", ext[2])
	assert.Equal(t, "l", ext[3])
}

func TestManyNZeroMatches(t *testing.T) {
	t.Parallel()

	rem, ext, err := chomp.ManyN(chomp.OneOf("eHl"), 0)("Good Morning")

	require.NoError(t, err)
	assert.Equal(t, "Good Morning", rem)
	assert.Empty(t, ext)
}

func TestPrefixed(t *testing.T) {
	t.Parallel()

	rem, ext, err := chomp.Prefixed(chomp.Tag("Hello"), chomp.Tag(`"`))(`"Hello, World"`)

	require.NoError(t, err)
	assert.Equal(t, `, World"`, rem)
	assert.Equal(t, "Hello", ext)
}

func TestSuffixed(t *testing.T) {
	t.Parallel()

	rem, ext, err := chomp.Suffixed(chomp.Tag("Hello"), chomp.Tag(","))("Hello, World")

	require.NoError(t, err)
	assert.Equal(t, " World", rem)
	assert.Equal(t, "Hello", ext)
}

func TestSeparatedList(t *testing.T) {
	t.Parallel()

	rem, ext, err := chomp.SeparatedList(chomp.Alpha(), chomp.Tag(","))("apple,banana,cherry,")

	require.NoError(t, err)
	assert.Equal(t, ",", rem)
	require.Len(t, ext, 3)
	assert.Equal(t, "apple", ext[0])
	assert.Equal(t, "banana", ext[1])
	assert.Equal(t, "cherry", ext[2])
}

func TestSeparatedListSingleElement(t *testing.T) {
	t.Parallel()

	rem, ext, err := chomp.SeparatedList(chomp.Alpha(), chomp.Tag(","))("apple")

	require.NoError(t, err)
	assert.Equal(t, "", rem)
	require.Len(t, ext, 1)
	assert.Equal(t, "apple", ext[0])
}

func TestSeparatedListNoMatch(t *testing.T) {
	t.Parallel()

	_, _, err := chomp.SeparatedList(chomp.Alpha(), chomp.Tag(","))("123")

	require.Error(t, err)
}

func TestSeparatedList0(t *testing.T) {
	t.Parallel()

	rem, ext, err := chomp.SeparatedList0(chomp.Alpha(), chomp.Tag(","))("apple,banana")

	require.NoError(t, err)
	assert.Equal(t, "", rem)
	require.Len(t, ext, 2)
	assert.Equal(t, "apple", ext[0])
	assert.Equal(t, "banana", ext[1])
}

func TestSeparatedList0NoMatch(t *testing.T) {
	t.Parallel()

	rem, ext, err := chomp.SeparatedList0(chomp.Alpha(), chomp.Tag(","))("123")

	require.NoError(t, err)
	assert.Equal(t, "123", rem)
	assert.Empty(t, ext)
}

func TestManyTill(t *testing.T) {
	t.Parallel()

	rem, ext, err := chomp.ManyTill(chomp.AnyChar(), chomp.Tag("END"))("abcEND rest")

	require.NoError(t, err)
	assert.Equal(t, " rest", rem)
	require.Len(t, ext, 3)
	assert.Equal(t, "a", ext[0])
	assert.Equal(t, "b", ext[1])
	assert.Equal(t, "c", ext[2])
}

func TestManyTillNoElementsBeforeTerminator(t *testing.T) {
	t.Parallel()

	_, _, err := chomp.ManyTill(chomp.AnyChar(), chomp.Tag("END"))("END")

	require.Error(t, err)
}

func TestManyTill0(t *testing.T) {
	t.Parallel()

	rem, ext, err := chomp.ManyTill0(chomp.AnyChar(), chomp.Tag("END"))("abEND")

	require.NoError(t, err)
	assert.Equal(t, "", rem)
	require.Len(t, ext, 2)
	assert.Equal(t, "a", ext[0])
	assert.Equal(t, "b", ext[1])
}

func TestManyTill0EmptyBeforeTerminator(t *testing.T) {
	t.Parallel()

	rem, ext, err := chomp.ManyTill0(chomp.AnyChar(), chomp.Tag("END"))("END")

	require.NoError(t, err)
	assert.Equal(t, "", rem)
	assert.Empty(t, ext)
}

func TestFoldMany(t *testing.T) {
	t.Parallel()

	rem, ext, err := chomp.FoldMany(chomp.AnyDigit(), 0, func(acc int, val string) int {
		return acc + int(val[0]-'0')
	})("123abc")

	require.NoError(t, err)
	assert.Equal(t, "abc", rem)
	assert.Equal(t, 6, ext)
}

func TestFoldManyNoMatch(t *testing.T) {
	t.Parallel()

	_, _, err := chomp.FoldMany(chomp.AnyDigit(), 0, func(acc int, val string) int {
		return acc + int(val[0]-'0')
	})("abc")

	require.Error(t, err)
}

func TestFoldMany0(t *testing.T) {
	t.Parallel()

	rem, ext, err := chomp.FoldMany0(chomp.AnyDigit(), 0, func(acc int, val string) int {
		return acc + int(val[0]-'0')
	})("12abc")

	require.NoError(t, err)
	assert.Equal(t, "abc", rem)
	assert.Equal(t, 3, ext)
}

func TestFoldMany0NoMatch(t *testing.T) {
	t.Parallel()

	rem, ext, err := chomp.FoldMany0(chomp.AnyDigit(), 0, func(acc int, val string) int {
		return acc + int(val[0]-'0')
	})("abc")

	require.NoError(t, err)
	assert.Equal(t, "abc", rem)
	assert.Equal(t, 0, ext)
}

func TestManyCount(t *testing.T) {
	t.Parallel()

	rem, count, err := chomp.ManyCount(chomp.AnyLetter())("abc123")

	require.NoError(t, err)
	assert.Equal(t, "123", rem)
	assert.Equal(t, uint(3), count)
}

func TestManyCountNoMatch(t *testing.T) {
	t.Parallel()

	_, _, err := chomp.ManyCount(chomp.AnyLetter())("123")

	require.Error(t, err)
}

func TestManyCount0(t *testing.T) {
	t.Parallel()

	rem, count, err := chomp.ManyCount0(chomp.AnyLetter())("ab123")

	require.NoError(t, err)
	assert.Equal(t, "123", rem)
	assert.Equal(t, uint(2), count)
}

func TestManyCount0NoMatch(t *testing.T) {
	t.Parallel()

	rem, count, err := chomp.ManyCount0(chomp.AnyLetter())("123")

	require.NoError(t, err)
	assert.Equal(t, "123", rem)
	assert.Equal(t, uint(0), count)
}

func TestLengthCount(t *testing.T) {
	t.Parallel()

	lengthParser := chomp.Map(chomp.AnyDigit(), func(s string) uint {
		return uint(s[0] - '0')
	})

	rem, ext, err := chomp.LengthCount(lengthParser, chomp.AnyLetter())("3abcdef")

	require.NoError(t, err)
	assert.Equal(t, "def", rem)
	require.Len(t, ext, 3)
	assert.Equal(t, "a", ext[0])
	assert.Equal(t, "b", ext[1])
	assert.Equal(t, "c", ext[2])
}

func TestFill(t *testing.T) {
	t.Parallel()

	rem, ext, err := chomp.Fill(chomp.AnyLetter(), 3)("abcdef")

	require.NoError(t, err)
	assert.Equal(t, "def", rem)
	require.Len(t, ext, 3)
	assert.Equal(t, "a", ext[0])
	assert.Equal(t, "b", ext[1])
	assert.Equal(t, "c", ext[2])
}

func TestFillNotEnoughElements(t *testing.T) {
	t.Parallel()

	_, _, err := chomp.Fill(chomp.AnyLetter(), 5)("abc")

	require.Error(t, err)
}
