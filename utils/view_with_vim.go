package utils

import (
	"errors"
	"strings"
	"os"
	"os/exec"
)

// ViewWithVim previews text as markdown using Vim installed on a machine
func ViewWithVim(text string, readonly bool) error {

	_, err := exec.LookPath("vim")
	if err != nil {
		return err
	}

	if len(text) == 0 {
		return errors.New("error: preview string is empty")
	}

	input := strings.NewReader(text)

	args := []string{
		"-",
		"-c", "set filetype=markdown",
	}

	if readonly {
		args = append(args, "-R")
	}

	cmd := exec.Command("vim", args...)
	cmd.Stdin = input
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
