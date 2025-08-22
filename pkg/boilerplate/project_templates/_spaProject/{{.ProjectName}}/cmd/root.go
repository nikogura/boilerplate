/*
Copyright © 2024 Nik Ogura <nik@example.com>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//nolint:gochecknoglobals // cobra requires global command vars
var cfgFile string

//nolint:gochecknoglobals // cobra requires global root command
var rootCmd = &cobra.Command{
	Use:   "{{.ProjectName}}",
	Short: "{{.ProjectShortDesc}}",
	Long:  `{{.ProjectLongDesc}}`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

//nolint:gochecknoinits // cobra requires init for command registration
func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.{{.ProjectName}}.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".{{.ProjectName}}" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".{{.ProjectName}}")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// Set key replacer to convert dots and hyphens to underscores for env vars
	viper.SetEnvKeyReplacer(nil) // We'll handle this in config.go with proper prefix

	// If a config file is found, read it in.
	err := viper.ReadInConfig()
	if err == nil {
		// Config file found and read successfully
		_ = err // We don't need to handle this specific case
	}
}
