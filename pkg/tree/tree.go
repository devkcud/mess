package tree

import (
	"fmt"
	"os"
	"path/filepath"
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
	if baseDir == "~" || strings.HasPrefix(baseDir, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "warning: could not get home dir: %v\n", err)
		} else {
			if baseDir == "~" {
				baseDir = home
			} else {
				baseDir = utils.JoinPaths(home, baseDir[2:])
			}
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

func (n *Node) AddDirectory(path string) *Node {
	if n.children == nil {
		return n
	}

	if filepath.IsAbs(path) || strings.HasPrefix(path, "~/") {
		if strings.HasPrefix(path, "~/") {
			home, err := os.UserHomeDir()
			if err == nil {
				path = filepath.Join(home, path[2:])
			}
		}
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

	if filepath.IsAbs(path) || strings.HasPrefix(path, "~/") {
		if strings.HasPrefix(path, "~/") {
			home, err := os.UserHomeDir()
			if err == nil {
				path = utils.JoinPaths(home, path[2:])
			}
		}
		path = utils.CleanPath(path)
		parts := utils.SplitPath(path)

		current, start := findStart(n, parts)

		for _, part := range parts[start : len(parts)-1] {
			if part == "" {
				continue
			}
			current = current.AddDirectory(part)
		}

		fileName := parts[len(parts)-1]
		node := &Node{
			relativePath: fileName,
			parent:       current,
			children:     nil,
		}
		current.children = append(current.children, node)
		return node
	}

	node := &Node{
		relativePath: path,
		parent:       n,
		children:     nil,
	}
	n.children = append(n.children, node)
	return node
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
	if len(node.children) > 0 {
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

func (n *Node) FullPath() string {
	if n.parent == nil {
		return n.relativePath
	}
	return filepath.Join(n.parent.FullPath(), n.relativePath)
}

func (n *Node) GetChildren() []*Node {
	return n.children
}

func (n *Node) PrintFiles() {
	n.walkFiles(func(file *Node) {
		p := file.FullPath()

		dir := filepath.Dir(p)
		for {
			info, err := os.Stat(dir)
			if err != nil {
				if os.IsNotExist(err) {
					dir = filepath.Dir(dir)
					continue
				}
			} else if !info.IsDir() {
				dir = filepath.Dir(dir)
				continue
			}
			break
		}

		needsElev := utils.NeedsElevation(dir)
		if needsElev {
			fmt.Printf("%s (sudo required)\n", p)
		} else {
			fmt.Println(p)
		}
	})
}

func (n *Node) walkFiles(fn func(*Node)) {
	if n.children == nil {
		fn(n)
		return
	}
	for _, c := range n.children {
		c.walkFiles(fn)
	}
}

func (n *Node) PrintShellCommands() {
	n.walkShell("")
}

func (n *Node) walkShell(prefixPath string) {
	node := n
	parts := []string{}
	for {
		parts = append(parts, node.relativePath)
		if len(node.children) == 1 && node.children[0].children != nil {
			node = node.children[0]
			continue
		}
		break
	}

	collapsed := filepath.Join(parts...)
	full := filepath.Join(prefixPath, collapsed)

	dir := full
	if node.children == nil {
		dir = filepath.Dir(full)
	}
	for {
		info, err := os.Stat(dir)
		if err != nil {
			if os.IsNotExist(err) {
				dir = filepath.Dir(dir)
				continue
			}
		} else if !info.IsDir() {
			dir = filepath.Dir(dir)
			continue
		}
		break
	}
	needsElev := utils.NeedsElevation(dir)
	sudo := ""
	if needsElev {
		sudo = "sudo "
	}

	if node.children != nil {
		if full != string(filepath.Separator) {
			fmt.Printf("%smkdir -p %s\n", sudo, full)
		}
	} else {
		fmt.Printf("%stouch %s\n", sudo, full)
	}

	for _, c := range node.children {
		c.walkShell(full)
	}
}
