package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

func rootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "store",
		Short: "store is a tool to manages data with JSON",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.AddCommand(newAddCmd())
	cmd.AddCommand(newSetCmd())
	cmd.AddCommand(newRmCmd())
	cmd.AddCommand(newGetCmd())
	cmd.AddCommand(newListCmd())

	return cmd
}

func Execute() {
	cmd := rootCmd()
	cmd.SetOutput(os.Stdout)
	if err := cmd.Execute(); err != nil {
		cmd.SetOutput(os.Stderr)
		cmd.Println(err)
		os.Exit(1)
	}
}
