package testutil

import (
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

var update = flag.Bool("update", false, "update golden files")

func AssertGolden(t *testing.T, name string, actual []byte) {
	t.Helper()
	goldenPath := filepath.Join("..", "testdata", name+".golden")

	if *update {
		err := os.MkdirAll(filepath.Dir(goldenPath), 0755)
		require.NoError(t, err)
		err = os.WriteFile(goldenPath, actual, 0644)
		require.NoError(t, err)
		t.Logf("Updated golden file: %s", goldenPath)
	}

	expected, err := os.ReadFile(goldenPath)
	require.NoError(t, err, "golden file not found; run with -update to generate")

	require.Equal(t, string(expected), string(actual), "output does not match golden file")
}
