package shell

import (
	"fmt"
	"os"

	"github.com/c-bata/go-prompt"
)

type session struct {
	shell       *Shell
	prompter    *prompt.Prompt
	cursorState CursorState
	// TODO: Implement History
	// history     *History
}

func (s *session) start() {
	promptOptions := []prompt.Option{
		prompt.OptionTitle(s.shell.title),
		prompt.OptionPrefix(s.shell.promptPrefix),
		prompt.OptionPreviewSuggestionTextColor(prompt.Cyan),
		prompt.OptionSelectedSuggestionTextColor(prompt.White),
		prompt.OptionSelectedSuggestionBGColor(prompt.DarkBlue),
		prompt.OptionSelectedDescriptionBGColor(prompt.Blue),
		prompt.OptionSuggestionBGColor(prompt.DarkGray),
		prompt.OptionDescriptionBGColor(prompt.LightGray),
		prompt.OptionScrollbarBGColor(prompt.DarkGray),
		prompt.OptionScrollbarThumbColor(prompt.White),
	}

	for _, keyBind := range s.shell.keyBinds {
		if keyBind.Fn != nil {
			promptOptions = append(promptOptions, prompt.OptionAddKeyBind(prompt.KeyBind{
				Key: keyBind.Key,
				Fn:  keyBind.Fn,
			}))
		}
	}

	s.prompter = prompt.New(
		s.execute,
		s.complete,
		promptOptions...,
	)
	s.prompter.Run()
}

func (s *session) execute(cmdline string) {
	err := s.shell.DispatchString(cmdline)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func (s *session) complete(in prompt.Document) []prompt.Suggest {
	s.cursorState.Update(in)

	if in.TextBeforeCursor() == "" {
		return []prompt.Suggest{}
	}

	var suggest []prompt.Suggest
	if s.cursorState.PreviousWordsN() == 0 {
		suggest = s.shell.Suggestions()
	} else {
		suggest = s.shell.Completer(s.cursorState)
	}
	return prompt.FilterHasPrefix(suggest, in.GetWordBeforeCursor(), true)
}
