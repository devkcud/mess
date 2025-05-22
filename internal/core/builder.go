package core

import (
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/devkcud/mess/pkg/messlog"
	"github.com/devkcud/mess/pkg/tree"
	"github.com/devkcud/mess/pkg/utils"
)

type Summary struct {
	DirectoriesCreated int
	FilesCreated       int
	Failures           int
	Successes          int
}

type builder struct {
	logger *messlog.Logger

	dryRun bool
	echo   bool

	rootTree    *tree.Node
	currentTree *tree.Node
}

func NewBuilder(base string, logger *messlog.Logger, dry, echo bool) *builder {
	root, current := tree.New(base)

	return &builder{
		logger:      logger,
		dryRun:      dry,
		echo:        echo,
		rootTree:    root,
		currentTree: current,
	}
}

func (b *builder) addDirectory(path string) {
	newRoot := b.currentTree.AddDirectory(utils.CleanPath(path))
	b.currentTree = newRoot
}

func (b *builder) addFile(path string) {
	b.currentTree.AddFile(utils.CleanPath(path))
}

func (b *builder) ProcessToken(token string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			b.logger.Trace("Panic detected: %v\n%s", r, debug.Stack())
			err = fmt.Errorf("panic occurred: %v", r)
		}
	}()

	b.logger.Info("Processing token %q", token)
	b.logger.Trace("Tokenized: %s", token)

	switch {
	case token == "..":
		b.currentTree = b.currentTree.Up()
		b.logger.Trace("Stack tree moved up one parent: %s", token)

	case strings.HasSuffix(token, utils.OSPathSeparator):
		b.addDirectory(token)
		b.logger.Trace("Stack tree added one directory: %s", token)

	case strings.Contains(token, utils.OSPathSeparator):
		dir, file := utils.SeparatePath(token)

		cur := b.currentTree
		b.addDirectory(dir)
		b.addFile(file)
		b.currentTree = cur
		b.logger.Trace("Stack tree added one directory and one file: %s", token)
		b.logger.Trace("`currentTree` remains intact")

	default:
		b.addFile(token)
		b.logger.Trace("Stack tree added one file: %s", token)
	}

	return
}

func (b *builder) PrintDryRunTree() {
	b.rootTree.PrintTree()
}

func (b *builder) PrintEchoFiles() {
	b.rootTree.PrintShellCommands()
}

func (b *builder) BuildFiles() error {
	var walk func(n *tree.Node) error
	walk = func(n *tree.Node) error {
		p := n.FullPath()
		b.logger.Trace("Building: %s", p)

		if n.GetChildren() != nil {
			if err := utils.WriteDirectories(p); err != nil {
				return fmt.Errorf("mkdir %s: %w", p, err)
			}
			for _, c := range n.GetChildren() {
				if err := walk(c); err != nil {
					return err
				}
			}
		} else {
			if err := utils.WriteFile(p, ""); err != nil {
				return fmt.Errorf("write %s: %w", p, err)
			}
		}

		return nil
	}

	return walk(b.rootTree)
}
