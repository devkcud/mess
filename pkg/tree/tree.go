package tree

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/devkcud/mess/pkg/utils"
)

type NodeType int

type Node struct {
	relativePath string

	parent   *Node `json:"-"`
	children []*Node
}

const (
	File NodeType = iota
	Dir
)

func New(baseDir string) (root *Node, current *Node) {
	if !filepath.IsAbs(baseDir) || baseDir == "~" {
		home, _ := os.UserHomeDir()

		if baseDir == "~" {
			baseDir = home
		} else {
			baseDir = utils.JoinPaths(home, baseDir)
		}
	}

	baseDir = utils.CleanPath(baseDir)

	parts := utils.SplitPath(baseDir)
	if len(parts) == 0 {
		parts = []string{baseDir}
	}

	root = &Node{
		relativePath: parts[0],
		parent:       nil,
		children:     []*Node{},
	}
	current = root

	for _, part := range parts[1:] {
		if part == "" {
			continue
		}
		current = current.AddDirectory(part)
	}

	return
}

func (n *Node) Up() *Node {
	return n.parent
}

func (n *Node) AddDirectory(path string) *Node {
	if n.children == nil {
		return n
	}

	if filepath.IsAbs(path) {
		path = utils.CleanPath(path)
		parts := utils.SplitPath(path)

		current, start := findStart(n, parts)

		for _, part := range parts[start:] {
			if part == "" {
				continue
			}
			current = current.AddDirectory(part)
		}
		return current
	}

	for _, c := range n.children {
		if c.relativePath == path {
			return c
		}
	}

	node := &Node{
		relativePath: path,
		parent:       n,
		children:     []*Node{},
	}
	n.children = append(n.children, node)
	return node
}

func (n *Node) AddFile(path string) *Node {
	if n.children == nil {
		return n
	}

	node := &Node{
		relativePath: path,
		parent:       n,
		children:     nil,
	}

	n.children = append(n.children, node)

	return node
}

func (n *Node) FullPath() string {
	if n.parent == nil {
		return n.relativePath
	}
	return filepath.Join(n.parent.FullPath(), n.relativePath)
}

func (n *Node) GetChildren() []*Node {
	return n.children
}

func findStart(root *Node, parts []string) (current *Node, start int) {
	for root.parent != nil {
		root = root.parent
	}
	current = root

	start = 0
	if len(parts) > 0 && parts[0] == root.relativePath {
		start = 1
	}

	return
}

func (n *Node) PrintTree() {
	n.print("", true, 0)
}

func (n *Node) print(prefix string, isTail bool, depth int) {
	node := n
	parts := []string{node.relativePath}
	for len(node.children) == 1 {
		node = node.children[0]
		parts = append(parts, node.relativePath)
	}
	label := utils.CleanPath(strings.Join(parts, "/"))
	if node.children != nil {
		if label != "/" {
			label += "/"
		}
	}

	if depth == 0 {
		fmt.Println(label)
	} else {
		branch := "├── "
		if isTail {
			branch = "└── "
		}
		fmt.Printf("%s%s%s\n", prefix, branch, label)
	}

	for i, child := range node.children {
		var newPref string
		if depth == 0 {
			newPref = ""
		} else if isTail {
			newPref = prefix + "    "
		} else {
			newPref = prefix + "│   "
		}
		child.print(newPref, i == len(node.children)-1, depth+1)
	}
}

func (n *Node) PrintShellCommands() {
	var (
		leafDirs  []string
		filePaths []string
	)

	var collect func(node *Node, curr string)
	collect = func(node *Node, curr string) {
		full := filepath.Join(curr, node.relativePath)

		if node.children == nil {
			filePaths = append(filePaths, full)
			return
		}

		hasSubDir := false
		for _, c := range node.children {
			if c.children != nil {
				hasSubDir = true
				break
			}
		}

		if !hasSubDir {
			leafDirs = append(leafDirs, full)
		}

		for _, c := range node.children {
			collect(c, full)
		}
	}

	collect(n, "")

	findParent := func(p string) string {
		for {
			info, err := os.Stat(p)
			if err != nil || !info.IsDir() {
				p = filepath.Dir(p)
				continue
			}
			break
		}
		return p
	}

	sort.Slice(leafDirs, func(i, j int) bool {
		di := strings.Count(leafDirs[i], string(filepath.Separator))
		dj := strings.Count(leafDirs[j], string(filepath.Separator))
		if di != dj {
			return di > dj
		}
		return strings.ToLower(leafDirs[i]) < strings.ToLower(leafDirs[j])
	})
	for _, dir := range leafDirs {
		parent := findParent(dir)
		needsElev := utils.NeedsElevation(parent)
		prefix := ""
		if needsElev {
			prefix = "sudo "
		}
		fmt.Printf("%smkdir -p %s\n", prefix, dir)
	}

	sort.Slice(filePaths, func(i, j int) bool {
		return strings.ToLower(filePaths[i]) < strings.ToLower(filePaths[j])
	})
	for _, file := range filePaths {
		parent := findParent(filepath.Dir(file))
		needsElev := utils.NeedsElevation(parent)
		prefix := ""
		if needsElev {
			prefix = "sudo "
		}
		fmt.Printf("%stouch %s\n", prefix, file)
	}
}
