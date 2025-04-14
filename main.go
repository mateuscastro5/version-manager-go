package main

import (
	"os"

	"github.com/be-tech/version-manager/internal/git"
	"github.com/be-tech/version-manager/internal/ui"
	"github.com/be-tech/version-manager/internal/utils"
)

func main() {
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
