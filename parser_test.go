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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			rem, ext, err := chomp.Crlf()(tt.input)

			require.NoError(t, err)
			assert.Equal(t, tt.rem, rem)
			assert.Equal(t, tt.ext, ext)
		})
	}
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			rem, ext, err := chomp.Eol()(tt.input)

			require.NoError(t, err)
			assert.Equal(t, tt.rem, rem)
			assert.Equal(t, tt.ext, ext)
		})
	}
}
