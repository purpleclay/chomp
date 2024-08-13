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
		tt := tt
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
