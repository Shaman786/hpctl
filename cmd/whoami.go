package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/Shaman786/hpctl/pkg/api"
	"github.com/Shaman786/hpctl/pkg/models"
	"github.com/spf13/cobra"
)

var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Get current user details",
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewClient()

		// 1. Fetch Data
		data, err := client.Get("/details")
		if err != nil {
			fmt.Printf("❌ API Error: %v\n", err)
			return
		}

		// 2. Parse JSON into Struct
		var response models.ClientResponse
		err = json.Unmarshal(data, &response)
		if err != nil {
			fmt.Printf("❌ Failed to parse response: %v\n", err)
			return
		}

		// 3. Print Pretty Output
		c := response.Client
		fmt.Println("\n👤  USER PROFILE")
		fmt.Println("------------------------------------------------")

		// Use Tabwriter for perfect alignment
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintf(w, "ID:\t%s\n", c.ID)
		fmt.Fprintf(w, "Name:\t%s %s\n", c.FirstName, c.LastName)
		fmt.Fprintf(w, "Email:\t%s\n", c.Email)
		fmt.Fprintf(w, "Phone:\t%s\n", c.Phone)
		fmt.Fprintf(w, "Location:\t%s, %s\n", c.City, c.Country)
		w.Flush()
		fmt.Println("------------------------------------------------")
	},
}

func init() {
	rootCmd.AddCommand(whoamiCmd)
}
