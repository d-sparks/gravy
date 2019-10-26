package main

import (
	"github.com/d-sparks/gravy/kaggle"
	"github.com/spf13/cobra"
)

var kaggleCmd = &cobra.Command{
	Use:   "kaggle",
	Short: "trading pipeline for kaggle data",
	Run:   kaggleFn,
}

var kaggleInput string
var kaggleOutput string

func init() {
	pipelinesCmd.AddCommand(kaggleCmd)

	kaggleCmd.Flags().StringVarP(&kaggleInput, "input", "s", "./data/kaggle/historical_stock_prices.csv", "Kaggle data input")
	kaggleCmd.Flags().StringVarP(&kaggleOutput, "output", "o", "./data/kaggle/historical_as_windows.json", "Normalized output")
}

func kaggleFn(cmd *cobra.Command, args []string) {
	kaggle.Pipeline(kaggleInput, kaggleOutput)
}
