package core

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/devkcud/mess/pkg/utils"
)

type node struct {
	Name     string
	Children map[string]*node
	IsDir    bool
}

func (root *node) addToTree(path string, isDir bool) {
	parts := strings.Split(filepath.Clean(path), utils.OSPathSeparator)
	curr := root

	for i, part := range parts {
		child, exists := curr.Children[part]
		if !exists {
			child = &node{
				Name:     part,
				Children: make(map[string]*node),
				IsDir:    i != len(parts)-1 || isDir,
			}
			curr.Children[part] = child
		}
		curr = child
	}
}

func (n *node) printTree(prefix string, isLast bool) {
	connector := "├── "
	if isLast {
		connector = "└── "
	}

	pref := "F"
	if n.IsDir {
		pref = "D"
	}

	if n.Name != "." {
		fmt.Printf("%s%s%s %s\n", prefix, connector, pref, n.Name)
	}
	children := make([]*node, 0, len(n.Children))
	for _, c := range n.Children {
		children = append(children, c)
	}
	sort.Slice(children, func(i, j int) bool {
		if children[i].IsDir != children[j].IsDir {
			return children[i].IsDir
		}
		return children[i].Name < children[j].Name
	})

	for i, child := range children {
		newPrefix := prefix
		if isLast {
			newPrefix += "    "
		} else {
			newPrefix += "│   "
		}
		child.printTree(newPrefix, i == len(children)-1)
	}
}
