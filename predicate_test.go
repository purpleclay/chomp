package chomp_test

import (
	"testing"

	"github.com/purpleclay/chomp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhile(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		rem   string
		ext   string
	}{
		{
			name:  "Ascii",
			input: "Happy New Year. Welcome 2024",
			rem:   " New Year. Welcome 2024",
			ext:   "Happy",
		},
		{
			name:  "Unicode",
			input: "あけましておめでとう。ようこそ 2024 年",
			rem:   "。ようこそ 2024 年",
			ext:   "あけましておめでとう",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			rem, ext, err := chomp.While(chomp.IsLetter)(tt.input)

			require.NoError(t, err)
			assert.Equal(t, tt.rem, rem)
			assert.Equal(t, tt.ext, ext)
		})
	}
}

func TestWhileNot(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		rem   string
		ext   string
	}{
		{
			name:  "Ascii",
			input: "Happy New Year. Welcome 2024",
			rem:   "2024",
			ext:   "Happy New Year. Welcome ",
		},
		{
			name:  "Unicode",
			input: "あけましておめでとう。ようこそ 2024 年",
			rem:   "2024 年",
			ext:   "あけましておめでとう。ようこそ ",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			rem, ext, err := chomp.WhileNot(chomp.IsDigit)(tt.input)

			require.NoError(t, err)
			assert.Equal(t, tt.rem, rem)
			assert.Equal(t, tt.ext, ext)
		})
	}
}
