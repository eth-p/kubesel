package loader

import (
	"io/fs"
	"os"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type fakeFileInfo struct {
	mode os.FileMode
}

func (f *fakeFileInfo) Name() string       { return "fake.file" }
func (f *fakeFileInfo) Size() int64        { return 0 }
func (f *fakeFileInfo) Mode() os.FileMode  { return f.mode }
func (f *fakeFileInfo) ModTime() time.Time { return time.Unix(0, 0) }
func (f *fakeFileInfo) IsDir() bool        { return f.Mode().IsDir() }
func (f *fakeFileInfo) Sys() any           { return nil }

func TestFindDefaultKubeDirWindows(t *testing.T) {
	testcases := map[string]struct {
		env             map[string]string
		files           map[string]os.FileMode
		expected        string
		expectedErrorIs error
	}{
		"HOME used 1st when config file exists": {
			files: map[string]os.FileMode{
				"C:/Home/Fake/.kube/config":        0o777,
				"C:/HomePath/Fake/.kube/config":    0o777,
				"C:/UserProfile/Fake/.kube/config": 0o777,
			},
			env: map[string]string{
				"HOME":        "C:/Home/Fake",
				"HOMEDRIVE":   "C:/",
				"HOMEPATH":    "HomePath/Fake",
				"USERPROFILE": "C:/UserProfile/Fake",
			},
			expected: "C:/Home/Fake/.kube",
		},
		"HOMEDRIVE\\HOMEPATH used 2nd when config file exists": {
			files: map[string]os.FileMode{
				"C:/HomePath/Fake/.kube/config":    0o777,
				"C:/UserProfile/Fake/.kube/config": 0o777,
			},
			env: map[string]string{
				"HOME":        "C:/Home/Fake",
				"HOMEDRIVE":   "C:/",
				"HOMEPATH":    "HomePath/Fake",
				"USERPROFILE": "C:/UserProfile/Fake",
			},
			expected: "C:/HomePath/Fake/.kube",
		},
		"USERPROFILE used 3rd when config file exists": {
			files: map[string]os.FileMode{
				"C:/UserProfile/Fake/.kube/config": 0o777,
			},
			env: map[string]string{
				"HOME":        "C:/Home/Fake",
				"HOMEDRIVE":   "C:/",
				"HOMEPATH":    "HomePath/Fake",
				"USERPROFILE": "C:/UserProfile/Fake",
			},
			expected: "C:/UserProfile/Fake/.kube",
		},
		"HOME used 1st when dir exists and writeable": {
			files: map[string]os.FileMode{
				"C:/Home/Fake":        0o777,
				"C:/HomePath/Fake":    0o777,
				"C:/UserProfile/Fake": 0o777,
			},
			env: map[string]string{
				"HOME":        "C:/Home/Fake",
				"HOMEDRIVE":   "C:/",
				"HOMEPATH":    "HomePath/Fake",
				"USERPROFILE": "C:/UserProfile/Fake",
			},
			expected: "C:/Home/Fake/.kube",
		},
		"USERPROFILE used 2nd when dir exists and writeable": {
			files: map[string]os.FileMode{
				"C:/Home/Fake":        0o000,
				"C:/HomePath/Fake":    0o777,
				"C:/UserProfile/Fake": 0o777,
			},
			env: map[string]string{
				"HOME":        "C:/Home/Fake",
				"HOMEDRIVE":   "C:/",
				"HOMEPATH":    "HomePath/Fake",
				"USERPROFILE": "C:/UserProfile/Fake",
			},
			expected: "C:/UserProfile/Fake/.kube",
		},
		"HOMEDRIVE\\HOMEPATH used 3rd when dir exists and writeable": {
			files: map[string]os.FileMode{
				"C:/Home/Fake":        0o000,
				"C:/HomePath/Fake":    0o777,
				"C:/UserProfile/Fake": 0o000,
			},
			env: map[string]string{
				"HOME":        "C:/Home/Fake",
				"HOMEDRIVE":   "C:/",
				"HOMEPATH":    "HomePath/Fake",
				"USERPROFILE": "C:/UserProfile/Fake",
			},
			expected: "C:/HomePath/Fake/.kube",
		},
		"HOME used 1st when dir exists": {
			files: map[string]os.FileMode{
				"C:/Home/Fake":        0o000,
				"C:/HomePath/Fake":    0o000,
				"C:/UserProfile/Fake": 0o000,
			},
			env: map[string]string{
				"HOME":        "C:/Home/Fake",
				"HOMEDRIVE":   "C:/",
				"HOMEPATH":    "HomePath/Fake",
				"USERPROFILE": "C:/UserProfile/Fake",
			},
			expected: "C:/Home/Fake/.kube",
		},
		"USERPROFILE used 2nd when dir exists": {
			files: map[string]os.FileMode{
				"C:/HomePath/Fake":    0o000,
				"C:/UserProfile/Fake": 0o000,
			},
			env: map[string]string{
				"HOME":        "C:/Home/Fake",
				"HOMEDRIVE":   "C:/",
				"HOMEPATH":    "HomePath/Fake",
				"USERPROFILE": "C:/UserProfile/Fake",
			},
			expected: "C:/UserProfile/Fake/.kube",
		},
		"HOMEDRIVE\\HOMEPATH used 3rd when dir exists": {
			files: map[string]os.FileMode{
				"C:/HomePath/Fake": 0o000,
			},
			env: map[string]string{
				"HOME":        "C:/Home/Fake",
				"HOMEDRIVE":   "C:/",
				"HOMEPATH":    "HomePath/Fake",
				"USERPROFILE": "C:/UserProfile/Fake",
			},
			expected: "C:/HomePath/Fake/.kube",
		},
		"Errors if no existing candidate": {
			files: map[string]os.FileMode{},
			env: map[string]string{
				"HOME":        "C:/Home/Fake",
				"HOMEDRIVE":   "C:/",
				"HOMEPATH":    "HomePath/Fake",
				"USERPROFILE": "C:/UserProfile/Fake",
			},
			expectedErrorIs: ErrNoKubeDir,
		},
		"Errors if no environment vars": {
			files: map[string]os.FileMode{
				"C:/Home/Fake":        0o000,
				"C:/HomePath/Fake":    0o000,
				"C:/UserProfile/Fake": 0o000,
			},
			env:             map[string]string{},
			expectedErrorIs: ErrNoKubeDir,
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			actual, err := doFindDefaultKubeDirWindows(
				func(s string) (string, bool) {
					value, ok := tc.env[s]
					return value, ok
				},
				func(s string) (fs.FileInfo, error) {
					mode, ok := tc.files[path.Clean(s)]
					if !ok {
						return nil, os.ErrNotExist
					}

					return &fakeFileInfo{
						mode: mode,
					}, nil
				},
			)

			if tc.expectedErrorIs != nil {
				require.ErrorIs(t, err, tc.expectedErrorIs)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expected, path.Clean(actual))
			}
		})
	}

}

func TestFindDefaultKubeDirPOSIX(t *testing.T) {
	t.Parallel()
	testcases := map[string]struct {
		env             map[string]string
		expected        string
		expectedErrorIs error
	}{
		"Found in HOME dir": {
			env: map[string]string{
				"HOME": "/home/fake",
			},

			expected: "/home/fake/.kube",
		},

		"Errors if HOME environment var is unset": {
			expectedErrorIs: ErrNoKubeDir,
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			actual, err := doFindDefaultKubeDirPOSIX(func(s string) (string, bool) {
				value, ok := tc.env[s]
				return value, ok
			})

			if tc.expectedErrorIs == nil {
				require.NoError(t, err)
				require.Equal(t, tc.expected, actual)
			} else {
				require.ErrorIs(t, err, tc.expectedErrorIs)
			}
		})
	}
}
