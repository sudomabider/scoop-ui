package scoop

import (
	"os/exec"
	"regexp"
	"strings"
)

// App represents a scoop app
type App struct {
	Name    string
	Version string
	Bucket  string
}

// List returns a list of apps
func List() ([]App, error) {
	cmd := exec.Command("scoop", "list")

	output, err := cmd.Output()
	if err != nil {
		return make([]App, 0), err
	}

	lines := strings.Split(string(output), "\n")

	reg := regexp.MustCompile(`^\ +([^\ ]+)\ ([^\ ]+)(\ \[([^\ ]+)\])?`)

	apps := make([]App, 0)

	for _, line := range lines {
		line = strings.TrimRight(line, " ")
		if !reg.MatchString(line) {
			continue
		}

		matches := reg.FindSubmatch([]byte(line))
		apps = append(apps, App{
			Name:    string(matches[1]),
			Version: string(matches[2]),
			Bucket:  string(matches[4]),
		})
	}

	return apps, nil
}
