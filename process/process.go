package process

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/be-tech/version-manager/questions"
	"github.com/briandowns/spinner"
	"github.com/fatih/color"
)

// Process gerencia a execução dos comandos git baseados nas respostas do usuário
type Process struct {
	questions *questions.Questions
}

// NewProcess cria uma nova instância do Process
func NewProcess(q *questions.Questions) *Process {
	return &Process{
		questions: q,
	}
}

// ExecuteCommands executa todos os comandos do processo
func (p *Process) ExecuteCommands() error {
	green := color.New(color.FgHiGreen, color.Bold).SprintFunc()
	fmt.Println(green("Starting deploy!"))

	// Checkout para a branch de destino
	if err := p.checkoutDestinationBranch(); err != nil {
		return err
	}

	// Merge das branches
	if err := p.mergeBranches(false); err != nil {
		return err
	}

	// Push para o repositório remoto
	if err := p.pushToRemote(false); err != nil {
		return err
	}

	// Criar tag
	if err := p.createTag(); err != nil {
		return err
	}

	// Atualizar tag no repositório remoto
	if err := p.updateTag(); err != nil {
		return err
	}

	if !p.questions.RemoveBranch {
		// Checkout para a branch de origem
		if err := p.checkoutSourceBranch(); err != nil {
			return err
		}

		// Merge das branches (volta)
		if err := p.mergeBranches(true); err != nil {
			return err
		}

		// Push para o repositório remoto (volta)
		if err := p.pushToRemote(true); err != nil {
			return err
		}
	} else {
		// Remover branch de origem
		if err := p.removeSourceBranch(); err != nil {
			return err
		}
	}

	return nil
}

// checkoutDestinationBranch faz checkout para a branch de destino
func (p *Process) checkoutDestinationBranch() error {
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Suffix = fmt.Sprintf(" Checking out to destination branch: %s", p.questions.DestinationBranch)
	s.Start()

	// Delay para dar sensação de processamento (similar ao delay no código JS)
	time.Sleep(2 * time.Second)

	cmd := exec.Command("git", "checkout", p.questions.DestinationBranch)
	output, err := cmd.CombinedOutput()
	if err != nil {
		s.Stop()
		return fmt.Errorf("failed to checkout destination branch: %v\n%s", err, output)
	}

	s.Stop()
	fmt.Printf("✓ Successfully checked out to destination branch: %s\n", p.questions.DestinationBranch)
	return nil
}

// checkoutSourceBranch faz checkout para a branch de origem
func (p *Process) checkoutSourceBranch() error {
	if p.questions.RemoveBranch {
		return nil
	}

	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Suffix = fmt.Sprintf(" Checking out to source branch: %s", p.questions.SourceBranch)
	s.Start()

	// Delay para dar sensação de processamento (similar ao delay no código JS)
	time.Sleep(2 * time.Second)

	cmd := exec.Command("git", "checkout", p.questions.SourceBranch)
	output, err := cmd.CombinedOutput()
	if err != nil {
		s.Stop()
		return fmt.Errorf("failed to checkout source branch: %v\n%s", err, output)
	}

	s.Stop()
	fmt.Printf("✓ Successfully checked out to source branch: %s\n", p.questions.SourceBranch)
	return nil
}

// mergeBranches realiza o merge entre branches
func (p *Process) mergeBranches(returning bool) error {
	var source, destination, mergeMessage string

	if !returning {
		source = p.questions.SourceBranch
		destination = p.questions.DestinationBranch
		mergeMessage = fmt.Sprintf("Merging %s into %s", source, destination)
	} else {
		source = p.questions.DestinationBranch
		destination = p.questions.SourceBranch
		mergeMessage = fmt.Sprintf("Merging %s into %s", source, destination)
	}

	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Suffix = fmt.Sprintf(" %s", mergeMessage)
	s.Start()

	// Delay para dar sensação de processamento (similar ao delay no código JS)
	time.Sleep(2 * time.Second)

	cmd := exec.Command("git", "merge", source)
	output, err := cmd.CombinedOutput()
	if err != nil {
		s.Stop()
		return fmt.Errorf("failed to merge branches: %v\n%s", err, output)
	}

	s.Stop()
	fmt.Printf("✓ Successfully merged %s into %s\n", source, destination)
	return nil
}

// pushToRemote faz push para o repositório remoto
func (p *Process) pushToRemote(returning bool) error {
	if !p.questions.Push {
		return nil
	}

	var branch string
	if !returning {
		branch = p.questions.DestinationBranch
	} else {
		if !p.questions.RemoveBranch {
			branch = p.questions.SourceBranch
		} else {
			return nil
		}
	}

	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Suffix = fmt.Sprintf(" Pushing %s to %s", branch, p.questions.Remote)
	s.Start()

	// Delay para dar sensação de processamento (similar ao delay no código JS)
	time.Sleep(2 * time.Second)

	cmd := exec.Command("git", "push", p.questions.Remote, branch)
	output, err := cmd.CombinedOutput()
	if err != nil {
		s.Stop()
		return fmt.Errorf("failed to push to remote: %v\n%s", err, output)
	}

	s.Stop()
	fmt.Printf("✓ Successfully pushed %s to %s\n", branch, p.questions.Remote)
	return nil
}

// createTag cria uma nova tag de versão
func (p *Process) createTag() error {
	if p.questions.Tag == "" {
		return nil
	}

	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Suffix = fmt.Sprintf(" Creating version tag: %s", p.questions.Tag)
	s.Start()

	// Delay para dar sensação de processamento (similar ao delay no código JS)
	time.Sleep(2 * time.Second)

	// Em Go, precisamos gerenciar a versão diretamente com git tag em vez de npm version
	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0")
	lastTag, _ := cmd.Output() // Ignora erro se não houver tag anterior

	var newTag string
	if len(lastTag) == 0 {
		// Se não há tag anterior, começamos com v1.0.0
		newTag = "v1.0.0"
	} else {
		// Remove nova linha e quaisquer caracteres de controle
		currentTag := strings.TrimSpace(string(lastTag))

		// Remove o 'v' inicial se houver
		tagVersion := currentTag
		if strings.HasPrefix(currentTag, "v") {
			tagVersion = currentTag[1:]
		}

		// Divida a versão em partes
		parts := strings.Split(tagVersion, ".")
		if len(parts) < 3 {
			// Se a versão não estiver no formato esperado, use v1.0.0
			newTag = "v1.0.0"
		} else {
			// Incrementa a versão com base no tipo selecionado
			switch p.questions.Tag {
			case "major":
				major := parseVersionPart(parts[0]) + 1
				newTag = fmt.Sprintf("v%d.0.0", major)
			case "minor":
				major := parseVersionPart(parts[0])
				minor := parseVersionPart(parts[1]) + 1
				newTag = fmt.Sprintf("v%d.%d.0", major, minor)
			case "patch":
				major := parseVersionPart(parts[0])
				minor := parseVersionPart(parts[1])
				patch := parseVersionPart(parts[2]) + 1
				newTag = fmt.Sprintf("v%d.%d.%d", major, minor, patch)
			case "premajor", "preminor", "prepatch", "prerelease":
				// Implementação simplificada de pre-releases
				newTag = currentTag + "-pre"
			default:
				newTag = currentTag
			}
		}
	}

	// Criar a tag com git tag
	cmd = exec.Command("git", "tag", "-a", newTag, "-m", fmt.Sprintf("Version %s", newTag))
	output, err := cmd.CombinedOutput()
	if err != nil {
		s.Stop()
		return fmt.Errorf("failed to create tag: %v\n%s", err, output)
	}

	s.Stop()
	fmt.Printf("✓ Successfully created version tag: %s\n", newTag)
	return nil
}

// parseVersionPart converte uma parte da versão de string para int
func parseVersionPart(part string) int {
	var version int
	fmt.Sscanf(part, "%d", &version)
	return version
}

// updateTag atualiza a tag no repositório remoto
func (p *Process) updateTag() error {
	if !p.questions.Push || p.questions.Tag == "" {
		return nil
	}

	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Suffix = " Updating tag in remote repository"
	s.Start()

	// Delay para dar sensação de processamento (similar ao delay no código JS)
	time.Sleep(5 * time.Second)

	// Obter a tag mais recente
	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0")
	output, err := cmd.Output()
	if err != nil {
		s.Stop()
		return fmt.Errorf("failed to get latest tag: %v", err)
	}
	newtag := strings.TrimSpace(string(output))

	// Push a tag para o remote
	cmd = exec.Command("git", "push", p.questions.Remote, newtag)
	_, err = cmd.CombinedOutput()
	if err != nil {
		s.Stop()
		return fmt.Errorf("failed to push tag to remote: %v", err)
	}

	// Push a branch destination para o remote
	cmd = exec.Command("git", "push", p.questions.Remote, p.questions.DestinationBranch)
	_, err = cmd.CombinedOutput()
	if err != nil {
		s.Stop()
		return fmt.Errorf("failed to push destination branch to remote: %v", err)
	}

	s.Stop()
	fmt.Printf("✓ Successfully updated tag %s in remote repository\n", newtag)
	return nil
}

// removeSourceBranch remove a branch de origem
func (p *Process) removeSourceBranch() error {
	if !p.questions.RemoveBranch {
		return nil
	}

	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Suffix = fmt.Sprintf(" Removing source branch: %s", p.questions.SourceBranch)
	s.Start()

	// Delay para dar sensação de processamento (similar ao delay no código JS)
	time.Sleep(2 * time.Second)

	cmd := exec.Command("git", "branch", "-D", p.questions.SourceBranch)
	output, err := cmd.CombinedOutput()
	if err != nil {
		s.Stop()
		return fmt.Errorf("failed to remove source branch: %v\n%s", err, output)
	}

	s.Stop()
	fmt.Printf("✓ Successfully removed source branch: %s\n", p.questions.SourceBranch)
	return nil
}
