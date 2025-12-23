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
	Use:   "create [command] [interval]",
	Short: "Create a new scheduled task",
	Long:  `Create a new scheduled task with the specified name, command, and interval. Command and interval can be provided as flags or positional arguments.`,
	Example: `
	tusk schedule create -n my-task -c "echo 'Hello, World!'" -i "1m"
	tusk schedule create -n my-task -c "echo 'Hello, World!'" -i "1h"
	tusk schedule create -n my-task "echo 'Hello, World!'" "1m"
	tusk schedule create "echo 'Hello, World!'" "1h"
	`,
	Run: func(cmd *cobra.Command, args []string) {
		// Use positional arguments if flags are not provided
		if command == "" && len(args) > 0 {
			command = args[0]
		}
		if interval == "" && len(args) > 1 {
			interval = args[1]
		}

		// Validate required fields
		if command == "" {
			cmd.Help()
			return
		}
		if interval == "" {
			cmd.Help()
			return
		}

		scheduleExecutor := executor.NewScheduleExecutor()
		scheduleExecutor.Create(name, command, interval)
	},
}

// scheduleListCmd represents the schedule list command
var scheduleListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all scheduled tasks",
	Long:  `List all scheduled tasks with their details.`,
	Example: `
	tusk schedule list
	`,
	Run: func(cmd *cobra.Command, args []string) {
		scheduleExecutor := executor.NewScheduleExecutor()
		scheduleExecutor.List()
	},
}

// scheduleDeleteCmd represents the schedule delete command
var scheduleDeleteCmd = &cobra.Command{
	Use:   "delete [name]",
	Short: "Delete a scheduled task",
	Long:  `Delete a scheduled task by name. This will stop and disable the timer, delete all associated files, and remove the entry from storage. Name can be provided as a flag or positional argument.`,
	Example: `
	tusk schedule delete -n my-task
	tusk schedule delete my-task
	`,
	Run: func(cmd *cobra.Command, args []string) {
		// Use positional argument if flag is not provided
		if name == "" && len(args) > 0 {
			name = args[0]
		}

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
	Use:   "logs [name]",
	Short: "View logs for a scheduled task",
	Long:  `View journal logs for a specific scheduled task by name. Use the -f flag to follow logs in real-time. Name can be provided as a flag or positional argument.`,
	Example: `
	tusk schedule logs -n my-task
	tusk schedule logs my-task
	tusk schedule logs -n my-task -f
	tusk schedule logs my-task -f
	`,
	Run: func(cmd *cobra.Command, args []string) {
		// Use positional argument if flag is not provided
		if name == "" && len(args) > 0 {
			name = args[0]
		}

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

	// Add list subcommand
	scheduleCmd.AddCommand(scheduleListCmd)

	// Add delete subcommand
	scheduleCmd.AddCommand(scheduleDeleteCmd)
	scheduleDeleteCmd.Flags().StringVarP(&name, "name", "n", "", "name of the schedule to delete")

	// Add logs subcommand
	scheduleCmd.AddCommand(scheduleLogsCmd)
	scheduleLogsCmd.Flags().StringVarP(&name, "name", "n", "", "name of the schedule to view logs for")
	scheduleLogsCmd.Flags().BoolVarP(&follow, "follow", "f", false, "follow log output (similar to tail -f)")
}
