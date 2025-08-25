package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
	exampleservice "{{.ProjectPackage}}/pkg/{{.ProjectPackageName}}"
)

//nolint:gochecknoglobals // Cobra boilerplate
var trustedUsersFile string

// serverCmd represents the server command.
//
//nolint:gochecknoglobals // Cobra boilerplate
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the gRPC server",
	Long: `Start the {{.ProjectName}} gRPC server.

The server will:
- Listen for gRPC requests on the configured address (default: 0.0.0.0:50001)
- Serve Prometheus metrics on port 8080 at /metrics
- Provide health checks at /healthz and /readyz
- Authenticate requests using JWT-SSH tokens
- Support optional gRPC reflection for tooling

Configuration is handled via environment variables with the prefix 
{{.EnvPrefix}}_. See the --help output for available options.`,
	Run: runServer,
}

//nolint:gochecknoinits // Cobra boilerplate
func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().StringVar(&trustedUsersFile, "users", "users.json",
		"Path to trusted users JSON file (can also be set via {{.EnvPrefix}}_TRUSTED_USERS_FILE)")
}

func runServer(cmd *cobra.Command, args []string) {
	ctx := context.Background()

	// Load configuration
	config, err := exampleservice.LoadConfig()
	if err != nil {
		logger.Fatal("Failed to load configuration", zap.Error(err))
	}

	// Override users file if provided via flag
	if trustedUsersFile != "" {
		config.TrustedUsersFile = trustedUsersFile
	}

	logger.Info("Starting {{.ProjectName}} server",
		zap.String("grpc_address", config.GRPCAddress),
		zap.String("metrics_address", config.MetricsAddress),
		zap.Bool("reflection", config.EnableReflection),
		zap.Bool("plaintext", config.PlainText))

	// Load trusted users
	trustedUsers, err := loadTrustedUsersFromFile(config.TrustedUsersFile)
	if err != nil {
		logger.Fatal("Failed to load trusted users",
			zap.String("file", config.TrustedUsersFile),
			zap.Error(err))
	}

	logger.Info("Loaded trusted users",
		zap.Int("count", len(trustedUsers)),
		zap.String("file", config.TrustedUsersFile))

	// Start metrics server
	metricsServer := exampleservice.NewMetricsServer(config.MetricsAddress, logger)
	go func() {
		metricsErr := metricsServer.StartServer(ctx)
		if metricsErr != nil {
			logger.Fatal("Metrics server failed", zap.Error(metricsErr))
		}
	}()

	// Create and start gRPC server
	server, err := exampleservice.NewServer(config, logger, trustedUsers)
	if err != nil {
		logger.Fatal("Failed to create gRPC server", zap.Error(err))
	}

	// Set up graceful shutdown
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		logger.Info("Received shutdown signal")
		cancel()
	}()

	// Start server
	err = server.Start(ctx)
	if err != nil {
		logger.Fatal("Server failed", zap.Error(err))
	}

	logger.Info("Server shutdown complete")
}

func loadTrustedUsersFromFile(filename string) (users []exampleservice.TrustedUser, err error) {
	if filename == "" {
		users = nil
		err = errors.New("trusted users file not specified")
		return users, err
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		users = nil
		err = fmt.Errorf("failed to read trusted users file: %w", err)
		return users, err
	}

	// Try parsing as wrapper structure first (users.json format)
	var wrapper struct {
		Users []exampleservice.TrustedUser `json:"users"`
	}
	err = json.Unmarshal(data, &wrapper)
	if err == nil && len(wrapper.Users) > 0 {
		users = wrapper.Users
	} else {
		// Fall back to direct array parsing (test-users.json format)
		err = json.Unmarshal(data, &users)
		if err != nil {
			users = nil
			err = fmt.Errorf("failed to parse trusted users JSON: %w", err)
			return users, err
		}
	}

	// Validate users
	_, err = exampleservice.LoadTrustedUsers(users)
	if err != nil {
		users = nil
		err = fmt.Errorf("invalid trusted users: %w", err)
		return users, err
	}

	return users, err
}
