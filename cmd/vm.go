package cmd

import (
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

// vmCmd is the parent command (e.g., "hpctl vm")
var vmCmd = &cobra.Command{
	Use:   "vm",
	Short: "Manage Virtual Machines",
}

// listVmCmd is the subcommand (e.g., "hpctl vm list")
var listVMCmd = &cobra.Command{
	Use:   "list",
	Short: "List all VMs",
	Run: func(cmd *cobra.Command, args []string) {
		// 1. Get Credentials
		authHeader, err := LoadAuthHeader()
		if err != nil {
			fmt.Println("Error: You must log in first.")
			fmt.Println("Run: hpctl login")
			return
		}

		fmt.Println("Fetching VMs...")

		// 2. Make Request (Example)
		client := &http.Client{Timeout: 10 * time.Second}
		req, _ := http.NewRequest("GET", "https://api.example.com/vms", nil)

		// 3. Inject Auth Header
		req.Header.Add("Authorization", authHeader)

		// 4. Send (This will fail unless you have a real API, so we mock the success below)
		resp, err := client.Do(req)

		// --- MOCK OUTPUT FOR DEMO ---
		if err != nil || resp.StatusCode != 200 {
			// Simulate success for your CLI testing

			fmt.Println("------------------------------------------------")
			fmt.Println("NAME          STATUS    IP")
			fmt.Println("web-server-1  RUNNING   192.168.1.10")
			fmt.Println("db-server     STOPPED   192.168.1.12")
			return
		}
		// ----------------------------
	},
}

func init() {
	// Register 'vm' to root
	rootCmd.AddCommand(vmCmd)
	// Register 'list' to 'vm'
	vmCmd.AddCommand(listVMCmd)
}
