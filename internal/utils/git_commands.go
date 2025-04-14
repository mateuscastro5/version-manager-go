package utils

import (
	"fmt"
	"strings"
)

type GitCommands struct {
	runner CommandRunner
}

type CommandRunner interface {
	Run(name string, args ...string) ([]byte, error)
	Output(name string, args ...string) ([]byte, error)
}

func NewGitCommands(runner CommandRunner) *GitCommands {
	return &GitCommands{
		runner: runner,
	}
}

func (g *GitCommands) Checkout(branch string) error {
	output, err := g.runner.Run("git", "checkout", branch)
	if err != nil {
		return fmt.Errorf("erro ao fazer checkout para a branch %s: %v\n%s", branch, err, output)
	}
	return nil
}

func (g *GitCommands) Merge(sourceBranch string) error {
	output, err := g.runner.Run("git", "merge", sourceBranch)
	if err != nil {
		return fmt.Errorf("erro ao fazer merge da branch %s: %v\n%s", sourceBranch, err, output)
	}
	return nil
}

func (g *GitCommands) Push(remote string, branch string) error {
	output, err := g.runner.Run("git", "push", remote, branch)
	if err != nil {
		return fmt.Errorf("erro ao fazer push da branch %s para %s: %v\n%s", branch, remote, err, output)
	}
	return nil
}

func (g *GitCommands) RemoveBranch(branch string) error {
	output, err := g.runner.Run("git", "branch", "-D", branch)
	if err != nil {
		return fmt.Errorf("erro ao remover a branch %s: %v\n%s", branch, err, output)
	}
	return nil
}

func (g *GitCommands) CreateTag(tag string, message string) error {
	output, err := g.runner.Run("git", "tag", "-a", tag, "-m", message)
	if err != nil {
		return fmt.Errorf("erro ao criar a tag %s: %v\n%s", tag, err, output)
	}
	return nil
}

func (g *GitCommands) GetLatestTag() (string, error) {
	output, err := g.runner.Output("git", "describe", "--tags", "--abbrev=0")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func (g *GitCommands) PushTag(remote string, tag string) error {
	output, err := g.runner.Run("git", "push", remote, tag)
	if err != nil {
		return fmt.Errorf("erro ao fazer push da tag %s para %s: %v\n%s", tag, remote, err, output)
	}
	return nil
}

func (g *GitCommands) GetRemotes() ([]string, error) {
	output, err := g.runner.Output("git", "remote")
	if err != nil {
		return nil, fmt.Errorf("erro ao listar remotos: %v", err)
	}

	remotes := strings.Split(strings.TrimSpace(string(output)), "\n")
	return filterEmptyStrings(remotes), nil
}

func (g *GitCommands) GetBranches() ([]string, error) {
	output, err := g.runner.Output("git", "branch", "--sort=-worktreepath")
	if err != nil {
		return nil, fmt.Errorf("erro ao listar branches: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	branches := make([]string, 0, len(lines))

	for _, line := range lines {
		branch := strings.Replace(line, "*", "", 1)
		branch = strings.TrimSpace(branch)
		if branch != "" {
			branches = append(branches, branch)
		}
	}

	return branches, nil
}

func filterEmptyStrings(slice []string) []string {
	result := make([]string, 0, len(slice))
	for _, s := range slice {
		if s != "" {
			result = append(result, s)
		}
	}
	return result
}
