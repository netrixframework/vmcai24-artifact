package cmd

import "github.com/spf13/cobra"

var iterations int

func RootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "tendermint-testing",
	}
	cmd.CompletionOptions.DisableDefaultCmd = true
	cmd.PersistentFlags().IntVarP(&iterations, "iterations", "i", 100, "number of strategy iterations to run")
	cmd.AddCommand(PCTStrategy())
	cmd.AddCommand(PCTTestStrategy())
	return cmd
}
