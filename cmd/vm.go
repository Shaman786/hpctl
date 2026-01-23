package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"text/tabwriter" // Makes the output look like a pro table

	"github.com/spf13/cobra"
)

// CONFIG: Point to your Mock API
const API_URL = "http://localhost:5000/api/v1"

// --- STRUCTS (Data Models) ---
type Server struct {
	Name   string `json:"name"`
	Image  string `json:"image"`
	Status string `json:"status"`
}

type ApiResponse struct {
	Message  string   `json:"message"`
	ServerID string   `json:"server_id"`
	Servers  []Server `json:"servers"`
	Logs     string   `json:"logs"`
	Error    string   `json:"error"`
}

// --- COMMAND: PARENT 'vm' ---
var vmCmd = &cobra.Command{
	Use:   "vm",
	Short: "Manage virtual machines",
	Long:  `Create, list, inspect, and destroy virtual machine instances.`,
}

// --- COMMAND: LIST ---
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all active servers",
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := http.Get(API_URL + "/servers")
		if err != nil {
			fmt.Printf("❌ Connection Failed: Is hp-cloud-api running on port 5000?\n")
			return
		}
		defer resp.Body.Close()

		var res ApiResponse
		if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
			fmt.Println("❌ Error parsing response")
			return
		}

		// Print nicer table
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "NAME\tIMAGE\tSTATUS")
		fmt.Fprintln(w, "----\t-----\t------")
		for _, s := range res.Servers {
			fmt.Fprintf(w, "%s\t%s\t%s\n", s.Name, s.Image, s.Status)
		}
		w.Flush()
	},
}

// --- COMMAND: CREATE ---
var createCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Provision a new server",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		image, _ := cmd.Flags().GetString("image") // This is safe to ignore

		fmt.Printf("📡 Provisioning %s (%s)...\n", name, image)

		// FIX: Use a struct + json.Marshal for safety
		reqBody := map[string]string{
			"name":  name,
			"image": image,
		}
		jsonPayload, _ := json.Marshal(reqBody)

		resp, err := http.Post(API_URL+"/servers", "application/json", strings.NewReader(string(jsonPayload)))
		if err != nil {
			fmt.Printf("❌ API Error: %v\n", err)
			return
		}
		defer resp.Body.Close()

		var res ApiResponse
		// FIX: Handle JSON decode errors (e.g., if API crashes and returns HTML)
		if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
			fmt.Println("❌ Error: Received invalid response from Cloud API")
			return
		}

		if resp.StatusCode == 201 {
			fmt.Printf("✔ Success: %s (ID: %s)\n", res.Message, res.ServerID)
		} else {
			fmt.Printf("❌ Error: %s\n", res.Error)
		}
	},
}

// --- COMMAND: LOGS ---
var logsCmd = &cobra.Command{
	Use:   "logs [name]",
	Short: "Fetch logs from a server",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		resp, err := http.Get(fmt.Sprintf("%s/servers/%s/logs", API_URL, name))
		if err != nil {
			fmt.Printf("❌ API Error: %v\n", err)
			return
		}
		defer resp.Body.Close()

		var res ApiResponse
		if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
			fmt.Println("❌Error parsing logs")
			return
		}

		if resp.StatusCode != 200 {
			fmt.Printf("❌ Server Error: %s\n", res.Error)
			return
		}

		fmt.Println("\033[36m--- START LOGS ---\033[0m") // Cyan Color Header
		fmt.Println(res.Logs)
		fmt.Println("\033[36m--- END LOGS ---\033[0m")
	},
}

// --- COMMAND: DESTROY ---
var destroyCmd = &cobra.Command{
	Use:   "destroy [name]",
	Short: "Decommission a server",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		// FIX 1: Handle the error!
		req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/servers/%s", API_URL, name), nil)
		if err != nil {
			fmt.Printf("❌ Internal Error creating request: %v\n", err)
			return
		}

		client := &http.Client{}
		resp, err := client.Do(req)

		if err != nil {
			fmt.Printf("❌ Connection Error: %v\n", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == 200 {
			fmt.Printf("✔ Server '%s' destroyed successfully.\n", name)
		} else {
			// FIX 2: Decode the error message from the API instead of guessing
			var res ApiResponse
			json.NewDecoder(resp.Body).Decode(&res)
			fmt.Printf("❌ Failed: %s\n", res.Error)
		}
	},
}

// --- INIT ---
func init() {
	// Register subcommands
	rootCmd.AddCommand(vmCmd)
	vmCmd.AddCommand(listCmd)
	vmCmd.AddCommand(createCmd)
	vmCmd.AddCommand(logsCmd)
	vmCmd.AddCommand(destroyCmd)

	// Flags
	createCmd.Flags().StringP("image", "i", "alpine", "OS Image (alpine/nginx)")
}
