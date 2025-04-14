package git

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/be-tech/version-manager/internal/utils"
	"github.com/be-tech/version-manager/pkg/config"
	"github.com/be-tech/version-manager/pkg/release"
	"github.com/be-tech/version-manager/pkg/version"
)

type DefaultCommandRunner struct{}

func (r *DefaultCommandRunner) Run(name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	return cmd.CombinedOutput()
}

func (r *DefaultCommandRunner) Output(name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	return cmd.Output()
}

type Manager struct {
	config    *config.Config
	gitCmd    *utils.GitCommands
	logger    *utils.Logger
	delayTime time.Duration
}

func NewManager(config *config.Config) *Manager {
	return &Manager{
		config:    config,
		gitCmd:    utils.NewGitCommands(&DefaultCommandRunner{}),
		logger:    utils.NewLogger(),
		delayTime: 2 * time.Second,
	}
}

func NewManagerWithRunner(config *config.Config, runner utils.CommandRunner) *Manager {
	return &Manager{
		config:    config,
		gitCmd:    utils.NewGitCommands(runner),
		logger:    utils.NewLogger(),
		delayTime: 2 * time.Second,
	}
}

func (m *Manager) ExecuteVersionFlow() error {
	m.logger.Title("Starting deploy!")

	if err := m.checkoutDestinationBranch(); err != nil {
		return err
	}

	if err := m.mergeBranches(false); err != nil {
		return err
	}

	if err := m.pushToRemote(false); err != nil {
		return err
	}

	var newTagVersion string
	var err error

	if m.config.Tag != "" {
		newTagVersion, err = m.createVersionTag()
		if err != nil {
			return err
		}

		if err := m.updateTagOnRemote(); err != nil {
			return err
		}

		if m.config.CreateRelease {
			if err := m.createRelease(newTagVersion); err != nil {
				return err
			}
		}
	}

	if !m.config.RemoveBranch {
		if err := m.checkoutSourceBranch(); err != nil {
			return err
		}

		if err := m.mergeBranches(true); err != nil {
			return err
		}

		if err := m.pushToRemote(true); err != nil {
			return err
		}
	} else {
		if err := m.removeSourceBranch(); err != nil {
			return err
		}
	}

	return nil
}

func (m *Manager) checkoutDestinationBranch() error {
	spinner := utils.NewProgressSpinner(fmt.Sprintf("Checking out to destination branch: %s", m.config.DestinationBranch))

	err := spinner.WithDelay(func() error {
		return m.gitCmd.Checkout(m.config.DestinationBranch)
	}, m.delayTime)

	if err != nil {
		return err
	}

	m.logger.Success("Successfully checked out to destination branch: %s", m.config.DestinationBranch)
	return nil
}

func (m *Manager) checkoutSourceBranch() error {
	if m.config.RemoveBranch {
		return nil
	}

	spinner := utils.NewProgressSpinner(fmt.Sprintf("Checking out to source branch: %s", m.config.SourceBranch))

	err := spinner.WithDelay(func() error {
		return m.gitCmd.Checkout(m.config.SourceBranch)
	}, m.delayTime)

	if err != nil {
		return err
	}

	m.logger.Success("Successfully checked out to source branch: %s", m.config.SourceBranch)
	return nil
}

func (m *Manager) mergeBranches(returning bool) error {
	var source, destination, mergeMessage string

	if !returning {
		source = m.config.SourceBranch
		destination = m.config.DestinationBranch
		mergeMessage = fmt.Sprintf("Merging %s into %s", source, destination)
	} else {
		source = m.config.DestinationBranch
		destination = m.config.SourceBranch
		mergeMessage = fmt.Sprintf("Merging %s into %s", source, destination)
	}

	spinner := utils.NewProgressSpinner(mergeMessage)

	err := spinner.WithDelay(func() error {
		return m.gitCmd.Merge(source)
	}, m.delayTime)

	if err != nil {
		return err
	}

	m.logger.Success("Successfully merged %s into %s", source, destination)
	return nil
}

func (m *Manager) pushToRemote(returning bool) error {
	if !m.config.Push {
		return nil
	}

	var branch string
	if !returning {
		branch = m.config.DestinationBranch
	} else {
		if !m.config.RemoveBranch {
			branch = m.config.SourceBranch
		} else {
			return nil
		}
	}

	spinner := utils.NewProgressSpinner(fmt.Sprintf("Pushing %s to %s", branch, m.config.Remote))

	err := spinner.WithDelay(func() error {
		return m.gitCmd.Push(m.config.Remote, branch)
	}, m.delayTime)

	if err != nil {
		return err
	}

	m.logger.Success("Successfully pushed %s to %s", branch, m.config.Remote)
	return nil
}

func (m *Manager) createVersionTag() (string, error) {
	if m.config.Tag == "" {
		return "", nil
	}

	spinner := utils.NewProgressSpinner(fmt.Sprintf("Creating version tag: %s", m.config.Tag))

	var newTag string

	err := spinner.WithDelay(func() error {
		lastTag, err := m.gitCmd.GetLatestTag()
		if err != nil {
			lastTag = ""
		}

		versionHandler := version.NewHandler()
		generatedTag, err := versionHandler.GenerateNewTag(lastTag, m.config.Tag)
		if err != nil {
			return fmt.Errorf("failed to generate new tag: %v", err)
		}

		newTag = generatedTag

		return m.gitCmd.CreateTag(newTag, fmt.Sprintf("Version %s", newTag))
	}, m.delayTime)

	if err != nil {
		return "", err
	}

	m.logger.Success("Successfully created version tag: %s", newTag)
	return newTag, nil
}

func (m *Manager) updateTagOnRemote() error {
	if !m.config.Push || m.config.Tag == "" {
		return nil
	}

	spinner := utils.NewProgressSpinner("Updating tag in remote repository")

	err := spinner.WithDelay(func() error {
		newtag, err := m.gitCmd.GetLatestTag()
		if err != nil {
			return fmt.Errorf("failed to get latest tag: %v", err)
		}

		if err := m.gitCmd.PushTag(m.config.Remote, newtag); err != nil {
			return err
		}

		if err := m.gitCmd.Push(m.config.Remote, m.config.DestinationBranch); err != nil {
			return err
		}

		return nil
	}, 5*time.Second)

	if err != nil {
		return err
	}

	m.logger.Success("Successfully updated tag in remote repository")
	return nil
}

func (m *Manager) removeSourceBranch() error {
	if !m.config.RemoveBranch {
		return nil
	}

	spinner := utils.NewProgressSpinner(fmt.Sprintf("Removing source branch: %s", m.config.SourceBranch))

	err := spinner.WithDelay(func() error {
		return m.gitCmd.RemoveBranch(m.config.SourceBranch)
	}, m.delayTime)

	if err != nil {
		return err
	}

	m.logger.Success("Successfully removed source branch: %s", m.config.SourceBranch)
	return nil
}

func (m *Manager) createRelease(tagVersion string) error {
	if !m.config.CreateRelease || tagVersion == "" {
		return nil
	}

	spinner := utils.NewProgressSpinner(fmt.Sprintf("Creating release for tag %s on %s", tagVersion, m.config.RepoType))

	err := spinner.WithDelay(func() error {
		releaseManager := release.NewReleaseManager(m.config)
		return releaseManager.CreateRelease(tagVersion)
	}, m.delayTime)

	if err != nil {
		return fmt.Errorf("falha ao criar release: %v", err)
	}

	m.logger.Success("Release criada com sucesso!")
	return nil
}
