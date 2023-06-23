package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/ivantias/commits-assistant-cli/style"
	"github.com/ivantias/commits-assistant-cli/utils"
	"github.com/manifoldco/promptui"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	isRepoInitialized := utils.CheckGitRepoInitialized()

	if !isRepoInitialized {
		style.Yellow.Println("⚠️ Git repo not initialized.")
		style.Yellow.Print("Want to initialize git repo? (Y/N): ")

		utils.ExecCommandLoop(scanner, "git init")
	}

	stagedFiles := utils.OutputCommand("git diff --cached --name-only")
	modifiedFiles := utils.OutputCommand("git status --porcelain")

	if len(stagedFiles) == 0 && len(modifiedFiles) > 0 {
		style.Yellow.Println("⚠️ There are no staged files for committing.")
		style.Yellow.Println("List of non-staged modified files:")

		var filesToOutput string

		for _, file := range modifiedFiles {
			filesToOutput += file + "\n"
		}

		style.Cyan.Println(filesToOutput)
		style.Cyan.Print("Do you want to add all files to staging? (Y/N): ")

		utils.ExecCommandLoop(scanner, "git add .")
	}

	if len(modifiedFiles) == 0 && len(stagedFiles) == 0 {
		style.Red.Println("⚠️ There are no changes to commit.")
		style.Red.Println("Exiting commits assistant...")
		os.Exit(0)
	}

	commitTypePrompt := promptui.Select{
		Label: "Select commit type",
		Items: utils.FormatCommitOptions(utils.CommitTypes),
	}

	selectedIndex, _, err := commitTypePrompt.Run()

	if err != nil {
		log.Fatal(err)
	}

	commitPrefix := utils.CommitTypes[selectedIndex].Name

	descriptionPrompt := promptui.Prompt{
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

	description, err := descriptionPrompt.Run()

	if err != nil {
		log.Fatal(err)
	}

	bodyPrompt := promptui.Prompt{
		Label: "Enter commit body (optional)",
	}

	body, err := bodyPrompt.Run()

	if err != nil {
		log.Fatal(err)
	}

	isBreakingChangePrompt := promptui.Select{
		Label: "Does this commit break backwards compatibility?",
		Items: []string{"No", "Yes"},
	}

	_, isBreakingChange, err := isBreakingChangePrompt.Run()

	if err != nil {
		log.Fatal(err)
	}

	var footer string
	var commitCommand string
	var commitOutput string

	if isBreakingChange == "Yes" {
		footerPrompt := promptui.Prompt{
			Label: "Enter breaking change description",
			Validate: func(input string) error {
				if len(input) == 0 {
					return errors.New("breaking change description cannot be empty")
				}
				return nil
			},
		}

		footer, err = footerPrompt.Run()

		if err != nil {
			log.Fatal(err)
		}

		footer = fmt.Sprintf("\nBREAKING CHANGE: %s", footer)
	}

	if len(body) > 0 {
		commitCommand = fmt.Sprintf(`git commit -m "%s: %s" -m "%s" -m "%s"`, commitPrefix, description, body, footer)
		commitOutput = fmt.Sprintf("%s: %s\n%s%s", commitPrefix, description, body, footer)
	} else {
		commitCommand = fmt.Sprintf(`git commit -m "%s: %s" -m "%s"`, commitPrefix, description, footer)
		commitOutput = fmt.Sprintf("%s: %s%s", commitPrefix, description, footer)
	}

	style.Cyan.Println("Committing 👇")
	fmt.Println(commitOutput + "\n")

	utils.Commit(commitCommand)

	if isBreakingChange == "Yes" {
		style.Yellow.Println("IMPORTANT: This commit should trigger a major release.")
	}

	style.Green.Println("Successful commit. Thanks for using the assistant!")
}
