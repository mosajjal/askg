package main

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/charmbracelet/glamour"
	"github.com/manifoldco/promptui"
	"github.com/mosajjal/bard-cli/bard"
)

// RunInteractive runs the interactive mode of the CLI
func RunInteractive(bard *bard.Bard) {
	fmt.Println("Press Ctrl+D to exit")
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.UpdateCharSet(spinner.CharSets[11])
	for {
		q := promptui.Prompt{
			Label: "Q",
		}
		text, err := q.Run()
		if err != nil {
			return
		}
		s.Start()
		answer, err := bard.Ask(text)
		if err != nil {
			fmt.Println(err)
			return
		}
		s.Stop()
		out, _ := glamour.RenderWithEnvironmentConfig(answer)
		fmt.Println(out)
	}
}
