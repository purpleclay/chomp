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
