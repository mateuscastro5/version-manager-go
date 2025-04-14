package main

import (
	"os"

	"github.com/be-tech/version-manager/internal/git"
	"github.com/be-tech/version-manager/internal/ui"
	"github.com/be-tech/version-manager/internal/utils"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file right at the start of the application
	loadEnvFile()

	logger := utils.NewLogger()
	logger.Title("Version Manager - Git version control tool")

	userInterface := ui.NewUI()

	config, err := userInterface.CollectUserInput()
	if err != nil {
		logger.Error("Erro ao processar entrada do usuário: %v", err)
		os.Exit(1)
	}

	gitManager := git.NewManager(config)

	if err := gitManager.ExecuteVersionFlow(); err != nil {
		logger.Error("Erro ao executar operações Git: %v", err)
		os.Exit(1)
	}

	logger.Success("Gerenciamento de versão concluído com sucesso!")
}

// loadEnvFile tries to load environment variables from .env files
// searching in multiple locations to ensure it works regardless of where
// the application is executed from
func loadEnvFile() {
	// Try different locations for the .env file
	paths := []string{
		".env",    // Current directory
		"../.env", // Parent directory
	}

	// Try each path and stop at the first successful load
	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			_ = godotenv.Load(path)
			return
		}
	}

	// Fall back to the default load which tries the current directory
	_ = godotenv.Load()
}
