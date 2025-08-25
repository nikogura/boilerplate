package cmd

import (
	"fmt"
	"os"
	"os/user"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
	exampleservice "{{.ProjectPackage}}/pkg/{{.ProjectPackageName}}"
)

//nolint:gochecknoglobals // Cobra boilerplate
var debug bool

//nolint:gochecknoglobals // Cobra boilerplate
var verbose bool

//nolint:gochecknoglobals // Cobra boilerplate
var serverAddress string

//nolint:gochecknoglobals // Cobra boilerplate
var username string

//nolint:gochecknoglobals // Cobra boilerplate
var pubKeyFile string

//nolint:gochecknoglobals // Cobra boilerplate
var plaintext bool

//nolint:gochecknoglobals // Cobra boilerplate
var showToken bool

//nolint:gochecknoglobals // Cobra boilerplate
var logger *zap.Logger

// rootCmd represents the base command when called without any subcommands.
//
//nolint:gochecknoglobals // Cobra boilerplate
var rootCmd = &cobra.Command{
	Use:   "{{.ProjectName}}",
	Short: "A gRPC service demonstrating indirect selection patterns",
	Long: `{{.ProjectName}} is a gRPC-based service that demonstrates
indirect selection patterns for use in code generation tools.

The service provides:
- gRPC server with JWT-SSH authentication
- Method registry for dynamic method discovery
- Prometheus metrics and health checks
- Configurable reflection support

Example usage:
  {{.ProjectName}} server                    # Start the server
  {{.ProjectName}} client list-methods      # List available methods
  {{.ProjectName}} client call echo --message="hello"  # Call a method`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

//nolint:gochecknoinits // Cobra boilerplate
func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Enable debug logging")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose Output")
	rootCmd.PersistentFlags().StringVarP(&serverAddress, "address", "a", "0.0.0.0:50001", "Server address")
	rootCmd.PersistentFlags().BoolVarP(&plaintext, "plaintext", "", false, "connect without tls")
	rootCmd.PersistentFlags().StringVarP(&pubKeyFile, "pubkey-file", "f", "~/.ssh/id_ed25519.pub", "File containing SSH public key to use for authentication.")
	rootCmd.PersistentFlags().StringVarP(&username, "username", "u", "", "Username for authentication")
	rootCmd.PersistentFlags().BoolVarP(&showToken, "show-token", "", false, "Dump JWT to stdout")
}

func initConfig() {
	// Initialize logger
	var err error
	if debug {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}

	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
}

// clientCmd represents the client command.
//
//nolint:gochecknoglobals // Cobra boilerplate
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Client commands for testing indirect selection",
	Long:  `Simple client commands to test the indirect selection service.`,
}

//nolint:gochecknoinits // Cobra boilerplate
func init() {
	rootCmd.AddCommand(clientCmd)
}

// Common helper functions for client commands.
func createClient() (client *exampleservice.Client, err error) {
	// Load base configuration
	config, err := exampleservice.LoadConfig()
	if err != nil {
		client = nil
		err = fmt.Errorf("failed to load configuration: %w", err)
		return client, err
	}

	// Override with command line flags if provided
	if serverAddress != "" {
		config.GRPCAddress = serverAddress
	}
	if plaintext {
		config.PlainText = plaintext
	}

	// Create client
	client, err = exampleservice.NewClient(config, logger, username, pubKeyFile)
	if err != nil {
		client = nil
		err = fmt.Errorf("failed to create client: %w", err)
		return client, err
	}

	return client, err
}

// zapError is a helper function to create zap error fields.
func zapError(err error) zap.Field {
	return zap.Error(err)
}

// getEffectiveUsername returns the username from flags or current user as fallback.
func getEffectiveUsername() (effectiveUsername string, err error) {
	if username != "" {
		effectiveUsername = username
		return effectiveUsername, err
	}

	// Get current user as fallback
	currentUser, err := user.Current()
	if err != nil {
		effectiveUsername = ""
		err = fmt.Errorf("failed to get current user: %w", err)
		return effectiveUsername, err
	}

	effectiveUsername = currentUser.Username
	return effectiveUsername, err
}
