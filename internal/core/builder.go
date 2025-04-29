package core

import (
	"fmt"
	"path/filepath"
	"runtime/debug"
	"strings"

	"github.com/devkcud/mess/pkg/messlog"
	"github.com/devkcud/mess/pkg/utils"
)

type Summary struct {
	DirectoriesCreated int
	FilesCreated       int
	Failures           int
	Successes          int
}

type builder struct {
	wd      string
	logger  *messlog.Logger
	dryRun  bool
	echo    bool
	summary *Summary
	tree    *node
}

func NewBuilder(base string, logger *messlog.Logger, dry bool, summary *Summary, echo bool) *builder {
	return &builder{
		logger:  logger,
		dryRun:  dry,
		echo:    echo,
		summary: summary,
		tree: &node{
			Name:     base,
			Children: make(map[string]*node),
			IsDir:    true,
		},
	}
}

func (b *builder) addFolder(wd string, path string, setWD bool) (string, error) {
	full := utils.JoinPaths(wd, path)
	b.logger.Info("Creating folder %q (%s)", path, full)

	if b.echo {
		fmt.Printf("mkdir -p %s\n", utils.JoinPaths(wd, path))
		if setWD {
			wd = utils.JoinPaths(wd, path)
		}
		return wd, nil
	}

	if b.dryRun {
		relPath := utils.JoinPaths(wd, path)
		b.tree.addToTree(relPath, true)

		if setWD {
			b.logger.Info("[DRY-RUN] Adding folder %q to working directory", path)
			b.logger.Debug("[DRY-RUN] Current working directory %q", relPath)

			wd = relPath
		} else {
			b.logger.Info("Skipping folder %q append to working directory", path)
		}
	} else {
		newPath, err := utils.WriteDirectory(full)
		if err != nil {
			return wd, err
		}

		if setWD {
			b.logger.Info("Adding folder %q to working directory", path)
			b.logger.Debug("Current working directory %q", newPath)

			wd = newPath
		} else {
			b.logger.Info("Skipping folder %q append to working directory", path)
		}
	}

	if b.summary != nil {
		b.summary.DirectoriesCreated++
	}

	return wd, nil
}

func (b *builder) addFile(wd string, path string) error {
	full := utils.JoinPaths(wd, path)

	b.logger.Debug("Creating file %s (%s)", path, full)

	if b.echo {
		fmt.Printf("touch %s\n", utils.JoinPaths(wd, path))
		return nil
	}

	if b.dryRun {
		relPath := utils.JoinPaths(wd, path)
		b.tree.addToTree(relPath, false)
	} else {
		err := utils.WriteFile(full, "")
		if err != nil {
			return err
		}
	}

	if b.summary != nil {
		b.summary.FilesCreated++
	}

	return nil
}

func (b *builder) ProcessToken(token string) (err error) {
	originalWD := b.wd

	defer func() {
		if r := recover(); r != nil {
			b.logger.Trace("Panic detected: %v\n%s", r, debug.Stack())
			err = fmt.Errorf("panic occurred: %v", r)
		}

		if err != nil {
			b.wd = originalWD
			b.logger.Error("Token %q failed: %v", token, err)
			b.logger.Trace("Restoring working directory from %q to original state: %q", b.wd, originalWD)

			if b.summary != nil {
				b.summary.Failures++
			}
		} else {
			b.logger.Info("Token %q succeeded", token)

			if b.summary != nil {
				b.summary.Successes++
			}
		}
	}()

	b.logger.Info("Processing token %q", token)

	switch {
	case token == "..":
		if b.echo {
			fmt.Println("cd ..")
		}
		b.wd = filepath.Dir(b.wd)

	case strings.HasSuffix(token, utils.OSPathSeparator):
		b.wd, err = b.addFolder(b.wd, token, true)

	case strings.Contains(token, utils.OSPathSeparator):
		dirPart, filePart := utils.SplitPath(token)

		folderPath := utils.JoinPaths(b.wd, dirPart)
		b.logger.Info("Creating folder %q (%s)", dirPart, folderPath)
		if b.dryRun {
			b.tree.addToTree(folderPath, true)
		} else {
			if _, err = utils.WriteDirectory(folderPath); err != nil {
				return err
			}
		}
		if b.summary != nil {
			b.summary.DirectoriesCreated++
		}

		err = b.addFile(folderPath, filePart)

	default:
		err = b.addFile(b.wd, token)
	}

	return
}

func (b *builder) PrintDryRunTree() {
	b.tree.printTree("", true)
}
