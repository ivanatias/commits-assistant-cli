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
		utils.Yellow.Println("⚠️ There are no staged files for committing.")
		utils.Yellow.Println("List of non-staged modified files:")

		var filesToOutput string

		for _, file := range modifiedFiles {
			filesToOutput += file + "\n"
		}

		utils.Cyan.Println(filesToOutput)
		utils.Yellow.Print("Do you want to add all files to staging? (Y/N): ")

		utils.ExecCommandLoop(scanner, "git add .")
	}

	if len(modifiedFiles) == 0 && len(stagedFiles) == 0 {
		utils.Yellow.Println("There are no changes to commit.")
		utils.Red.Println("Exiting commits assistant...")
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

	utils.Cyan.Printf("Committing ---> %s\n", strings.TrimSpace(fullCommitMsg))

	utils.Commit(commitCommand)

	utils.Green.Println("Successful commit. Thanks for using the assistant!")
}
