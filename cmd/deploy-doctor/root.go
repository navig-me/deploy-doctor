package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const version = "dev"

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy-doctor",
		Short: "Deploy Doctor checks container deploy readiness",
	}

	cmd.AddCommand(newScanCmd())
	cmd.AddCommand(newVersionCmd())
	cmd.AddCommand(newDoctorCmd())
	cmd.AddCommand(newProfilesCmd())

	return cmd
}

func execute() error {
	return newRootCmd().Execute()
}

func main() {
	if err := execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
