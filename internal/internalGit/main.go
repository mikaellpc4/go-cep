package internalGit

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func GitClone(repo string, location string) error {
	return GitCloneWithDepth(repo, location, 0)
}

func GitCloneWithDepth(repo string, location string, depth int) error {
	args := []string{"clone", repo, location}

	if _, err := os.Stat(location); err != nil {
		_, err := issueCommand("mkdir", []string{location})
		return err
	}

	if depth > 0 {
		args = append(args, strconv.Itoa(depth))
	}

	_, err := issueCommand("git", args)

	if err != nil {
		return err
	}

	return nil
}

func GitPull(repo string, location string) error {
	args := []string{"clone", repo, location}

	_, err := issueCommand("git", args)

	if err != nil {
		return err
	}

	return nil
}

func GitLog(location string) error {
	args := []string{"-C", location, "log"}

	data, err := issueCommand("git", args)

	fmt.Println(data)

	if err != nil {
		return err
	}

	return nil
}

func issueCommand(command string, args []string) ([]string, error) {
	cmd := exec.Command(command, args...)

	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(out), "\n")
	return lines, nil
}
