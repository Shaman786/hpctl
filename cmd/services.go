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

var servicesCmd = &cobra.Command{
    Use:   "services",
    Short: "List all active services",
    Run: func(cmd *cobra.Command, args []string) {
        client := api.NewClient()

        fmt.Println("📦 Fetching your services...")

        // 1. Get Data
        data, err := client.Get("/service")
        if err != nil {
            fmt.Printf("❌ API Error: %v\n", err)
            return
        }

        // 2. Parse JSON
        var response models.ServiceResponse
        err = json.Unmarshal(data, &response)
        if err != nil {
            fmt.Printf("❌ Failed to parse services: %v\n", err)
            return
        }

        // 3. Print Table
        if len(response.Services) == 0 {
            fmt.Println("No services found.")
            return
        }

        fmt.Println("\nYOUR INFRASTRUCTURE")
        fmt.Println("--------------------------------------------------------------------------------")

        // TabWriter formats text into aligned columns
        w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)

        // Header
        fmt.Fprintln(w, "ID\tNAME\tDOMAIN\tPRICE\tSTATUS\tNEXT DUE")

        // Rows
        for _, s := range response.Services {
            fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n", 
                s.ID, 
                s.Name, 
                s.Domain, 
                s.Price, 
                s.Status, 
                s.NextDue,
            )
        }
        w.Flush()
        fmt.Println("--------------------------------------------------------------------------------")
    },
}

func init() {
    rootCmd.AddCommand(servicesCmd)
}