// Copyright (c) 2020 Sorint.lab S.p.A.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package repo

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/1set/gut/yos"
	"github.com/ercole-io/ercole/utils"
)

// ArtifactInfo contains info about all files in repository
type ArtifactInfo struct {
	Repository            string
	Installed             bool `json:"-"`
	Version               string
	ReleaseDate           string
	Filename              string
	Name                  string
	OperatingSystemFamily string
	OperatingSystem       string
	Arch                  string
	UpstreamType          string
	UpstreamInfo          map[string]interface{}
	Install               func(ai *ArtifactInfo)              `json:"-"`
	Uninstall             func(ai *ArtifactInfo)              `json:"-"`
	Download              func(ai *ArtifactInfo, dest string) `json:"-"`
}

//Regex for filenames
var agentRHEL5Regex *regexp.Regexp = regexp.MustCompile("^ercole-agent-(?P<version>.*)-1.(?P<arch>x86_64).rpm$")
var agentRHELRegex *regexp.Regexp = regexp.MustCompile("^ercole-agent-(?P<version>.*)-1.el(?P<dist>\\d+).(?P<arch>x86_64).rpm$")
var agentVirtualizationRHELRegex *regexp.Regexp = regexp.MustCompile("^ercole-agent-virtualization-(?P<version>.*)-1.el(?P<dist>\\d+).(?P<arch>x86_64).rpm$")
var agentExadataRHELRegex *regexp.Regexp = regexp.MustCompile("^ercole-agent-exadata-(?P<version>.*)-1.el(?P<dist>\\d+).(?P<arch>x86_64).rpm$")
var agentWinRegex *regexp.Regexp = regexp.MustCompile("^ercole-agent-setup-(?P<version>.*).exe$")
var agentHpuxRegex *regexp.Regexp = regexp.MustCompile("^ercole-agent-hpux-(?P<version>.*).tar.gz")
var agentAixRegexRpm *regexp.Regexp = regexp.MustCompile("^ercole-agent-aix-(?P<version>.*)-1.(?P<dist>.*).(?P<arch>noarch).rpm$")
var agentAixRegexTarGz *regexp.Regexp = regexp.MustCompile("^ercole-agent-aix-(?P<version>.*).tar.gz$")
var ercoleRHELRegex *regexp.Regexp = regexp.MustCompile("^ercole-(?P<version>.*)-1.el(?P<dist>\\d+).(?P<arch>x86_64).rpm$")
var ercoleWebRHELRegex *regexp.Regexp = regexp.MustCompile("^ercole-web-(?P<version>.*)-1.el(?P<dist>\\d+).(?P<arch>noarch).rpm$")

// GetFullName return the fullname of the file
func (artifact *ArtifactInfo) GetFullName() string {
	return fmt.Sprintf("%s/%s@%s", artifact.Repository, artifact.Name, artifact.Version)
}

// IsInstalled return true if file is detected in the distribution directory
func (artifact *ArtifactInfo) IsInstalled(distributedFiles string) bool {
	_, err := os.Stat(filepath.Join(distributedFiles, "all", artifact.Filename))

	return !os.IsNotExist(err)
}

// SetInfoFromFileName sets to fileInfo informations taken from filename
func (artifact *ArtifactInfo) SetInfoFromFileName(filename string) {
	switch {
	case agentVirtualizationRHELRegex.MatchString(filename): //agent virtualization RHEL
		data := utils.FindNamedMatches(agentVirtualizationRHELRegex, filename)
		artifact.Name = "ercole-agent-virtualization-rhel" + data["dist"]
		artifact.Version = data["version"]
		artifact.Arch = data["arch"]
		artifact.OperatingSystemFamily = "rhel"
		artifact.OperatingSystem = "rhel" + data["dist"]
	case agentExadataRHELRegex.MatchString(filename): //agent exadata RHEL
		data := utils.FindNamedMatches(agentExadataRHELRegex, filename)
		artifact.Name = "ercole-agent-exadata-rhel" + data["dist"]
		artifact.Version = data["version"]
		artifact.Arch = data["arch"]
		artifact.OperatingSystemFamily = "rhel"
		artifact.OperatingSystem = "rhel" + data["dist"]
	case agentRHEL5Regex.MatchString(filename): //agent RHEL5
		data := utils.FindNamedMatches(agentRHEL5Regex, filename)
		artifact.Name = "ercole-agent-rhel5"
		artifact.Version = data["version"]
		artifact.Arch = data["arch"]
		artifact.OperatingSystemFamily = "rhel"
		artifact.OperatingSystem = "rhel5"
	case agentRHELRegex.MatchString(filename): //agent RHEL
		data := utils.FindNamedMatches(agentRHELRegex, filename)
		artifact.Name = "ercole-agent-rhel" + data["dist"]
		artifact.Version = data["version"]
		artifact.Arch = data["arch"]
		artifact.OperatingSystemFamily = "rhel"
		artifact.OperatingSystem = "rhel" + data["dist"]
	case ercoleRHELRegex.MatchString(filename): //ercole RHEL
		data := utils.FindNamedMatches(ercoleRHELRegex, filename)
		artifact.Name = "ercole-" + data["dist"]
		artifact.Version = data["version"]
		artifact.Arch = data["arch"]
		artifact.OperatingSystemFamily = "rhel"
		artifact.OperatingSystem = "rhel" + data["dist"]
	case ercoleWebRHELRegex.MatchString(filename): //ercole-web RHEL
		data := utils.FindNamedMatches(ercoleWebRHELRegex, filename)
		artifact.Name = "ercole-web" + data["dist"]
		artifact.Version = data["version"]
		artifact.Arch = data["arch"]
		artifact.OperatingSystemFamily = "rhel"
		artifact.OperatingSystem = "rhel" + data["dist"]
	case agentWinRegex.MatchString(filename): //agent WIN
		data := utils.FindNamedMatches(agentWinRegex, filename)
		artifact.Name = "ercole-agent-win"
		artifact.Version = data["version"]
		artifact.Arch = "x86_64"
		artifact.OperatingSystemFamily = "win"
		artifact.OperatingSystem = "win"
	case agentHpuxRegex.MatchString(filename): //agent HPUX
		data := utils.FindNamedMatches(agentHpuxRegex, filename)
		artifact.Name = "ercole-agent-hpux"
		artifact.Version = data["version"]
		artifact.Arch = "noarch"
		artifact.OperatingSystemFamily = "hpux"
		artifact.OperatingSystem = "hpux"
	case agentAixRegexRpm.MatchString(filename): //agent AIX
		data := utils.FindNamedMatches(agentAixRegexRpm, filename)
		artifact.Name = "ercole-agent-aix"
		artifact.Version = data["version"]
		artifact.Arch = "noarch"
		artifact.OperatingSystemFamily = "aix"
		artifact.OperatingSystem = data["dist"]
	case agentAixRegexTarGz.MatchString(filename): //agent AIX
		data := utils.FindNamedMatches(agentAixRegexTarGz, filename)
		artifact.Name = "ercole-agent-aix-targz"
		artifact.Version = data["version"]
		artifact.Arch = "noarch"
		artifact.OperatingSystemFamily = "aix-tar-gz"
		artifact.OperatingSystem = "aix6.1"
	default:
		panic(fmt.Errorf("Filename %s is not supported. Please check that is correct", filename))
	}
}

// SetDownloader set the downloader of the artifact
func (artifact *ArtifactInfo) SetDownloader(verbose bool) {
	switch artifact.UpstreamType {
	case "github-release":
		artifact.Download = func(ai *ArtifactInfo, dest string) {
			utils.DownloadFile(dest, ai.UpstreamInfo["DownloadUrl"].(string))
		}
	case "directory":
		artifact.Download = func(ai *ArtifactInfo, dest string) {
			if verbose {
				fmt.Printf("Copying file from %s to %s\n", ai.UpstreamInfo["Filename"].(string), dest)
			}
			err := yos.CopyFile(ai.UpstreamInfo["Filename"].(string), dest)
			if err != nil {
				panic(err)
			}
			err = os.Chmod(dest, 0755)
			if err != nil {
				panic(err)
			}
		}
	case "ercole-reposervice":
		artifact.Download = func(ai *ArtifactInfo, dest string) {
			utils.DownloadFile(dest, ai.UpstreamInfo["DownloadUrl"].(string))
		}
	default:
		panic(artifact)
	}
}

// SetInstaller set the installer of the artifact
func (artifact *ArtifactInfo) SetInstaller(verbose bool, distributedFiles string) {
	switch {
	case strings.HasSuffix(artifact.Filename, ".rpm"):
		artifact.Install = func(ai *ArtifactInfo) {
			//Create missing directories
			if verbose {
				fmt.Printf("Creating the directories (if missing) %s, %s\n",
					filepath.Join(distributedFiles, ai.OperatingSystemFamily, ai.OperatingSystem, ai.Arch),
					filepath.Join(distributedFiles, "all"),
				)
			}
			err := os.MkdirAll(filepath.Join(distributedFiles, ai.OperatingSystemFamily, ai.OperatingSystem, ai.Arch), 0755)
			if err != nil {
				panic(err)
			}
			err = os.MkdirAll(filepath.Join(distributedFiles, "all"), 0755)
			if err != nil {
				panic(err)
			}

			//Download the file in the right location
			if verbose {
				fmt.Printf("Downloading the artifact %s to %s\n", ai.Filename, filepath.Join(distributedFiles, ai.OperatingSystemFamily, ai.OperatingSystem, ai.Arch, ai.Filename))
			}
			ai.Download(ai, filepath.Join(distributedFiles, ai.OperatingSystemFamily, ai.OperatingSystem, ai.Arch, ai.Filename))

			//Create a link to all
			if verbose {
				fmt.Printf("Linking the artifact to %s\n", filepath.Join(distributedFiles, "all", ai.Filename))
			}
			err = os.Link(filepath.Join(distributedFiles, ai.OperatingSystemFamily, ai.OperatingSystem, ai.Arch, ai.Filename), filepath.Join(distributedFiles, "all", ai.Filename))
			if err != nil {
				panic(err)
			}

			//Launch the createrepo command
			if verbose {
				fmt.Printf("Executing createrepo %s\n", filepath.Join(distributedFiles, ai.OperatingSystemFamily, ai.OperatingSystem, ai.Arch))
			}
			cmd := exec.Command("createrepo", filepath.Join(distributedFiles, ai.OperatingSystemFamily, ai.OperatingSystem, ai.Arch))
			if verbose {
				cmd.Stdout = os.Stdout
			}
			cmd.Stderr = os.Stderr
			cmd.Run()

			//Settint it to installed
			ai.Installed = true
		}
	default:
		artifact.Install = func(ai *ArtifactInfo) {
			//Create missing directories
			if verbose {
				fmt.Printf("Creating the directories (if missing) %s, %s\n",
					filepath.Join(distributedFiles, ai.OperatingSystemFamily, ai.OperatingSystem, ai.Arch),
					filepath.Join(distributedFiles, "all"),
				)
			}
			err := os.MkdirAll(filepath.Join(distributedFiles, ai.OperatingSystemFamily, ai.OperatingSystem, ai.Arch), 0755)
			if err != nil {
				panic(err)
			}
			err = os.MkdirAll(filepath.Join(distributedFiles, "all"), 0755)
			if err != nil {
				panic(err)
			}

			//Download the file in the right location
			if verbose {
				fmt.Printf("Downloading the artifact %s to %s\n", ai.Filename, filepath.Join(distributedFiles, ai.OperatingSystemFamily, ai.OperatingSystem, ai.Arch, ai.Filename))
			}
			ai.Download(ai, filepath.Join(distributedFiles, ai.OperatingSystemFamily, "/", ai.OperatingSystem, ai.Arch, ai.Filename))

			//Create a link to all
			if verbose {
				fmt.Printf("Linking the artifact to %s\n", filepath.Join(distributedFiles, "all", ai.Filename))
			}
			err = os.Link(filepath.Join(distributedFiles, ai.OperatingSystemFamily, "/", ai.OperatingSystem, ai.Arch, ai.Filename), filepath.Join(distributedFiles, "all", ai.Filename))
			if err != nil {
				panic(err)
			}

			//Setting it to installed
			ai.Installed = true
		}
	}
}

// SetUninstaller set the uninstaller of the artifact
func (artifact *ArtifactInfo) SetUninstaller(verbose bool, distributedFiles string) {
	switch {
	case strings.HasSuffix(artifact.Filename, ".rpm"):
		artifact.Uninstall = func(ai *ArtifactInfo) {
			//Removing the link to all
			if verbose {
				fmt.Printf("Removing Linking the artifact to %s\n", filepath.Join(distributedFiles, "all", ai.Filename))
			}
			err := os.Remove(filepath.Join(distributedFiles, "all", ai.Filename))
			if err != nil {
				panic(err)
			}

			//Removing the file
			if verbose {
				fmt.Printf("Removing the file %s\n", filepath.Join(distributedFiles, ai.OperatingSystemFamily, ai.OperatingSystem, ai.Arch, ai.Filename))
			}
			err = os.Remove(filepath.Join(distributedFiles, ai.OperatingSystemFamily, ai.OperatingSystem, ai.Arch, ai.Filename))
			if err != nil {
				panic(err)
			}

			//Launch the createrepo command
			if verbose {
				fmt.Printf("Executing createrepo %s\n", filepath.Join(distributedFiles, ai.OperatingSystemFamily, ai.OperatingSystem, ai.Arch))
			}
			cmd := exec.Command("createrepo", filepath.Join(distributedFiles, ai.OperatingSystemFamily, ai.OperatingSystem, ai.Arch))
			if verbose {
				cmd.Stdout = os.Stdout
			}
			cmd.Stderr = os.Stderr
			err = cmd.Run()
			if err != nil {
				panic(err)
			}

			//Set it to not installed
			ai.Installed = false
		}
	default:
		artifact.Uninstall = func(ai *ArtifactInfo) {
			//Removing the link to all
			if verbose {
				fmt.Printf("Removing Linking the artifact to %s\n", filepath.Join(distributedFiles, "all", ai.Filename))
			}
			err := os.Remove(filepath.Join(distributedFiles, "all", ai.Filename))
			if err != nil {
				panic(err)
			}

			//Removing the file
			if verbose {
				fmt.Printf("Removing the file %s\n", filepath.Join(distributedFiles, ai.OperatingSystemFamily, ai.OperatingSystem, ai.Arch, ai.Filename))
			}
			err = os.Remove(filepath.Join(distributedFiles, ai.OperatingSystemFamily, ai.OperatingSystem, ai.Arch, ai.Filename))
			if err != nil {
				panic(err)
			}

			//Set it to not installed
			ai.Installed = false
		}
	}
}