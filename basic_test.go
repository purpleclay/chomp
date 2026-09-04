package chomp_test

import (
	"errors"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/purpleclay/chomp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// zeroWidthNonEmptyExt always succeeds without consuming any input, but
// reports a non-empty ext - violating the assumption that a non-empty ext
// implies real progress. Used to prove Escaped/EscapedTransform measure
// progress from the remainder, not ext, and so never loop forever on it.
func zeroWidthNonEmptyExt(s chomp.State) (chomp.State, string, error) {
	return s, "ok", nil
}

// consumesButEmptyExt consumes one byte but reports an empty ext -
// violating the assumption that ext length reflects what was consumed.
// Used to prove Escaped/EscapedTransform still credit real progress even
// when ext is empty.
func consumesButEmptyExt(s chomp.State) (chomp.State, string, error) {
	if s.Rest() == "" {
		return s, "", errors.New("empty")
	}
	return s.Advance(1), "", nil
}

func TestChar(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		char  rune
		rem   string
		ext   string
	}{
		{
			name:  "Ascii",
			input: ",,rest",
			char:  ',',
			rem:   ",rest",
			ext:   ",",
		},
		{
			name:  "Unicode",
			input: "★星空",
			char:  '★',
			rem:   "星空",
			ext:   "★",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			rem, ext, err := chomp.Char(tt.char).Run(tt.input)

			require.NoError(t, err)
			assert.Equal(t, tt.rem, rem)
			assert.Equal(t, tt.ext, ext)
		})
	}
}

func TestAnyChar(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		rem   string
		ext   string
	}{
		{
			name:  "Ascii",
			input: "Hello",
			rem:   "ello",
			ext:   "H",
		},
		{
			name:  "Unicode",
			input: "こんにちは",
			rem:   "んにちは",
			ext:   "こ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			rem, ext, err := chomp.AnyChar().Run(tt.input)

			require.NoError(t, err)
			assert.Equal(t, tt.rem, rem)
			assert.Equal(t, tt.ext, ext)
		})
	}
}

func TestSatisfy(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		pred  func(rune) bool
		rem   string
		ext   string
	}{
		{
			name:  "UppercaseAscii",
			input: "Hello",
			pred:  func(r rune) bool { return r >= 'A' && r <= 'Z' },
			rem:   "ello",
			ext:   "H",
		},
		{
			name:  "UnicodeHiragana",
			input: "あいうえお",
			pred:  func(r rune) bool { return r >= 'あ' && r <= 'ん' },
			rem:   "いうえお",
			ext:   "あ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			rem, ext, err := chomp.Satisfy(tt.pred).Run(tt.input)

			require.NoError(t, err)
			assert.Equal(t, tt.rem, rem)
			assert.Equal(t, tt.ext, ext)
		})
	}
}

func TestTag(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		tag   string
		rem   string
	}{
		{
			name:  "Ascii",
			input: "hello and good morning",
			tag:   "hello",
			rem:   " and good morning",
		},
		{
			name:  "Unicode",
			input: "こんにちは、おはよう",
			tag:   "こんにちは",
			rem:   "、おはよう",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			rem, tag, err := chomp.Tag(tt.tag).Run(tt.input)

			require.NoError(t, err)
			assert.Equal(t, tt.rem, rem)
			assert.Equal(t, tt.tag, tag)
		})
	}
}

func TestIsA(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		any   string
		rem   string
		ext   string
	}{
		{
			name:  "Ascii",
			input: "dark knight",
			any:   "krad ",
			rem:   "night",
			ext:   "dark k",
		},
		{
			name:  "Unicode",
			input: "ダークナイト",
			any:   "ダー",
			rem:   "クナイト",
			ext:   "ダー",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			rem, ext, err := chomp.IsA(tt.any).Run(tt.input)

			require.NoError(t, err)
			assert.Equal(t, tt.rem, rem)
			assert.Equal(t, tt.ext, ext)
		})
	}
}

func TestIsNot(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		not   string
		rem   string
		ext   string
	}{
		{
			name:  "Ascii",
			input: "dark knight",
			not:   "tighn",
			rem:   "night",
			ext:   "dark k",
		},
		{
			name:  "Unicode",
			input: "ダークナイト",
			not:   "トイ",
			rem:   "イト",
			ext:   "ダークナ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			rem, ext, err := chomp.IsNot(tt.not).Run(tt.input)

			require.NoError(t, err)
			assert.Equal(t, tt.rem, rem)
			assert.Equal(t, tt.ext, ext)
		})
	}
}

func TestUntil(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		until string
		input string
		rem   string
		ext   string
	}{
		{
			name:  "Ascii",
			until: "jumps",
			input: "the quick brown fox jumps over the lazy dog",
			rem:   "jumps over the lazy dog",
			ext:   "the quick brown fox ",
		},
		{
			name:  "Unicode",
			until: "の",
			input: "素早い茶色のキツネが怠惰な犬を飛び越える",
			rem:   "のキツネが怠惰な犬を飛び越える",
			ext:   "素早い茶色",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			rem, ext, err := chomp.Until(tt.until).Run(tt.input)

			require.NoError(t, err)
			assert.Equal(t, tt.rem, rem)
			assert.Equal(t, tt.ext, ext)
		})
	}
}

func TestOneOf(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		oneOf string
		input string
		rem   string
		ext   string
	}{
		{
			name:  "Ascii",
			oneOf: "!,eH",
			input: "Hello, World!",
			rem:   "ello, World!",
			ext:   "H",
		},
		{
			name:  "Unicode",
			oneOf: "はおうこ、",
			input: "こんにちは、おはよう",
			rem:   "んにちは、おはよう",
			ext:   "こ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			rem, ext, err := chomp.OneOf(tt.oneOf).Run(tt.input)

			require.NoError(t, err)
			assert.Equal(t, tt.rem, rem)
			assert.Equal(t, tt.ext, ext)
		})
	}
}

func TestNoneOf(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		noneOf string
		input  string
		rem    string
		ext    string
	}{
		{
			name:   "Ascii",
			noneOf: "eqzygoqui",
			input:  "the quick brown fox jumps over the lazy dog",
			rem:    "he quick brown fox jumps over the lazy dog",
			ext:    "t",
		},
		{
			name:   "Unicode",
			noneOf: "が早越ネをのる",
			input:  "素早い茶色のキツネが怠惰な犬を飛び越える",
			rem:    "早い茶色のキツネが怠惰な犬を飛び越える",
			ext:    "素",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			rem, ext, err := chomp.NoneOf(tt.noneOf).Run(tt.input)

			require.NoError(t, err)
			assert.Equal(t, tt.rem, rem)
			assert.Equal(t, tt.ext, ext)
		})
	}
}

func TestTagNoCase(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		tag   string
		rem   string
		ext   string
	}{
		{
			name:  "MatchExact",
			input: "hello and good morning",
			tag:   "hello",
			rem:   " and good morning",
			ext:   "hello",
		},
		{
			name:  "MatchUppercase",
			input: "HELLO and good morning",
			tag:   "hello",
			rem:   " and good morning",
			ext:   "HELLO",
		},
		{
			name:  "MatchMixedCase",
			input: "HeLLo and good morning",
			tag:   "hello",
			rem:   " and good morning",
			ext:   "HeLLo",
		},
		{
			name:  "Unicode",
			input: "ΓΕΙΑ and good morning",
			tag:   "γεια",
			rem:   " and good morning",
			ext:   "ΓΕΙΑ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			rem, ext, err := chomp.TagNoCase(tt.tag).Run(tt.input)

			require.NoError(t, err)
			assert.Equal(t, tt.rem, rem)
			assert.Equal(t, tt.ext, ext)
		})
	}
}

func TestTagNoCaseFoldLengthMismatch(t *testing.T) {
	t.Parallel()

	const (
		kelvinSign = "K" // KELVIN SIGN, 3 bytes, folds with 'k'/'K'
		longS      = "ſ" // LATIN SMALL LETTER LONG S, 2 bytes, folds with 's'/'S'
	)

	tests := []struct {
		name  string
		input string
		tag   string
		rem   string
		ext   string
	}{
		{
			name:  "KelvinSignInInputShorterTag",
			input: kelvinSign + "elvin",
			tag:   "k",
			rem:   "elvin",
			ext:   kelvinSign,
		},
		{
			name:  "KelvinSignInTagShorterInput",
			input: "Kelvin",
			tag:   kelvinSign,
			rem:   "elvin",
			ext:   "K",
		},
		{
			name:  "LongSInInputShorterTag",
			input: longS + "ong",
			tag:   "s",
			rem:   "ong",
			ext:   longS,
		},
		{
			name:  "LongSInTagShorterInput",
			input: "Song",
			tag:   longS,
			rem:   "ong",
			ext:   "S",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			rem, ext, err := chomp.TagNoCase(tt.tag).Run(tt.input)

			require.NoError(t, err)
			assert.Equal(t, tt.rem, rem)
			assert.Equal(t, tt.ext, ext)
		})
	}
}

func TestTagNoCaseRejectsMalformedUTF8AsReplacementChar(t *testing.T) {
	t.Parallel()

	rem, ext, err := chomp.TagNoCase("�").Run("\xff")

	require.Error(t, err)
	assert.Equal(t, "\xff", rem)
	assert.Equal(t, "", ext)
}

func TestTagNoCaseMatchesGenuineReplacementChar(t *testing.T) {
	t.Parallel()

	rem, ext, err := chomp.TagNoCase("�").Run("�x")

	require.NoError(t, err)
	assert.Equal(t, "x", rem)
	assert.Equal(t, "�", ext)
}

func TestTagNoCaseRejectsMalformedTagAsReplacementChar(t *testing.T) {
	t.Parallel()

	rem, ext, err := chomp.TagNoCase("\xff").Run("�x")

	require.Error(t, err)
	assert.Equal(t, "�x", rem)
	assert.Equal(t, "", ext)
}

// referenceTagNoCase is an independent rune-by-rune reimplementation of
// TagNoCase's matching logic, deferring the actual fold-equality check to
// strings.EqualFold (applied per rune) rather than chomp's own foldEqual.
// Used as the oracle for FuzzTagNoCase.
func referenceTagNoCase(tag, input string) (rem, ext string, failed bool) {
	pos, tagPos := 0, 0
	for tagPos < len(tag) {
		want, wantSize := utf8.DecodeRuneInString(tag[tagPos:])
		if pos >= len(input) || (want == utf8.RuneError && wantSize <= 1) {
			return "", "", true
		}

		got, size := utf8.DecodeRuneInString(input[pos:])
		if (got == utf8.RuneError && size <= 1) || !strings.EqualFold(string(got), string(want)) {
			return "", "", true
		}
		pos += size
		tagPos += wantSize
	}

	return input[pos:], input[:pos], false
}

func FuzzTagNoCase(f *testing.F) {
	f.Add("k", "Kelvin sign: Kelvin")
	f.Add("K", "Kelvin")
	f.Add("s", "long s: ſong")
	f.Add("ſ", "Song")
	f.Add("hello", "HELLO, World!")
	f.Add("γεια", "ΓΕΙΑ and good morning")
	f.Add("k", "")
	f.Add("hello", "hell")
	f.Add("�", "\xff")
	f.Add("�", "�x")
	f.Add("\xff", "�x")

	f.Fuzz(func(t *testing.T, tag, input string) {
		if tag == "" {
			t.Skip()
		}

		rem, ext, err := chomp.TagNoCase(tag).Run(input)
		wantRem, wantExt, wantFailed := referenceTagNoCase(tag, input)

		if wantFailed {
			require.Error(t, err)
			assert.Equal(t, input, rem)
			assert.Equal(t, "", ext)
			return
		}

		require.NoError(t, err)
		assert.Equal(t, wantRem, rem)
		assert.Equal(t, wantExt, ext)
	})
}

func TestTake(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		n     int
		rem   string
		ext   string
	}{
		{
			name:  "Ascii",
			input: "Hello, World!",
			n:     5,
			rem:   ", World!",
			ext:   "Hello",
		},
		{
			name:  "Unicode",
			input: "こんにちは、おはよう",
			n:     5,
			rem:   "、おはよう",
			ext:   "こんにちは",
		},
		{
			name:  "EntireInput",
			input: "Hello",
			n:     5,
			rem:   "",
			ext:   "Hello",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			rem, ext, err := chomp.Take(tt.n).Run(tt.input)

			require.NoError(t, err)
			assert.Equal(t, tt.rem, rem)
			assert.Equal(t, tt.ext, ext)
		})
	}
}

func TestTakeUntil1(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		until string
		input string
		rem   string
		ext   string
	}{
		{
			name:  "Ascii",
			until: ",",
			input: "Hello, World!",
			rem:   ", World!",
			ext:   "Hello",
		},
		{
			name:  "Unicode",
			until: "、",
			input: "こんにちは、おはよう",
			rem:   "、おはよう",
			ext:   "こんにちは",
		},
		{
			name:  "MultiCharDelimiter",
			until: "World",
			input: "Hello, World!",
			rem:   "World!",
			ext:   "Hello, ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			rem, ext, err := chomp.TakeUntil1(tt.until).Run(tt.input)

			require.NoError(t, err)
			assert.Equal(t, tt.rem, rem)
			assert.Equal(t, tt.ext, ext)
		})
	}
}

func TestEscaped(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		rem   string
		ext   string
	}{
		{
			name:  "WithEscapedQuote",
			input: `Hello\"World`,
			rem:   "",
			ext:   `Hello\"World`,
		},
		{
			name:  "WithEscapedBackslash",
			input: `Hello\\World`,
			rem:   "",
			ext:   `Hello\\World`,
		},
		{
			name:  "WithEscapedNewline",
			input: `Hello\nWorld`,
			rem:   "",
			ext:   `Hello\nWorld`,
		},
		{
			name:  "NoEscape",
			input: "HelloWorld",
			rem:   "",
			ext:   "HelloWorld",
		},
		{
			name:  "MultipleEscapes",
			input: `Hello\"World\nTest`,
			rem:   "",
			ext:   `Hello\"World\nTest`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			rem, ext, err := chomp.Escaped(
				chomp.While(chomp.IsLetter),
				'\\',
				chomp.OneOf(`"n\`),
			).Run(tt.input)

			require.NoError(t, err)
			assert.Equal(t, tt.rem, rem)
			assert.Equal(t, tt.ext, ext)
		})
	}
}

func TestEscapedUnicode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		rem   string
		ext   string
	}{
		{
			name:  "UnicodeWithEscape",
			input: "こんにちは★は世界",
			rem:   "",
			ext:   "こんにちは★は世界",
		},
		{
			name:  "UnicodeMultipleEscapes",
			input: "始まり★に中間★は終わり",
			rem:   "",
			ext:   "始まり★に中間★は終わり",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			rem, ext, err := chomp.Escaped(
				chomp.While(chomp.IsLetter),
				'★',
				chomp.OneOf("はに"),
			).Run(tt.input)

			require.NoError(t, err)
			assert.Equal(t, tt.rem, rem)
			assert.Equal(t, tt.ext, ext)
		})
	}
}

func TestEscapedTransform(t *testing.T) {
	t.Parallel()

	transform := func(s chomp.State) (chomp.State, string, error) {
		rest := s.Rest()
		if rest == "" {
			return s, "", chomp.CombinatorParseError{State: s, Type: "transform"}
		}
		switch rest[0] {
		case 'n':
			return s.Advance(1), "\n", nil
		case '"':
			return s.Advance(1), "\"", nil
		case '\\':
			return s.Advance(1), "\\", nil
		}
		return s, "", chomp.CombinatorParseError{State: s, Type: "transform"}
	}

	tests := []struct {
		name  string
		input string
		rem   string
		ext   string
	}{
		{
			name:  "WithEscapedNewline",
			input: `Hello\nWorld`,
			rem:   "",
			ext:   "Hello\nWorld",
		},
		{
			name:  "WithEscapedQuote",
			input: `Hello\"World`,
			rem:   "",
			ext:   `Hello"World`,
		},
		{
			name:  "WithEscapedBackslash",
			input: `Hello\\World`,
			rem:   "",
			ext:   `Hello\World`,
		},
		{
			name:  "NoEscape",
			input: "HelloWorld",
			rem:   "",
			ext:   "HelloWorld",
		},
		{
			name:  "MultipleEscapes",
			input: `Hello\nWorld\"Test`,
			rem:   "",
			ext:   "Hello\nWorld\"Test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			rem, ext, err := chomp.EscapedTransform(
				chomp.While(chomp.IsLetter),
				'\\',
				transform,
			).Run(tt.input)

			require.NoError(t, err)
			assert.Equal(t, tt.rem, rem)
			assert.Equal(t, tt.ext, ext)
		})
	}
}

func TestEscapedTransformUnicode(t *testing.T) {
	t.Parallel()

	transform := func(s chomp.State) (chomp.State, string, error) {
		rest := s.Rest()
		if rest == "" {
			return s, "", chomp.CombinatorParseError{State: s, Type: "transform"}
		}
		r, size := utf8.DecodeRuneInString(rest)
		switch r {
		case 'は':
			return s.Advance(size), "【", nil
		case 'に':
			return s.Advance(size), "】", nil
		}
		return s, "", chomp.CombinatorParseError{State: s, Type: "transform"}
	}

	tests := []struct {
		name  string
		input string
		rem   string
		ext   string
	}{
		{
			name:  "UnicodeContent",
			input: "こんにちは★は世界",
			rem:   "",
			ext:   "こんにちは【世界",
		},
		{
			name:  "UnicodeMultipleEscapes",
			input: "こんにちは★は世界★に終",
			rem:   "",
			ext:   "こんにちは【世界】終",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			rem, ext, err := chomp.EscapedTransform(
				chomp.While(chomp.IsLetter),
				'★',
				transform,
			).Run(tt.input)

			require.NoError(t, err)
			assert.Equal(t, tt.rem, rem)
			assert.Equal(t, tt.ext, ext)
		})
	}
}

func TestEscapedTransformingNormal(t *testing.T) {
	t.Parallel()

	rem, ext, err := chomp.Escaped(chomp.Parentheses(), '\\', chomp.OneOf("n")).Run("(ab)(cd)")

	require.NoError(t, err)
	assert.Equal(t, "", rem)
	assert.Equal(t, "(ab)(cd)", ext)
}

func TestEscapedTransformingEscapable(t *testing.T) {
	t.Parallel()

	rem, ext, err := chomp.Escaped(chomp.Tag("a"), '\\', chomp.Parentheses()).Run(`a\(x)a`)

	require.NoError(t, err)
	assert.Equal(t, "", rem)
	assert.Equal(t, `a\(x)a`, ext)
}

func TestEscapedNeverStallsOnEmptyExt(t *testing.T) {
	t.Parallel()

	rem, ext, err := chomp.Escaped(consumesButEmptyExt, '\\', chomp.OneOf("n")).Run("abc")

	require.NoError(t, err)
	assert.Equal(t, "", rem)
	assert.Equal(t, "abc", ext)
}

func TestEscapedNeverLoopsOnZeroWidthNonEmptyExt(t *testing.T) {
	t.Parallel()

	withTimeout(t, hangTimeout, func() {
		_, _, err := chomp.Escaped(zeroWidthNonEmptyExt, '\\', chomp.OneOf("n")).Run("abc")
		require.Error(t, err)
	})
}

func TestEscapedTransformNeverStallsOnEmptyTransform(t *testing.T) {
	t.Parallel()

	rem, ext, err := chomp.EscapedTransform(chomp.While(chomp.IsLetter), '\\', consumesButEmptyExt).Run(`ab\xcd`)

	require.NoError(t, err)
	assert.Equal(t, "", rem)
	assert.Equal(t, "abcd", ext)
}

func TestEscapedTransformNeverLoopsOnZeroWidthNonEmptyExt(t *testing.T) {
	t.Parallel()

	withTimeout(t, hangTimeout, func() {
		transform := func(s chomp.State) (chomp.State, string, error) {
			return s, "", errors.New("boom")
		}
		_, _, err := chomp.EscapedTransform(zeroWidthNonEmptyExt, '\\', transform).Run("abc")
		require.Error(t, err)
	})
}

func TestCombinatorError(t *testing.T) {
	t.Parallel()

	_, _, err := chomp.OneOf("!h").Run("Happy Monday")

	assert.EqualError(t, err, `chomp: parse error at line 1, column 1 (offset 0): expected a character in "!h"`)
}

func TestParserCombinatorError(t *testing.T) {
	t.Parallel()

	_, _, err := chomp.All(
		chomp.Tag("the legend of batman"),
		chomp.Tag(":"),
		chomp.Tag("marvel")).Run("the legend of batman:dc:9781801260336:£19.99")

	assert.EqualError(t, err, `chomp: parse error at line 1, column 22 (offset 21): expected "marvel"`)
}

func TestEmptyPatternBehaviour(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		run     chomp.Combinator[string]
		input   string
		rem     string
		ext     string
		wantErr bool
	}{
		{
			name:  "Tag",
			run:   chomp.Tag(""),
			input: "abc",
			rem:   "abc",
			ext:   "",
		},
		{
			name:  "TagNoCase",
			run:   chomp.TagNoCase(""),
			input: "abc",
			rem:   "abc",
			ext:   "",
		},
		{
			name:  "Until",
			run:   chomp.Until(""),
			input: "abc",
			rem:   "abc",
			ext:   "",
		},
		{
			name:    "TakeUntil1",
			run:     chomp.TakeUntil1(""),
			input:   "abc",
			wantErr: true,
		},
		{
			name:    "IsA",
			run:     chomp.IsA(""),
			input:   "abc",
			wantErr: true,
		},
		{
			name:  "IsNot",
			run:   chomp.IsNot(""),
			input: "abc",
			rem:   "",
			ext:   "abc",
		},
		{
			name:    "OneOf",
			run:     chomp.OneOf(""),
			input:   "abc",
			wantErr: true,
		},
		{
			name:  "NoneOf",
			run:   chomp.NoneOf(""),
			input: "abc",
			rem:   "bc",
			ext:   "a",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			rem, ext, err := tt.run.Run(tt.input)

			if tt.wantErr {
				require.Error(t, err)
				assert.Equal(t, tt.input, rem)
				assert.Equal(t, "", ext)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.rem, rem)
			assert.Equal(t, tt.ext, ext)
		})
	}
}

func TestEmptyPatternZeroWidthCombinatorsDoNotHang(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		c    chomp.Combinator[string]
	}{
		{name: "Tag", c: chomp.Tag("")},
		{name: "TagNoCase", c: chomp.TagNoCase("")},
		{name: "Until", c: chomp.Until("")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			withTimeout(t, hangTimeout, func() {
				rem, ext, err := chomp.ManyN(tt.c, 0).Run("abc")
				require.NoError(t, err)
				assert.Equal(t, "abc", rem)
				assert.Empty(t, ext)
			})
		})
	}
}

func TestIsNotEmptySequenceOnEmptyInputFails(t *testing.T) {
	t.Parallel()

	rem, ext, err := chomp.IsNot("").Run("")

	require.Error(t, err)
	assert.Equal(t, "", rem)
	assert.Equal(t, "", ext)
}
