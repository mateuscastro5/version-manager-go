package questions

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/AlecAivazis/survey/v2"
)

type Questions struct {
	Remote            string
	SourceBranch      string
	DestinationBranch string
	Push              bool
	RemoveBranch      bool
	Tag               string
}

func NewQuestions() *Questions {
	return &Questions{}
}

func (q *Questions) PublishQuestions() error {
	var err error

	if err = q.selectOrigin(); err != nil {
		return err
	}

	if err = q.selectSourceBranch(); err != nil {
		return err
	}

	if err = q.selectDestinationBranch(); err != nil {
		return err
	}

	if err = q.askWantsPush(); err != nil {
		return err
	}

	if err = q.askRemoveFromBranch(); err != nil {
		return err
	}

	if err = q.handleTagCreation(); err != nil {
		return err
	}

	return nil
}

func (q *Questions) selectOrigin() error {
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

func (q *Questions) selectSourceBranch() error {
	cmd := exec.Command("git", "branch", "--sort=-worktreepath")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to list git branches: %v", err)
	}

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

func (q *Questions) selectDestinationBranch() error {
	cmd := exec.Command("git", "branch", "--sort=-worktreepath")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to list git branches: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	branches := make([]string, 0, len(lines))

	for _, line := range lines {
		branch := strings.Replace(line, "*", "", 1)
		branch = strings.TrimSpace(branch)
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

func (q *Questions) askRemoveFromBranch() error {
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

func (q *Questions) handleTagCreation() error {
	if !q.Push {
		q.Tag = ""
		return nil
	}

	cmd := exec.Command("git", "ls-remote", "--tags", "origin")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Run()

	existingTags := out.String()

	if existingTags == "" {
		q.Tag = "1.0.0"
		return nil
	}

	var options []string

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

	versionType := strings.Split(versionChoice, " ")[0]

	if versionType == "Not" {
		q.Tag = ""
	} else {
		q.Tag = versionType
	}

	fmt.Printf("Selected release version: %s\n", q.Tag)

	return nil
}
