package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type Manifest struct {
	Version   string `yaml:"version"`
	Resources []struct {
		Type       string `yaml:"type"`
		Name       string `yaml:"name"`
		Properties struct {
			Image string `yaml:"image"`
		} `yaml:"properties"`
	} `yaml:"resources"`
}

var applyCmd = &cobra.Command{
	Use:   "apply -f [file]",
	Short: "Apply configuration from a file (IaC)",
	Run: func(cmd *cobra.Command, args []string) {
		file, _ := cmd.Flags().GetString("file")
		if file == "" {
			fmt.Println("❌ Error: flag --file is required")
			return
		}

		data, err := os.ReadFile(file)
		if err != nil {
			fmt.Printf("❌ Read Error: %v\n", err)
			return
		}

		var manifest Manifest
		// LINT FIX: Handle Unmarshal error
		if err := yaml.Unmarshal(data, &manifest); err != nil {
			fmt.Printf("❌ YAML Parse Error: %v\n", err)
			return
		}

		client := NewClient()

		for _, res := range manifest.Resources {
			if res.Type == "vm" {
				fmt.Printf("⚙️  Syncing %s...\n", res.Name)

				reqBody := map[string]string{
					"name":  res.Name,
					"image": res.Properties.Image,
				}

				// LINT FIX: Handle Post error
				_, err := client.Post("/servers", reqBody)
				if err != nil {
					fmt.Printf("   ⚠️  Failed: %v\n", err)
				} else {
					fmt.Printf("   ✅ Active: %s\n", res.Name)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)
	applyCmd.Flags().StringP("file", "f", "", "Path to YAML configuration")
}
