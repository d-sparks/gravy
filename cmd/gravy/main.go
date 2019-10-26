package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gravy",
	Short: "Gravy command line tool",
	Run:   rootFn,
}

func rootFn(cmd *cobra.Command, args []string) {
	fmt.Println("hi")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}
