package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"text/template" // <--- Added this for templates

	"github.com/Shaman786/hpctl/pkg/api"
	"github.com/Shaman786/hpctl/pkg/models"
	"github.com/spf13/cobra"
)

var inspectCmd = &cobra.Command{
	Use:   "inspect [SERVICE_ID]",
	Short: "Show detailed information about a specific service",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serviceID := args[0]
		client := api.NewClient()

		// 1. Fetch Data
		endpoint := fmt.Sprintf("/service/%s", serviceID)
		data, err := client.Get(endpoint)
		if err != nil {
			fmt.Printf("❌ API Error: %v\n", err)
			return
		}

		// 2. Parse JSON
		var response models.ServiceDetailResponse
		err = json.Unmarshal(data, &response)
		if err != nil {
			// Fallback: Try unmarshaling directly if API returns raw object
			var directResponse models.ServiceDetail
			if err2 := json.Unmarshal(data, &directResponse); err2 == nil {
				response.Service = directResponse
			} else {
				fmt.Printf("❌ Failed to parse details: %v\n", err)
				return
			}
		}

		s := response.Service

		// 3. Define the Template (The "View")
		const layout = `
📋 SERVICE DETAILS
------------------------------------------------
Name:        {{ .Name }}
Domain:      {{ .Domain }}
Status:      {{ .Status }}
IP Address:  {{ .DedicatedIP }}
Price:       {{ .Amount }} {{ .BillingCycle }}
Next Due:    {{ .NextDueDate }}
Payment:     {{ .PaymentMethod }}
------------------------------------------------
`

		// 4. Render the Output
		t := template.Must(template.New("service").Parse(layout))
		err = t.Execute(os.Stdout, s)
		if err != nil {
			fmt.Println("Error formatting output:", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(inspectCmd)
}
