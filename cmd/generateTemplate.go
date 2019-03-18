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
	"log"
	"os/exec"

	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

// generateTemplateCmd represents the generateTemplate command
var generateTemplateCmd = &cobra.Command{
	Use:   "generateTemplate",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		fromCluster, _ := cmd.Flags().GetBool("from-cluster")
		if fromCluster {
			cmd := exec.Command("oc", "get", args[0], "--export", "-o", "yaml")
			// cmd.Stdin = strings.NewReader("some input")
			var out bytes.Buffer
			cmd.Stdout = &out
			err := cmd.Run()
			if err != nil {
				log.Fatal("Could not export the desired resource from the cluster. Are you sure that it exists?")
			}
			var anyYaml map[string]interface{}
			yaml.Unmarshal(out.Bytes(), &anyYaml)
			fmt.Printf("kind: %s\n", anyYaml["kind"].(string))
			fmt.Printf("replicas: %d\n", anyYaml["spec"].(map[interface{}]interface{})["replicas"].(int))
			var outByte []byte
			outByte, _ = yaml.Marshal(anyYaml)
			fmt.Printf("%s\n", outByte)
			//fmt.Printf("%s\n", out.String())
		}

	},
}

func init() {
	rootCmd.AddCommand(generateTemplateCmd)

	generateTemplateCmd.Flags().BoolP("from-cluster", "c", false, "Use an existing cluster as a source for this resource")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generateTemplateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateTemplateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
