package core

import (
	"fmt"
	flag "github.com/spf13/pflag"
	"os"
)

type flagWrapper struct {
	fs *flag.FlagSet
}

func simpleHelp(fs *flag.FlagSet) {
	fmt.Fprintf(fs.Output(), "Usage: %s [-flags] <..|dir/|dir/file|file>[@<owner>|%%<perms>]...\n", fs.Name())
}

func NewCLI() *flagWrapper {
	fs := flag.NewFlagSet("mess", flag.ExitOnError)
	fs.Usage = func() {
		fs.SetOutput(os.Stderr)

		simpleHelp(fs)
		fmt.Fprint(fs.Output(), "\nFlags:\n")
		fs.PrintDefaults()
	}

	return &flagWrapper{fs}
}

func (fw *flagWrapper) Parse() ([]string, error) {
	if err := fw.fs.Parse(os.Args[1:]); err != nil {
		return nil, err
	}

	return fw.Args(), nil
}

func (fw *flagWrapper) Args() []string {
	return fw.fs.Args()
}

func (fw *flagWrapper) HelpExit(simple bool) {
	if simple {
		simpleHelp(fw.fs)
	} else {
		fw.fs.Usage()
	}
	os.Exit(1)
}

func (fw *flagWrapper) Bool(name string, def bool, usage string) *bool {
	return fw.fs.Bool(name, def, usage)
}

func (fw *flagWrapper) Int(name string, def int, usage string) *int {
	return fw.fs.Int(name, def, usage)
}

func (fw *flagWrapper) String(name, def, usage string) *string {
	return fw.fs.String(name, def, usage)
}

func (fw *flagWrapper) BoolP(name, shorthand string, def bool, usage string) *bool {
	return fw.fs.BoolP(name, shorthand, def, usage)
}

func (fw *flagWrapper) IntP(name, short string, def int, usage string) *int {
	return fw.fs.IntP(name, short, def, usage)
}

func (fw *flagWrapper) StringP(name, short, def, usage string) *string {
	return fw.fs.StringP(name, short, def, usage)
}
