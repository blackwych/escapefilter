package main

import (
	"fmt"
	"github.com/blackwych/escapefilter/escapefilter"
	"github.com/jessevdk/go-flags"
	"os"
)

var (
	// Application name
	appname = "escapefilter"

	// Version number starting with "v"
	// The value will be injected by build script.
	version = "(unknown version)"
)

// arguments represents positional arguments
type arguments struct {
	Infiles []string `positional-arg-name:"INFILE"`
}

// options represents command-line options (and positional arguments)
type options struct {
	Help    bool      `short:"h" long:"help" description:"Print this help and exit"`
	Version bool      `short:"v" long:"version" description:"Print version information and exit"`
	Args    arguments `positional-args:"true"`
}

// versionInfo returns version information.
func versionInfo() string {
	return fmt.Sprintf("%s %s", appname, version)
}

// parseOptions parses command-line options.
// (nil, nil) means "Do nothing and exit successfully."
func parseOptions() (*options, error) {
	parser := flags.NewNamedParser(appname, flags.PassDoubleDash)

	opts := &options{}
	parser.AddGroup("Options", "Options", opts)

	if _, err := parser.Parse(); err != nil {
		return nil, err
	}

	if opts.Help {
		parser.WriteHelp(os.Stdout)
		return nil, nil
	}

	if opts.Version {
		fmt.Println(versionInfo())
		return nil, nil
	}

	// Use standard input if no files specified
	if len(opts.Args.Infiles) == 0 {
		opts.Args.Infiles = []string{"-"}
	}

	return opts, nil
}

// exitWithError reports error and exits with status 1 (= error).
func exitWithError(err error) {
	fmt.Fprintf(os.Stderr, "%s: %v\n", appname, err)
	os.Exit(1)
}

// load loads an input file.
func load(filter *escapefilter.EscapeFilter, filename string) error {
	var file *os.File

	if filename == "-" {
		file = os.Stdin
	} else {
		f, err := os.Open(filename)
		if err != nil {
			return err
		}

		file = f
		defer file.Close()
	}

	if err := filter.Load(file); err != nil {
		return err
	}

	return nil
}

func main() {
	opts, err := parseOptions()
	if err != nil {
		exitWithError(err)
		return
	}

	if opts == nil {
		os.Exit(0)
		return
	}

	filter := escapefilter.New()

	for _, infile := range opts.Args.Infiles {
		if err := load(filter, infile); err != nil {
			exitWithError(err)
			return
		}
	}

	fmt.Print(filter)
}
