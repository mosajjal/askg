package main

import (
	"fmt"
	"net/url"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/briandowns/spinner"
	"github.com/c-bata/go-prompt"
	"github.com/charmbracelet/glamour"
	"github.com/mosajjal/bard-cli/bard"
)

func completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "!quit", Description: "Quit the application"},
		{Text: "!editor", Description: "use $EDITOR to write the question"},
		{Text: "!reset", Description: "Reset the bard conversation"},
	}
	return prompt.FilterHasPrefix(s, d.Text, true)
}

func renderToMD(text string) string {
	out, _ := glamour.RenderWithEnvironmentConfig(text)
	return out
}

func santizeQuestion(question string) string {
	// url encode the question
	return url.QueryEscape(question)
}

func normalQ(bard *bard.Bard, question string) string {
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.UpdateCharSet(spinner.CharSets[11])
	s.Start()
	answer, err := bard.Ask(santizeQuestion(question))
	if err != nil {
		fmt.Println(err)
		answer = ""
	}
	s.Stop()
	return answer
}

// RunInteractive runs the interactive mode of the CLI
func RunInteractive(bard *bard.Bard) {
	fmt.Println("Press Ctrl+D to exit")
	fmt.Println("press <tab> to see the list of commands")
	for {
		text := prompt.Input("> ", completer)

		switch text {
		case "!quit":
			return
		case "!editor":
			editorQ := ""
			prompt := &survey.Editor{
				Message:  "Question:",
				FileName: "*.md",
			}
			if err := survey.AskOne(prompt, &editorQ); err != nil {
				fmt.Println(err)
				continue
			}
			// trim the last newline
			editorQ = editorQ[:len(editorQ)-1]
			fmt.Println(renderToMD(normalQ(bard, editorQ)))
			continue
		case "!reset":
			bard.Clear()
			continue
		case "":
			continue
		default:
			fmt.Println(renderToMD(normalQ(bard, text)))
			continue
		}

	}
}
