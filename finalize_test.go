package chomp_test

import (
	"errors"
	"runtime"
	"strings"
	"testing"
	"time"
	"unsafe"

	"github.com/purpleclay/chomp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFinalizeDropsOriginalInput(t *testing.T) {
	buf := make([]byte, 8*1024*1024)
	for i := range buf {
		buf[i] = 'x'
	}
	input := string(buf)
	buf = nil

	collected := make(chan struct{})
	runtime.SetFinalizer(unsafe.StringData(input), func(*byte) {
		close(collected)
	})

	_, _, err := chomp.Tag("nomatch").Run(input)
	require.Error(t, err)

	input = "" // drop the only other reference to the backing array

	for range 20 {
		runtime.GC()
		select {
		case <-collected:
			runtime.KeepAlive(err)
			return
		case <-time.After(20 * time.Millisecond):
		}
	}

	t.Fatal("original input's backing array was not collected; the escaped error still retains it")
}

func TestFinalizeIsIdempotent(t *testing.T) {
	t.Parallel()

	_, _, rawErr := chomp.Tag("Hello")(chomp.NewState("Goodbye"))
	require.Error(t, rawErr)

	once := chomp.Finalize(rawErr)
	twice := chomp.Finalize(once)

	assert.Equal(t, rawErr.Error(), once.Error(), "finalising must not change the rendered message")
	assert.Equal(t, once.Error(), twice.Error())

	var pe1, pe2 chomp.CombinatorParseError
	require.True(t, errors.As(once, &pe1))
	require.True(t, errors.As(twice, &pe2))
	assert.Equal(t, pe1, pe2, "applying Finalize twice must be a no-op")
	assert.Equal(t, pe1.Snippet(), pe2.Snippet())
}

// assertSnippetStableAcrossFinalize proves Snippet renders identically
// whether read from a raw, unfinalised error (c invoked directly against a
// State) or from the finalised error Run returns for the same failure.
func assertSnippetStableAcrossFinalize[T chomp.Result](t *testing.T, input string, c chomp.Combinator[T]) {
	t.Helper()

	_, _, rawErr := c(chomp.NewState(input))
	require.Error(t, rawErr)
	var rawPE chomp.CombinatorParseError
	require.True(t, errors.As(rawErr, &rawPE))

	_, _, runErr := c.Run(input)
	require.Error(t, runErr)
	var finalPE chomp.CombinatorParseError
	require.True(t, errors.As(runErr, &finalPE))

	assert.Equal(t, rawPE.Error(), finalPE.Error())
	assert.Equal(t, rawPE.Snippet(), finalPE.Snippet())
}

func TestSnippetIdenticalBeforeAndAfterFinalize(t *testing.T) {
	t.Parallel()

	t.Run("MidLine", func(t *testing.T) {
		t.Parallel()
		assertSnippetStableAcrossFinalize(t, "version = 1.4.x", chomp.All(
			chomp.Until("."), chomp.Tag("."), chomp.Digit(), chomp.Tag("."), chomp.Digit(),
		))
	})

	t.Run("Tabbed", func(t *testing.T) {
		t.Parallel()
		assertSnippetStableAcrossFinalize(t, "a\tb\tX", chomp.Pair(chomp.Tag("a\tb\t"), chomp.Tag("Y")))
	})

	t.Run("LongLineTruncated", func(t *testing.T) {
		t.Parallel()
		input := strings.Repeat("x", 150) + "!"
		assertSnippetStableAcrossFinalize(t, input, chomp.Pair(chomp.Take(120), chomp.Tag("Y")))
	})

	t.Run("MultiByte", func(t *testing.T) {
		t.Parallel()
		assertSnippetStableAcrossFinalize(t, "こんにちは、World",
			chomp.Pair(chomp.Tag("こんにちは、"), chomp.Tag("Earth")))
	})
}
