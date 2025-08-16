package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/nyasuto/pit/cmd"
)

type CLI struct {
	Init       cmd.InitCmd       `cmd:"" help:"Initialize a new pit repository"`
	HashObject cmd.HashObjectCmd `cmd:"" help:"Compute hash of a file"`
}

func main() {
	var cli CLI
	parser := kong.Must(&cli,
		kong.Name("pit"),
		kong.Description("A tiny, educational Git implementation in Go"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
		}),
	)

	ctx, err := parser.Parse(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", os.Args[0], err)
		os.Exit(1)
	}

	if err := ctx.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", os.Args[0], err)
		os.Exit(1)
	}
}
