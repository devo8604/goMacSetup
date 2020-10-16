package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
)

// Status enum
const (
	STS_FILE_DOWNLOAD_ERR = 1
	STS_ERR_PERM_SET_ERR
	STS_XCODE_TOOLS_INSTALL_ERR
	STS_BREW_INSTALL_ERR
)

// Downloads a file from a given url to a given filepath
func downloadFile(url string, filepath string) error {

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}

	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// Installs my env on MacOS
func main() {

	var pythonPackages []string
	var homeBrewPackages []string
	var scripts []string
	var err error
	var cmd *exec.Cmd

	pythonPackages = []string{
		"install", // This takes care of the install arg for us
		"matplotlib",
		"scikit-learn",
		"tensorflow",
	}

	homeBrewPackages = []string{
		"install", // This takes care of the install arg for us
		"python3.9",
		"ansible",
		"redis",
		"wget",
		"youtube-dl",
		"1password",
		"devdocs",
		"hyper",
		"parallels",
		"steam",
		"stellarium",
		"visual-studio-code",
		"wireshark",
		"openra",
		"google-chrome",
		"yarn",
		"node",
		"go",
		"git",
	}

	// Download the Homebrew install script
	url := "https://raw.githubusercontent.com/Homebrew/install/master/install.sh"
	filename := "./install.sh"
	err = downloadFile(url, filename)
	if err != nil {
		os.Exit(STS_BREW_INSTALL_ERR)
	}

	// Make sure to delete the file before exiting
	defer os.Remove(filename)

	// chmod it...
	err = os.Chmod(filename, 0777)
	if err != nil {
		os.Exit(STS_ERR_PERM_SET_ERR)
	}

	bashPath, bashErr := exec.LookPath("bash")
	scripts = []string{
		bashPath,
		filename,
	}

	if bashErr == nil {

		// Install xcode command line tools
		xcodePath, xcodeErr := exec.LookPath("xcode-select")
		if xcodeErr != nil {
			os.Exit(STS_XCODE_TOOLS_INSTALL_ERR)
		}
		xcodeTools := []string{
			"--install",
		}

		cmd := exec.Command(xcodePath, xcodeTools...)
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		err = cmd.Run()
		if err != nil {
			fmt.Println("Unable to install XCode Commandline tools:", err)
			os.Exit(STS_XCODE_TOOLS_INSTALL_ERR)
		}

		scriptsCmd := &exec.Cmd{
			Path:   bashPath,
			Args:   scripts,
			Stdin:  os.Stdin,
			Stderr: os.Stderr,
		}

		fmt.Println("Installing HomeBrew...")
		err = scriptsCmd.Run()
		if err != nil {
			fmt.Printf("\nError running script: %s\n", err)
			os.Exit(STS_BREW_INSTALL_ERR)
		}
	}

	brewBinaryPath, brewErr := exec.LookPath("brew")

	if brewErr == nil {

		fmt.Println("Installing HomeBrew Packages...")
		// Install the homebrew packages next:
		cmd = exec.Command(brewBinaryPath, homeBrewPackages...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			fmt.Printf("\nError installing HomeBrew packages: %s\n", err)
		}
	}

	pip3Path, pipErr := exec.LookPath("pip3")

	if pipErr == nil {
		fmt.Println("Installing Python Packages...")
		// Install the python packagers I want
		cmd = exec.Command(pip3Path, pythonPackages...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			fmt.Printf("\nError installing Python packages: %s\n", err)
		}
	}
}
