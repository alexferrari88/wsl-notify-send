package cmd

import (
	"fmt"
	"wsl-notify-send/internal/config"
	"wsl-notify-send/internal/notify"

	"github.com/spf13/cobra"
)

// Version can be set at build time with: go build -ldflags "-X wsl-notify-send/cmd.Version=x.y.z"
var Version = "0.1.0"

var cfg config.Config

var rootCmd = &cobra.Command{
	Use:   "wsl-notify-send [flags] <title> [message]",
	Short: "Send desktop notifications on Windows and WSL2",
	Long: `wsl-notify-send is a cross-platform notification tool for Windows and WSL2.
It provides a clean interface to send desktop notifications, alerts, and beeps
using the beeep library.

Examples:
  wsl-notify-send "Hello" "World"
  wsl-notify-send --alert "Warning" "Something happened"
  wsl-notify-send --beep
  wsl-notify-send --icon icon.png "Info" "With custom icon"
  wsl-notify-send --app-name "MyApp" "Custom" "From MyApp"`,
	Args: func(cmd *cobra.Command, args []string) error {
		// If version mode, no args required
		if cfg.Version {
			return nil
		}
		
		// If beep mode, no args required
		if cfg.BeepMode {
			return nil
		}
		
		// Otherwise, need at least title
		if len(args) < 1 {
			return fmt.Errorf("requires at least a title argument")
		}
		
		// Maximum 2 args (title and message)
		if len(args) > 2 {
			return fmt.Errorf("too many arguments, expected: <title> [message]")
		}
		
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// Handle version flag
		if cfg.Version {
			cmd.Printf("wsl-notify-send version %s\n", Version)
			return nil
		}
		
		// Validate configuration
		if err := cfg.Validate(); err != nil {
			return fmt.Errorf("invalid configuration: %w", err)
		}
		
		// Handle beep mode
		if cfg.BeepMode {
			return notify.Beep(cfg.Frequency, cfg.Duration)
		}
		
		// Parse title and message
		title := args[0]
		message := ""
		if len(args) > 1 {
			message = args[1]
		}
		
		// Send notification
		if cfg.AlertMode {
			return notify.Alert(title, message, cfg.Icon, cfg.AppName)
		}
		
		return notify.Notify(title, message, cfg.Icon, cfg.AppName)
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func IsQuietMode() bool {
	return cfg.Quiet
}

func init() {
	// Notification mode flags
	rootCmd.Flags().BoolVarP(&cfg.AlertMode, "alert", "a", false, "Send alert notification with sound")
	rootCmd.Flags().BoolVarP(&cfg.BeepMode, "beep", "b", false, "Just beep (no notification)")
	
	// Content flags
	rootCmd.Flags().StringVarP(&cfg.Icon, "icon", "i", "", "Icon file path or stock icon name")
	rootCmd.Flags().StringVar(&cfg.AppName, "app-name", "wsl-notify-send", "Application name")
	
	// Beep customization flags
	rootCmd.Flags().Float64Var(&cfg.Frequency, "freq", 587.0, "Beep frequency in Hz")
	rootCmd.Flags().IntVar(&cfg.Duration, "duration", 500, "Beep duration in milliseconds")
	
	// Utility flags
	rootCmd.Flags().BoolVarP(&cfg.Quiet, "quiet", "q", false, "Suppress error output")
	rootCmd.Flags().BoolVar(&cfg.Version, "version", false, "Show version information")
	
	// Mark beep-related flags as hidden when not in beep mode
	rootCmd.Flags().MarkHidden("freq")
	rootCmd.Flags().MarkHidden("duration")
}