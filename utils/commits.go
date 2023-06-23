package utils

import (
	"fmt"
	"strings"
)

type CommitType struct {
	Name        string
	Description string
}

type BreakingChangeOption struct {
	Answer         string
	BreakingChange bool
}

var BreakingChangeOptions = []BreakingChangeOption{
	{"No", false},
	{"Yes", true},
}

var CommitTypes = []CommitType{
	{
		"📦 feat",
		"A new feature",
	},
	{
		"🐛 fix",
		"A bug fix",
	},
	{
		"📃 docs",
		"Documentation only changes",
	},
	{
		"💅 style",
		"Changes that do not affect the meaning of the code",
	},
	{
		"🔧 refactor",
		"A code change that neither fixes a bug nor adds a feature",
	},
	{
		"🚀 perf",
		"A code change that improves performance",
	},
	{
		"🧪 test",
		"Adding missing tests",
	},
	{
		"👀 chore",
		"Changes to the build process or auxiliary tools",
	},
	{
		"👈 revert",
		"Reverts a previous commit",
	},
}

func FormatCommitOptions() []string {
	var options = make([]string, len(CommitTypes))
	var longestPrefixLength int

	for i := range CommitTypes {
		prefixLength := len(CommitTypes[i].Name)
		if prefixLength > longestPrefixLength {
			longestPrefixLength = prefixLength
		}
	}

	for i, commitType := range CommitTypes {
		var descPadding = longestPrefixLength - len(commitType.Name)
		descPad := strings.Repeat(" ", descPadding)
		options[i] = fmt.Sprintf("%s: %s%s",
			commitType.Name,
			descPad,
			commitType.Description,
		)
	}

	return options
}

func FormatBreakingChangeOptions() []string {
	var options = make([]string, len(BreakingChangeOptions))

	for i, option := range BreakingChangeOptions {
		options[i] = option.Answer
	}

	return options
}
