package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/nordlys-labs/nordlys-installer/internal/config"
	"github.com/nordlys-labs/nordlys-installer/internal/constants"
	"github.com/nordlys-labs/nordlys-installer/internal/runtime"
	"github.com/nordlys-labs/nordlys-installer/internal/tools"
	"github.com/nordlys-labs/nordlys-installer/internal/ui"
	"github.com/nordlys-labs/nordlys-installer/internal/updater"
	"github.com/spf13/cobra"
)

var (
	apiKey         string
	model          string
	toolsList      []string
	nonInteractive bool
)

var rootCmd = &cobra.Command{
	Use:   "nordlys-installer",
	Short: "Configure developer tools to use Nordlys",
	Long: `Nordlys Installer configures your developer tools to use Nordlys's
Mixture of Models for intelligent model selection and cost optimization.

Supported tools: claude-code, opencode, codex, gemini-cli, grok-cli, qwen-code, zed`,
	Run: runInstaller,
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update nordlys-installer to the latest version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Checking for updates...")
		if err := updater.SelfUpdate(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

var uninstallCmd = &cobra.Command{
	Use:   "uninstall [tool...]",
	Short: "Remove Nordlys configuration from tools",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, toolName := range args {
			tool := tools.GetToolByName(toolName)
			if tool == nil {
				fmt.Fprintf(os.Stderr, "Unknown tool: %s\n", toolName)
				continue
			}

			if err := tool.Uninstall(); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to uninstall %s: %v\n", toolName, err)
				continue
			}

			fmt.Printf("✅ Removed Nordlys configuration from %s\n", toolName)
		}
	},
}

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate Nordlys configuration",
	Run: func(cmd *cobra.Command, args []string) {
		apiKey := os.Getenv("NORDLYS_API_KEY")
		if apiKey == "" {
			fmt.Fprintln(os.Stderr, "NORDLYS_API_KEY environment variable not set")
			os.Exit(1)
		}

		fmt.Println("Validating API key format...")
		if err := config.ValidateAPIKey(apiKey); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Invalid API key: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✅ API key format is valid")

		fmt.Println("Testing API connection...")
		if err := config.ValidateAPIConnection(apiKey); err != nil {
			fmt.Fprintf(os.Stderr, "❌ API connection failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✅ API connection successful")

		fmt.Println("\nValidating tool configurations...")
		for _, tool := range tools.GetAllTools() {
			if err := tool.Validate(); err != nil {
				fmt.Printf("⚠️  %s: %v\n", tool.Name(), err)
			} else {
				fmt.Printf("✅ %s: configured\n", tool.Name())
			}
		}
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("nordlys-installer v%s\n", constants.Version)
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all supported tools",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Supported tools:")
		for _, tool := range tools.GetAllTools() {
			installed := "✗"
			if tool.IsInstalled() {
				installed = "✓"
			}
			node := ""
			if tool.RequiresNode() {
				node = " (requires Node.js)"
			}
			fmt.Printf("  [%s] %s - %s%s\n", installed, tool.Name(), tool.Description(), node)
		}
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&apiKey, "api-key", "k", "", "Nordlys API key")
	rootCmd.PersistentFlags().StringVarP(&model, "model", "m", constants.DefaultModel, "Model to use")
	rootCmd.Flags().StringSliceVarP(&toolsList, "tools", "t", nil, "Tools to configure (comma-separated)")
	rootCmd.Flags().BoolVar(&nonInteractive, "non-interactive", false, "Run without TUI")

	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(uninstallCmd)
	rootCmd.AddCommand(validateCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(listCmd)

	for _, tool := range tools.GetAllTools() {
		toolCmd := &cobra.Command{
			Use:   tool.Name(),
			Short: tool.Description(),
		}
		toolCmd.AddCommand(&cobra.Command{
			Use:   "update",
			Short: "Update Nordlys configuration",
			Run:   runToolUpdate(tool),
		})
		rootCmd.AddCommand(toolCmd)
	}
}

func runInstaller(cmd *cobra.Command, args []string) {
	if apiKey == "" {
		apiKey = os.Getenv("NORDLYS_API_KEY")
	}

	if !nonInteractive && len(toolsList) == 0 && apiKey == "" {
		if err := ui.RunTUI(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "Error: API key required (--api-key or NORDLYS_API_KEY)")
		os.Exit(1)
	}

	if err := config.ValidateAPIKey(apiKey); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	var selectedTools []tools.Tool
	if len(toolsList) > 0 {
		for _, name := range toolsList {
			name = strings.TrimSpace(name)
			tool := tools.GetToolByName(name)
			if tool == nil {
				fmt.Fprintf(os.Stderr, "Unknown tool: %s\n", name)
				os.Exit(1)
			}
			selectedTools = append(selectedTools, tool)
		}
	} else {
		selectedTools = tools.GetAllTools()
	}

	for _, tool := range selectedTools {
		fmt.Printf("Configuring %s...\n", tool.Name())

		if tool.RequiresNode() {
			if err := runtime.EnsureNodeJS(); err != nil {
				fmt.Fprintf(os.Stderr, "❌ %s: %v\n", tool.Name(), err)
				continue
			}
		}

		if err := tool.UpdateConfig(apiKey, model, constants.APIBaseURL); err != nil {
			fmt.Fprintf(os.Stderr, "❌ %s: %v\n", tool.Name(), err)
			continue
		}

		fmt.Printf("✅ %s configured successfully\n", tool.Name())
	}
}

func runToolUpdate(tool tools.Tool) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		resolvedAPIKey := apiKey
		if resolvedAPIKey == "" {
			resolvedAPIKey = os.Getenv("NORDLYS_API_KEY")
		}
		if resolvedAPIKey == "" {
			existingKey, _ := tool.GetExistingConfig()
			resolvedAPIKey = existingKey
		}
		if resolvedAPIKey == "" {
			fmt.Fprintln(os.Stderr, "Error: API key required (--api-key, NORDLYS_API_KEY, or existing config)")
			os.Exit(1)
		}
		if err := config.ValidateAPIKey(resolvedAPIKey); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		resolvedModel := model
		if resolvedModel == "" {
			_, existingModel := tool.GetExistingConfig()
			resolvedModel = existingModel
		}
		if resolvedModel == "" {
			resolvedModel = constants.DefaultModel
		}

		fmt.Printf("Updating %s...\n", tool.Name())

		if tool.RequiresNode() {
			if err := runtime.EnsureNodeJS(); err != nil {
				fmt.Fprintf(os.Stderr, "❌ %s: %v\n", tool.Name(), err)
				os.Exit(1)
			}
		}

		if err := tool.UpdateConfig(resolvedAPIKey, resolvedModel, constants.APIBaseURL); err != nil {
			fmt.Fprintf(os.Stderr, "❌ %s: %v\n", tool.Name(), err)
			os.Exit(1)
		}

		fmt.Printf("✅ %s updated successfully\n", tool.Name())
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
