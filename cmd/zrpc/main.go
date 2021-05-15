package main

import (
	"log"

	"github.com/mike955/zrpc/cmd/zrpc/new"

	"github.com/spf13/cobra"
)

var (
	version string = "v0.0.1-alpha1"

	rootCmd = &cobra.Command{
		Use:     "zrpc",
		Short:   "zrpc: An cli tookkit for zrpc framework (a mini go framework)",
		Long:    "zrpc: An cli tookkit for zrpc framework (a mini go framework)",
		Version: version,
	}
)

func init() {
	rootCmd.AddCommand(new.CmdNew)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
