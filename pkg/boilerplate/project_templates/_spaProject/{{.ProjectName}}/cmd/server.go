/*
Copyright © 2024 Nik Ogura <nik@example.com>
*/
package cmd

import (
	"context"
	"log"

	"github.com/spf13/cobra"
	"{{.ProjectPackage}}/pkg/{{.ProjectPackageName}}"
)

//nolint:gochecknoglobals // cobra requires global command vars
var address string

// serverCmd represents the server command.
//
//nolint:gochecknoglobals // cobra requires global command vars
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run the {{.ProjectName}} server",
	Long: `Run the {{.ProjectName}} server.

This starts both the main HTTP server (port 9999 by default) for serving 
the Single-Page Application and the metrics server (port 8080) for 
Prometheus metrics and health checks.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		config, err := {{.ProjectPackageName}}.LoadConfig()
		if err != nil {
			log.Fatalf("Failed to load configuration: %s", err)
		}

		if address != "" {
			// Override config with command line flag
			config.ServerAddress = address
		}

		server, err := {{.ProjectPackageName}}.NewServer(ctx, config)
		if err != nil {
			log.Fatalf("Failed creating server: %s", err)
		}

		err = server.Run(ctx)
		if err != nil {
			log.Fatalf("Failed running server: %s", err)
		}
	},
}

//nolint:gochecknoinits // cobra requires init for command registration
func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().StringVarP(&address, "address", "a", "", "address on which to run (overrides config)")
}
