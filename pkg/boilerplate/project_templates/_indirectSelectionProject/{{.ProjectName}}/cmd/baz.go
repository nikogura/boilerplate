package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

// bazCmd calls the baz method
//
//nolint:gochecknoglobals // Cobra boilerplate
var bazCmd = &cobra.Command{
	Use:   "baz",
	Short: "Call the baz method",
	Long:  `Calls the baz method on the server.`,
	Run:   runBaz,
}

//nolint:gochecknoinits // Cobra boilerplate
func init() {
	rootCmd.AddCommand(bazCmd)
}

func runBaz(cmd *cobra.Command, args []string) {
	client, err := createClient()
	if err != nil {
		logger.Fatal("Failed to create client", zapError(err))
	}
	defer client.Close()

	// Get effective username (from flag or current user)
	effectiveUsername, err := getEffectiveUsername()
	if err != nil {
		logger.Fatal("Failed to get username", zapError(err))
	}

	ctx := context.Background()
	response, err := client.CallBaz(ctx, effectiveUsername)
	if err != nil {
		logger.Fatal("Failed to call baz", zapError(err))
	}

	fmt.Println(response)
}
