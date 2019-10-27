package main

import (
	"github.com/spf13/cobra"
)

var pipelinesCmd = &cobra.Command{
	Use:   "pipelines",
	Short: "trading data pipelines",
}

func init() {
	rootCmd.AddCommand(pipelinesCmd)
}
