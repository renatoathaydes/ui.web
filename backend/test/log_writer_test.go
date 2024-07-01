package test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"ui.web/server/src"
)

func CanWriteSingleLine(t *testing.T) {
	var lines []string
	writer := src.NewWriter(func(line string) {
		lines = append(lines, line)
	})

	c, err := writer.Write([]byte{65, 66, 67})

	require.Nil(t, err)
	require.Equal(t, 3, c)
	require.Equal(t, []string{"ABC"}, lines)
}

func CanWriteMultiLines(t *testing.T) {
	var lines []string
	writer := src.NewWriter(func(line string) {
		lines = append(lines, line)
	})

	c, err := writer.Write([]byte{65, 66})
	require.Nil(t, err)
	require.Equal(t, 2, c)

	c, err = writer.Write([]byte{67, 10, 68, 69, 10, 70})
	require.Nil(t, err)
	require.Equal(t, 6, c)

	c, err = writer.Write([]byte{71, 72, 10})
	require.Nil(t, err)
	require.Equal(t, 3, c)

	writer.Close()

	require.Equal(t, []string{"ABC", "DE", "FGH"}, lines)
}
