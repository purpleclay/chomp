package chomp_test

import (
	"testing"

	"github.com/purpleclay/chomp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
		tt := tt
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
		tt := tt
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
		tt := tt
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
		tt := tt
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
		tt := tt
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			rem, ext, err := chomp.NoneOf(tt.noneOf)(tt.input)

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
