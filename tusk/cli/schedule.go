package cli

import (
	"fmt"

	"github.com/sboy99/projects.go/tusk/internal/executor"
	"github.com/spf13/cobra"
)

var command string
var interval string

// scheduleCmd represents the schedule command
var scheduleCmd = &cobra.Command{
	Use:   "schedule",
	Short: "Schedule a task",
	Long:  `Schedule a task with the specified name.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if command == "" {
			return fmt.Errorf("schedule command is required")
		}
		if interval == "" {
			return fmt.Errorf("interval is required")
		}

		scheduleExecutor := executor.NewScheduleExecutor(command, interval)
		if err := scheduleExecutor.Execute(); err != nil {
			return fmt.Errorf("failed to execute schedule: %w", err)
		}

		fmt.Printf("Scheduled task with ID '%s'\n", scheduleExecutor.GetID())
		return nil
	},
}

func init() {
	rootCmd.AddCommand(scheduleCmd)
	scheduleCmd.Flags().StringVarP(&command, "command", "c", "", "command to schedule (e.g. echo 'Hello, World!')")
	scheduleCmd.Flags().StringVarP(&interval, "interval", "i", "", "interval for scheduling (e.g. 1m, 1h, 1d, 1w)")
	scheduleCmd.MarkFlagRequired("command")
	scheduleCmd.MarkFlagRequired("interval")
}
