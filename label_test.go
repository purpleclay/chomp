package chomp_test

import (
	"errors"
	"log/slog"
	"testing"

	"github.com/purpleclay/chomp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLabel_Success(t *testing.T) {
	t.Parallel()

	rem, ext, err := chomp.Label("digits", chomp.Digit()).Run("123abc")
	require.NoError(t, err)
	assert.Equal(t, "abc", rem)
	assert.Equal(t, "123", ext)
}

func TestLabel_SingleLabel(t *testing.T) {
	t.Parallel()

	_, _, err := chomp.Label("version", chomp.Digit()).Run("abc")
	require.Error(t, err)

	var pe chomp.CombinatorParseError
	require.True(t, errors.As(err, &pe))
	assert.Equal(t, []string{"version"}, pe.Labels)
	assert.Equal(t,
		"chomp: parse error at line 1, column 1 (offset 0): expected digit while parsing version",
		err.Error())
}

func TestLabel_NestedOutermostFirst(t *testing.T) {
	t.Parallel()

	_, _, err := chomp.Label("manifest", chomp.Label("version", chomp.Digit())).Run("abc")
	require.Error(t, err)

	var pe chomp.CombinatorParseError
	require.True(t, errors.As(err, &pe))
	assert.Equal(t, []string{"manifest", "version"}, pe.Labels)
	assert.Equal(t,
		"chomp: parse error at line 1, column 1 (offset 0): expected digit while parsing manifest > version",
		err.Error())
}

func TestLabel_NameWithNewlineDoesNotBreakSingleLine(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		label string
		want  string
	}{
		{name: "LF", label: "manifest\nversion", want: "manifest version"},
		{name: "CRLF", label: "manifest\r\nversion", want: "manifest version"},
		{name: "BareCR", label: "manifest\rversion", want: "manifest version"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, _, err := chomp.Label(tt.label, chomp.Digit()).Run("abc")
			require.Error(t, err)

			assert.NotContains(t, err.Error(), "\n")
			assert.NotContains(t, err.Error(), "\r")

			var pe chomp.CombinatorParseError
			require.True(t, errors.As(err, &pe))
			assert.Equal(t, []string{tt.want}, pe.Labels)
			assert.Equal(t,
				"chomp: parse error at line 1, column 1 (offset 0): expected digit while parsing "+tt.want,
				err.Error())
		})
	}
}

func TestLabel_SnippetIncludesLabelChain(t *testing.T) {
	t.Parallel()

	_, _, err := chomp.Label("manifest", chomp.Label("version", chomp.Digit())).Run("abc")
	require.Error(t, err)

	var pe chomp.CombinatorParseError
	require.True(t, errors.As(err, &pe))
	assert.Equal(t,
		"  |\n1 | abc\n  | ^ expected digit while parsing manifest > version",
		pe.Snippet())
}

func TestLabel_LogValueContext(t *testing.T) {
	t.Parallel()

	_, _, err := chomp.Label("manifest", chomp.Label("version", chomp.Digit())).Run("abc")
	require.Error(t, err)

	var pe chomp.CombinatorParseError
	require.True(t, errors.As(err, &pe))

	attrs := map[string]slog.Value{}
	for _, a := range pe.LogValue().Group() {
		attrs[a.Key] = a.Value
	}
	assert.Equal(t, "manifest > version", attrs["context"].String())
}

func TestLabel_UnwrapChainIntact(t *testing.T) {
	t.Parallel()

	t.Run("ParserError", func(t *testing.T) {
		t.Parallel()
		_, _, err := chomp.Label("greeting",
			chomp.Pair(chomp.Tag("Hello, "), chomp.Tag("World"))).
			Run("Hello, Mars")
		require.Error(t, err)

		var wrapped chomp.ParserError
		require.True(t, errors.As(err, &wrapped))
		assert.Equal(t, "pair", wrapped.Type)

		var pe chomp.CombinatorParseError
		require.True(t, errors.As(err, &pe))
		assert.Equal(t, []string{"greeting"}, pe.Labels)
	})

	t.Run("RangedParserError", func(t *testing.T) {
		t.Parallel()
		_, _, err := chomp.Label("triple", chomp.Repeat(chomp.Tag("ab"), 3)).Run("abxy")
		require.Error(t, err)

		var wrapped chomp.RangedParserError
		require.True(t, errors.As(err, &wrapped))
		assert.Equal(t, "repeat", wrapped.Type)

		var pe chomp.CombinatorParseError
		require.True(t, errors.As(err, &pe))
		assert.Equal(t, []string{"triple"}, pe.Labels)
	})

	t.Run("CutError", func(t *testing.T) {
		t.Parallel()
		_, _, err := chomp.Label("committed", chomp.Cut(chomp.Tag("Hello"))).Run("World")
		require.Error(t, err)

		var wrapped chomp.CutError
		require.True(t, errors.As(err, &wrapped))

		var pe chomp.CombinatorParseError
		require.True(t, errors.As(err, &pe))
		assert.Equal(t, []string{"committed"}, pe.Labels)
	})
}

func TestLabel_NonCombinatorParseErrorUnaffected(t *testing.T) {
	t.Parallel()

	labelled := chomp.Label("index",
		chomp.I(chomp.All(chomp.Tag("Hello"), chomp.Tag(", World!")), 5))
	unlabelled := chomp.I(chomp.All(chomp.Tag("Hello"), chomp.Tag(", World!")), 5)

	_, _, labelledErr := labelled.Run("Hello, World!")
	_, _, unlabelledErr := unlabelled.Run("Hello, World!")

	require.Error(t, labelledErr)
	require.Error(t, unlabelledErr)
	assert.Equal(t, unlabelledErr.Error(), labelledErr.Error(),
		"Label must not alter an error with no CombinatorParseError to attach to")

	var pe chomp.CombinatorParseError
	assert.False(t, errors.As(labelledErr, &pe))
}
