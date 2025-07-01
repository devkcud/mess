package utils

import (
	"os"
	"os/user"
	"strconv"
	"strings"
	"syscall"

	"golang.org/x/sys/unix"
)

const (
	DirPerm  = os.FileMode(0o755)
	FilePerm = os.FileMode(0o644)
)

const OSPathSeparator = string(os.PathSeparator)

var UserHomeDirectory = func() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return home
}()

func SplitPath(path string) []string {
	parts := strings.Split(path, OSPathSeparator)
	if len(parts) > 0 && parts[0] == "" {
		parts[0] = OSPathSeparator
	}

	return parts
}

func DoesPathExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func NeedsElevation(path string) bool {
	if os.Geteuid() == 0 {
		return false
	}

	if err := unix.Access(path, unix.W_OK); err == nil {
		return false
	} else if err == syscall.EACCES {
		return true
	} else {
		return false
	}
}

func GetOwnerInfo(path string) (uid uint32, username string) {
	info, _ := os.Stat(path)
	stat := info.Sys().(*syscall.Stat_t)
	uid = stat.Uid

	u, err := user.LookupId(strconv.Itoa(int(uid)))
	if err != nil {
		return 0, ""
	}
	username = u.Username

	return
}
