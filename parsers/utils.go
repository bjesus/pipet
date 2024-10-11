package parsers

import (
	"os/exec"
	"strings"
)

func ExecutePipe(input string, command string) (string, error) {
	cmd := exec.Command("bash", "-c", command)
	cmd.Stdin = strings.NewReader(input)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func CalculateIndentation(s string) int {
	count := 0
	for _, char := range s {
		if char != ' ' {
			break
		}
		count++
	}
	return count
}

func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}
