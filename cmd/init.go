package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	githubapi "github.com/jacobsee/applier-gen/pkg/github_api"
	yamlresources "github.com/jacobsee/applier-gen/pkg/yaml_resources"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes an empty OpenShift-Applier inventory",
	Long: `Scaffolds an empty OpenShift-Applier inventory, including:

An inventory directory, with host_vars and group_vars
A templates directory
A params directory
A files directory
An Ansible Galaxy requirements file with the latest release of OpenShift-Applier on GitHub
	
In addition, the Ansible Galaxy requirements are installed if the ansible-galaxy bin is available.`,
	Run: func(cmd *cobra.Command, args []string) {
		makeAllDirectories()
		latestReleasedVersion := getLatestApplierReleaseTag()
		writeGalaxyRequirements(latestReleasedVersion)
		installGalaxyRequirements()
	},
}

func makeAllDirectories() {
	os.Mkdir("inventory", 0766)
	os.Mkdir("inventory/host_vars", 0766)
	os.Mkdir("inventory/group_vars", 0766)
	os.Mkdir("templates", 0766)
	os.Mkdir("params", 0766)
	os.Mkdir("files", 0766)
}

func writeConfigs() {
	hostVars := []byte("ansible_connection: local")
	err := ioutil.WriteFile("inventory/host_vars/localhost.yml", hostVars, 0766)
	if err != nil {
		log.Fatal("Could not write inventory/host_vars/localhost.yml")
	}
		
}

func getLatestApplierReleaseTag() string {
	currentApplierVersion := githubapi.GetLatestVersionInfo()
	return currentApplierVersion.TagName
}

func writeGalaxyRequirements(version string) {
	requirements := &yamlresources.Requirements{{
		Src:     "https://github.com/redhat-cop/openshift-applier",
		SCM:     "git",
		Version: version,
		Name:    "openshift-applier",
	}}
	yamlRequirements, err := yaml.Marshal(requirements)
	if err != nil {
		log.Fatal("Could not generate requirements file.")
	}
	err = ioutil.WriteFile("requirements.yml", yamlRequirements, 0766)
}

func installGalaxyRequirements() {
	galaxy := exec.Command("ansible-galaxy", "install", "-r", "requirements.yml", "--roles-path=roles", "--f")
	galaxy.Dir, _ = os.Getwd()
	_, err := galaxy.Output()
	fmt.Println("Initialized empty openshift-applier inventory")
	if err != nil {
		fmt.Println("Could not invoke ansible-galaxy to install requirements. Please check your installation and install ansible-galaxy requirements manually.")
	} else {
		fmt.Println("Successfully installed ansible-galaxy requirements")
	}
}

func init() {
	rootCmd.AddCommand(initCmd)
}
