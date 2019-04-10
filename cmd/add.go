package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	clusterinterface "github.com/jacobsee/applier-cli/pkg/cluster_interface"
	fileinterface "github.com/jacobsee/applier-cli/pkg/file_interface"
	yamlresources "github.com/jacobsee/applier-cli/pkg/yaml_resources"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	yaml "gopkg.in/yaml.v2"
)

type runFlags struct {
	fromCluster  bool
	fromFile     bool
	makeTemplate bool
	edit         bool
}

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

		flags := runFlags{
			fromCluster:  fromCluster,
			fromFile:     fromFile,
			makeTemplate: makeTemplate,
			edit:         edit,
		}

		var clusterInterface *clusterinterface.OCClusterInterface
		var fileInterface *fileinterface.FileSystemInterface

		add(flags, args, clusterInterface, fileInterface)
	},
}

func add(flags runFlags, args []string, clusterInterface clusterinterface.ClusterInterface, fileInterface fileinterface.FileInterface) {
	var resource map[string]interface{}
	var err error

	if flags.fromCluster {
		resource, err = clusterInterface.GetResource(args[0])
	} else if flags.fromFile {
		resource, err = fileInterface.ReadResource(args[0])
	} else {
		log.Fatal("Unclear where to get the resource. Please use --from-cluster (-c) or --from-file (-f).")
	}
	if err != nil {
		log.Fatal("Unable to get the resource requested.")
	}

	resourceKind := resource["kind"].(string)
	name := resource["metadata"].(map[interface{}]interface{})["name"].(string)
	var fileCreated string

	filterUndesireableKeys(resource)

	if flags.makeTemplate || resourceKind == "Template" {
		fileCreated = writeTemplateToInventory(resource, resourceKind, name, fileInterface)
	} else {
		fileCreated = writeFileToInventory(resource, resourceKind, name, fileInterface)
	}

	if flags.edit {
		cmd := exec.Command(viper.GetString("editor"), fileCreated)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func filterUndesireableKeys(resource map[string]interface{}) {
	delete(resource["metadata"].(map[interface{}]interface{})["annotations"].(map[interface{}]interface{}), "kubectl.kubernetes.io/last-applied-configuration")
	delete(resource["metadata"].(map[interface{}]interface{})["annotations"].(map[interface{}]interface{}), "deployment.kubernetes.io/revision")
	delete(resource["metadata"].(map[interface{}]interface{}), "selfLink")
	delete(resource["metadata"].(map[interface{}]interface{}), "generation")
	delete(resource["metadata"].(map[interface{}]interface{}), "creationTimestamp")
	delete(resource["metadata"].(map[interface{}]interface{}), "resourceVersion")
	delete(resource, "status")
}

func writeTemplateToInventory(resource map[string]interface{}, kind string, name string, fileInterface fileinterface.FileInterface) string {

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
	paramsFile := fmt.Sprintf("params/%s", name)
	err := ioutil.WriteFile(outFile, outByte, 0766)

	fileInterface.TouchParamsFile(paramsFile)

	clusterContents, err := fileInterface.ReadClusterContents()
	if err != nil {
		log.Fatal("Unable to read the inventory in the current directory. Check the path or run \"applier-cli init\".")
	}
	clusterContents.OpenShiftClusterContent = append(clusterContents.OpenShiftClusterContent, yamlresources.ClusterContentObject{
		Object: name,
		Content: []yamlresources.ClusterContent{{
			Name:     name,
			Template: outFile,
			Params:   paramsFile,
		}},
	})
	err = fileInterface.WriteClusterContents(clusterContents)

	if err != nil {
		log.Fatal("Unable to add the template to the current inventory.")
	} else {
		fmt.Println("Template added to the current inventory.")
	}

	return outFile

}

func writeFileToInventory(resource map[string]interface{}, kind string, name string, fileInterface fileinterface.FileInterface) string {

	outByte, _ := yaml.Marshal(resource)

	outFile := fmt.Sprintf("files/%s.yml", name)
	err := fileInterface.WriteFile(outFile, outByte, 0766)
	if err != nil {
		log.Fatal("Unable to write the file to the current inventory.")
	}

	clusterContents, err := fileInterface.ReadClusterContents()
	if err != nil {
		log.Fatal("Unable to read the inventory in the current directory. Check the path or run \"applier-cli init\".")
	}

	clusterContents.OpenShiftClusterContent = append(clusterContents.OpenShiftClusterContent, yamlresources.ClusterContentObject{
		Object: name,
		Content: []yamlresources.ClusterContent{{
			Name: name,
			File: outFile,
		}},
	})
	fileInterface.WriteClusterContents(clusterContents)

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
