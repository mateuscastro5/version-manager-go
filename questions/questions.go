package questions

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/AlecAivazis/survey/v2"
)

// Questions gerencia todas as perguntas feitas ao usuário e armazena as respostas
type Questions struct {
	Remote            string
	SourceBranch      string
	DestinationBranch string
	Push              bool
	RemoveBranch      bool
	Tag               string
}

// NewQuestions cria uma nova instância de Questions
func NewQuestions() *Questions {
	return &Questions{}
}

// PublishQuestions executa todas as perguntas para o usuário
func (q *Questions) PublishQuestions() error {
	var err error

	// Seleção do repositório remoto
	if err = q.selectOrigin(); err != nil {
		return err
	}

	// Seleção da branch de origem
	if err = q.selectSourceBranch(); err != nil {
		return err
	}

	// Seleção da branch de destino
	if err = q.selectDestinationBranch(); err != nil {
		return err
	}

	// Pergunta se deseja fazer push
	if err = q.askWantsPush(); err != nil {
		return err
	}

	// Pergunta se deseja remover a branch de origem
	if err = q.askRemoveFromBranch(); err != nil {
		return err
	}

	// Criação de tag
	if err = q.handleTagCreation(); err != nil {
		return err
	}

	return nil
}

// selectOrigin pergunta qual repositório remoto usar
func (q *Questions) selectOrigin() error {
	// Executar comando git para listar os repositórios remotos
	cmd := exec.Command("git", "remote")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to list git remotes: %v", err)
	}

	remotes := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(remotes) == 0 || (len(remotes) == 1 && remotes[0] == "") {
		return fmt.Errorf("no git remotes found")
	}

	prompt := &survey.Select{
		Message: "Which remote do you want to use?",
		Options: remotes,
	}

	var remote string
	if err := survey.AskOne(prompt, &remote); err != nil {
		return fmt.Errorf("failed to get remote selection: %v", err)
	}

	q.Remote = remote
	fmt.Printf("Remote selected: %s\n", q.Remote)

	return nil
}

// selectSourceBranch pergunta qual branch de origem usar
func (q *Questions) selectSourceBranch() error {
	// Executar comando git para listar as branches
	cmd := exec.Command("git", "branch", "--sort=-worktreepath")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to list git branches: %v", err)
	}

	// Processar a saída para obter a lista de branches
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	branches := make([]string, 0, len(lines))

	for _, line := range lines {
		branch := strings.Replace(line, "*", "", 1)
		branch = strings.TrimSpace(branch)
		if branch != "" {
			branches = append(branches, branch)
		}
	}

	if len(branches) == 0 {
		return fmt.Errorf("no git branches found")
	}

	prompt := &survey.Select{
		Message: "Which source branch do you want to merge into the destination branch?",
		Options: branches,
	}

	var sourceBranch string
	if err := survey.AskOne(prompt, &sourceBranch); err != nil {
		return fmt.Errorf("failed to get source branch selection: %v", err)
	}

	q.SourceBranch = sourceBranch
	fmt.Printf("Source Branch selected: %s\n", q.SourceBranch)

	return nil
}

// selectDestinationBranch pergunta qual branch de destino usar
func (q *Questions) selectDestinationBranch() error {
	// Executar comando git para listar as branches
	cmd := exec.Command("git", "branch", "--sort=-worktreepath")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to list git branches: %v", err)
	}

	// Processar a saída para obter a lista de branches
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	branches := make([]string, 0, len(lines))

	for _, line := range lines {
		branch := strings.Replace(line, "*", "", 1)
		branch = strings.TrimSpace(branch)
		// Filtrar a branch de origem
		if branch != "" && branch != q.SourceBranch {
			branches = append(branches, branch)
		}
	}

	if len(branches) == 0 {
		return fmt.Errorf("no other git branches found")
	}

	prompt := &survey.Select{
		Message: "Which destination branch do you want to merge your source branch into?",
		Options: branches,
	}

	var destinationBranch string
	if err := survey.AskOne(prompt, &destinationBranch); err != nil {
		return fmt.Errorf("failed to get destination branch selection: %v", err)
	}

	q.DestinationBranch = destinationBranch
	fmt.Printf("Destination Branch selected: %s\n", q.DestinationBranch)

	return nil
}

// askWantsPush pergunta se o usuário deseja fazer push para o repositório remoto
func (q *Questions) askWantsPush() error {
	prompt := &survey.Confirm{
		Message: "Do you want to push your changes to the remote repository?",
		Default: false,
	}

	var push bool
	if err := survey.AskOne(prompt, &push); err != nil {
		return fmt.Errorf("failed to get push confirmation: %v", err)
	}

	q.Push = push
	fmt.Printf("It will push: %t\n", q.Push)

	return nil
}

// askRemoveFromBranch pergunta se o usuário deseja remover a branch de origem
func (q *Questions) askRemoveFromBranch() error {
	// Verificar se a branch de origem é uma das branches protegidas
	protectedBranches := []string{"master", "main", "develop", "stage"}
	for _, protected := range protectedBranches {
		if q.SourceBranch == protected {
			q.RemoveBranch = false
			return nil
		}
	}

	prompt := &survey.Confirm{
		Message: "Do you want to remove the source branch?",
		Default: false,
	}

	var remove bool
	if err := survey.AskOne(prompt, &remove); err != nil {
		return fmt.Errorf("failed to get remove branch confirmation: %v", err)
	}

	q.RemoveBranch = remove
	fmt.Printf("It will remove the source branch: %t\n", q.RemoveBranch)

	return nil
}

// handleTagCreation gerencia a criação de tags
func (q *Questions) handleTagCreation() error {
	if !q.Push {
		q.Tag = ""
		return nil
	}

	// Verificar se já existem tags
	cmd := exec.Command("git", "ls-remote", "--tags", "origin")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Run() // Ignorando erros aqui, pois pode não haver tags

	existingTags := out.String()

	if existingTags == "" {
		q.Tag = "1.0.0"
		return nil
	}

	// Opções de versionamento
	var options []string

	// Determinar se estamos em uma branch stage
	isStage := q.DestinationBranch == "stage"

	if isStage {
		options = []string{
			"premajor - version before a major release that is still in development. vX.x.x-x",
			"preminor - version before a minor release that is still in development. vx.X.x-x",
			"prepatch - version before a patch release that is still in development. vx.x.X-x",
			"prerelease - version before a stable release that is still in development. vx.x.x-X",
			"Not a version - Select this to not create a tag",
		}
	} else {
		options = []string{
			"major - significant changes, compatibility impact. vX.x.x",
			"minor - small changes, new features, improvements. vx.X.x",
			"patch - bug fixes, minor changes. vx.x.X",
			"Not a version - Select this to not create a tag",
		}
	}

	prompt := &survey.Select{
		Message: "Which release version do you want to tag?",
		Options: options,
	}

	var versionChoice string
	if err := survey.AskOne(prompt, &versionChoice); err != nil {
		return fmt.Errorf("failed to get version selection: %v", err)
	}

	// Extrair o tipo de versão da opção selecionada
	versionType := strings.Split(versionChoice, " ")[0]

	if versionType == "Not" {
		q.Tag = ""
	} else {
		q.Tag = versionType
	}

	fmt.Printf("Selected release version: %s\n", q.Tag)

	return nil
}
