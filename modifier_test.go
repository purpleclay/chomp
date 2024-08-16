package chomp_test

import (
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
	)("1,2")

	require.NoError(t, err)
	assert.Equal(t, 1, out.X)
	assert.Equal(t, 2, out.Y)
}

func TestOpt(t *testing.T) {
	t.Parallel()

	rem, ext, err := chomp.Opt(chomp.Tag("the"))("dark knight")

	require.NoError(t, err)
	assert.Equal(t, "dark knight", rem)
	assert.Equal(t, "", ext)
}

func TestS(t *testing.T) {
	t.Parallel()

	rem, ext, err := chomp.S(chomp.Tag("hello"))("hello and good morning")

	require.NoError(t, err)
	assert.Equal(t, " and good morning", rem)
	require.Len(t, ext, 1)
	assert.Equal(t, "hello", ext[0])
}

func TestI(t *testing.T) {
	t.Parallel()

	rem, ext, err := chomp.I(
		chomp.Repeat(chomp.All(chomp.Until(" "), chomp.Tag(" ")), 3),
		2)("hello and good morning")

	require.NoError(t, err)
	assert.Equal(t, "morning", rem)
	assert.Equal(t, "and", ext)
}

func TestPeek(t *testing.T) {
	t.Parallel()
	rem, ext, err := chomp.Peek(chomp.Tag("Hello"))("Hello and Good Morning!")

	require.NoError(t, err)
	assert.Equal(t, "Hello and Good Morning!", rem)
	assert.Equal(t, "Hello", ext)
}

func TestPeekUsingSequence(t *testing.T) {
	t.Parallel()
	rem, ext, err := chomp.Peek(
		chomp.Many(chomp.Suffixed(chomp.Until(" "), chomp.Tag(" "))),
	)("Hello and Good Morning!")

	require.NoError(t, err)
	assert.Equal(t, "Hello and Good Morning!", rem)
	assert.Equal(t, []string{"Hello", "and", "Good"}, ext)
}
