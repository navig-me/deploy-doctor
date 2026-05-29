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
				cmd.Println("docker: available")
			} else {
				cmd.Println("docker: unavailable")
				cmd.Println("hint: run scan --static-only for Dockerless checks")
			}
			return nil
		},
	}
}
