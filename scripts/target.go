package scripts

import (
	"fmt"
	"strings"
)

func parseTarget(input string) (string, string, error) {
	var target string
	var targetName string
	if input == "" {
		target = "host"
	} else {
		targetParts := strings.Split(input, " ")
		target = targetParts[0]
		switch target {
		case "host":
			if len(targetParts) > 1 {
				return "", "", fmt.Errorf("target %q must not have a target name", target)
			}
		case "container", "service":
			if len(targetParts) == 0 {
				return "", "", fmt.Errorf("missing target name")
			}
			targetName = targetParts[1]
		default:
			return "", "", fmt.Errorf("invalid target %q", input)
		}
	}
	return target, targetName, nil
}
