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
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"

	//templatev1 "github.com/openshift/api/template/v1"
	yamlresources "github.com/jacobsee/applier-gen/pkg/yaml_resources"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
	//"k8s.io/apimachinery/pkg/runtime"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		fromCluster, _ := cmd.Flags().GetBool("from-cluster")
		fromFile, _ := cmd.Flags().GetBool("from-file")
		makeTemplate, _ := cmd.Flags().GetBool("make-template")

		var resourceYaml map[string]interface{}
		var out bytes.Buffer

		if fromCluster {

			cmd := exec.Command("oc", "get", args[0], "--export", "-o", "yaml")

			cmd.Stdout = &out
			err := cmd.Run()
			if err != nil {
				log.Fatal("Could not export the desired resource from the cluster. Are you sure that it exists?")
			}
			err = yaml.Unmarshal(out.Bytes(), &resourceYaml)
			if err != nil {
				log.Fatal("Unable to interpret the resource.")
			}

		} else if fromFile {

			fileContents, err := ioutil.ReadFile(args[0])
			if err != nil {
				log.Fatal("Could not read the specified file.")
			}
			yaml.Unmarshal(fileContents, &resourceYaml)
			if err != nil {
				log.Fatal("Unable to interpret the file as valid YAML.")
			}

		} else {

			log.Fatal("Unclear where to get the resource. Please use --from-cluster (-c) or --from-file (-f).")
			return

		}

		resourceKind := resourceYaml["kind"].(string)

		if makeTemplate || resourceKind == "Template" {

			outByte := []byte{}

			name := resourceYaml["metadata"].(map[interface{}]interface{})["name"].(string)

			if resourceKind != "Template" {
				template := yamlresources.Template{
					APIVersion: "v1",
					Kind:       "Template",
					Metadata: yamlresources.TemplateMetadata{
						Name: name,
					},
					Objects: []map[string]interface{}{
						resourceYaml,
					},
					Parameters: nil,
				}
				outByte, _ = yaml.Marshal(template)
			} else {
				outByte, _ = yaml.Marshal(resourceYaml)
			}

			err := ioutil.WriteFile(fmt.Sprintf("templates/%s.yml", name), outByte, 0766)
			if err != nil {
				log.Fatal(err)
				fmt.Println("Could not add template to the current inventory.")
			} else {
				fmt.Println("Template added to the current inventory.")
			}

		} else {

			outByte, _ := yaml.Marshal(resourceYaml)

			name := resourceYaml["metadata"].(map[interface{}]interface{})["name"].(string)

			err := ioutil.WriteFile(fmt.Sprintf("files/%s.yml", name), outByte, 0766)
			if err != nil {
				log.Fatal(err)
				fmt.Println("Could not add file to the current inventory.")
			} else {
				fmt.Println("File added to the current inventory.")
			}

		}

	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().BoolP("from-cluster", "c", false, "Use an existing cluster as a source for this resource")
	addCmd.Flags().BoolP("from-file", "f", false, "Use a yaml file as a source for this resource")
	addCmd.Flags().BoolP("make-template", "t", false, "Convert the resource into a template (if it isn't already)")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
