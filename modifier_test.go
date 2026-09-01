package chomp_test

import (
	"errors"
	"strconv"
	"testing"

	"github.com/purpleclay/chomp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMap(t *testing.T) {
	t.Parallel()

	type Coord struct {
		X int
		Y int
	}

	_, out, err := chomp.Map(
		chomp.SepPair(chomp.While(chomp.IsDigit), chomp.Tag(","), chomp.While(chomp.IsDigit)),
		func(in []string) Coord {
			x, _ := strconv.Atoi(in[0])
			y, _ := strconv.Atoi(in[1])

			return Coord{X: x, Y: y}
		},
	).Run("1,2")

	require.NoError(t, err)
	assert.Equal(t, 1, out.X)
	assert.Equal(t, 2, out.Y)
}

func TestOpt(t *testing.T) {
	t.Parallel()

	rem, ext, err := chomp.Opt(chomp.Tag("the")).Run("dark knight")

	require.NoError(t, err)
	assert.Equal(t, "dark knight", rem)
	assert.Equal(t, "", ext)
}

// TestOptNeverTrustsInnerRemOnFailure asserts that Opt never advances the
// input when its inner combinator fails, even when that combinator itself
// partially consumes before failing. This must hold both for combinators
// that already conform to the contract, and defensively for one that
// deliberately doesn't.
func TestOptNeverTrustsInnerRemOnFailure(t *testing.T) {
	t.Parallel()

	t.Run("ReproducesRoadmapIssue", func(t *testing.T) {
		t.Parallel()

		rem, ext, err := chomp.Opt(chomp.Pair(chomp.Tag("a"), chomp.Tag("b"))).Run("ac")

		require.NoError(t, err)
		assert.Equal(t, "ac", rem)
		assert.Nil(t, ext)
	})

	t.Run("SepPair", func(t *testing.T) {
		t.Parallel()

		rem, ext, err := chomp.Opt(chomp.SepPair(chomp.Tag("Hello"), chomp.Tag(", "), chomp.Tag("World"))).Run("HelloWorld")

		require.NoError(t, err)
		assert.Equal(t, "HelloWorld", rem)
		assert.Nil(t, ext)
	})

	t.Run("Many", func(t *testing.T) {
		t.Parallel()

		rem, ext, err := chomp.Opt(chomp.Many(chomp.Tag("a"))).Run("xyz")

		require.NoError(t, err)
		assert.Equal(t, "xyz", rem)
		assert.Nil(t, ext)
	})

	t.Run("NonConformingInner", func(t *testing.T) {
		t.Parallel()

		// Deliberately violates the combinator contract by returning a
		// partially-consumed remainder on failure, to prove Opt itself
		// never trusts an inner combinator's rem when it errors, rather
		// than relying on the inner combinator to behave.
		var nonConforming chomp.Combinator[string] = func(s chomp.State) (chomp.State, string, error) {
			return s.Advance(1), "", errors.New("boom")
		}

		rem, ext, err := chomp.Opt(nonConforming).Run("abc")

		require.NoError(t, err)
		assert.Equal(t, "abc", rem)
		assert.Equal(t, "", ext)
	})
}

func TestS(t *testing.T) {
	t.Parallel()

	rem, ext, err := chomp.S(chomp.Tag("hello")).Run("hello and good morning")

	require.NoError(t, err)
	assert.Equal(t, " and good morning", rem)
	require.Len(t, ext, 1)
	assert.Equal(t, "hello", ext[0])
}

func TestI(t *testing.T) {
	t.Parallel()

	rem, ext, err := chomp.I(
		chomp.Repeat(chomp.All(chomp.Until(" "), chomp.Tag(" ")), 3),
		2).Run("hello and good morning")

	require.NoError(t, err)
	assert.Equal(t, "morning", rem)
	assert.Equal(t, "and", ext)
}

func TestPeek(t *testing.T) {
	t.Parallel()
	rem, ext, err := chomp.Peek(chomp.Tag("Hello")).Run("Hello and Good Morning!")

	require.NoError(t, err)
	assert.Equal(t, "Hello and Good Morning!", rem)
	assert.Equal(t, "Hello", ext)
}

func TestPeekUsingSequence(t *testing.T) {
	t.Parallel()
	rem, ext, err := chomp.Peek(
		chomp.Many(chomp.Suffixed(chomp.Until(" "), chomp.Tag(" "))),
	).Run("Hello and Good Morning!")

	require.NoError(t, err)
	assert.Equal(t, "Hello and Good Morning!", rem)
	assert.Equal(t, []string{"Hello", "and", "Good"}, ext)
}

func TestFlatten(t *testing.T) {
	t.Parallel()

	rem, ext, err := chomp.Flatten(
		chomp.Many(chomp.Parentheses()),
	).Run("(H)(el)(lo) and Good Morning!")

	require.NoError(t, err)
	assert.Equal(t, " and Good Morning!", rem)
	assert.Equal(t, "Hello", ext)
}
