package gitinfo

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"
)

// GitInfo contains info about the current directory's git checkout
type GitInfo struct {
	RevisionCount  string // How many "commits" are in the log
	ProjectVersion string // Latest tag name
	Version        string // Synthetic tag name + revision count version number
	Date           string // Date of last commit (in UTC)
	Repository     string // Repository location
	Hash           string // Hash of the current commit/checkout
}

// GetGitInfo will gather all the git related information
// and return a single object containing the details.
func GetGitInfo() *GitInfo {
	gi := &GitInfo{}
	gi.RevisionCount = GitRevisionCount()
	gi.ProjectVersion = GitProjectVersion()
	gi.Version = GitVersion()
	gi.Date = GitDate()
	gi.Repository = GitRepository()
	gi.Hash = GitHash()
	return gi
}

// GitRevisionCount determines the current revision count.
func GitRevisionCount() string {
	cmd := exec.Command("git", "log", "--oneline")
	b, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("running %#v %#v: %v", cmd.Path, cmd.Args, err)
	}
	s := strings.TrimSpace(string(b))
	lines := strings.Split(s, "\n")
	return fmt.Sprintf("%v", len(lines))
}

// GitHash finds the current git commit hash.
func GitHash() string {
	cmd := exec.Command("git", "log", "--oneline", "-1")
	b, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("running %#v %#v: %v", cmd.Path, cmd.Args, err)
	}
	parts := strings.Split(string(b), " ")
	return parts[0]

}

// GitProjectVersion gets the latest git tag.
func GitProjectVersion() string {
	cmd := exec.Command("git", "describe", "--tags", "--long")
	b, err := cmd.CombinedOutput()
	if err != nil {
		return "x.notags"
		//log.Fatalf("running %#v %#v: %v", cmd.Path, cmd.Args, err)
	}
	s := strings.TrimSpace(string(b))
	return s
}

// GitVersion combines GitProjectVersion with GitRevisionCount
func GitVersion() string {
	s := GitProjectVersion()
	parts := strings.Split(s, "-")
	version := fmt.Sprintf("%v.%v", parts[0], GitRevisionCount())
	return version
}

// GitDate gets the latest git commit date
func GitDate() string {
	cmd := exec.Command("env", "TZ=UTC", "git", "log", "-1", `--format=%cd`)
	b, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("running %#v %#v: %v", cmd.Path, cmd.Args, err)
	}
	s := strings.TrimSpace(string(b))
	return s
}

// GitRepository reports the current repo name
// (useful when people fork the project)
func GitRepository() string {
	cmd := exec.Command("git", "remote", "-v")
	b, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("running %#v %#v: %v", cmd.Path, cmd.Args, err)
	}
	lines := strings.Split(string(b), "\n")
	re := regexp.MustCompile(`(\S+)\s+\(fetch\)$`)
	for _, line := range lines {
		m := re.FindString(line)
		if len(m) > 0 {
			return m
		}
	}
	return "unparseable"
}
