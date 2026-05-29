package main

import (
	"os/exec"

	"github.com/spf13/cobra"
)

func dockerAvailable() bool {
	cmd := exec.Command("docker", "version", "--format", "{{.Server.Version}}")
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

func newDoctorCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "doctor",
		Short: "Diagnose local tool/runtime prerequisites",
		RunE: func(cmd *cobra.Command, args []string) error {
			if dockerAvailable() {
				cmd.Println("Docker daemon: available")
				cmd.Println("Runtime probes are enabled.")
			} else {
				cmd.Println("Docker daemon: unavailable")
				cmd.Println("Fallback: run `deploy-doctor scan --static-only` for static checks only.")
			}
			return nil
		},
	}
}
