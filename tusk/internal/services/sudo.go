package services

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/sboy99/projects.go/tusk/pkg/logger"
)

// SudoService handles sudo privilege management
type SudoService struct {
	logger *logger.Logger
}

// NewSudoService creates a new SudoService instance
func NewSudoService() *SudoService {
	return &SudoService{
		logger: logger.NewLogger(""),
	}
}

// RequestPrivileges checks if the user has sudo privileges and prompts for password if needed
func (s *SudoService) RequestPrivileges() error {
	// Check if already running as root
	if os.Geteuid() == 0 {
		s.logger.Info("running as root, no sudo required")
		return nil
	}

	// Check if sudo is available
	if _, err := exec.LookPath("sudo"); err != nil {
		return fmt.Errorf("sudo command not found: %w", err)
	}

	// Validate sudo privileges by running sudo -v
	// This will prompt for password if needed and cache credentials
	s.logger.Info("requesting sudo privileges...")
	cmd := exec.Command("sudo", "-v")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to obtain sudo privileges: %w", err)
	}

	s.logger.Success("sudo privileges obtained")
	return nil
}
