package test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"ui.web/server/src"
)

func TestChangeExt(t *testing.T) {
	require.Equal(t, "hi.png", src.ChangExtension("hi", ".png"))
	require.Equal(t, "hi.png", src.ChangExtension("hi.txt", ".png"))
	require.Equal(t, "foo/bar.exe", src.ChangExtension("foo/bar", ".exe"))
	require.Equal(t, "foo/bar.exe", src.ChangExtension("foo/bar.txt", ".exe"))
}
