// Package launch is the registry of supported AI agents and how to invoke
// them with an initial prompt.
package launch

import (
	"fmt"
	"os/exec"
)

// Agent describes a launchable AI agent.
type Agent struct {
	Slug        string
	DisplayName string
	// Binary is the executable name probed on PATH.
	Binary string
	// BuildCmd produces the argv to run, given a prompt to inject.
	BuildCmd func(prompt string) []string
}

// Registry is the canonical list of agents algotutor can auto-launch.
var Registry = []Agent{
	{
		Slug:        "claude",
		DisplayName: "Claude Code",
		Binary:      "claude",
		BuildCmd: func(prompt string) []string {
			return []string{"claude", "--dangerously-skip-permissions", prompt}
		},
	},
	{
		Slug:        "codex",
		DisplayName: "OpenAI Codex CLI",
		Binary:      "codex",
		BuildCmd: func(prompt string) []string {
			return []string{"codex", prompt}
		},
	},
	{
		Slug:        "opencode",
		DisplayName: "OpenCode",
		Binary:      "opencode",
		BuildCmd: func(prompt string) []string {
			return []string{"opencode", prompt}
		},
	},
	{
		Slug:        "gemini",
		DisplayName: "Gemini CLI",
		Binary:      "gemini",
		BuildCmd: func(prompt string) []string {
			return []string{"gemini", prompt}
		},
	},
}

// DetectAvailable returns the agents whose binary is on PATH.
func DetectAvailable() []Agent {
	out := make([]Agent, 0, len(Registry))
	for _, a := range Registry {
		if _, err := exec.LookPath(a.Binary); err == nil {
			out = append(out, a)
		}
	}
	return out
}

// Lookup returns the agent for slug, or false if unknown.
func Lookup(slug string) (Agent, bool) {
	for _, a := range Registry {
		if a.Slug == slug {
			return a, true
		}
	}
	return Agent{}, false
}

// Command builds the argv to invoke agent slug with prompt as the initial
// message.
func Command(slug, prompt string) ([]string, error) {
	a, ok := Lookup(slug)
	if !ok {
		return nil, fmt.Errorf("unknown agent: %s", slug)
	}
	if _, err := exec.LookPath(a.Binary); err != nil {
		return nil, fmt.Errorf("%s not on PATH", a.Binary)
	}
	return a.BuildCmd(prompt), nil
}
