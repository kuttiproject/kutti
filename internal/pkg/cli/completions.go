package cli

import (
	"strings"

	"github.com/spf13/cobra"
)

// StringCompletions filters a slice of strings to those members that start with
// the specified string.
func StringCompletions(possibilities []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if toComplete == "" {
		return possibilities, cobra.ShellCompDirectiveNoFileComp
	}

	result := make([]string, 0, len(possibilities))
	for _, value := range possibilities {
		if strings.HasPrefix(value, toComplete) {
			result = append(result, value)
		}
	}

	return result, cobra.ShellCompDirectiveNoFileComp | cobra.ShellCompDirectiveNoSpace
}
