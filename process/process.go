package process

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
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
	defer s.Stop()

	// Get last tag
	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0")
	lastTagBytes, err := cmd.CombinedOutput()
	if err != nil && !strings.Contains(err.Error(), "no tags can be found") {
		return fmt.Errorf("error getting last tag: %v\n%s", err, lastTagBytes)
	}

	currentTag := strings.TrimSpace(string(lastTagBytes))
	var newVersion *semver.Version

	// Parse or initialize version
	if currentTag == "" {
		newVersion = semver.MustParse("v1.0.0")
	} else {
		v, err := semver.NewVersion(currentTag)
		if err != nil {
			return fmt.Errorf("invalid version format: %s", currentTag)
		}
		newVersion = v
	}

	// Apply version bump
	switch p.questions.Tag {
	case "major":
		v := newVersion.IncMajor()
		newVersion = &v
	case "minor":
		v := newVersion.IncMinor()
		newVersion = &v
	case "patch":
		v := newVersion.IncPatch()
		newVersion = &v
	case "premajor":
		v := newVersion.IncMajor()
		v, _ = v.SetPrerelease("beta")
		newVersion = &v
	case "preminor":
		v := newVersion.IncMinor()
		v, _ = v.SetPrerelease("beta")
		newVersion = &v
	case "prepatch":
		v := newVersion.IncPatch()
		v, _ = v.SetPrerelease("beta")
		newVersion = &v
	case "prerelease":
		pre := newVersion.Prerelease()
		if pre == "" {
			v := *newVersion
			v, _ = v.SetPrerelease("beta.0")
			newVersion = &v
		} else {
			// Increment prerelease version
			v := *newVersion
			v, _ = v.SetPrerelease(pre)
			newVersion = &v
		}
	default:
		return fmt.Errorf("invalid version bump type: %s", p.questions.Tag)
	}

	newTag := "v" + newVersion.String()

	// NPM version
	npmCmd := exec.Command("npm", "version", newTag, "--no-git-tag-version")
	if output, err := npmCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("npm version failed: %v\n%s", err, output)
	}

	// Git operations
	gitCmds := []*exec.Cmd{
		exec.Command("git", "add", "package.json", "package-lock.json"),
		exec.Command("git", "commit", "-m", fmt.Sprintf("chore: release %s", newTag)),
		exec.Command("git", "tag", "-a", newTag, "-m", fmt.Sprintf("Release %s", newTag)),
		exec.Command("git", "push", p.questions.Remote, p.questions.DestinationBranch),
		exec.Command("git", "push", p.questions.Remote, newTag),
	}

	for _, cmd := range gitCmds {
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("git command failed: %v\n%s", err, output)
		}
	}

	s.Stop()
	fmt.Printf("✓ Successfully created and pushed version tag: %s\n", newTag)
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
