package node

import (
	"path/filepath"

	"github.com/devkcud/mess/pkg/utils"
)

func (n *Node) insertChild(path string, nodeType NodeType) *Node {
	current := n

	if filepath.IsAbs(path) {
		current = n.Root()
		path = path[1:]
	}
	if current.Type != TypeDirectory {
		return nil
	}

	parts := utils.SplitPath(path)
	for i, part := range parts {
		if part == "/" {
			current = current.Root()
			continue
		}

		if part == "" || part == "." {
			continue
		}

		if i == 0 {
			part = ExpandUserHome(part)
		}

		if part == ".." {
			if current.Parent != nil {
				current = current.Parent
			}
			continue
		}

		information, err := ParsePathPart(part)
		if err != nil {
			panic(err) // TODO: Handle errors better
		}
		part = information.Name

		found := false

		for _, child := range current.Children {
			if child.Name == part {
				current = child
				found = true
				break
			}
		}

		if !found {
			newType := TypeDirectory
			if i == len(parts)-1 {
				newType = nodeType
			}

			perm := utils.FilePerm
			if newType == TypeDirectory {
				perm = utils.DirPerm
			}

			if information.Permission != nil {
				perm = *information.Permission
			}

			newNode := &Node{
				Name:       part,
				Type:       newType,
				Permission: perm,
				Parent:     current,
				Children:   []*Node{},
			}

			path := newNode.BuildPathBackwards()
			pathExists := utils.DoesPathExist(path)

			if pathExists {
				newNode.NeedsElevation = utils.NeedsElevation(path)
			} else {
				newNode.NeedsElevation = newNode.Parent.NeedsElevation
			}

			if pathExists {
				_, newNode.Owner = utils.GetOwnerInfo(path)
			} else {
				if newNode.NeedsElevation {
					newNode.Owner = utils.RootUser
				} else {
					newNode.Owner = utils.CurrentUser
				}
			}

			if information.Owner != "" {
				newNode.Owner = information.Owner
			}

			current.Children = append(current.Children, newNode)
			current = newNode
		}
	}

	return current
}

func (n *Node) AddFile(file string) *Node {
	return n.insertChild(file, TypeFile)
}

func (n *Node) AddDirectory(directory string) *Node {
	return n.insertChild(directory, TypeDirectory)
}
