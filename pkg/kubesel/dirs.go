package kubesel

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/adrg/xdg"
)

// findDataHomeDir returns the XDG_DATA_HOME directory.
// If the environment variable is set, it will be used.
//
// Otherwise, it should be ~/.local/share
func findDataHomeDir() string {
	dataHome, ok := os.LookupEnv("XDG_DATA_HOME")
	if ok {
		return dataHome
	}

	// Command-line tools on Darwin should be similar to Linux.
	if runtime.GOOS == "darwin" {
		return filepath.Join(xdg.Home, ".local", "share")
	}

	// Use the library.
	return xdg.DataHome
}
