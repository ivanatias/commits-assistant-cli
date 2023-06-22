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

	selectPrompt := promptui.Select{
		Label: "Select type of commit",
		Items: utils.FormatCommitOptions(utils.CommitTypes),
	}

	selectedIndex, _, err := selectPrompt.Run()

	if err != nil {
		log.Fatal(err)
	}

	commitPrefix := utils.CommitTypes[selectedIndex].Name

	messagePrompt := promptui.Prompt{
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

	message, err := messagePrompt.Run()

	if err != nil {
		log.Fatal(err)
	}

	descPrompt := promptui.Prompt{
		Label: "Enter commit description (optional)",
	}

	desc, err := descPrompt.Run()

	if err != nil {
		log.Fatal(err)
	}

	var commitCommand string
	var commitOutput string

	if len(desc) > 0 {
		commitCommand = fmt.Sprintf(`git commit -m "%s: %s" -m "%s"`, commitPrefix, message, desc)
		commitOutput = fmt.Sprintf("%s: %s\n%s", commitPrefix, message, desc)
	} else {
		commitCommand = fmt.Sprintf(`git commit -m "%s: %s"`, commitPrefix, message)
		commitOutput = fmt.Sprintf("%s: %s", commitPrefix, message)
	}

	style.Cyan.Println("Committing 👇")
	fmt.Println(commitOutput + "\n")

	utils.Commit(commitCommand)

	style.Green.Println("Successful commit. Thanks for using the assistant!")
}
