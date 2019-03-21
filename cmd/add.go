package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	yamlresources "github.com/jacobsee/applier-gen/pkg/yaml_resources"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	yaml "gopkg.in/yaml.v2"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add new resources to the current inventory",
	Long: `Add new resources to the current inventory. Resources that
are templates are added to the templates directory, while all
others are added to the files directory. Non-template resources
can also be converted into templates. Generated resources can immediately
be opened in your default editor for further tuning.`,
	Run: func(cmd *cobra.Command, args []string) {

		// Get all flags
		fromCluster, _ := cmd.Flags().GetBool("from-cluster")
		fromFile, _ := cmd.Flags().GetBool("from-file")
		makeTemplate, _ := cmd.Flags().GetBool("make-template")
		edit, _ := cmd.Flags().GetBool("edit")

		var resourceYaml map[string]interface{}

		if fromCluster {
			resourceYaml = getResourceFromCluster(args[0])
		} else if fromFile {
			resourceYaml = getResourceFromFile(args[0])
		} else {
			log.Fatal("Unclear where to get the resource. Please use --from-cluster (-c) or --from-file (-f).")
		}

		resourceKind := resourceYaml["kind"].(string)
		name := resourceYaml["metadata"].(map[interface{}]interface{})["name"].(string)
		var fileCreated string

		if makeTemplate || resourceKind == "Template" {
			fileCreated = writeTemplate(resourceYaml, resourceKind, name)
		} else {
			fileCreated = writeFile(resourceYaml, resourceKind, name)
		}

		if edit {
			cmd := exec.Command(viper.GetString("editor"), fileCreated)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Run()
		}

	},
}

func getResourceFromCluster(resourceName string) map[string]interface{} {

	var out bytes.Buffer
	var resource map[string]interface{}

	cmd := exec.Command("oc", "get", resourceName, "--export", "-o", "yaml")
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal("Could not export the desired resource from the cluster. Are you sure that you are logged in and it exists?")
	}
	err = yaml.Unmarshal(out.Bytes(), &resource)
	if err != nil {
		log.Fatal("Unable to interpret the resource.")
	}

	return resource

}

func getResourceFromFile(fileName string) map[string]interface{} {

	var resource map[string]interface{}

	fileContents, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal("Could not read the specified file.")
	}
	yaml.Unmarshal(fileContents, &resource)
	if err != nil {
		log.Fatal("Unable to interpret the file as valid YAML.")
	}

	return resource

}

func writeTemplate(resource map[string]interface{}, kind string, name string) string {

	outByte := []byte{}

	if kind != "Template" {
		template := yamlresources.Template{
			APIVersion: "v1",
			Kind:       "Template",
			Metadata: yamlresources.TemplateMetadata{
				Name: name,
			},
			Objects: []map[string]interface{}{
				resource,
			},
			Parameters: nil,
		}
		outByte, _ = yaml.Marshal(template)
	} else {
		outByte, _ = yaml.Marshal(resource)
	}

	outFile := fmt.Sprintf("templates/%s.yml", name)
	err := ioutil.WriteFile(outFile, outByte, 0766)
	if err != nil {
		log.Fatal(err)
		fmt.Println("Could not add template to the current inventory.")
	} else {
		fmt.Println("Template added to the current inventory.")
	}

	return outFile

}

func writeFile(resource map[string]interface{}, kind string, name string) string {

	outByte, _ := yaml.Marshal(resource)

	outFile := fmt.Sprintf("files/%s.yml", name)
	err := ioutil.WriteFile(outFile, outByte, 0766)
	if err != nil {
		log.Fatal(err)
		fmt.Println("Could not add file to the current inventory.")
	} else {
		fmt.Println("File added to the current inventory.")
	}

	return outFile

}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().BoolP("from-cluster", "c", false, "Use an existing cluster as a source for this resource")
	addCmd.Flags().BoolP("from-file", "f", false, "Use a yaml file as a source for this resource")
	addCmd.Flags().BoolP("make-template", "t", false, "Convert the resource into a template (if it isn't already)")
	addCmd.Flags().BoolP("edit", "e", false, "Immediately open the file in your default editor once created")
}
