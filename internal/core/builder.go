package core

import (
	"fmt"
	"path/filepath"
	"runtime/debug"
	"strings"

	"github.com/devkcud/mess/pkg/messlog"
	"github.com/devkcud/mess/pkg/node"
	"github.com/devkcud/mess/pkg/utils"
)

type builder struct {
	logger *messlog.Logger

	dryRun bool
	echo   bool

	root *node.Node
}

func NewBuilder(base string, logger *messlog.Logger, dry, echo bool) *builder {
	return &builder{
		logger: logger,
		dryRun: dry,
		echo:   echo,
		root:   node.New(base),
	}
}

func (b *builder) addDirectory(path string) {
	b.logger.Info("Added directory %s", path)
	b.root = b.root.AddDirectory(path)
}

func (b *builder) addFile(path string) {
	b.logger.Info("Added file %s", path)
	b.root.AddFile(path)
}

func (b *builder) ProcessToken(token string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			b.logger.Trace("Panic detected: %v\n%s", r, debug.Stack())
			err = fmt.Errorf("panic occurred: %v", r)
		}
	}()

	switch {
	case token == "..":
		b.logger.Debug("Rule found: ..")
		b.root = b.root.Up()
		b.logger.Trace("Stack tree moved up one parent: %s", token)

	case strings.HasSuffix(token, utils.OSPathSeparator):
		b.logger.Debug("Rule found: dir/")
		b.addDirectory(token)
		b.logger.Trace("Stack tree added one directory: %s", token)

	case strings.Contains(token, utils.OSPathSeparator):
		b.logger.Debug("Rule found: dir/file")
		dir, file := filepath.Split(token)

		cur := b.root
		b.addDirectory(dir)
		b.addFile(file)
		b.root = cur

		b.logger.Trace("Stack tree added one directory and one file: %s", token)
		b.logger.Trace("`currentTree` remains intact")

	default:
		b.logger.Debug("Rule found: file")
		b.addFile(token)
		b.logger.Trace("Stack tree added one file: %s", token)
	}

	return
}

func (b *builder) PrintDryRunTree() {
	b.root.Root().PrintNodeTree()
}

func (b *builder) PrintEchoFiles() {
	b.root.Root().PrintCommands()
}
func (b *builder) PrintJSON() error {
	j, err := b.root.Root().PrintJSON("    ")
	defer fmt.Println(j)
	if err != nil {
		return err
	}
	return nil
}

func (b *builder) BuildFiles() error {
	b.logger.Debug("Building files...")
	defer b.logger.Debug("Build done!")

	return b.root.Root().BuildFiles()
}
