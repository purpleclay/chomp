package chomp_test

import (
	"testing"
	"time"

	"github.com/purpleclay/chomp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// withTimeout runs fn and fails the test if it does not complete within d.
// Used to turn a suspected infinite loop into a fast, reported test failure
// instead of hanging the test binary.
func withTimeout(t *testing.T, d time.Duration, fn func()) {
	t.Helper()

	done := make(chan struct{})
	go func() {
		fn()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(d):
		t.Fatalf("did not terminate within %s (suspected infinite loop on zero-width match)", d)
	}
}

const hangTimeout = 500 * time.Millisecond

func TestManyNZeroWidthTerminates(t *testing.T) {
	t.Parallel()
	withTimeout(t, hangTimeout, func() {
		rem, ext, err := chomp.ManyN(chomp.Digit0(), 0).Run("abc")
		require.NoError(t, err)
		assert.Equal(t, "abc", rem)
		assert.Empty(t, ext)
	})
}

func TestManyNZeroWidthUntilTerminates(t *testing.T) {
	t.Parallel()
	withTimeout(t, hangTimeout, func() {
		rem, ext, err := chomp.ManyN(chomp.Until(""), 0).Run("abc")
		require.NoError(t, err)
		assert.Equal(t, "abc", rem)
		assert.Empty(t, ext)
	})
}

func TestManyZeroWidthTerminates(t *testing.T) {
	t.Parallel()
	withTimeout(t, hangTimeout, func() {
		_, _, err := chomp.Many(chomp.Digit0()).Run("abc")
		require.Error(t, err)
	})
}

func TestSeparatedListZeroWidthTerminates(t *testing.T) {
	t.Parallel()
	withTimeout(t, hangTimeout, func() {
		rem, ext, err := chomp.SeparatedList(chomp.Digit0(), chomp.Opt(chomp.Tag(","))).Run("abc")
		require.NoError(t, err)
		assert.Equal(t, "abc", rem)
		require.Len(t, ext, 1)
		assert.Equal(t, "", ext[0])
	})
}

func TestSeparatedList0ZeroWidthTerminates(t *testing.T) {
	t.Parallel()
	withTimeout(t, hangTimeout, func() {
		rem, ext, err := chomp.SeparatedList0(chomp.Digit0(), chomp.Opt(chomp.Tag(","))).Run("abc")
		require.NoError(t, err)
		assert.Equal(t, "abc", rem)
		require.Len(t, ext, 1)
		assert.Equal(t, "", ext[0])
	})
}

func TestManyTillZeroWidthTerminates(t *testing.T) {
	t.Parallel()
	withTimeout(t, hangTimeout, func() {
		_, _, err := chomp.ManyTill(chomp.Digit0(), chomp.Tag("END")).Run("abc")
		require.Error(t, err)
	})
}

func TestManyTill0ZeroWidthTerminates(t *testing.T) {
	t.Parallel()
	withTimeout(t, hangTimeout, func() {
		_, _, err := chomp.ManyTill0(chomp.Digit0(), chomp.Tag("END")).Run("abc")
		require.Error(t, err)
	})
}

func TestFoldManyZeroWidthTerminates(t *testing.T) {
	t.Parallel()
	withTimeout(t, hangTimeout, func() {
		_, _, err := chomp.FoldMany(chomp.Digit0(), 0, func(acc int, _ string) int { return acc }).Run("abc")
		require.Error(t, err)
	})
}

func TestFoldMany0ZeroWidthTerminates(t *testing.T) {
	t.Parallel()
	withTimeout(t, hangTimeout, func() {
		rem, ext, err := chomp.FoldMany0(chomp.Digit0(), 0, func(acc int, _ string) int { return acc }).Run("abc")
		require.NoError(t, err)
		assert.Equal(t, "abc", rem)
		assert.Equal(t, 0, ext)
	})
}

func TestManyCountZeroWidthTerminates(t *testing.T) {
	t.Parallel()
	withTimeout(t, hangTimeout, func() {
		_, _, err := chomp.ManyCount(chomp.Digit0()).Run("abc")
		require.Error(t, err)
	})
}

func TestManyCount0ZeroWidthTerminates(t *testing.T) {
	t.Parallel()
	withTimeout(t, hangTimeout, func() {
		rem, count, err := chomp.ManyCount0(chomp.Digit0()).Run("abc")
		require.NoError(t, err)
		assert.Equal(t, "abc", rem)
		assert.Equal(t, uint(0), count)
	})
}
