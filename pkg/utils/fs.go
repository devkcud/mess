package utils

import (
	"os"
	"path/filepath"
)

const (
	DirPerm  = 0o755
	FilePerm = 0o644
)

const OSPathSeparator = string(os.PathSeparator)

func JoinPaths(paths ...string) string {
	return filepath.Clean(filepath.Join(paths...))
}

func SplitPath(path string) (dir, file string) {
	return filepath.Split(filepath.Clean(path))
}

func WriteFile(path string, data string) error {
	return os.WriteFile(filepath.Clean(path), []byte(data), FilePerm)
}

func WriteDirectory(paths ...string) (string, error) {
	full := JoinPaths(paths...)
	return full, os.MkdirAll(full, DirPerm)
}
