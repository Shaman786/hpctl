package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/Shaman786/hpctl/pkg/api"
	"github.com/Shaman786/hpctl/pkg/models"
	"github.com/Shaman786/hpctl/pkg/templates"
	"github.com/spf13/cobra"
)

var inspectCmd = &cobra.Command{
	Use:   "inspect [SERVICE_ID]",
	Short: "Show detailed information about a specific service",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serviceID := args[0]
		client := api.NewClient()

		// 1. Fetch
		endpoint := fmt.Sprintf("/service/%s", serviceID)
		data, err := client.Get(endpoint)
		if err != nil {
			fmt.Printf("❌ API Error: %v\n", err)
			return
		}

		// 2. Parse (Handles potential API inconsistencies)
		var response models.ServiceDetailResponse
		if err := json.Unmarshal(data, &response); err != nil {
			// Fallback if API returns direct object instead of wrapper
			var directResponse models.ServiceDetail
			if err2 := json.Unmarshal(data, &directResponse); err2 == nil {
				response.Service = directResponse
			} else {
				fmt.Printf("❌ Failed to parse details: %v\n", err)
				return
			}
		}

		// 3. Render
		templates.Render(templates.ServiceDetail, response.Service)
	},
}

func init() {
	rootCmd.AddCommand(inspectCmd)
}
