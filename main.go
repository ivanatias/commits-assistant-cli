package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ivantias/commits-assistant-cli/utils"
	"github.com/manifoldco/promptui"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	isRepoInitialized := utils.CheckGitRepoInitialized()

	if !isRepoInitialized {
		fmt.Println("Git repo not initialized.")
		fmt.Print("Want to initialize git repo? (Y/N): ")

		utils.ExecCommandLoop(scanner, "git init")
	}

	stagedFiles := utils.OutputCommand("git diff --cached --name-only")
	modifiedFiles := utils.OutputCommand("git status --porcelain")

	if len(stagedFiles) == 0 && len(modifiedFiles) > 0 {
		fmt.Println("There are no staged files for committing.")
		fmt.Println("List of non-staged modified files:")

		var filesToOutput string

		for _, file := range modifiedFiles {
			filesToOutput += file + "\n"
		}

		fmt.Println(filesToOutput)
		fmt.Print("Do you want to add all files to staging? (Y/N): ")

		utils.ExecCommandLoop(scanner, "git add .")
	}

	if len(modifiedFiles) == 0 && len(stagedFiles) == 0 {
		fmt.Println("There are no changes to commit.")
		fmt.Println("Exiting commits assistant...")
		os.Exit(0)
	}

	selectPrompt := promptui.Select{
		Label: "Select type of commit",
		Items: utils.FormatCommitOptions(utils.CommitTypes),
	}

	selectedIndex, _, err := selectPrompt.Run()

	commitPrefix := utils.CommitTypes[selectedIndex].Name

	if err != nil {
		log.Fatal(err)
	}

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

	message, _ := messagePrompt.Run()

	fullCommitMsg := fmt.Sprintf("%s: %s", commitPrefix, message)
	commitCommand := fmt.Sprintf(`git commit -m "%s"`, fullCommitMsg)

	fmt.Printf("Committing ---> %s\n", strings.TrimSpace(fullCommitMsg))

	utils.Commit(commitCommand)

	fmt.Println("Successful commit. Thanks for using the assistant!")
}
