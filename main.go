package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/ivantias/commits-assistant-cli/prompts"
	"github.com/ivantias/commits-assistant-cli/style"
	"github.com/ivantias/commits-assistant-cli/utils"
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

	selectedIndex, _, err := prompts.CommitTypePrompt.Run()

	if err != nil {
		log.Fatal(err)
	}

	commitPrefix := utils.CommitTypes[selectedIndex].Name

	description, err := prompts.DescriptionPrompt.Run()

	if err != nil {
		log.Fatal(err)
	}

	body, err := prompts.BodyPrompt.Run()

	if err != nil {
		log.Fatal(err)
	}

	selectedIndex, _, err = prompts.IsBreakingChangePrompt.Run()

	if err != nil {
		log.Fatal(err)
	}

	isBreakingChange := utils.BreakingChangeOptions[selectedIndex].BreakingChange

	var footer, commitCommand, commitOutput string

	if isBreakingChange {
		footer, err = prompts.FooterPrompt.Run()

		if err != nil {
			log.Fatal(err)
		}

		footer = fmt.Sprintf("\nBREAKING CHANGE: %s", footer)
	}

	if len(body) > 0 {
		commitCommand = fmt.Sprintf(
			`git commit -m "%s: %s" -m "%s" -m "%s"`,
			commitPrefix,
			description,
			body,
			footer,
		)
		commitOutput = fmt.Sprintf("%s: %s\n%s%s",
			commitPrefix,
			description,
			body,
			footer,
		)
	} else {
		commitCommand = fmt.Sprintf(`git commit -m "%s: %s" -m "%s"`,
			commitPrefix,
			description,
			footer,
		)
		commitOutput = fmt.Sprintf("%s: %s%s", commitPrefix, description, footer)
	}

	style.Cyan.Println("Committing 👇")
	fmt.Println(commitOutput + "\n")

	utils.Commit(commitCommand)

	if isBreakingChange {
		style.Yellow.Println("IMPORTANT: This commit should trigger a major release.")
	}

	style.Green.Println("Successful commit. Thanks for using the assistant!")
}
