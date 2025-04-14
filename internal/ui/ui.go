package ui

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/be-tech/version-manager/internal/utils"
	"github.com/be-tech/version-manager/pkg/config"
)

type UI struct {
	config *config.Config
	logger *utils.Logger
	gitCmd *utils.GitCommands
}

type DefaultCommandRunner struct{}

func (r *DefaultCommandRunner) Run(name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	return cmd.CombinedOutput()
}

func (r *DefaultCommandRunner) Output(name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	return cmd.Output()
}

func NewUI() *UI {
	return &UI{
		config: config.NewConfig(),
		logger: utils.NewLogger(),
		gitCmd: utils.NewGitCommands(&DefaultCommandRunner{}),
	}
}

func (u *UI) CollectUserInput() (*config.Config, error) {
	var err error

	if err = u.selectOrigin(); err != nil {
		return nil, err
	}

	if err = u.selectSourceBranch(); err != nil {
		return nil, err
	}

	if err = u.selectDestinationBranch(); err != nil {
		return nil, err
	}

	if err = u.askWantsPush(); err != nil {
		return nil, err
	}

	if err = u.askRemoveFromBranch(); err != nil {
		return nil, err
	}

	if err = u.handleVersionTag(); err != nil {
		return nil, err
	}

	if u.config.Tag != "" {
		if err = u.askCreateRelease(); err != nil {
			return nil, err
		}

		if u.config.CreateRelease {
			if err = u.selectRepoType(); err != nil {
				return nil, err
			}

			if err = u.collectReleaseInfo(); err != nil {
				return nil, err
			}
		}
	}

	return u.config, nil
}

func (u *UI) logChoice(description string, choice interface{}) {
	u.logger.Info("%s: %v", description, choice)
}

func (u *UI) selectOrigin() error {
	remotes, err := u.gitCmd.GetRemotes()
	if err != nil {
		return fmt.Errorf("falha ao listar repositórios remotos: %v", err)
	}

	if len(remotes) == 0 {
		return fmt.Errorf("nenhum repositório remoto encontrado")
	}

	prompt := &survey.Select{
		Message: "Qual repositório remoto você deseja usar?",
		Options: remotes,
	}

	var remote string
	if err := survey.AskOne(prompt, &remote); err != nil {
		return fmt.Errorf("falha na seleção do repositório remoto: %v", err)
	}

	u.config.Remote = remote
	u.logChoice("Repositório remoto selecionado", u.config.Remote)

	return nil
}

func (u *UI) selectSourceBranch() error {
	branches, err := u.gitCmd.GetBranches()
	if err != nil {
		return err
	}

	prompt := &survey.Select{
		Message: "Qual branch de origem você deseja mesclar na branch de destino?",
		Options: branches,
	}

	var sourceBranch string
	if err := survey.AskOne(prompt, &sourceBranch); err != nil {
		return fmt.Errorf("falha na seleção da branch de origem: %v", err)
	}

	u.config.SourceBranch = sourceBranch
	u.logChoice("Branch de origem selecionada", u.config.SourceBranch)

	return nil
}

func (u *UI) selectDestinationBranch() error {
	branches, err := u.gitCmd.GetBranches()
	if err != nil {
		return err
	}

	filteredBranches := make([]string, 0, len(branches))
	for _, branch := range branches {
		if branch != u.config.SourceBranch {
			filteredBranches = append(filteredBranches, branch)
		}
	}

	if len(filteredBranches) == 0 {
		return fmt.Errorf("nenhuma outra branch encontrada")
	}

	prompt := &survey.Select{
		Message: "Qual branch de destino você deseja para mesclar sua branch de origem?",
		Options: filteredBranches,
	}

	var destinationBranch string
	if err := survey.AskOne(prompt, &destinationBranch); err != nil {
		return fmt.Errorf("falha na seleção da branch de destino: %v", err)
	}

	u.config.DestinationBranch = destinationBranch
	u.logChoice("Branch de destino selecionada", u.config.DestinationBranch)

	return nil
}

func (u *UI) askWantsPush() error {
	prompt := &survey.Confirm{
		Message: "Você deseja enviar as alterações para o repositório remoto?",
		Default: false,
	}

	var push bool
	if err := survey.AskOne(prompt, &push); err != nil {
		return fmt.Errorf("falha ao obter confirmação de push: %v", err)
	}

	u.config.Push = push
	u.logChoice("Enviará para o repositório remoto", u.config.Push)

	return nil
}

func (u *UI) askRemoveFromBranch() error {
	protectedBranches := []string{"master", "main", "develop", "stage"}
	for _, protected := range protectedBranches {
		if u.config.SourceBranch == protected {
			u.config.RemoveBranch = false
			return nil
		}
	}

	prompt := &survey.Confirm{
		Message: "Você deseja remover a branch de origem?",
		Default: false,
	}

	var remove bool
	if err := survey.AskOne(prompt, &remove); err != nil {
		return fmt.Errorf("falha ao obter confirmação para remover branch: %v", err)
	}

	u.config.RemoveBranch = remove
	u.logChoice("Removerá a branch de origem", u.config.RemoveBranch)

	return nil
}

func (u *UI) handleVersionTag() error {
	if !u.config.Push {
		u.config.Tag = ""
		return nil
	}

	output, err := u.gitCmd.GetLatestTag()
	hasExistingTags := err == nil && output != ""

	if !hasExistingTags {
		u.config.Tag = "1.0.0"
		u.logChoice("Tag de versão inicial", u.config.Tag)
		return nil
	}

	var options []string

	isStage := u.config.DestinationBranch == "stage"

	if isStage {
		options = []string{
			"premajor - versão antes de um lançamento principal que ainda está em desenvolvimento. vX.x.x-x",
			"preminor - versão antes de um lançamento secundário que ainda está em desenvolvimento. vx.X.x-x",
			"prepatch - versão antes de um lançamento de correção que ainda está em desenvolvimento. vx.x.X-x",
			"prerelease - versão antes de um lançamento estável que ainda está em desenvolvimento. vx.x.x-X",
			"Not a version - Selecione para não criar uma tag",
		}
	} else {
		options = []string{
			"major - mudanças significativas, impacto na compatibilidade. vX.x.x",
			"minor - pequenas mudanças, novos recursos, melhorias. vx.X.x",
			"patch - correção de bugs, pequenas mudanças. vx.x.X",
			"Not a version - Selecione para não criar uma tag",
		}
	}

	prompt := &survey.Select{
		Message: "Qual versão de lançamento você deseja para a tag?",
		Options: options,
	}

	var versionChoice string
	if err := survey.AskOne(prompt, &versionChoice); err != nil {
		return fmt.Errorf("falha ao obter seleção de versão: %v", err)
	}

	versionParts := strings.SplitN(versionChoice, " ", 2)
	versionType := versionParts[0]

	if versionType == "Not" {
		u.config.Tag = ""
	} else {
		u.config.Tag = versionType
	}

	u.logChoice("Versão de lançamento selecionada", u.config.Tag)

	return nil
}

func (u *UI) askCreateRelease() error {
	prompt := &survey.Confirm{
		Message: "Você deseja criar uma release no GitHub/GitLab com esta tag?",
		Default: false,
	}

	var createRelease bool
	if err := survey.AskOne(prompt, &createRelease); err != nil {
		return fmt.Errorf("falha ao obter confirmação para criar release: %v", err)
	}

	u.config.CreateRelease = createRelease
	u.logChoice("Criar release", u.config.CreateRelease)

	return nil
}

func (u *UI) selectRepoType() error {
	prompt := &survey.Select{
		Message: "Qual é o tipo do seu repositório?",
		Options: []string{"GitHub", "GitLab"},
	}

	var repoType string
	if err := survey.AskOne(prompt, &repoType); err != nil {
		return fmt.Errorf("falha na seleção do tipo de repositório: %v", err)
	}

	u.config.RepoType = strings.ToLower(repoType)
	u.logChoice("Tipo de repositório", u.config.RepoType)

	return nil
}

func (u *UI) collectReleaseInfo() error {
	titlePrompt := &survey.Input{
		Message: "Título da release (deixe em branco para usar a tag):",
	}

	var title string
	if err := survey.AskOne(titlePrompt, &title); err != nil {
		return fmt.Errorf("falha ao obter título da release: %v", err)
	}

	if title == "" {
		title = "Release v" + u.config.Tag
	}
	u.config.ReleaseTitle = title
	u.logChoice("Título da release", u.config.ReleaseTitle)

	notesPrompt := &survey.Multiline{
		Message: "Notas da release (descrição das mudanças):",
	}

	var notes string
	if err := survey.AskOne(notesPrompt, &notes); err != nil {
		return fmt.Errorf("falha ao obter notas da release: %v", err)
	}

	u.config.ReleaseNotes = notes
	u.logChoice("Notas da release foram preenchidas", len(notes) > 0)

	return nil
}
