package installer

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"womm/internal/config"

	"github.com/charmbracelet/huh"
)

type PackageManager struct {
	Name       string
	Command    string
	InstallCmd string
	OS         []string
}

func GetAvailablePackageManagers() []PackageManager {
	osName := runtime.GOOS
	all := []PackageManager{
		{"Homebrew (brew)", "brew", "brew install %s", []string{"darwin", "linux"}},
		{"apt (Debian/Ubuntu)", "apt-get", "sudo apt-get install -y %s", []string{"linux"}},
		{"yum (RHEL/CentOS)", "yum", "sudo yum install -y %s", []string{"linux"}},
		{"Chocolatey (choco)", "choco", "choco install %s -y", []string{"windows"}},
		{"Scoop (scoop)", "scoop", "scoop install %s", []string{"windows"}},
		{"Windows Package Manager (winget)", "winget", "winget install --id %s --accept-package-agreements --accept-source-agreements", []string{"windows"}},
	}

	var available []PackageManager
	for _, pm := range all {
		if contains(pm.OS, osName) && commandExists(pm.Command) {
			available = append(available, pm)
		}
	}
	return available
}

func Install(tool config.Tool) {
	availablePMs := GetAvailablePackageManagers()
	if len(availablePMs) == 0 {
		fmt.Println("No supported package managers found.")
		return
	}

	selectedPM := availablePMs[0]
	if len(availablePMs) == 1 {
		fmt.Printf("Using detected package manager: %s\n", selectedPM.Name)
	} else {
		options := make([]huh.Option[int], len(availablePMs))
		for i, pm := range availablePMs {
			options[i] = huh.NewOption(pm.Name, i)
		}

		var selectedIndex int
		prompt := huh.NewSelect[int]().
			Title(fmt.Sprintf("Which package manager should I use for %s?", tool.Name)).
			Options(options...).
			Value(&selectedIndex)

		if err := prompt.Run(); err != nil {
			fmt.Println("Install cancelled.")
			return
		}
		selectedPM = availablePMs[selectedIndex]
	}

	packageName := tool.GetPackageName()
	if override, exists := tool.InstallOverrides[selectedPM.Command]; exists {
		packageName = override
	}

	installCmd := fmt.Sprintf(selectedPM.InstallCmd, packageName)
	executeInstallCommand(installCmd, selectedPM)
}

func executeInstallCommand(cmdStr string, pm PackageManager) {
	fmt.Printf("▶ Installing via %s...\n", pm.Name)

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", cmdStr)
	} else {
		cmd = exec.Command("sh", "-c", cmdStr)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("❌ Failed to install via %s\n", pm.Name)
		return
	}

	fmt.Printf("✅ Successfully installed via %s!\n", pm.Name)
}

func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
