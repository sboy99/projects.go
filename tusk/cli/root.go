package cli

import (
	"github.com/sboy99/projects.go/tusk/internal/config"
	"github.com/sboy99/projects.go/tusk/internal/version"
	"github.com/sboy99/projects.go/tusk/pkg/logger"
	"github.com/spf13/cobra"
)

var log *logger.Logger

func init() {
	log = logger.NewLogger("")
}

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tusk",
	Short: "Tusk CLI - A command-line tool",
	Long: `Tusk is a CLI application that provides various commands
for managing services and tasks.`,
	Version: version.GetVersion(),
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.tusk/config.yaml)")
	rootCmd.PersistentFlags().Bool("debug", false, "enable debug mode")

	// Add version command
	rootCmd.AddCommand(versionCmd)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		config.Load(cfgFile)
	} else {
		config.Load("")
	}
}

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Long:  `Print the version number and build information.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println(version.GetBuildInfo())
	},
}
