package ui

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"kcc-cli/m/v2/resources"
	"strings"
)

// SwitchContext Switch the current-context based on the selection made
func SwitchContext(config *resources.KubeConfig, contextName string) (string, error) {
	var selectedContextName = contextName
	for _, context := range config.Contexts {
		if context.Name == selectedContextName {
			config.CurrentContext = selectedContextName
			return selectedContextName, nil
		}
	}
	return "", fmt.Errorf("context %s not found", contextName)
}

func CussorPositionPointer(config *resources.KubeConfig) (int, []resources.Context) {
	cursorPosition := -1
	contexts := config.Contexts
	currentContext := config.CurrentContext
	if currentContext != " " {
		for i, context := range contexts {
			if currentContext == context.Name {
				cursorPosition = i
			}
		}
	}
	return cursorPosition, contexts
}

// ShowSelector render selector
func ShowSelector(options []resources.Context, currentPos int) (string, error) {

	modifiedOptions := make([]resources.Context, len(options))
	copy(modifiedOptions, options)
	modifiedOptions[currentPos].Name = options[currentPos].Name + " (*)"

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}? {{ `/ to search` | faint }}",
		Active:   ">    {{ .Name | cyan | bold }}",
		Inactive: "     {{ .Name | white}}",
		Selected: "     {{ .Name | cyan }}",
		Details: `{{ "CONTEXT:" | green | bold }}	{{ .Name | white  }}
{{ "CLUSTER:" | green | bold  }}	{{ .Context.cluster | white }}
{{ "AUTH INFO:" | green | bold  }}	{{ .Context.user | white }}
`,
	}

	// Search contexts in the selector
	searcher := func(input string, index int) bool {
		option := options[index]
		context := strings.Replace(strings.ToLower(option.Name), " ", "", -1)

		input = strings.Replace(strings.ToLower(input), " ", "", -1)
		return strings.Contains(context, input)
	}
	prompt := promptui.Select{
		Label:        "Select Kubernetes cluster context",
		Items:        modifiedOptions,
		Templates:    templates,
		Size:         5,
		Searcher:     searcher,
		CursorPos:    currentPos,
		HideSelected: true,
		HideHelp:     true,
	}

	i, _, err := prompt.RunCursorAt(currentPos, currentPos-3)
	if err != nil {
		return "", err
	}

	return options[i].Name, nil
}
