package prompts

import (
	"errors"

	"github.com/manifoldco/promptui"
)

var (
	CommitTypePrompt = promptui.Select{
		Label: "Select commit type",
		Items: commitTypeSelection,
	}

	DescriptionPrompt = promptui.Prompt{
		Label: "Enter commit message",
		Validate: func(input string) error {
			if len(input) == 0 {
				return errors.New("commit message cannot be empty")
			} else if len(input) > 50 {
				return errors.New("commit message cannot be longer than 50 characters")
			}
			return nil
		},
	}

	IsBreakingChangePrompt = promptui.Select{
		Label: "Does this commit break backwards compatibility?",
		Items: breakingChangeSelection,
	}

	BodyPrompt = promptui.Prompt{
		Label: "Enter commit body (optional)",
	}

	FooterPrompt = promptui.Prompt{
		Label: "Enter breaking change description",
		Validate: func(input string) error {
			if len(input) == 0 {
				return errors.New("breaking change description cannot be empty")
			}
			return nil
		},
	}
)
