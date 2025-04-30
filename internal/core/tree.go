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
	p := filepath.Clean(path)
	parts := strings.Split(p, utils.OSPathSeparator)
	if parts[0] == "" {
		parts = parts[1:]
	}
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

func (root *node) PrintCollapsedTree() {
	kids := sortedChildren(root)
	for i, child := range kids {
		printCollapsed(child, "", i == len(kids)-1)
	}
}

func printCollapsed(n *node, prefix string, isLast bool) {
	parts := []string{n.Name}
	curr := n
	for curr.IsDir && len(curr.Children) == 1 {
		only := firstChild(curr)
		if !only.IsDir {
			break
		}
		parts = append(parts, only.Name)
		curr = only
	}

	connector := "├── "
	if isLast {
		connector = "└── "
	}
	name := filepath.Join(parts...)
	if curr.IsDir {
		name += "/"
	}
	fmt.Printf("%s%s%s\n", prefix, connector, name)

	kids := sortedChildren(curr)
	newPrefix := prefix
	if isLast {
		newPrefix += "    "
	} else {
		newPrefix += "│   "
	}
	for i, child := range kids {
		printCollapsed(child, newPrefix, i == len(kids)-1)
	}
}

func sortedChildren(n *node) []*node {
	out := make([]*node, 0, len(n.Children))
	for _, c := range n.Children {
		out = append(out, c)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].IsDir != out[j].IsDir {
			return out[i].IsDir
		}
		return out[i].Name < out[j].Name
	})
	return out
}

func firstChild(n *node) *node {
	for _, c := range n.Children {
		return c
	}
	return nil
}
