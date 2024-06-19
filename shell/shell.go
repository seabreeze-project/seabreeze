package shell

import (
	"fmt"

	prompt "github.com/c-bata/go-prompt"
	"github.com/mattn/go-shellwords"
	"github.com/spf13/cobra"
)

type CommandsMap map[string]*Command
type CommandDispatcher func(args []string) (bool, error)
type Completer func(CursorState) []prompt.Suggest

type CommandlineParser interface {
	Parse(line string) ([]string, error)
	ParseWithEnvs(line string) (envs []string, args []string, err error)
}

type Shell struct {
	title        string
	appVersion   string
	promptPrefix string
	// TODO: Implement HistoryLoader
	// historyLoader HistoryLoader
	Completer     Completer
	HelpGenerator HelpGenerator
	commands      CommandsMap
	aliases       CommandsMap
	dispatchers   []CommandDispatcher
	keyBinds      []KeyBind
}

func New(title string, commands ...*Command) *Shell {
	keyBinds := []KeyBind{
		{Key: prompt.ControlA, Desc: "Go to the beginning of the line (Home)"},
		{Key: prompt.ControlE, Desc: "Go to the end of the line (End)"},
		{Key: prompt.ControlP, Desc: "Previous command (Up arrow)"},
		{Key: prompt.ControlN, Desc: "Next command (Down arrow)"},
		{Key: prompt.ControlF, Desc: "Forward one character"},
		{Key: prompt.ControlB, Desc: "Backward one character"},
		{Key: prompt.ControlD, Desc: "Delete character under the cursor"},
		{Key: prompt.ControlH, Desc: "Delete character before the cursor (Backspace)"},
		{Key: prompt.ControlW, Desc: "Cut the word before the cursor to the clipboard"},
		{Key: prompt.ControlK, Desc: "Cut the line after the cursor to the clipboard"},
		{Key: prompt.ControlU, Desc: "Cut the line before the cursor to the clipboard"},
		{Key: prompt.ControlL, Desc: "Clear the screen"},
	}
	sh := &Shell{
		title:         title,
		commands:      make(CommandsMap),
		aliases:       make(CommandsMap),
		keyBinds:      keyBinds,
		HelpGenerator: &DefaultHelpGenerator{},
	}
	sh.AddCommand(commands...)
	return sh
}

func (sh *Shell) SetAppVersion(version string) *Shell {
	sh.appVersion = version
	return sh
}

func (sh *Shell) AddCommand(commands ...*Command) {
	for _, command := range commands {
		sh.checkRegisterCommand(command)
		sh.commands[command.Name] = command

		for _, alias := range command.Aliases {
			sh.checkRegisterAlias(alias, command)
			sh.aliases[alias] = command
		}
	}
}

func (sh *Shell) AddVirtualCommand(commands ...*Command) {
	for _, command := range commands {
		if command.Run != nil {
			panic(fmt.Sprintf("cannot register virtual shell command %q: virtual command must not have a Run function", command.Name))
		}
		sh.checkRegisterCommand(command)
		sh.commands[command.Name] = command

		for _, alias := range command.Aliases {
			sh.checkRegisterAlias(alias, command)
			sh.aliases[alias] = command
		}
	}
}

func (sh *Shell) AddExternalCommand(cobraCmds ...*cobra.Command) {
	for _, cobraCmd := range cobraCmds {
		command := &Command{
			Name:    cobraCmd.Name(),
			Desc:    cobraCmd.Short,
			Long:    cobraCmd.Long,
			Aliases: cobraCmd.Aliases,
			Run:     nil,
		}
		sh.AddVirtualCommand(command)
	}
}

func (sh *Shell) AddDispatcher(dispatcher CommandDispatcher) *Shell {
	sh.dispatchers = append(sh.dispatchers, dispatcher)
	return sh
}

func (sh *Shell) AddKeyBind(keyBinds ...KeyBind) *Shell {
	sh.keyBinds = append(sh.keyBinds, keyBinds...)
	return sh
}

func (sh *Shell) SetPromptPrefix(prefix string) *Shell {
	sh.promptPrefix = prefix
	return sh
}

func (sh *Shell) Help() string {
	help, err := sh.HelpGenerator.GenerateHelp(TemplateData{
		Title:    sh.title,
		Version:  sh.appVersion,
		Commands: sh.commands,
		KeyBinds: sh.keyBinds,
	})
	if err != nil {
		panic(err)
	}
	return help
}

func (sh *Shell) Suggestions() []prompt.Suggest {
	suggestions := []prompt.Suggest{}
	for _, command := range sh.commands {
		suggestions = append(suggestions, prompt.Suggest{
			Text:        command.Name,
			Description: command.Desc,
		})
	}
	for alias, command := range sh.aliases {
		suggestions = append(suggestions, prompt.Suggest{
			Text:        alias,
			Description: command.Desc + " (alias)",
		})
	}
	return suggestions
}

func (sh *Shell) Run() {
	s := &session{shell: sh}
	s.start()
}

func (sh *Shell) Dispatch(args []string) error {
	if len(args) == 0 {
		return nil
	}

	command := FindIn(args[0], sh.commands, sh.aliases)
	if command != nil {
		return command.Run(args[1:])
	}

	for _, dispatcher := range sh.dispatchers {
		ok, err := dispatcher(args)
		if ok || err != nil {
			return err
		}
	}

	return fmt.Errorf("unknown command %q", args[0])
}

func (sh *Shell) DispatchString(cmdline string) error {
	args, _ := shellwords.Parse(cmdline)
	return sh.Dispatch(args)
}

func FindIn(name string, registries ...CommandsMap) *Command {
	for _, registry := range registries {
		if command, ok := registry[name]; ok && command.Run != nil {
			return command
		}
	}
	return nil
}

func (sh *Shell) checkRegisterCommand(command *Command) {
	if command.Name == "" {
		panic("cannot register shell command without name")
	}
	if _, ok := sh.commands[command.Name]; ok {
		panic(fmt.Sprintf("cannot register shell command %q: command with same name already registered", command.Name))
	}
}

func (sh *Shell) checkRegisterAlias(alias string, command *Command) {
	if alias == "" {
		panic(fmt.Sprintf("cannot register shell alias with empty name for command %q", command.Name))
	}
	if _, ok := sh.commands[alias]; ok {
		panic(fmt.Sprintf("cannot register shell alias %q for command %q: command with same name already registered", alias, command.Name))
	}
	if _, ok := sh.aliases[alias]; ok {
		panic(fmt.Sprintf("cannot register shell alias %q for command %q: alias with same name already registered", alias, command.Name))
	}
}
