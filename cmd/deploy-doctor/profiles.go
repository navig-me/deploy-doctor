package main

import (
	"fmt"

	"docker-doctor/internal/profiles"
	"github.com/spf13/cobra"
)

func newProfilesCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "profiles", Short: "Work with deployment profiles"}
	cmd.AddCommand(newProfilesListCmd())
	cmd.AddCommand(newProfilesExplainCmd())
	return cmd
}

func newProfilesListCmd() *cobra.Command {
	var recommended bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List available profiles",
		RunE: func(cmd *cobra.Command, args []string) error {
			if recommended {
				for _, p := range []string{"render", "flyio"} {
					cmd.Println(p)
				}
				return nil
			}
			for _, p := range profiles.List() {
				cmd.Println(p.Name)
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&recommended, "recommended", false, "List recommended profiles based on repo detection")
	return cmd
}

func newProfilesExplainCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "explain <profile>",
		Short: "Explain a profile",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			p, err := profiles.Get(args[0])
			if err != nil {
				return err
			}
			cmd.Println(p.Name)
			cmd.Println(p.Description)
			cmd.Printf("thresholds: image_warn=%d image_critical=%d layer_warn=%d mem_warn=%d\n", p.Thresholds.ImageSizeWarnMB, p.Thresholds.ImageSizeCriticalMB, p.Thresholds.LayerWarnCount, p.Thresholds.StartupMemoryWarnMB)
			return nil
		},
	}
}

func formatProfileLine(name string) string { return fmt.Sprintf("%s", name) }
