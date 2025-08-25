package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

// barCmd calls the bar method
//
//nolint:gochecknoglobals // Cobra boilerplate
var barCmd = &cobra.Command{
	Use:   "bar",
	Short: "Call the bar method",
	Long:  `Calls the bar method on the server.`,
	Run:   runBar,
}

//nolint:gochecknoinits // Cobra boilerplate
func init() {
	rootCmd.AddCommand(barCmd)
}

func runBar(cmd *cobra.Command, args []string) {
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
	response, err := client.CallBar(ctx, effectiveUsername)
	if err != nil {
		logger.Fatal("Failed to call bar", zapError(err))
	}

	fmt.Println(response)
}
