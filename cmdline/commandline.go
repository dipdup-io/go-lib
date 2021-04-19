package cmdline

import "flag"

// Args - struct contains all command-line arguments
type Args struct {
	Config string
	Help   bool
}

// Parse -
func Parse() (args Args) {
	flag.BoolVar(&args.Help, "h", false, "Show usage")
	flag.StringVar(&args.Config, "f", "config.yaml", "Path to YAML config file")

	flag.Parse()

	if flag.NFlag() == 0 {
		args.Help = true
	}

	if args.Help {
		flag.Usage()
	}
	return
}
