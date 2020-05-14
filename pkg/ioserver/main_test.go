package ioserver

import (
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"testing"
)

func NewSocketPath(t *testing.T) (path string, cleanup func()) {
	dir, err := ioutil.TempDir("", "sumpboi-ioserver")
	require.NoError(t, err, "socket directory should have created successfully")
	return dir + "/gpio.sock", func() {
		require.NoError(t, os.RemoveAll(dir), "socket directory cleanup should have succeeded")
	}
}
