package cmd

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in and save credentials",
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)

		// 1. Get Username
		fmt.Print("Username (Email): ")
		username, _ := reader.ReadString('\n')
		username = strings.TrimSpace(username)

		// 2. Get Password (Hidden)
		fmt.Print("Password: ")
		bytePassword, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			fmt.Println("\nError reading password")
			return
		}
		password := string(bytePassword)
		fmt.Println() // Newline after input

		// 3. Save Logic
		if err := saveCredentials(username, password); err != nil {
			fmt.Printf("Error saving credentials: %v\n", err)
			return
		}
		fmt.Println("Success! Credentials verified and saved. 🚀")
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}

// ------------------------------------------------------
// AUTH HELPER FUNCTIONS (Shared by other commands)
// ------------------------------------------------------

// saveCredentials encodes user:pass and saves to ~/.hpctl.yaml
func saveCredentials(username, password string) error {
	authStr := fmt.Sprintf("%s:%s", username, password)
	encodedAuth := base64.StdEncoding.EncodeToString([]byte(authStr))
	content := fmt.Sprintf("auth: %s", encodedAuth)

	home, _ := os.UserHomeDir()
	configPath := filepath.Join(home, ".hpctl.yaml")
	return os.WriteFile(configPath, []byte(content), 0o644)
}

// LoadAuthHeader reads the config and returns the "Basic ..." string
// This is used by vm.go and other commands.
func LoadAuthHeader() (string, error) {
	home, _ := os.UserHomeDir()
	configPath := filepath.Join(home, ".hpctl.yaml")

	// Read file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return "", fmt.Errorf("not logged in (config file missing)")
	}

	// Parse "auth: XXXXX"
	content := string(data)
	if !strings.HasPrefix(content, "auth: ") {
		return "", fmt.Errorf("invalid config format")
	}

	encodedAuth := strings.TrimSpace(strings.TrimPrefix(content, "auth: "))
	return "Basic " + encodedAuth, nil
}
