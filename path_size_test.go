package code_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"code"
)

// Fixture layout under testdata (sizes in bytes):
//
//	test.txt               6
//	empty.txt              0
//	dir/a.txt              4
//	dir/b.txt             10
//	dir/nested/inner.txt  11   (second level, must not be counted)
//	onlydirs/sub/.gitkeep  0   (onlydirs has no files on the first level)
//	hidden/visible.txt     4
//	hidden/.secret         7   (hidden, counted only with all=true)
//	hidden/.git/config     2   (inside a hidden directory, second level)
func TestGetPathSize(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		path      string
		recursive bool
		all       bool
		want      string
		wantErr   error
	}{
		{
			name: "file",
			path: filepath.Join("testdata", "test.txt"),
			want: "6B",
		},
		{
			name: "empty file",
			path: filepath.Join("testdata", "empty.txt"),
			want: "0B",
		},
		{
			name: "dir counts first level files only",
			path: filepath.Join("testdata", "dir"),
			want: "14B",
		},
		{
			// a.txt (4B) + b.txt (10B) + nested/inner.txt (11B).
			name:      "dir recursive",
			path:      filepath.Join("testdata", "dir"),
			recursive: true,
			want:      "25B",
		},
		{
			name: "dir with trailing separator",
			path: filepath.Join("testdata", "dir") + string(filepath.Separator),
			want: "14B",
		},
		{
			name: "dir with only subdirectories",
			path: filepath.Join("testdata", "onlydirs"),
			want: "0B",
		},
		{
			name: "dir skips hidden files by default",
			path: filepath.Join("testdata", "hidden"),
			want: "4B",
		},
		{
			name: "dir includes hidden files with all",
			path: filepath.Join("testdata", "hidden"),
			all:  true,
			want: "11B",
		},
		{
			// Hidden directories are not descended into without all.
			name:      "dir recursive skips hidden dirs by default",
			path:      filepath.Join("testdata", "hidden"),
			recursive: true,
			want:      "4B",
		},
		{
			// visible.txt (4B) + .secret (7B) + .git/config (2B).
			name:      "dir recursive includes hidden dirs with all",
			path:      filepath.Join("testdata", "hidden"),
			recursive: true,
			all:       true,
			want:      "13B",
		},
		{
			// An explicitly named hidden path is reported as 0B unless all is set.
			name: "hidden file passed explicitly is ignored by default",
			path: filepath.Join("testdata", "hidden", ".secret"),
			want: "0B",
		},
		{
			name: "hidden file passed explicitly is measured with all",
			path: filepath.Join("testdata", "hidden", ".secret"),
			all:  true,
			want: "7B",
		},
		{
			name: "hidden dir passed explicitly is ignored by default",
			path: filepath.Join("testdata", "hidden", ".git"),
			want: "0B",
		},
		{
			name: "hidden dir passed explicitly is measured with all",
			path: filepath.Join("testdata", "hidden", ".git"),
			all:  true,
			want: "2B",
		},
		{
			// "." and ".." start with a dot but are not hidden entries.
			name: "current dir reference is not hidden",
			path: filepath.Join("testdata", "dir") + string(filepath.Separator) + ".",
			want: "14B",
		},
		{
			// testdata/dir/.. resolves to testdata: test.txt (6B) + empty.txt (0B).
			name: "parent dir reference is not hidden",
			path: filepath.Join("testdata", "dir") + string(filepath.Separator) + "..",
			want: "6B",
		},
		{
			name:    "missing path",
			path:    filepath.Join("testdata", "missing.txt"),
			wantErr: fs.ErrNotExist,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := code.GetPathSize(tt.path, tt.recursive, false, tt.all)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				require.Empty(t, got)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestGetPathSize_EmptyDir(t *testing.T) {
	t.Parallel()

	got, err := code.GetPathSize(t.TempDir(), false, false, false)

	require.NoError(t, err)
	require.Equal(t, "0B", got)
}

func TestGetPathSize_LargeFile(t *testing.T) {
	t.Parallel()

	path := writeTempFile(t, "large.bin", 1500)

	got, err := code.GetPathSize(path, false, false, false)

	require.NoError(t, err)
	require.Equal(t, "1500B", got)
}

func TestGetPathSize_Human(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		size int
		want string
	}{
		{name: "zero", size: 0, want: "0B"},
		{name: "bytes stay integer", size: 1023, want: "1023B"},
		{name: "exactly one kilobyte", size: 1024, want: "1.0KB"},
		{name: "fractional kilobytes", size: 1536, want: "1.5KB"},
		{name: "rounds to one decimal", size: 1100, want: "1.1KB"},
		{name: "megabytes", size: 3 * 1024 * 1024, want: "3.0MB"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			path := writeTempFile(t, "file.bin", tt.size)

			got, err := code.GetPathSize(path, false, true, false)

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

// writeTempFile creates a zero-filled file of the given size in a temporary
// directory that is removed when the test finishes and returns its path.
func writeTempFile(t *testing.T, name string, size int) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), name)
	require.NoError(t, os.WriteFile(path, make([]byte, size), 0o600))

	return path
}
