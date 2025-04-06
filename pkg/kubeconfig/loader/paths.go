package loader

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"

	"github.com/eth-p/kubesel/internal/parallel"
)

const DefaultKubeDirName = ".kube"
const DefaultKubeconfigFileName = "config"

// FindKubeConfigFiles finds all the kubectl configuration files.
//
// This uses the same logic as kubectl, splitting the `KUBECONFIG` environment
// variable as a list of paths and falling back to `$HOME/.kube/config` if
// the variable is not defined.
//
// REF: https://kubernetes.io/docs/concepts/configuration/organize-cluster-access-kubeconfig/#the-kubeconfig-environment-variable
func FindKubeConfigFiles() ([]string, error) {
	kubeconfigVar, ok := os.LookupEnv("KUBECONFIG")
	if ok {
		var files []string

		for _, file := range filepath.SplitList(kubeconfigVar) {
			if file == "" {
				continue
			}

			files = append(files, file)
		}
		return files, nil
	}

	// Search for the default kubeconfig file.
	defaultKubeconfig, err := FindDefaultKubeconfigFile()
	if err != nil {
		return nil, err
	}

	return []string{defaultKubeconfig}, nil
}

// FindDefaultKubeconfigFile returns the path to the default `.kube/config`
// file.
//
// On Linux/Mac, this will be `$HOME/.kube/config`.
//
// On Windows, this will one of:
//   - `%HOME%`
//   - `%HOMEDRIVE%/%HOMEPATH%`
//   - `%USERPROFILE%`
//
// See [kubernetes source] comments for more info.
//
// [kubernetes source]: https://github.com/kubernetes/client-go/blob/2086688a727d00268de695a9c701e52458939701/util/homedir/homedir.go#L25-L30
func FindDefaultKubeconfigFile() (string, error) {
	kubedir, err := FindDefaultKubeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(kubedir, DefaultKubeconfigFileName), nil
}

// FindDefaultKubeDir returns the path to the default `.kube` directory.
//
// See: [FindDefaultKubeDirPOSIX] and [FindDefaultKubeDirWindows]
func FindDefaultKubeDir() (string, error) {
	if runtime.GOOS == "windows" {
		return FindDefaultKubeDirWindows()
	} else {
		return FindDefaultKubeDirPOSIX()
	}
}

// FindDefaultKubeDirPOSIX returns the path to the kubectl config directory
// for a Linux or MacOS machine.
//
// This will be `$HOME/.kube/`.
//
// Contrary to the [Kubernetes implementation], if `$HOME` is unset, an error
// will be returned instead of an empty string.
//
// [Kubernetes implementation]: https://github.com/kubernetes/client-go/blob/2086688a727d00268de695a9c701e52458939701/util/homedir/homedir.go#L25-L30
func FindDefaultKubeDirPOSIX() (string, error) {
	return doFindDefaultKubeDirPOSIX(os.LookupEnv)
}

func doFindDefaultKubeDirPOSIX(lookupEnv func(string) (string, bool)) (string, error) {
	homedir, ok := lookupEnv("HOME")
	if !ok {
		return "", fmt.Errorf("%w: $HOME environment variable undefined", ErrNoKubeDir)
	}

	return filepath.Join(homedir, DefaultKubeDirName), nil
}

// FindDefaultKubeDirWindows returns the path to the kubectl config directory
// for a Windows machine.
//
// This will be one of:
//   - `%HOME%/.kube/`
//   - `%HOMEDRIVE%/%HOMEPATH%/.kube/`
//   - `%USERPROFILE%/.kube/`
//
// See [kubernetes implementation] comments for more info.
//
// [kubernetes implementation]: https://github.com/kubernetes/client-go/blob/2086688a727d00268de695a9c701e52458939701/util/homedir/homedir.go#L25-L30
func FindDefaultKubeDirWindows() (string, error) {
	return doFindDefaultKubeDirWindows(os.LookupEnv, os.Stat)
}

func doFindDefaultKubeDirWindows(
	lookupEnv func(string) (string, bool),
	stat func(string) (fs.FileInfo, error),
) (string, error) {
	type candidate struct {
		path          string
		hasKubeConfig bool
		isExist       bool
		isWriteable   bool
	}

	var (
		fromHome        candidate
		fromDriveHome   candidate
		fromUserProfile candidate
	)

	// Create candidate paths, depending on which environment variables are set.
	if homeDirVar, ok := lookupEnv("HOME"); ok && len(homeDirVar) > 0 {
		fromHome.path = homeDirVar
	}

	if homeDriveVar, ok := lookupEnv("HOMEDRIVE"); ok && len(homeDriveVar) > 0 {
		if homePathVar, ok := lookupEnv("HOMEPATH"); ok && len(homePathVar) > 0 {
			fromDriveHome.path = filepath.Join(homeDriveVar, homePathVar)
		}
	}

	if userProfileVar, ok := lookupEnv("USERPROFILE"); ok && len(userProfileVar) > 0 {
		fromUserProfile.path = userProfileVar
	}

	// Return the first candidate with a `.kube/config` file.
	// While we're at it, determine if the directory exists or is writeable.
	//  * https://github.com/kubernetes/client-go/blob/2086688a727d00268de695a9c701e52458939701/util/homedir/homedir.go#L40-L41
	checkConfigOrder := []*candidate{&fromHome, &fromDriveHome, &fromUserProfile}
	parallel.Run(checkConfigOrder, func(c *candidate) {
		if c.path == "" {
			return
		}

		_, err := stat(filepath.Join(c.path, DefaultKubeDirName, DefaultKubeconfigFileName))
		if err == nil {
			c.hasKubeConfig = true
			return
		}

		const permOwnerWrite = 0o200
		if stat, err := stat(c.path); err == nil {
			c.isExist = true
			c.isWriteable = (stat.Mode().Perm() & permOwnerWrite) > 0
		}
	})

	for _, c := range checkConfigOrder {
		if c.hasKubeConfig {
			return filepath.Join(c.path, DefaultKubeDirName), nil
		}
	}

	// Return the first writeable candidate.
	//  * https://github.com/kubernetes/client-go/blob/2086688a727d00268de695a9c701e52458939701/util/homedir/homedir.go#L55-L56
	checkDirOrder := []*candidate{&fromHome, &fromUserProfile, &fromDriveHome}
	for _, c := range checkDirOrder {
		if c.isWriteable {
			return filepath.Join(c.path, DefaultKubeDirName), nil
		}
	}

	// Return the first existent candidate.
	for _, c := range checkDirOrder {
		if c.isExist {
			return filepath.Join(c.path, DefaultKubeDirName), nil
		}
	}

	return "", ErrNoKubeDir
}
