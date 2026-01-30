package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

// Server represents a virtual machine instance running on the infrastructure.
type Server struct {
	Name   string `json:"name"`
	Image  string `json:"image"`
	Status string `json:"status"`
}

// APIResponse represents the standard JSON response format from the backend API.
type APIResponse struct {
	Message  string   `json:"message"`
	ServerID string   `json:"server_id"`
	Servers  []Server `json:"servers"`
	Logs     string   `json:"logs"`
	Error    string   `json:"error"`
}

var vmCmd = &cobra.Command{
	Use:   "vm",
	Short: "Manage virtual machines",
}

// LIST COMMAND
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all active servers",
	Run: func(cmd *cobra.Command, args []string) {
		client := NewClient()

		data, err := client.Get("/servers")
		if err != nil {
			fmt.Printf("❌ Connection Error: %v\n", err)
			return
		}

		var res APIResponse
		if err := json.Unmarshal(data, &res); err != nil {
			fmt.Println("❌ Error parsing API response")
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		// LINT FIX: Explicitly ignore write errors (standard for CLI output)
		_, _ = fmt.Fprintln(w, "NAME\tIMAGE\tSTATUS")
		_, _ = fmt.Fprintln(w, "----\t-----\t------")
		for _, s := range res.Servers {
			_, _ = fmt.Fprintf(w, "%s\t%s\t%s\n", s.Name, s.Image, s.Status)
		}
		_ = w.Flush()
	},
}

// CREATE COMMAND
var createCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Provision a new server",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := NewClient()
		name := args[0]
		image, _ := cmd.Flags().GetString("image")

		fmt.Printf("📡 Provisioning %s (%s)...\n", name, image)

		reqBody := map[string]string{"name": name, "image": image}

		data, err := client.Post("/servers", reqBody)
		if err != nil {
			fmt.Printf("❌ Provisioning Failed: %v\n", err)
			return
		}

		var res APIResponse
		// LINT FIX: Check error
		if err := json.Unmarshal(data, &res); err != nil {
			fmt.Printf("❌ Error parsing response: %v\n", err)
			return
		}
		fmt.Printf("✔ Success: %s (ID: %s)\n", res.Message, res.ServerID)
	},
}

// DESTROY COMMAND
var destroyCmd = &cobra.Command{
	Use:   "destroy [name]",
	Short: "Decommission a server",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := NewClient()
		name := args[0]

		_, err := client.Delete("/servers/" + name)
		if err != nil {
			fmt.Printf("❌ Failed: %v\n", err)
			return
		}
		fmt.Printf("✔ Server '%s' destroyed successfully.\n", name)
	},
}

// LOGS COMMAND
var logsCmd = &cobra.Command{
	Use:   "logs [name]",
	Short: "Fetch logs from a server",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := NewClient()
		name := args[0]

		data, err := client.Get(fmt.Sprintf("/servers/%s/logs", name))
		if err != nil {
			fmt.Printf("❌ Failed to fetch logs: %v\n", err)
			return
		}

		var res APIResponse
		// LINT FIX: Check error
		if err := json.Unmarshal(data, &res); err != nil {
			fmt.Printf("❌ Error parsing logs: %v\n", err)
			return
		}

		fmt.Println("\033[36m--- START LOGS ---\033[0m")
		fmt.Println(res.Logs)
		fmt.Println("\033[36m--- END LOGS ---\033[0m")
	},
}

func init() {
	rootCmd.AddCommand(vmCmd)
	vmCmd.AddCommand(listCmd)
	vmCmd.AddCommand(createCmd)
	vmCmd.AddCommand(logsCmd)
	vmCmd.AddCommand(destroyCmd)
	createCmd.Flags().StringP("image", "i", "alpine", "OS Image (alpine/nginx)")
}
