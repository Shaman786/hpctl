package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/Shaman786/hpctl/pkg/api"
	"github.com/Shaman786/hpctl/pkg/models"
	"github.com/Shaman786/hpctl/pkg/templates"
	"github.com/spf13/cobra"
)

var servicesCmd = &cobra.Command{
	Use:   "services",
	Short: "List all active services",
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewClient()
		fmt.Println("📦 Fetching your services...")

		// 1. Fetch
		data, err := client.Get("/service")
		if err != nil {
			fmt.Printf("❌ API Error: %v\n", err)
			return
		}

		// 2. Parse
		var response models.ServiceResponse
		if err := json.Unmarshal(data, &response); err != nil {
			fmt.Printf("❌ JSON Error: %v\n", err)
			return
		}

		if len(response.Services) == 0 {
			fmt.Println("No services found.")
			return
		}

		// 3. Render
		templates.Render(templates.ServicesList, response)
	},
}

func init() {
	rootCmd.AddCommand(servicesCmd)
}
