package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

// fooCmd calls the foo method
//
//nolint:gochecknoglobals // Cobra boilerplate
var fooCmd = &cobra.Command{
	Use:   "foo",
	Short: "Call the foo method",
	Long:  `Calls the foo method on the server.`,
	Run:   runFoo,
}

//nolint:gochecknoinits // Cobra boilerplate
func init() {
	rootCmd.AddCommand(fooCmd)
}

func runFoo(cmd *cobra.Command, args []string) {
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
	response, err := client.CallFoo(ctx, effectiveUsername)
	if err != nil {
		logger.Fatal("Failed to call foo", zapError(err))
	}

	fmt.Println(response)
}
