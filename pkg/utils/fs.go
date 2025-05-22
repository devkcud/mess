package utils

import (
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"golang.org/x/sys/unix"
)

const (
	DirPerm  = 0o755
	FilePerm = 0o644
)

const OSPathSeparator = string(os.PathSeparator)

func CleanPath(path string) string {
	return filepath.Clean(path)
}

func JoinPaths(paths ...string) string {
	return filepath.Join(paths...)
}

func SplitPath(path string) []string {
	parts := strings.Split(path, OSPathSeparator)
	if len(parts) > 0 && parts[0] == "" {
		parts[0] = OSPathSeparator
	}

	return parts
}

func SeparatePath(path string) (dir, file string) {
	return filepath.Split(path)
}

func WriteFile(path string, data string) error {
	return os.WriteFile(path, []byte(data), FilePerm)
}

func WriteDirectories(paths ...string) error {
	return os.MkdirAll(JoinPaths(paths...), DirPerm)
}

func NeedsElevation(dir string) bool {
	if os.Geteuid() == 0 {
		return false
	}

	if err := unix.Access(dir, unix.W_OK); err == nil {
		return false
	} else if err == syscall.EACCES {
		return true
	} else {
		return false
	}
}
