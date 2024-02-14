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
