package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3" // Run: go get gopkg.in/yaml.v3
)

type Manifest struct {
	Resources []struct {
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
		data, _ := os.ReadFile(file)

		var manifest Manifest
		yaml.Unmarshal(data, &manifest)

		client := NewClient()
		for _, res := range manifest.Resources {
			fmt.Printf("⚙️  Syncing %s...\n", res.Name)
			client.Post("/servers", map[string]string{
				"name":  res.Name,
				"image": res.Properties.Image,
			})
		}
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)
	applyCmd.Flags().StringP("file", "f", "", "Config file")
}
