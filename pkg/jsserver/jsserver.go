package jsserver

import (
	"errors"
	"fmt"
	"os/exec"
)

func CheckNpm() error {
	nodeCmd := exec.Command("cmd.exe", "/C", "node", "-v")

	_, err := nodeCmd.Output()
	if err != nil {
		errText := fmt.Sprintf("NodeJs is not installed: %s", err)
		return errors.New(errText)
	}

	npmCmd := exec.Command("cmd.exe", "/C", "npm", "-v")
	_, err = npmCmd.Output()
	if err != nil {
		errText := fmt.Sprintf("Npm is not installed: %s", err)
		return errors.New(errText)
	}

	return nil
}
