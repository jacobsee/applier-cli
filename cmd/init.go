// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	apireponses "github.com/jacobsee/applier-gen/pkg/api_responses"
	yamlresources "github.com/jacobsee/applier-gen/pkg/yaml_resources"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		currentApplierVersion := apireponses.GetLatestVersionInfo()
		requirements := &yamlresources.Requirements{{
			Src:     "https://github.com/redhat-cop/openshift-applier",
			SCM:     "git",
			Version: currentApplierVersion.TagName,
			Name:    "openshift-applier",
		}}
		yamlRequirements, err := yaml.Marshal(requirements)
		if err != nil {
			//oops
		}
		err = ioutil.WriteFile("requirements.yml", yamlRequirements, 0766)
		os.Mkdir("inventory", 0766)
		os.Mkdir("inventory/host_vars", 0766)
		os.Mkdir("inventory/group_vars", 0766)
		os.Mkdir("templates", 0766)
		os.Mkdir("params", 0766)
		os.Mkdir("files", 0766)
		galaxy := exec.Command("ansible-galaxy", "install", "-r", "requirements.yml", "--roles-path=roles", "--f")
		galaxy.Dir, _ = os.Getwd()
		_, err = galaxy.Output()
		fmt.Println("Initialized empty openshift-applier inventory")
		if err != nil {
			fmt.Println("Could not invoke ansible-galaxy to install requirements. Please check your installation and install ansible-galaxy reqiurements manually.")
		} else {
			fmt.Println("Successfully installed ansible-galaxy requirements")
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
