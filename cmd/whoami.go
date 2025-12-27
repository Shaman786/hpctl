package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/Shaman786/hpctl/pkg/api"
	"github.com/Shaman786/hpctl/pkg/models"
	"github.com/Shaman786/hpctl/pkg/templates" // <--- Import the UI package
	"github.com/spf13/cobra"
)

var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Get current user details",
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewClient()

		// 1. Fetch
		data, err := client.Get("/details")
		if err != nil {
			fmt.Printf("❌ API Error: %v\n", err)
			return
		}

		// 2. Parse
		var response models.ClientResponse
		if err := json.Unmarshal(data, &response); err != nil {
			fmt.Printf("❌ JSON Error: %v\n", err)
			return
		}

		// 3. Render (Delegated to templates package)
		templates.Render(templates.WhoAmI, response.Client)
	},
}

func init() {
	rootCmd.AddCommand(whoamiCmd)
}
