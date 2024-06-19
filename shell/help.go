package shell

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/c-bata/go-prompt"
)

type TemplateData struct {
	Title    string
	Version  string
	Commands CommandsMap
	KeyBinds []KeyBind
}

type HelpGenerator interface {
	GenerateHelp(data TemplateData) (string, error)
}

type DefaultHelpGenerator struct {
}

func (hg *DefaultHelpGenerator) GenerateHelp(data TemplateData) (string, error) {
	templateSource := `{{.Title}}{{if .Version}} {{.Version}}{{end}}

Available Commands:{{range .Commands}}
  {{.Name | pad 10}}  {{if .Aliases}}({{.Aliases | join ", "}}) {{end}}{{.Desc}}{{end}}

Keybindings:{{range .KeyBinds}}
  {{.Key | friendlyKey | pad 10}}  {{.Desc}}{{end}}`

	funcs := hg.GetFunctions(data)
	t := template.Must(template.New("help").Funcs(funcs).Parse(templateSource))

	var sb strings.Builder
	err := t.Execute(&sb, data)
	if err != nil {
		return "", fmt.Errorf("error executing template: %w", err)
	}

	return sb.String(), nil
}

func (hg *DefaultHelpGenerator) GetFunctions(data TemplateData) template.FuncMap {
	return template.FuncMap{
		"join": func(sep string, elems []string) string {
			return strings.Join(elems, sep)
		},
		"pad": func(maxLength int, str string) string {
			return fmt.Sprintf(fmt.Sprintf("%%-%ds", maxLength), str)
		},
		// "pad2": func(str string) string {
		// 	maxLength := 0
		// 	for k, _ := range data.Commands {
		// 		if len(k) > maxLength {
		// 			maxLength = len(k)
		// 		}
		// 	}
		// 	return fmt.Sprintf(fmt.Sprintf("%%-%ds", maxLength), str)
		// },
		// "padded": func(cmds CommandsMap) CommandsMap {
		// 	maxKeyLength := 0
		// 	for k, _ := range cmds {
		// 		if len(k) > maxKeyLength {
		// 			maxKeyLength = len(k)
		// 		}
		// 	}
		// 	for k, v := range cmds {
		// 		cmds[k].Name = fmt.Sprintf(fmt.Sprintf("%%-%ds", maxKeyLength), v.Name)
		// 	}
		// 	return cmds
		// },
		"friendlyKey": func(key prompt.Key) string {
			if s, ok := KeyNames[key]; ok {
				return s
			}
			return fmt.Sprintf("%v", key)
		},
	}
}
