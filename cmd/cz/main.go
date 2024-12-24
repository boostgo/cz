package main

import (
	"cz/internal/git"
	"fmt"
	"github.com/boostgo/lite/log"
	"github.com/boostgo/lite/system/trace"
	"github.com/manifoldco/promptui"
	"os"
)

// CommitType represents a type of commit
type CommitType struct {
	Label       string
	Description string
}

var commitTypes = []CommitType{
	{Label: "feat", Description: "A new feature"},
	{Label: "fix", Description: "A bug fix"},
	{Label: "docs", Description: "Documentation only changes"},
	{Label: "ref", Description: "A code change that neither fixes a bug nor adds a feature"},
	{Label: "perf", Description: "A code change that improves performance"},
	{Label: "test", Description: "Adding missing tests or correcting existing tests"},
	{Label: "ci", Description: "Changes to our CI configuration files and scripts"},
}

func main() {
	trace.IAmMaster()
	log.PrettyLog()

	if len(os.Args) < 2 || os.Args[1] != "commit" {
		fmt.Println("Usage: cz commit")
		os.Exit(1)
	}

	commitMsg, err := generateCommitMessage()
	if err != nil {
		log.Error().Err(err).Msg("Generate commit message")
		os.Exit(1)
	}

	if err = git.Commit(commitMsg); err != nil {
		log.Error().Err(err).Msg("Commit")
		os.Exit(1)
	}

	log.Info().Msg("cz")
}

func generateCommitMessage() (string, error) {
	// 1. Select commit type
	typeSelect := promptui.Select{
		Label: "Select commit type",
		Items: commitTypes,
		Templates: &promptui.SelectTemplates{
			Label:    "{{ . }}?",
			Active:   "\U0001F449 {{ .Label | cyan }} - {{ .Description | cyan }}",
			Inactive: "  {{ .Label | white }} - {{ .Description | white }}",
			Selected: "\U0001F44D {{ .Label | green }} - {{ .Description | green }}",
		},
	}

	typeIdx, _, err := typeSelect.Run()
	if err != nil {
		return "", fmt.Errorf("prompt failed: %v", err)
	}

	// 2. Enter scope (optional)
	scopePrompt := promptui.Prompt{
		Label:     "Scope (optional)",
		Default:   "",
		AllowEdit: true,
	}

	scope, err := scopePrompt.Run()
	if err != nil {
		return "", fmt.Errorf("prompt failed: %v", err)
	}

	// 3. Enter commit message
	messagePrompt := promptui.Prompt{
		Label: "Commit message",
		Validate: func(input string) error {
			if len(input) < 3 {
				return fmt.Errorf("commit message must be at least 3 characters")
			}
			return nil
		},
	}

	message, err := messagePrompt.Run()
	if err != nil {
		return "", fmt.Errorf("prompt failed: %v", err)
	}

	// 4. Format commit message
	commitMsg := commitTypes[typeIdx].Label
	if scope != "" {
		commitMsg += "(" + scope + ")"
	}
	commitMsg += ": " + message

	// 5. Optional longer description
	bodyPrompt := promptui.Prompt{
		Label:     "Longer description (optional, press Enter to skip)",
		Default:   "",
		AllowEdit: true,
	}

	body, err := bodyPrompt.Run()
	if err != nil {
		return "", fmt.Errorf("prompt failed: %v", err)
	}

	if body != "" {
		commitMsg += "\n\n" + body
	}

	return commitMsg, nil
}
