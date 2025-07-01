package node

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/devkcud/mess/pkg/utils"
)

type simpleNode struct {
	fpath string
	owner string
	perms os.FileMode
}

var (
	ErrNotDirectory = errors.New("path is a file not a directory")
	ErrIsDirectory  = errors.New("path is a directory not a file")
)

func (n *Node) Up() *Node {
	parent := n.Parent
	if parent == nil {
		return n
	}
	return parent
}

func (n *Node) Root() *Node {
	if n.Parent != nil {
		return n.Parent.Root()
	}

	return n
}

func (n *Node) UserHome() *Node {
	return n.Root().AddDirectory(utils.UserHomeDirectory)
}

func ExpandUserHome(path string) string {
	if path == "~" {
		path = utils.UserHomeDirectory
	} else if strings.HasPrefix(path, "~/") {
		path = filepath.Join(utils.UserHomeDirectory, path[2:])
	}

	return path
}

func (n *Node) BuildPathBackwards() string {
	path := n.Name
	current := n
	for {
		if current.Parent == nil {
			break
		}
		path = filepath.Join(current.Parent.Name, path)
		current = current.Parent
	}

	return path
}

func (n *Node) Collapse() (string, *Node) {
	name := n.Name
	for len(n.Children) == 1 {
		n = n.Children[0]
		name = filepath.Join(name, n.Name)
	}
	return name, n
}

func (n *Node) BuildFiles() error {
	dirs := make([]simpleNode, 0)
	files := make([]simpleNode, 0)

	var walk func(node *Node) error
	walk = func(node *Node) error {
		sn := simpleNode{
			fpath: node.BuildPathBackwards(),
			owner: node.Owner,
			perms: node.Permission,
		}

		info, err := os.Stat(sn.fpath)
		if node.Type == TypeDirectory {
			if err == nil {
				if !info.IsDir() {
					return fmt.Errorf("%w: %s", ErrNotDirectory, sn.fpath)
				}
			} else if !os.IsNotExist(err) {
				return err
			} else {
				dirs = append(dirs, sn)
			}

			for _, child := range node.Children {
				if err := walk(child); err != nil {
					return err
				}
			}
			return nil
		}

		if err == nil {
			if info.IsDir() {
				return fmt.Errorf("%w: %s", ErrIsDirectory, sn.fpath)
			}
		} else if !os.IsNotExist(err) {
			return err
		} else {
			files = append(files, sn)
		}

		return nil
	}

	if err := walk(n); err != nil {
		return err
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir.fpath, dir.perms); err != nil {
			return fmt.Errorf("%w: %s", err, dir.fpath)
		}
	}

	for _, file := range files {
		if err := os.WriteFile(file.fpath, []byte(""), file.perms); err != nil {
			return fmt.Errorf("%w: %s", err, file.fpath)
		}
	}

	return nil
}
