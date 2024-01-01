/*
Copyright (c) 2023 Purple Clay

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

package chomp_test

import (
	"testing"

	"github.com/purpleclay/chomp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTag(t *testing.T) {
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
			rem, tag, err := chomp.Tag(tt.tag)(tt.input)

			assert.Equal(t, tt.rem, rem)
			assert.Equal(t, tt.tag, tag)
			assert.NoError(t, err)
		})
	}
}

func TestCombinatorError(t *testing.T) {
	_, _, err := chomp.Tag("missing")("does not contain text")

	require.EqualError(t, err, "tag combinator failed to parse text using input 'missing'")
}

func TestAny(t *testing.T) {
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
			rem, ext, err := chomp.Any(tt.any)(tt.input)

			assert.Equal(t, tt.rem, rem)
			assert.Equal(t, tt.ext, ext)
			assert.NoError(t, err)
		})
	}
}

func TestNot(t *testing.T) {
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
			rem, ext, err := chomp.Not(tt.not)(tt.input)

			assert.Equal(t, tt.rem, rem)
			assert.Equal(t, tt.ext, ext)
			assert.NoError(t, err)
		})
	}
}

func TestUntil(t *testing.T) {
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
			rem, ext, err := chomp.Until(tt.until)(tt.input)

			assert.Equal(t, tt.rem, rem)
			assert.Equal(t, tt.ext, ext)
			assert.NoError(t, err)
		})
	}
}

func TestCrlf(t *testing.T) {
	tests := []struct {
		name  string
		input string
		rem   string
		ext   string
	}{
		{
			name:  "LF",
			input: "\nHello",
			rem:   "Hello",
			ext:   "\n",
		},
		{
			name:  "CRLF",
			input: "\r\nこんにちは",
			rem:   "こんにちは",
			ext:   "\r\n",
		},
		{
			name:  "LFOnly",
			input: "\n",
			rem:   "",
			ext:   "\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rem, ext, err := chomp.Crlf()(tt.input)

			assert.Equal(t, tt.rem, rem)
			assert.Equal(t, tt.ext, ext)
			assert.NoError(t, err)
		})
	}
}

func TestOneOf(t *testing.T) {
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
			rem, ext, err := chomp.OneOf(tt.oneOf)(tt.input)

			assert.Equal(t, tt.rem, rem)
			assert.Equal(t, tt.ext, ext)
			assert.NoError(t, err)
		})
	}
}

func TestNoneOf(t *testing.T) {
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
			rem, ext, err := chomp.NoneOf(tt.noneOf)(tt.input)

			assert.Equal(t, tt.rem, rem)
			assert.Equal(t, tt.ext, ext)
			assert.NoError(t, err)
		})
	}
}
