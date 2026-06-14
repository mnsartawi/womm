package checker

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"

	"womm/internal/config"
	"womm/internal/installer"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

type DepResult struct {
	Tool      config.Tool
	Installed bool
	Version   string
	MeetsReq  bool
}

type EnvResult struct {
	Name   string
	Set    bool
	Value  string
	Expect string
}

type FileResult struct {
	Path   string
	Exists bool
}

var platformNames = map[string]string{
	"windows": "windows",
	"linux":   "linux",
	"macos":   "darwin",
}

// ── Dependencies ──

func CheckDeps(tools map[string]config.Tool) []DepResult {
	var results []DepResult
	for _, t := range tools {
		results = append(results, checkTool(t))
	}
	return results
}

func checkTool(t config.Tool) DepResult {
	cmd := shellCmd(t.GetCheckCmd())
	out, err := cmd.Output()
	if err != nil {
		return DepResult{Tool: t, Installed: false}
	}

	re := regexp.MustCompile(t.GetVersionRegex())
	matches := re.FindStringSubmatch(string(out))
	if len(matches) < 1 {
		return DepResult{Tool: t, Installed: false}
	}

	version := matches[0]
	meetsReq := strings.Contains(version, t.Version) || version >= strings.TrimPrefix(t.Version, ">=")

	return DepResult{Tool: t, Installed: true, Version: version, MeetsReq: meetsReq}
}

func PrintDepResults(results []DepResult) {
	for _, r := range results {
		fmt.Println(depStatus(r))
	}
}

func depStatus(r DepResult) string {
	if !r.Installed {
		return fmt.Sprintf("❌ %s: not found (requires %s)", r.Tool.Name, r.Tool.Version)
	}
	if r.MeetsReq {
		return fmt.Sprintf("✓ %s: %s", r.Tool.Name, r.Version)
	}
	return fmt.Sprintf("⚠️  %s: %s (requires %s)", r.Tool.Name, r.Version, r.Tool.Version)
}

func GetMissing(results []DepResult) []DepResult {
	var missing []DepResult
	for _, r := range results {
		if !r.Installed || !r.MeetsReq {
			missing = append(missing, r)
		}
	}
	return missing
}

func PromptInstall(missing []DepResult) {
	for _, m := range missing {
		var confirm bool
		prompt := huh.NewConfirm().
			Title(fmt.Sprintf("%s is missing. Would you like to install it?", m.Tool.Name)).
			Affirmative("Yes").
			Negative("No").
			Value(&confirm).
			WithTheme(blueTheme())

		_ = prompt.Run()

		if confirm {
			installer.Install(m.Tool)
		}
	}
}

// ── Env Vars ──

func CheckEnv(env map[string]string) []EnvResult {
	var results []EnvResult
	for name, expect := range env {
		val, set := os.LookupEnv(name)
		results = append(results, EnvResult{
			Name:   name,
			Set:    set && val != "",
			Value:  val,
			Expect: expect,
		})
	}
	return results
}

func PrintEnvResults(results []EnvResult) {
	if len(results) == 0 {
		return
	}
	fmt.Println()
	for _, r := range results {
		if r.Set {
			fmt.Printf("  ✓ %s\n", r.Name)
		} else {
			fmt.Printf("  ❌ %s (expected: %s)\n", r.Name, r.Expect)
		}
	}
}

// ── Files ──

func CheckFiles(files []string) []FileResult {
	var results []FileResult
	for _, f := range files {
		_, err := os.Stat(f)
		results = append(results, FileResult{Path: f, Exists: err == nil})
	}
	return results
}

func PrintFileResults(results []FileResult) {
	if len(results) == 0 {
		return
	}
	fmt.Println()
	for _, r := range results {
		if r.Exists {
			fmt.Printf("  ✓ %s\n", r.Path)
		} else {
			fmt.Printf("  ❌ %s not found\n", r.Path)
		}
	}
}

// ── Platforms ──

func CheckPlatforms(platforms map[string]bool) bool {
	current := runtime.GOOS
	for name, goos := range platformNames {
		if goos == current {
			if supported, ok := platforms[name]; ok && supported {
				return true
			}
			return false
		}
	}

	for name, supported := range platforms {
		expectedGOOS, ok := platformNames[name]
		if !ok || !supported {
			continue
		}
		if expectedGOOS == current {
			return true
		}
	}
	return false
}

func PrintPlatformResult(supported bool) {
	if supported {
		fmt.Printf("✓ Platform: %s\n", runtime.GOOS)
	} else {
		fmt.Printf("❌ Platform %s is not supported by this project\n", runtime.GOOS)
	}
}

// ── Commands ──

func RunCommands(commands map[string][]string) {
	if len(commands) == 0 {
		return
	}
	fmt.Println()
	for name, steps := range commands {
		var confirm bool
		prompt := huh.NewConfirm().
			Title(fmt.Sprintf("Run %q commands?", name)).
			Affirmative("Yes").
			Negative("Skip").
			Value(&confirm).
			WithTheme(blueTheme())

		_ = prompt.Run()
		if !confirm {
			continue
		}

		for _, step := range steps {
			fmt.Printf("  ▶ %s\n", step)
			cmd := shellCmd(step)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				fmt.Printf("  ❌ command failed: %s\n", step)
				break
			}
			fmt.Printf("  ✓ %s\n", step)
		}
	}
}

// ── Helpers ──

func shellCmd(command string) *exec.Cmd {
	if runtime.GOOS == "windows" {
		return exec.Command("cmd", "/C", command)
	}
	return exec.Command("sh", "-c", command)
}

func blueTheme() *huh.Theme {
	theme := huh.ThemeCharm()
	theme.Focused.FocusedButton = theme.Focused.FocusedButton.
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#2563EB")).
		Bold(true)
	return theme
}
