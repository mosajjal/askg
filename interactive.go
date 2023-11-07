package main

import (
	"fmt"
	"os"
	"strings"
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
		{Text: "!write <file>", Description: "Write the conversation to a file"},
	}
	return prompt.FilterHasPrefix(s, d.Text, true)
}

func renderToMD(f *os.File, text string) {
	out := ""
	if f == os.Stdout {
		o, _ := os.Stdout.Stat()
		if (o.Mode()&os.ModeCharDevice) == os.ModeCharDevice && !nocolor { //Terminal
			out, _ = glamour.RenderWithEnvironmentConfig(text)
		} else {
			out = text
		}
	} else {
		out = text
	}
	fmt.Fprintln(f, out)
}

func sanitizeQuestion(question string) string {
	// url encode the question. BUG: this breaks UTF-8. for now returning the question
	return question
	// return url.QueryEscape(question)
}

func normalQ(bard *bard.Bard, question string) string {
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.UpdateCharSet(spinner.CharSets[11])
	s.Start()
	answer, err := bard.Ask(sanitizeQuestion(question))
	if err != nil {
		fmt.Println(err)
		answer = ""
	}
	s.Stop()
	return answer
}

// RunInteractive runs the interactive mode of the CLI
func RunInteractive(bard *bard.Bard) {
	fmt.Println("press <tab> to see the list of commands")
	outFile := os.Stdout
	for {
		text := prompt.Input("> ", completer)
		switch {
		case strings.HasPrefix(text, "!quit"):
			return
		case strings.HasPrefix(text, "!editor"):
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
			fmt.Fprintln(outFile, "Q: "+text)
			renderToMD(outFile, normalQ(bard, editorQ))
			continue
		case strings.HasPrefix(text, "!reset"):
			bard.Clear()
			continue
		case strings.TrimSpace(text) == "":
			continue
		case strings.HasPrefix(text, "!write"):
			// get the file name by spiliting the text after !write
			fname := strings.Split(text, "!write")
			if len(fname) < 2 {
				fmt.Println("Please provide a file name")
				continue
			}
			outFileName := strings.TrimSpace(fname[1])
			if outFileName == "" {
				fmt.Println("Please provide a file name")
				continue
			}
			var err error
			outFile, err = os.OpenFile(outFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				fmt.Println(err)
				continue
			}
		default:
			// Write the question back as well
			fmt.Fprintln(outFile, "Q: "+text)
			renderToMD(outFile, normalQ(bard, text))
			continue
		}

	}
}
