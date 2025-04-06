package loader

import (
	"fmt"
	"io"
	"os"

	"github.com/eth-p/kubesel/internal/parallel"
	"github.com/eth-p/kubesel/pkg/kubeconfig"
	"gopkg.in/yaml.v3"
)

// LoadedKubeconfig represents a loaded kubectl configuration file.
type LoadedKubeconfig struct {
	Path   string
	Config kubeconfig.Config
	Errors []error
}

type LoadedKubeconfigCollection struct {
	Configs []*LoadedKubeconfig
	Merged  *kubeconfig.Config
}

func LoadMultipleFiles(files []string) *LoadedKubeconfigCollection {
	result := new(LoadedKubeconfigCollection)
	result.Configs = make([]*LoadedKubeconfig, len(files))
	result.Merged = new(kubeconfig.Config)

	// Load each config file in parallel, merging them iteratively.
	for i, kc := range parallel.Ordered(files, LoadFromFile) {
		result.Configs[i] = kc
		result.Merged = kubeconfig.MergeConfig(result.Merged, &kc.Config)
	}

	return result
}

// LoadFromFile reads and parses a [kubeconfig.Config] file from the filesystem,
// returning a [LoadedKubeconfig] with its contents.
func LoadFromFile(file string) *LoadedKubeconfig {
	handle, err := os.OpenFile(file, os.O_RDONLY, 0)
	if err != nil {
		return &LoadedKubeconfig{
			Errors: []error{fmt.Errorf("%w: %w", ErrReading, err)},
		}
	}

	defer handle.Close()

	result := LoadFromReader(handle)
	result.Path = file
	return result
}

// LoadFromReader reads and parses a [kubeconfig.Config] file from a
// [io.Reader], returning a [LoadedKubeconfig] with its contents.
func LoadFromReader(reader io.Reader) *LoadedKubeconfig {
	buffer, err := io.ReadAll(reader)
	if err != nil {
		return &LoadedKubeconfig{
			Errors: []error{fmt.Errorf("%w: %w", ErrReading, err)},
		}
	}

	result := new(LoadedKubeconfig)
	err = yaml.Unmarshal(buffer, &result.Config)
	if err != nil {
		result.Errors = []error{fmt.Errorf("%w: %w", ErrParsing, err)}
	}

	return result
}
