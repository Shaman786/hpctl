// Package cmd handles the CLI commands and authentication logic for the hpctl application.
package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "hpctl",
	Short: "HPCTL is a CLI tool for managing resources",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// We don't define commands here anymore.
	// We just register them in their respective files (auth.go, vm.go).
}
