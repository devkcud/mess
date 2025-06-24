package node

import (
	"path/filepath"
	"strings"

	"github.com/devkcud/mess/pkg/utils"
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
