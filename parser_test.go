package chomp_test

import (
	"testing"

	"github.com/purpleclay/chomp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCrlf(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		rem   string
		ext   string
	}{
		{
			name:  "CRLF",
			input: "\r\nこんにちは",
			rem:   "こんにちは",
			ext:   "\r\n",
		},
		{
			name:  "CRLFOnly",
			input: "\r\n",
			rem:   "",
			ext:   "\r\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			rem, ext, err := chomp.Crlf().Run(tt.input)

			require.NoError(t, err)
			assert.Equal(t, tt.rem, rem)
			assert.Equal(t, tt.ext, ext)
		})
	}
}

func TestCrlfRejectsBareLineEndings(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
	}{
		{name: "BareLF", input: "\nHello"},
		{name: "BareCR", input: "\rHello"},
		{name: "NoLineEnding", input: "Hello"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			rem, ext, err := chomp.Crlf().Run(tt.input)

			require.Error(t, err)
			assert.Equal(t, tt.input, rem)
			assert.Equal(t, "", ext)
		})
	}
}

func TestLineEnding(t *testing.T) {
	t.Parallel()

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
			t.Parallel()
			rem, ext, err := chomp.LineEnding().Run(tt.input)

			require.NoError(t, err)
			assert.Equal(t, tt.rem, rem)
			assert.Equal(t, tt.ext, ext)
		})
	}
}

func TestLineEndingRejectsBareCR(t *testing.T) {
	t.Parallel()

	rem, ext, err := chomp.LineEnding().Run("\rHello")

	require.Error(t, err)
	assert.Equal(t, "\rHello", rem)
	assert.Equal(t, "", ext)
}

func TestEol(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		rem   string
		ext   string
	}{
		{
			name: "LF",
			input: `Hello, World!
It's a great day`,
			rem: "It's a great day",
			ext: "Hello, World!",
		},
		{
			name:  "NoLF",
			input: "こんにちは、おはよう",
			rem:   "",
			ext:   "こんにちは、おはよう",
		},
		{
			name: "EmptyLineBeforeLF",
			input: `
こんにちは、おはよう`,
			rem: "こんにちは、おはよう",
			ext: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			rem, ext, err := chomp.Eol().Run(tt.input)

			require.NoError(t, err)
			assert.Equal(t, tt.rem, rem)
			assert.Equal(t, tt.ext, ext)
		})
	}
}
