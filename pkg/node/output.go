package node

import (
	"encoding/json"
	"fmt"

	"github.com/devkcud/mess/pkg/utils"
)

func (n *Node) PrintNodeTree() {
	rootPath, rootNode := n.Collapse()

	if rootNode.Type == TypeDirectory {
		rootPath += "/"
	}

	fmt.Println(rootPath)

	for i, child := range rootNode.Children {
		last := i == len(rootNode.Children)-1
		child.print("", last)
	}
}

func (n *Node) print(prefix string, isLast bool) {
	collapsed, node := n.Collapse()

	branch := "├── "
	if isLast {
		branch = "└── "
	}

	if node.Type == TypeDirectory {
		collapsed += "/"
	}

	fmt.Printf("%s%s%s\n", prefix, branch, collapsed)

	nextPrefix := prefix
	if isLast {
		nextPrefix += "    "
	} else {
		nextPrefix += "│   "
	}

	for i, child := range node.Children {
		last := i == len(node.Children)-1
		child.print(nextPrefix, last)
	}
}

func (n *Node) PrintCommands() {
	currentUser := utils.CurrentUser

	var (
		sudoMkdirs, mkdirs   []string
		sudoTouches, touches []string
		sudoChmods, chmods   []string
		sudoChowns, chowns   []string
	)

	var walkDirs func(node *Node)
	walkDirs = func(node *Node) {
		if node.Type != TypeDirectory {
			return
		}
		if node.Parent != nil && node.Parent.Type == TypeDirectory && len(node.Parent.Children) == 1 {
			return
		}

		_, deepest := node.Collapse()
		if deepest.Type == TypeFile {
			deepest = deepest.Up()
		}
		fullPath := ExpandUserHome(deepest.BuildPathBackwards())

		if deepest.Parent == nil {
			for _, c := range deepest.Children {
				walkDirs(c)
			}
			return
		}

		if !utils.DoesPathExist(fullPath) {
			cmd := fmt.Sprintf("mkdir -p %s", fullPath)
			if deepest.NeedsElevation {
				sudoMkdirs = append(sudoMkdirs, "sudo "+cmd)
			} else {
				mkdirs = append(mkdirs, cmd)
			}
		}

		if deepest.Permission != utils.DirPerm {
			cmd := fmt.Sprintf("chmod %o %s", deepest.Permission, fullPath)
			if deepest.NeedsElevation {
				sudoChmods = append(sudoChmods, "sudo "+cmd)
			} else {
				chmods = append(chmods, cmd)
			}
		}

		if deepest.Owner != "" && deepest.Owner != currentUser {
			cmd := fmt.Sprintf("chown %s %s", deepest.Owner, fullPath)
			if deepest.NeedsElevation {
				sudoChowns = append(sudoChowns, "sudo "+cmd)
			} else {
				chowns = append(chowns, cmd)
			}
		}

		for _, c := range deepest.Children {
			walkDirs(c)
		}
	}
	walkDirs(n)

	var walkFiles func(node *Node)
	walkFiles = func(node *Node) {
		if node.Type == TypeFile {
			fullPath := ExpandUserHome(node.BuildPathBackwards())

			if node.Parent == nil {
				return
			}

			if !utils.DoesPathExist(fullPath) {
				cmd := fmt.Sprintf("touch %s", fullPath)
				if node.NeedsElevation {
					sudoTouches = append(sudoTouches, "sudo "+cmd)
				} else {
					touches = append(touches, cmd)
				}
			}

			if node.Permission != utils.FilePerm {
				cmd := fmt.Sprintf("chmod %o %s", node.Permission, fullPath)
				if node.NeedsElevation {
					sudoChmods = append(sudoChmods, "sudo "+cmd)
				} else {
					chmods = append(chmods, cmd)
				}
			}

			if node.Owner != "" && node.Owner != currentUser {
				cmd := fmt.Sprintf("chown %s %s", node.Owner, fullPath)
				if node.NeedsElevation {
					sudoChowns = append(sudoChowns, "sudo "+cmd)
				} else {
					chowns = append(chowns, cmd)
				}
			}
			return
		}
		for _, c := range node.Children {
			walkFiles(c)
		}
	}
	walkFiles(n)

	for _, cmd := range sudoMkdirs {
		fmt.Println(cmd)
	}
	for _, cmd := range mkdirs {
		fmt.Println(cmd)
	}
	for _, cmd := range sudoTouches {
		fmt.Println(cmd)
	}
	for _, cmd := range touches {
		fmt.Println(cmd)
	}
	for _, cmd := range sudoChmods {
		fmt.Println(cmd)
	}
	for _, cmd := range chmods {
		fmt.Println(cmd)
	}
	for _, cmd := range sudoChowns {
		fmt.Println(cmd)
	}
	for _, cmd := range chowns {
		fmt.Println(cmd)
	}
}

func (n *Node) PrintJSON(indent string) (string, error) {
	bytes, err := json.MarshalIndent(n, "", indent)
	if err != nil {
		return "", err
	}
	return string(bytes), err
}
