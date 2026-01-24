package chomp_test

import (
	"testing"

	"github.com/purpleclay/chomp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
			rem, ext, err := chomp.Char(tt.char)(tt.input)

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
			rem, ext, err := chomp.AnyChar()(tt.input)

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
			rem, ext, err := chomp.Satisfy(tt.pred)(tt.input)

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
			rem, tag, err := chomp.Tag(tt.tag)(tt.input)

			require.NoError(t, err)
			assert.Equal(t, tt.rem, rem)
			assert.Equal(t, tt.tag, tag)
		})
	}
}

func TestAny(t *testing.T) {
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
			rem, ext, err := chomp.Any(tt.any)(tt.input)

			require.NoError(t, err)
			assert.Equal(t, tt.rem, rem)
			assert.Equal(t, tt.ext, ext)
		})
	}
}

func TestNot(t *testing.T) {
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
			rem, ext, err := chomp.Not(tt.not)(tt.input)

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
			rem, ext, err := chomp.Until(tt.until)(tt.input)

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
			rem, ext, err := chomp.OneOf(tt.oneOf)(tt.input)

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
			rem, ext, err := chomp.NoneOf(tt.noneOf)(tt.input)

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
			rem, ext, err := chomp.TagNoCase(tt.tag)(tt.input)

			require.NoError(t, err)
			assert.Equal(t, tt.rem, rem)
			assert.Equal(t, tt.ext, ext)
		})
	}
}

func TestTake(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		n     uint
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
			rem, ext, err := chomp.Take(tt.n)(tt.input)

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
			rem, ext, err := chomp.TakeUntil1(tt.until)(tt.input)

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
			)(tt.input)

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
			)(tt.input)

			require.NoError(t, err)
			assert.Equal(t, tt.rem, rem)
			assert.Equal(t, tt.ext, ext)
		})
	}
}

func TestEscapedTransform(t *testing.T) {
	t.Parallel()

	transform := func(s string) (string, string, error) {
		if s == "" {
			return s, "", chomp.CombinatorParseError{Text: s, Type: "transform"}
		}
		switch s[0] {
		case 'n':
			return s[1:], "\n", nil
		case '"':
			return s[1:], "\"", nil
		case '\\':
			return s[1:], "\\", nil
		}
		return s, "", chomp.CombinatorParseError{Text: s, Type: "transform"}
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
			)(tt.input)

			require.NoError(t, err)
			assert.Equal(t, tt.rem, rem)
			assert.Equal(t, tt.ext, ext)
		})
	}
}

func TestEscapedTransformUnicode(t *testing.T) {
	t.Parallel()

	transform := func(s string) (string, string, error) {
		if s == "" {
			return s, "", chomp.CombinatorParseError{Text: s, Type: "transform"}
		}
		runes := []rune(s)
		switch runes[0] {
		case 'は':
			return string(runes[1:]), "【", nil
		case 'に':
			return string(runes[1:]), "】", nil
		}
		return s, "", chomp.CombinatorParseError{Text: s, Type: "transform"}
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
			)(tt.input)

			require.NoError(t, err)
			assert.Equal(t, tt.rem, rem)
			assert.Equal(t, tt.ext, ext)
		})
	}
}

func TestCombinatorError(t *testing.T) {
	t.Parallel()

	_, _, err := chomp.OneOf("!h")("Happy Monday")

	assert.EqualError(t, err, "(one_of) combinator failed to parse text 'Happy Monday' with input '!h'")
}

func TestParserCombinatorError(t *testing.T) {
	t.Parallel()

	_, _, err := chomp.All(
		chomp.Tag("the legend of batman"),
		chomp.Tag(":"),
		chomp.Tag("marvel"))("the legend of batman:dc:9781801260336:£19.99")

	assert.EqualError(t, err, "(all) parser failed. (tag) combinator failed to parse text 'dc:9781801260336:£19.99' with input 'marvel'")
}
