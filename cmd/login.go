package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/term"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with Host-Palace",
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)

		// 1. Get Username
		fmt.Print("Enter Username: ")
		username, _ := reader.ReadString('\n')
		username = strings.TrimSpace(username)

		// 2. Get Password (Hidden)
		fmt.Print("Enter Password: ")
		bytePassword, _ := term.ReadPassword(int(syscall.Stdin))
		password := string(bytePassword)
		fmt.Println("\nLogging in...")

		// 3. Save to Config
		viper.Set("auth.username", username)
		viper.Set("auth.password", password) // In production, we would use a keyring

		err := viper.WriteConfig()
		if err != nil {
			// If file doesn't exist, create it safely
			viper.SafeWriteConfig()
		}

		fmt.Println("✅ Credentials saved to ~/.hpctl.yaml")
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
