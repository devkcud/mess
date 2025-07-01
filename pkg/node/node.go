package node

import (
	"os"
	"path/filepath"

	"github.com/devkcud/mess/pkg/utils"
)

type NodeType int

type Node struct {
	Name string

	Type NodeType

	Permission     os.FileMode
	NeedsElevation bool
	Owner          string

	Parent   *Node `json:"-"`
	Children []*Node
}

const (
	TypeDirectory NodeType = iota
	TypeFile
)

func New(baseDirectory string) *Node {
	if !filepath.IsAbs(baseDirectory) {
		baseDirectory = filepath.Join(utils.UserHomeDirectory, baseDirectory)
	}

	baseDirectory = filepath.Clean(baseDirectory)

	root := &Node{
		Name:           "/",
		Type:           TypeDirectory,
		Permission:     utils.DirPerm,
		NeedsElevation: true,
		Owner:          utils.RootUser,
		Parent:         nil,
		Children:       []*Node{},
	}
	current := root
	for _, part := range utils.SplitPath(baseDirectory)[1:] {
		if part == "" {
			continue
		}
		part = ExpandUserHome(part)
		current = current.AddDirectory(part)
	}

	return current
}

func (nt NodeType) String() (name string) {
	switch nt {
	case TypeDirectory:
		name = "directory"
	case TypeFile:
		name = "file"
	}
	return
}
