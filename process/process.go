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

type Process struct {
	questions *questions.Questions
}

func NewProcess(q *questions.Questions) *Process {
	return &Process{
		questions: q,
	}
}

func (p *Process) ExecuteCommands() error {
	green := color.New(color.FgHiGreen, color.Bold).SprintFunc()
	fmt.Println(green("Starting deploy!"))

	if err := p.checkoutDestinationBranch(); err != nil {
		return err
	}

	if err := p.mergeBranches(false); err != nil {
		return err
	}

	if err := p.pushToRemote(false); err != nil {
		return err
	}

	if err := p.createTag(); err != nil {
		return err
	}

	if err := p.updateTag(); err != nil {
		return err
	}

	if !p.questions.RemoveBranch {
		if err := p.checkoutSourceBranch(); err != nil {
			return err
		}

		if err := p.mergeBranches(true); err != nil {
			return err
		}

		if err := p.pushToRemote(true); err != nil {
			return err
		}
	} else {
		if err := p.removeSourceBranch(); err != nil {
			return err
		}
	}

	return nil
}

func (p *Process) checkoutDestinationBranch() error {
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Suffix = fmt.Sprintf(" Checking out to destination branch: %s", p.questions.DestinationBranch)
	s.Start()

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

func (p *Process) checkoutSourceBranch() error {
	if p.questions.RemoveBranch {
		return nil
	}

	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Suffix = fmt.Sprintf(" Checking out to source branch: %s", p.questions.SourceBranch)
	s.Start()

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

func (p *Process) createTag() error {
	if p.questions.Tag == "" {
		return nil
	}

	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Suffix = fmt.Sprintf(" Creating version tag: %s", p.questions.Tag)
	s.Start()

	time.Sleep(2 * time.Second)

	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0")
	lastTag, _ := cmd.Output()

	var newTag string
	if len(lastTag) == 0 {
		newTag = "v1.0.0"
	} else {
		currentTag := strings.TrimSpace(string(lastTag))

		tagVersion := currentTag
		if strings.HasPrefix(currentTag, "v") {
			tagVersion = currentTag[1:]
		}

		parts := strings.Split(tagVersion, ".")
		if len(parts) < 3 {
			newTag = "v1.0.0"
		} else {
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
				newTag = currentTag + "-pre"
			default:
				newTag = currentTag
			}
		}
	}

	cmd = exec.Command("npm", "version", newTag, "--allow-same-version=true", fmt.Sprintf("Version %s", newTag))
	updateTag, err := cmd.CombinedOutput()
	if err != nil {
		s.Stop()
		return fmt.Errorf("failed to create tag: %v\n%s", err, updateTag)
	}

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

func parseVersionPart(part string) int {
	var version int
	fmt.Sscanf(part, "%d", &version)
	return version
}

func (p *Process) updateTag() error {
	if !p.questions.Push || p.questions.Tag == "" {
		return nil
	}

	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Suffix = " Updating tag in remote repository"
	s.Start()

	time.Sleep(5 * time.Second)

	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0")
	output, err := cmd.Output()
	if err != nil {
		s.Stop()
		return fmt.Errorf("failed to get latest tag: %v", err)
	}
	newtag := strings.TrimSpace(string(output))

	cmd = exec.Command("git", "push", p.questions.Remote, newtag)
	_, err = cmd.CombinedOutput()
	if err != nil {
		s.Stop()
		return fmt.Errorf("failed to push tag to remote: %v", err)
	}

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

func (p *Process) removeSourceBranch() error {
	if !p.questions.RemoveBranch {
		return nil
	}

	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Suffix = fmt.Sprintf(" Removing source branch: %s", p.questions.SourceBranch)
	s.Start()

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
