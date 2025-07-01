package node

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

type NodeInformation struct {
	Name string

	Owner      string
	Permission *os.FileMode
}

var ErrEmptyName = errors.New("name is empty")

func ParsePathPart(part string) (*NodeInformation, error) {
	info := new(NodeInformation)

	indexAt := strings.LastIndex(part, "@")
	indexPercentage := strings.LastIndex(part, "%")

	endName := len(part)
	if indexAt != -1 && (indexPercentage == -1 || indexAt < indexPercentage) {
		endName = indexAt
	} else if indexPercentage != -1 && (indexAt == -1 || indexPercentage < indexAt) {
		endName = indexPercentage
	}
	info.Name = part[:endName]

	if info.Name == "" {
		return nil, ErrEmptyName
	}

	if indexAt != -1 {
		start := indexAt + 1
		end := len(part)

		if indexPercentage != -1 && indexPercentage > indexAt {
			end = indexPercentage
		}

		info.Owner = part[start:end]
	}

	if indexPercentage != -1 {
		start := indexPercentage + 1
		end := len(part)

		if indexAt != -1 && indexAt > indexPercentage {
			end = indexAt
		}

		permissionString := part[start:end]

		m, err := strconv.ParseUint(permissionString, 8, 32)
		if err != nil {
			return info, err
		}
		perm := os.FileMode(m)
		info.Permission = &perm
	}

	return info, nil
}
