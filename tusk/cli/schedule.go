package cli

import (
	"github.com/sboy99/projects.go/tusk/internal/executor"
	"github.com/spf13/cobra"
)

var name string
var command string
var interval string
var follow bool

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
		scheduleExecutor := executor.NewScheduleExecutor()
		scheduleExecutor.Create(name, command, interval)
	},
}

// scheduleListCmd represents the schedule list command
var scheduleListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all scheduled tasks",
	Long:  `List all scheduled tasks with their details.`,
	Run: func(cmd *cobra.Command, args []string) {
		scheduleExecutor := executor.NewScheduleExecutor()
		scheduleExecutor.List()
	},
}

// scheduleDeleteCmd represents the schedule delete command
var scheduleDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a scheduled task",
	Long:  `Delete a scheduled task by name. This will stop and disable the timer, delete all associated files, and remove the entry from storage.`,
	Run: func(cmd *cobra.Command, args []string) {
		if name == "" {
			cmd.Help()
			return
		}
		scheduleExecutor := executor.NewScheduleExecutor()
		scheduleExecutor.Delete(name)
	},
}

// scheduleLogsCmd represents the schedule logs command
var scheduleLogsCmd = &cobra.Command{
	Use:   "logs",
	Short: "View logs for a scheduled task",
	Long:  `View journal logs for a specific scheduled task by name. Use the -f flag to follow logs in real-time.`,
	Run: func(cmd *cobra.Command, args []string) {
		if name == "" {
			cmd.Help()
			return
		}
		scheduleExecutor := executor.NewScheduleExecutor()
		scheduleExecutor.Logs(name, follow)
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

	// Add delete subcommand
	scheduleCmd.AddCommand(scheduleDeleteCmd)
	scheduleDeleteCmd.Flags().StringVarP(&name, "name", "n", "", "name of the schedule to delete")
	scheduleDeleteCmd.MarkFlagRequired("name")

	// Add logs subcommand
	scheduleCmd.AddCommand(scheduleLogsCmd)
	scheduleLogsCmd.Flags().StringVarP(&name, "name", "n", "", "name of the schedule to view logs for")
	scheduleLogsCmd.Flags().BoolVarP(&follow, "follow", "f", false, "follow log output (similar to tail -f)")
	scheduleLogsCmd.MarkFlagRequired("name")
}
