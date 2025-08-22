/*
Copyright © {{.ProjectVersion}} {{.OwnerName}} <{{.OwnerEmail}}>
*/
package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"{{.ProjectPackage}}/pkg/{{.ProjectPackageName}}"
)

// serverCmd represents the server command
//
//nolint:gochecknoglobals // Cobra boilerplate
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "{{.ServerShortDesc}}",
	Long: `
{{.ServerLongDesc}}
`,
	RunE: runServer,
}

func runServer(cmd *cobra.Command, args []string) (err error) {
	// Create logger with simple production config
	var logger *zap.Logger
	logger, err = {{.ProjectPackageName}}.NewLogger("info", "json")
	if err != nil {
		return err
	}
	defer func() {
		_ = logger.Sync()
	}()

	// Create service with 5 second interval (default)
	service := {{.ProjectPackageName}}.NewService(logger, 0)

	// Set up context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle OS signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigCh
		logger.Info("Received shutdown signal")
		cancel()
	}()

	// Start service
	err = service.Start(ctx)
	return err
}

//nolint:gochecknoinits // Cobra boilerplate
func init() {
	rootCmd.AddCommand(serverCmd)
}