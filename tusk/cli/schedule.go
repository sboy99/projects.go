package cli

import (
	"github.com/sboy99/projects.go/tusk/internal/executor"
	"github.com/spf13/cobra"
)

var name string
var command string
var interval string

// scheduleCmd represents the schedule command
var scheduleCmd = &cobra.Command{
	Use:   "schedule",
	Short: "Manage scheduled tasks",
	Long:  `Manage scheduled tasks with create, list, and other operations.`,
}

// scheduleCreateCmd represents the schedule create command
var scheduleCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new scheduled task",
	Long:  `Create a new scheduled task with the specified name, command, and interval.`,
	Run: func(cmd *cobra.Command, args []string) {
		scheduleExecutor := executor.NewScheduleExecutor(name, command, interval)
		scheduleExecutor.Create()
	},
}

// scheduleListCmd represents the schedule list command
var scheduleListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all scheduled tasks",
	Long:  `List all scheduled tasks with their details.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: list all schedules from storage

	},
}

func init() {
	rootCmd.AddCommand(scheduleCmd)

	// Add create subcommand
	scheduleCmd.AddCommand(scheduleCreateCmd)
	scheduleCreateCmd.Flags().StringVarP(&name, "name", "n", "", "name for the schedule")
	scheduleCreateCmd.Flags().StringVarP(&command, "command", "c", "", "command to schedule (e.g. echo 'Hello, World!')")
	scheduleCreateCmd.Flags().StringVarP(&interval, "interval", "i", "", "interval for scheduling (e.g. 1m, 1h, 1d, 1w)")
	scheduleCreateCmd.MarkFlagRequired("command")
	scheduleCreateCmd.MarkFlagRequired("interval")

	// Add list subcommand
	scheduleCmd.AddCommand(scheduleListCmd)
}
