package chomp_test

import (
	"testing"

	"github.com/purpleclay/chomp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAlpha(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		rem   string
		ext   string
	}{
		{
			name:  "Ascii",
			input: "Hello123",
			rem:   "123",
			ext:   "Hello",
		},
		{
			name:  "Unicode",
			input: "こんにちは123",
			rem:   "123",
			ext:   "こんにちは",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			rem, ext, err := chomp.Alpha()(tt.input)

			require.NoError(t, err)
			assert.Equal(t, tt.rem, rem)
			assert.Equal(t, tt.ext, ext)
		})
	}
}

func TestAlpha0(t *testing.T) {
	t.Parallel()

	rem, ext, err := chomp.Alpha0()("123Hello")

	require.NoError(t, err)
	assert.Equal(t, "123Hello", rem)
	assert.Equal(t, "", ext)
}

func TestDigit(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		rem   string
		ext   string
	}{
		{
			name:  "Ascii",
			input: "123abc",
			rem:   "abc",
			ext:   "123",
		},
		{
			name:  "Unicode",
			input: "123あいう",
			rem:   "あいう",
			ext:   "123",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			rem, ext, err := chomp.Digit()(tt.input)

			require.NoError(t, err)
			assert.Equal(t, tt.rem, rem)
			assert.Equal(t, tt.ext, ext)
		})
	}
}

func TestDigit0(t *testing.T) {
	t.Parallel()

	rem, ext, err := chomp.Digit0()("abc123")

	require.NoError(t, err)
	assert.Equal(t, "abc123", rem)
	assert.Equal(t, "", ext)
}

func TestAlphanumeric(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		rem   string
		ext   string
	}{
		{
			name:  "Ascii",
			input: "Hello123!",
			rem:   "!",
			ext:   "Hello123",
		},
		{
			name:  "Unicode",
			input: "こんにちは123!",
			rem:   "!",
			ext:   "こんにちは123",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			rem, ext, err := chomp.Alphanumeric()(tt.input)

			require.NoError(t, err)
			assert.Equal(t, tt.rem, rem)
			assert.Equal(t, tt.ext, ext)
		})
	}
}

func TestAlphanumeric0(t *testing.T) {
	t.Parallel()

	rem, ext, err := chomp.Alphanumeric0()("!Hello123")

	require.NoError(t, err)
	assert.Equal(t, "!Hello123", rem)
	assert.Equal(t, "", ext)
}

func TestSpace(t *testing.T) {
	t.Parallel()

	rem, ext, err := chomp.Space()("   \tHello")

	require.NoError(t, err)
	assert.Equal(t, "Hello", rem)
	assert.Equal(t, "   \t", ext)
}

func TestSpace0(t *testing.T) {
	t.Parallel()

	rem, ext, err := chomp.Space0()("Hello")

	require.NoError(t, err)
	assert.Equal(t, "Hello", rem)
	assert.Equal(t, "", ext)
}

func TestMultispace(t *testing.T) {
	t.Parallel()

	rem, ext, err := chomp.Multispace()("  \n\t\rHello")

	require.NoError(t, err)
	assert.Equal(t, "Hello", rem)
	assert.Equal(t, "  \n\t\r", ext)
}

func TestMultispace0(t *testing.T) {
	t.Parallel()

	rem, ext, err := chomp.Multispace0()("Hello")

	require.NoError(t, err)
	assert.Equal(t, "Hello", rem)
	assert.Equal(t, "", ext)
}

func TestHexDigit(t *testing.T) {
	t.Parallel()

	rem, ext, err := chomp.HexDigit()("1a2B3cFf rest")

	require.NoError(t, err)
	assert.Equal(t, " rest", rem)
	assert.Equal(t, "1a2B3cFf", ext)
}

func TestHexDigit0(t *testing.T) {
	t.Parallel()

	rem, ext, err := chomp.HexDigit0()("xyz")

	require.NoError(t, err)
	assert.Equal(t, "xyz", rem)
	assert.Equal(t, "", ext)
}

func TestOctalDigit(t *testing.T) {
	t.Parallel()

	rem, ext, err := chomp.OctalDigit()("01onal")

	require.NoError(t, err)
	assert.Equal(t, "onal", rem)
	assert.Equal(t, "01", ext)
}

func TestOctalDigit0(t *testing.T) {
	t.Parallel()

	rem, ext, err := chomp.OctalDigit0()("89")

	require.NoError(t, err)
	assert.Equal(t, "89", rem)
	assert.Equal(t, "", ext)
}

func TestBinaryDigit(t *testing.T) {
	t.Parallel()

	rem, ext, err := chomp.BinaryDigit()("1010 rest")

	require.NoError(t, err)
	assert.Equal(t, " rest", rem)
	assert.Equal(t, "1010", ext)
}

func TestBinaryDigit0(t *testing.T) {
	t.Parallel()

	rem, ext, err := chomp.BinaryDigit0()("234")

	require.NoError(t, err)
	assert.Equal(t, "234", rem)
	assert.Equal(t, "", ext)
}

func TestNewline(t *testing.T) {
	t.Parallel()

	rem, ext, err := chomp.Newline()("\nHello")

	require.NoError(t, err)
	assert.Equal(t, "Hello", rem)
	assert.Equal(t, "\n", ext)
}

func TestTab(t *testing.T) {
	t.Parallel()

	rem, ext, err := chomp.Tab()("\tHello")

	require.NoError(t, err)
	assert.Equal(t, "Hello", rem)
	assert.Equal(t, "\t", ext)
}

func TestNotLineEnding(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		rem   string
		ext   string
	}{
		{
			name:  "WithNewline",
			input: "Hello, World!\nNext line",
			rem:   "\nNext line",
			ext:   "Hello, World!",
		},
		{
			name:  "WithCarriageReturn",
			input: "Hello, World!\rNext line",
			rem:   "\rNext line",
			ext:   "Hello, World!",
		},
		{
			name:  "Unicode",
			input: "こんにちは\n次の行",
			rem:   "\n次の行",
			ext:   "こんにちは",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			rem, ext, err := chomp.NotLineEnding()(tt.input)

			require.NoError(t, err)
			assert.Equal(t, tt.rem, rem)
			assert.Equal(t, tt.ext, ext)
		})
	}
}

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
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			rem, ext, err := chomp.While(chomp.IsLetter)(tt.input)

			require.NoError(t, err)
			assert.Equal(t, tt.rem, rem)
			assert.Equal(t, tt.ext, ext)
		})
	}
}

func TestWhileN(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		rem   string
		ext   string
	}{
		{
			name:  "Ascii",
			input: "2024 adventure awaits",
			rem:   " adventure awaits",
			ext:   "2024",
		},
		{
			name:  "Unicode",
			input: "2024 年の冒険が待っています",
			rem:   " 年の冒険が待っています",
			ext:   "2024",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			rem, ext, err := chomp.WhileN(chomp.IsDigit, 2)(tt.input)

			require.NoError(t, err)
			assert.Equal(t, tt.rem, rem)
			assert.Equal(t, tt.ext, ext)
		})
	}
}

func TestWhileNM(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		rem   string
		ext   string
	}{
		{
			name:  "Ascii",
			input: "2024 adventure awaits",
			rem:   " adventure awaits",
			ext:   "2024",
		},
		{
			name:  "Unicode",
			input: "2024 年の冒険が待っています",
			rem:   " 年の冒険が待っています",
			ext:   "2024",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			rem, ext, err := chomp.WhileNM(chomp.IsDigit, 1, 8)(tt.input)

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
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			rem, ext, err := chomp.WhileNot(chomp.IsDigit)(tt.input)

			require.NoError(t, err)
			assert.Equal(t, tt.rem, rem)
			assert.Equal(t, tt.ext, ext)
		})
	}
}

func TestWhileNotN(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		rem   string
		ext   string
	}{
		{
			name:  "Ascii",
			input: "2024 adventure awaits",
			rem:   "adventure awaits",
			ext:   "2024 ",
		},
		{
			name:  "Unicode",
			input: "2024 年の冒険が待っています",
			rem:   "年の冒険が待っています",
			ext:   "2024 ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			rem, ext, err := chomp.WhileNotN(chomp.IsLetter, 2)(tt.input)

			require.NoError(t, err)
			assert.Equal(t, tt.rem, rem)
			assert.Equal(t, tt.ext, ext)
		})
	}
}

func TestWhileNotNM(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		rem   string
		ext   string
	}{
		{
			name:  "Ascii",
			input: "2024 adventure awaits",
			rem:   "adventure awaits",
			ext:   "2024 ",
		},
		{
			name:  "Unicode",
			input: "2024 年の冒険が待っています",
			rem:   "年の冒険が待っています",
			ext:   "2024 ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			rem, ext, err := chomp.WhileNotNM(chomp.IsLetter, 1, 8)(tt.input)

			require.NoError(t, err)
			assert.Equal(t, tt.rem, rem)
			assert.Equal(t, tt.ext, ext)
		})
	}
}
