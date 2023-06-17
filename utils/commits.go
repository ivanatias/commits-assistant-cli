package utils

import (
	"fmt"
	"strings"
)

type CommitType struct {
	Name        string
	Description string
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

func FormatCommitOptions(commitTypes []CommitType) []string {
	var options = make([]string, len(commitTypes))
	var longestPrefixLength int

	for i := range commitTypes {
		prefixLength := len(commitTypes[i].Name)
		if prefixLength > longestPrefixLength {
			longestPrefixLength = prefixLength
		}
	}

	for i, commitType := range commitTypes {
		var descPadding = longestPrefixLength - len(commitType.Name)
		descPad := strings.Repeat(" ", descPadding)
		options[i] = fmt.Sprintf("%s: %s%s", commitType.Name, descPad, commitType.Description)
	}

	return options
}
