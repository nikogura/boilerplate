/*
Copyright © {{.ProjectVersion}} {{.OwnerName}} <{{.OwnerEmail}}>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
//
//nolint:gochecknoglobals // Cobra boilerplate
var rootCmd = &cobra.Command{
	Use:   "{{.ProjectName}}",
	Short: "{{.ProjectShortDesc}}",
	Long: `
{{.ProjectLongDesc}}
`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

//nolint:gochecknoinits // Cobra boilerplate
func init() {

}
