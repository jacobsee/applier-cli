package cmd

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Provides commands which can be eval'd to run the current OpenShift-Applier inventory.",
	Long: `Provides commands which can be eval'd to run Ansible 
with the current OpenShift-Applier inventory. Default is 
to run using local Ansible, but can also provide the command 
for running OpenShift-Applier in a Docker container.`,
	Run: func(cmd *cobra.Command, args []string) {
		docker, _ := cmd.Flags().GetBool("docker")

		if docker {
			if checkSELinux() {
				fmt.Println(`docker run -u $(id -u) \
-v $HOME/.kube/config:/openshift-applier/.kube/config:z
-v $HOME/src/inventory/:/tmp/inventory
-e INVENTORY_PATH=/tmp/inventory
-t redhatcop/openshift-applier`)
			} else {
				fmt.Println(`docker run -u $(id -u) \
-v $HOME/.kube/config:/openshift-applier/.kube/config
-v $HOME/src/inventory/:/tmp/inventory
-e INVENTORY_PATH=/tmp/inventory
-t redhatcop/openshift-applier`)
			}
		} else {
			fmt.Println("ansible-playbook apply.yml -i inventory/")
		}
	},
}

func checkSELinux() bool {
	var out bytes.Buffer
	cmd := exec.Command("getenforce")
	cmd.Stdout = &out
	err := cmd.Run()
	if err == nil && strings.Trim(out.String(), "\n") == "Enforcing" {
		return true
	}
	return false
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().BoolP("docker", "d", false, "Run using the OpenShift-Applier Docker image instead of local Ansible.")
}
