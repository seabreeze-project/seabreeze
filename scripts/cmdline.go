package scripts

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/seabreeze-project/seabreeze/util"
	"mvdan.cc/sh/v3/expand"
	"mvdan.cc/sh/v3/syntax"
)

func generateCmdline(command string, env map[string]string, args ...string) ([]string, error) {
	p := syntax.NewParser()
	var words []*syntax.Word
	err := p.Words(strings.NewReader(command), func(w *syntax.Word) bool {
		words = append(words, w)
		return true
	})
	if err != nil {
		return nil, err
	}

	fullEnv := map[string]string{}
	if len(env) > 0 {
		for k, v := range env {
			// disallow overriding @, *, #, and all numbered variables
			if util.StringIsNumber(k) || k == "@" || k == "*" || k == "#" {
				return nil, fmt.Errorf("cannot override $%s in environment", k)
			}
			fullEnv[k] = v
		}
	}

	getenv := func(s string) string {
		switch s {
		case "#":
			return fmt.Sprintf("%d", len(args))
		case "*":
			return strings.Join(args, " ")
		case "IFS":
			return " \t\n"
		}
		if util.StringIsNumber(s) {
			i, err := strconv.Atoi(s)
			if i < 1 || err != nil {
				// command name
				return ""
			}
			v := args[i-1]
			return v
		}
		if v, ok := fullEnv[s]; ok {
			return v
		}
		return "$" + s
	}
	cfg := &expand.Config{
		Env: expand.FuncEnviron(getenv),
	}

	var processed []string
	fields, err := expand.Fields(cfg, words...)
	if err != nil {
		return nil, err
	}

	for _, f := range fields {
		if f == "$@" {
			processed = append(processed, args...)
		} else {
			processed = append(processed, f)
		}
	}

	if len(processed) == 0 {
		return nil, fmt.Errorf("script has no effective command")
	}

	fmt.Println("processed:", strings.Join(processed, " "))
	return processed, nil
}
