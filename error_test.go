package chomp_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"log/slog"
	"strings"
	"testing"

	"github.com/purpleclay/chomp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCombinatorParseError_Error(t *testing.T) {
	t.Parallel()

	t.Run("LiteralQuoting", func(t *testing.T) {
		t.Parallel()
		_, _, err := chomp.Tag("Hello").Run("Goodbye")
		assert.Equal(t, `chomp: parse error at line 1, column 1 (offset 0): expected "Hello"`, err.Error())
	})

	t.Run("PredicateDerived", func(t *testing.T) {
		t.Parallel()
		_, _, err := chomp.Digit().Run("abc")
		assert.Equal(t, "chomp: parse error at line 1, column 1 (offset 0): expected digit", err.Error())
	})

	t.Run("FallbackMap", func(t *testing.T) {
		t.Parallel()
		_, _, err := chomp.Eof().Run("remaining")
		assert.Equal(t, "chomp: parse error at line 1, column 1 (offset 0): expected end of input", err.Error())
	})

	t.Run("SecondLine", func(t *testing.T) {
		t.Parallel()
		_, _, err := chomp.Pair(chomp.Tag("first\n"), chomp.Tag("third")).Run("first\nsecond")
		assert.Equal(t, `chomp: parse error at line 2, column 1 (offset 6): expected "third"`, err.Error())
	})

	t.Run("LabelsAppendWhileParsingClause", func(t *testing.T) {
		t.Parallel()
		err := chomp.CombinatorParseError{
			Expected: `"Hello"`,
			State:    chomp.NewState("Goodbye"),
			Labels:   []string{"manifest", "version"},
			Type:     "tag",
		}
		assert.Equal(t,
			`chomp: parse error at line 1, column 1 (offset 0): expected "Hello" while parsing manifest > version`,
			err.Error())
	})

	t.Run("NeverContainsNewline", func(t *testing.T) {
		t.Parallel()
		err := chomp.CombinatorParseError{
			Expected: "x",
			State:    chomp.NewState("y"),
			Labels:   []string{"a", "b"},
			Type:     "tag",
		}
		assert.NotContains(t, err.Error(), "\n")
	})
}

func TestCombinatorParseError_Snippet(t *testing.T) {
	t.Parallel()

	t.Run("MidLine", func(t *testing.T) {
		t.Parallel()
		_, _, err := chomp.All(
			chomp.Until("."),
			chomp.Tag("."),
			chomp.Digit(),
			chomp.Tag("."),
			chomp.Digit(),
		).Run("version = 1.4.x")
		require.Error(t, err)

		var pe chomp.CombinatorParseError
		require.True(t, errors.As(err, &pe))

		want := "  |\n" +
			"1 | version = 1.4.x\n" +
			"  | " + strings.Repeat(" ", 14) + "^ expected digit"
		assert.Equal(t, want, pe.Snippet())
	})

	t.Run("EOF", func(t *testing.T) {
		t.Parallel()
		_, _, err := chomp.Pair(chomp.Tag("Hello"), chomp.Tag(", World")).Run("Hello")
		require.Error(t, err)

		var pe chomp.CombinatorParseError
		require.True(t, errors.As(err, &pe))

		want := "  |\n" +
			"1 | Hello\n" +
			"  | " + strings.Repeat(" ", 5) + `^ expected ", World"`
		assert.Equal(t, want, pe.Snippet())
	})

	t.Run("Tabbed", func(t *testing.T) {
		t.Parallel()
		_, _, err := chomp.Pair(chomp.Tag("a\tb\t"), chomp.Tag("Y")).Run("a\tb\tX")
		require.Error(t, err)

		var pe chomp.CombinatorParseError
		require.True(t, errors.As(err, &pe))

		want := "  |\n" +
			"1 | a\tb\tX\n" +
			"  | " + " \t \t" + `^ expected "Y"`
		assert.Equal(t, want, pe.Snippet())
	})

	t.Run("LongLineTruncated", func(t *testing.T) {
		t.Parallel()
		input := strings.Repeat("x", 150) + "!"
		_, _, err := chomp.Pair(chomp.Take(120), chomp.Tag("Y")).Run(input)
		require.Error(t, err)

		var pe chomp.CombinatorParseError
		require.True(t, errors.As(err, &pe))

		wantLine := "1 | … " + strings.Repeat("x", 70) + "!"
		wantCaret := "  | " + strings.Repeat(" ", 42) + `^ expected "Y"`
		want := "  |\n" + wantLine + "\n" + wantCaret
		assert.Equal(t, want, pe.Snippet())
	})

	t.Run("MultiByte", func(t *testing.T) {
		t.Parallel()
		_, _, err := chomp.Pair(chomp.Tag("こんにちは、"), chomp.Tag("Earth")).Run("こんにちは、World")
		require.Error(t, err)

		var pe chomp.CombinatorParseError
		require.True(t, errors.As(err, &pe))

		assert.Equal(t, 18, pe.Offset(), "offset must be a byte offset")
		line, col := pe.Position()
		assert.Equal(t, 1, line)
		assert.Equal(t, 7, col, "column must be rune-counted, not byte-counted")

		want := "  |\n" +
			"1 | こんにちは、World\n" +
			"  | " + strings.Repeat(" ", 6) + `^ expected "Earth"`
		assert.Equal(t, want, pe.Snippet())
	})

	t.Run("BareCRLineBoundary", func(t *testing.T) {
		t.Parallel()
		_, _, err := chomp.All(chomp.Eol(), chomp.Eol(), chomp.Tag("GOODTAG")).
			Run("line1\rline2\rBADTAG")
		require.Error(t, err)

		var pe chomp.CombinatorParseError
		require.True(t, errors.As(err, &pe))

		assert.Equal(t, 12, pe.Offset())
		line, col := pe.Position()
		assert.Equal(t, 3, line, "two bare-CR line breaks precede the failure")
		assert.Equal(t, 1, col)

		want := "  |\n" +
			"3 | BADTAG\n" +
			"  | ^ expected \"GOODTAG\""
		assert.Equal(t, want, pe.Snippet())
	})
}

func TestCombinatorParseError_LogValue(t *testing.T) {
	t.Parallel()

	t.Run("GroupContents", func(t *testing.T) {
		t.Parallel()
		_, _, err := chomp.Digit().Run("abc")
		var pe chomp.CombinatorParseError
		require.True(t, errors.As(err, &pe))

		v := pe.LogValue()
		require.Equal(t, slog.KindGroup, v.Kind())

		attrs := map[string]slog.Value{}
		for _, a := range v.Group() {
			attrs[a.Key] = a.Value
		}

		assert.Equal(t, int64(0), attrs["offset"].Int64())
		assert.Equal(t, int64(1), attrs["line"].Int64())
		assert.Equal(t, int64(1), attrs["column"].Int64())
		assert.Equal(t, "digit", attrs["expected"].String())
		_, hasContext := attrs["context"]
		assert.False(t, hasContext, "context key must be omitted when Labels is empty")
	})

	t.Run("ContextPresentWhenLabelsSet", func(t *testing.T) {
		t.Parallel()
		err := chomp.CombinatorParseError{
			Expected: "digit",
			State:    chomp.NewState("version = 1.4.x"),
			Labels:   []string{"version value", "semver"},
			Type:     "is_digit",
		}

		attrs := map[string]slog.Value{}
		for _, a := range err.LogValue().Group() {
			attrs[a.Key] = a.Value
		}
		assert.Equal(t, "version value > semver", attrs["context"].String())
	})

	t.Run("JSONHandlerEndToEnd", func(t *testing.T) {
		t.Parallel()
		_, _, err := chomp.Digit().Run("abc")

		var buf bytes.Buffer
		logger := slog.New(slog.NewJSONHandler(&buf, nil))
		logger.Error("failed to parse manifest", "err", err)

		var record map[string]any
		require.NoError(t, json.Unmarshal(buf.Bytes(), &record))

		errField, ok := record["err"].(map[string]any)
		require.True(t, ok, "err field must be a structured object, not a string")
		assert.Equal(t, float64(0), errField["offset"])
		assert.Equal(t, float64(1), errField["line"])
		assert.Equal(t, float64(1), errField["column"])
		assert.Equal(t, "digit", errField["expected"])
	})
}

func TestWrapperErrors_DelegateRendering(t *testing.T) {
	t.Parallel()

	// Pair wraps its inner failure in a ParserError.
	_, _, pairErr := chomp.Pair(chomp.Tag("Hello, "), chomp.Tag("World")).Run("Hello, Mars")
	var wrappedParserErr chomp.ParserError
	require.True(t, errors.As(pairErr, &wrappedParserErr))

	var innerFromPair chomp.CombinatorParseError
	require.True(t, errors.As(pairErr, &innerFromPair))
	assert.Equal(t, innerFromPair.Error(), wrappedParserErr.Error(),
		"ParserError.Error() must equal its inner CombinatorParseError.Error()")
	assert.Equal(t, innerFromPair.LogValue(), wrappedParserErr.LogValue())

	// Repeat wraps its inner failure in a RangedParserError.
	_, _, repeatErr := chomp.Repeat(chomp.Tag("ab"), 3).Run("abxy")
	var wrappedRangedErr chomp.RangedParserError
	require.True(t, errors.As(repeatErr, &wrappedRangedErr))

	var innerFromRepeat chomp.CombinatorParseError
	require.True(t, errors.As(repeatErr, &innerFromRepeat))
	assert.Equal(t, innerFromRepeat.Error(), wrappedRangedErr.Error())
	assert.Equal(t, innerFromRepeat.LogValue(), wrappedRangedErr.LogValue())

	// Cut wraps its inner failure in a CutError.
	_, _, cutErr := chomp.Cut(chomp.Tag("Hello")).Run("World")
	var wrappedCutErr chomp.CutError
	require.True(t, errors.As(cutErr, &wrappedCutErr))

	var innerFromCut chomp.CombinatorParseError
	require.True(t, errors.As(cutErr, &innerFromCut))
	assert.Equal(t, innerFromCut.Error(), wrappedCutErr.Error())
	assert.Equal(t, innerFromCut.LogValue(), wrappedCutErr.LogValue())
}

func TestZeroValueWrappersDoNotPanic(t *testing.T) {
	t.Parallel()

	t.Run("ParserError", func(t *testing.T) {
		t.Parallel()
		var e chomp.ParserError
		e.Type = "pair"

		require.NotPanics(t, func() { _ = e.Error() })
		assert.Equal(t, "chomp: pair parser error with no underlying cause", e.Error())
		assert.NotContains(t, e.Error(), "\n")

		require.NotPanics(t, func() { _ = e.LogValue() })
		assert.Equal(t, e.Error(), e.LogValue().String())
	})

	t.Run("RangedParserError", func(t *testing.T) {
		t.Parallel()
		var e chomp.RangedParserError
		e.Type = "repeat"

		require.NotPanics(t, func() { _ = e.Error() })
		assert.Equal(t, "chomp: repeat parser error with no underlying cause", e.Error())
		assert.NotContains(t, e.Error(), "\n")

		require.NotPanics(t, func() { _ = e.LogValue() })
		assert.Equal(t, e.Error(), e.LogValue().String())
	})

	t.Run("CutError", func(t *testing.T) {
		t.Parallel()
		var e chomp.CutError

		require.NotPanics(t, func() { _ = e.Error() })
		assert.Equal(t, "chomp: cut error with no underlying cause", e.Error())
		assert.NotContains(t, e.Error(), "\n")

		require.NotPanics(t, func() { _ = e.LogValue() })
		assert.Equal(t, e.Error(), e.LogValue().String())
	})
}

func TestIOutOfBoundsError(t *testing.T) {
	t.Parallel()

	_, _, err := chomp.I(
		chomp.SepPair(chomp.Tag("Hello"), chomp.Tag(", "), chomp.Tag("World")), 5).
		Run("Hello, World!")
	require.Error(t, err)

	assert.True(t, strings.HasPrefix(err.Error(), "chomp: "))
	assert.NotContains(t, err.Error(), "\n")

	var pe chomp.CombinatorParseError
	assert.False(t, errors.As(err, &pe), "I's out-of-bounds error must not be a CombinatorParseError")
}
