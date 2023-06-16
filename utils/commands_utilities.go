package utils

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/ivantias/commits-assistant-cli/style"
)

func OutputCommand(command string) []string {
	cmdArgs := strings.Split(command, " ")
	output, err := exec.Command(cmdArgs[0], cmdArgs[1:]...).Output()

	if err != nil {
		log.Fatal(err)
	}

	strOutput := string(output)

	parsedOutput := strings.Split(strOutput, "\n")

	// remove the last element if it's empty
	if len(parsedOutput) > 0 && parsedOutput[len(parsedOutput)-1] == "" {
		parsedOutput = parsedOutput[:len(parsedOutput)-1]
	}

	return parsedOutput
}

func RunCommand(command string) {
	cmdArgs := strings.Split(command, " ")

	err := exec.Command(cmdArgs[0], cmdArgs[1:]...).Run()

	if err != nil {
		log.Fatal(err)
	}
}

func CheckGitRepoInitialized() bool {
	err := exec.Command("git", "rev-parse", "--is-inside-work-tree").Run()

	return err == nil
}

func ExecCommandLoop(scanner *bufio.Scanner, command string) {
	shouldContinueProcess := false

	commandResults := map[string]string{
		"git add .": "Files added to staging!",
		"git init":  "Git repo initialized!",
	}

	resultOutput, ok := commandResults[command]

	if !ok {
		log.Fatal("Invalid command.")
	}

	for !shouldContinueProcess {
		if !scanner.Scan() {
			break
		}

		answer := strings.TrimSpace(strings.ToLower(scanner.Text()))

		switch answer {
		case "y":
			shouldContinueProcess = true
			RunCommand(command)
			style.Green.Println(resultOutput)
		case "n":
			style.Red.Println("Exiting commits assistant...")
			os.Exit(0)
		default:
			style.Red.Print("Invalid input. Please enter Y or N: ")
			continue
		}
	}
}

func Commit(command string) {
	cmdArgs := []string{"-c", command}
	cmd := exec.Command("sh", cmdArgs...)

	output, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Println(string(output))
		panic(err)
	}
}
