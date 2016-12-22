package commands

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use: "aws-utility",
	Short: "Lightweight AWS utility to analyse and produce services and data",
	Run: func(cmd *cobra.Command, args []string) {},
}
